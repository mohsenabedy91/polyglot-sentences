//go:build !test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const RequestTimeout = 30 * time.Second

var AuthorizeURL string
var KongBaseURL string

func main() {
	fmt.Println("Starting...")

	configProvider := &config.Config{}
	conf := configProvider.GetConfig()
	KongBaseURL = conf.Kong.APIBaseUrl
	AuthorizeURL = conf.Kong.AuthorizeUrl

	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	kongRoutes, err := fetchKongRoutes(ctx)
	if err != nil {
		logErrorAndExit("Error fetching Kong routes:", err)
	}

	if len(os.Args) > 1 {
		handleArgs(os.Args[1:], kongRoutes, ctx)
		return
	}
	// TODO add services if doesn't exists on Kong
	kongServices, err := fetchKongServices(ctx)
	if err != nil {
		logErrorAndExit("Error fetching Kong services:", err)
	}

	data, err := os.ReadFile(conf.Kong.SwaggerFilePath)
	if err != nil {
		logErrorAndExit("Error reading swagger file:", err)
	}

	swaggerRoutes, err := parseSwaggerRoutes(data)
	if err != nil {
		logErrorAndExit("Error parsing swagger routes:", err)
	}

	// TODO update/sync routes and plugins on Kong
	newRoutes := filterNewRoutes(swaggerRoutes, kongRoutes, kongServices)
	processNewRoutes(newRoutes, ctx)
	fmt.Println("===============================================")
	fmt.Println("========> Added routes successfully <==========")
	fmt.Println("===============================================")
}

func logErrorAndExit(message string, err error) {
	fmt.Println(message, err)
	os.Exit(1)
}

func handleArgs(args []string, kongRoutes []Route, ctx context.Context) {
	for _, arg := range args {
		if arg == "--remove-all" {
			deleteKongRoutes(ctx, kongRoutes)
			fmt.Println("===============================================")
			fmt.Println("========> Deleted routes successfully <========")
			fmt.Println("===============================================")
		}
	}
}

func processNewRoutes(routes []Route, ctx context.Context) {
	var wg sync.WaitGroup
	for _, route := range routes {
		wg.Add(1)
		go func(route Route) {
			defer wg.Done()
			createdRoute, err := createKongRoute(ctx, route)
			if err != nil {
				fmt.Println("Error creating Kong route:", err)
				return
			}

			if err := addRequestTransformerPlugin(ctx, createdRoute); err != nil {
				fmt.Println("Error adding request transformer plugin:", err)
				return
			}

			if route.AuthBearer {
				if err := addAuthorizePlugin(ctx, createdRoute); err != nil {
					fmt.Println("Error adding authorize plugin:", err)
					return
				}
			}
		}(route)
	}
	wg.Wait()
}

// Service represents a Kong service
type Service struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// ServicesResponse represents a response from Kong for fetching services
type ServicesResponse struct {
	Data []Service `json:"data"`
}

// Route represents a Kong route
type Route struct {
	ID            string
	Service       Service
	Name          string
	Method        string
	Path          string
	AuthBearer    bool
	Permissions   []string
	Tags          []string
	RegexPriority int
}

// Security represents security settings in Swagger
type Security struct {
	AuthBearer []string `json:"AuthBearer,omitempty"`
}

// ExtraKey represents extra values in Swagger
type ExtraKey struct {
	Service *string `json:"service,omitempty"`
}

// Operation represents an operation in Swagger
type Operation struct {
	ID       *string    `json:"operationId,omitempty"`
	Tags     []string   `json:"tags,omitempty"`
	ExtraKey *ExtraKey  `json:"x-kong,omitempty"`
	Security []Security `json:"security,omitempty"`
}

// Paths maps paths to their operations in Swagger
type Paths map[string]map[string]Operation

// API represents the Swagger API
type API struct {
	Paths Paths `json:"paths"`
}

func filterNewRoutes(swaggerRoutes, existingRoutes []Route, services *ServicesResponse) []Route {
	for i, kongRoute := range existingRoutes {
		for _, service := range services.Data {
			if kongRoute.Service.ID != "" && kongRoute.Service.ID == service.ID {
				existingRoutes[i].Service = service
			}
		}
	}

	var result []Route
	for _, route := range swaggerRoutes {
		var exists bool
		for _, kongRoute := range existingRoutes {
			if route.Service.Name == kongRoute.Service.Name && route.Name == kongRoute.Name {
				exists = true
				break
			}
		}
		if !exists {
			for _, service := range services.Data {
				if route.Service.Name != "" && route.Service.Name == service.Name {
					route.Service = service
				}
			}
			result = append(result, route)
		}
	}

	return result
}

func parseSwaggerRoutes(file []byte) ([]Route, error) {
	var api API
	if err := json.Unmarshal(file, &api); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	routeMap := make(map[string]*Route)
	for path, methods := range api.Paths {
		for method, operation := range methods {
			if !isValidOperation(operation) {
				continue
			}

			var permissions []string
			var isAuthBearer bool
			for _, security := range operation.Security {
				if security.AuthBearer != nil {
					isAuthBearer = true
					permissions = append(permissions, security.AuthBearer...)
				}
			}

			newRoute := &Route{
				Name:          *operation.ID,
				Service:       Service{Name: *operation.ExtraKey.Service},
				Method:        strings.ToUpper(method),
				Path:          path,
				Permissions:   permissions,
				Tags:          operation.Tags,
				AuthBearer:    isAuthBearer,
				RegexPriority: 10,
			}

			for existingRouteName, route := range routeMap {
				newRoutePath := strings.Split(newRoute.Path, "/")
				existingRoutePath := strings.Split(route.Path, "/")

				if isSubRoute(newRoutePath, existingRoutePath) {

					diffPath := len(newRoutePath) - len(existingRoutePath)

					if diffPath > 0 {
						newRoute.RegexPriority += 10

					} else if diffPath == 0 {

						if strings.HasPrefix(newRoutePath[len(newRoutePath)-1], "{") {
							routeMap[existingRouteName].RegexPriority += 10

						} else if strings.HasPrefix(existingRoutePath[len(existingRoutePath)-1], "{") {
							newRoute.RegexPriority += 10
						}

					} else {
						routeMap[existingRouteName].RegexPriority += 10
					}
				}
			}

			routeMap[newRoute.Name] = newRoute
		}
	}

	routes := make([]Route, 0, len(routeMap))
	for _, route := range routeMap {
		routes = append(routes, *route)
	}

	return routes, nil
}

func isSubRoute(parts1, parts2 []string) bool {
	if len(parts1) == 0 || len(parts2) == 0 {
		return false
	}

	minLen := len(parts1)
	if len(parts2) < minLen {
		minLen = len(parts2)
	}

	if len(parts1) == len(parts2) {
		minLen -= 1
	}

	for i := 0; i < minLen; i++ {
		if parts1[i] != parts2[i] {
			return false
		}
	}

	return true
}

func preparePath(path string) string {
	paramRegex := regexp.MustCompile(`{([^{}]+)}`)
	replacedPath := paramRegex.ReplaceAllStringFunc(path, func(match string) string {
		paramName := strings.Trim(match, "{}")
		return fmt.Sprintf("(?<%s>\\S+)", paramName)
	})
	return strings.Trim(replacedPath, "/")
}

func preparePluginPath(originalURL string) string {
	paramRegex := regexp.MustCompile(`{([^{}]+)}`)
	replacedURL := paramRegex.ReplaceAllStringFunc(originalURL, func(match string) string {
		paramName := strings.Trim(match, "{}")
		return fmt.Sprintf("$(uri_captures['%s'])", paramName)
	})
	return strings.Trim(replacedURL, "/")
}

func isValidOperation(operation Operation) bool {
	return operation.ID != nil && operation.ExtraKey != nil && operation.ExtraKey.Service != nil
}

func fetchKongRoutes(ctx context.Context) ([]Route, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, KongBaseURL+"/routes", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %v", err)
	}

	var routesResponse struct {
		Data []struct {
			ID      string  `json:"id"`
			Name    string  `json:"name"`
			Service Service `json:"service"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &routesResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling routes: %v", err)
	}

	var routes []Route
	for _, route := range routesResponse.Data {
		routes = append(routes, Route{
			ID:   route.ID,
			Name: route.Name,
			Service: Service{
				ID: route.Service.ID,
			},
		})
	}

	return routes, nil
}

func fetchKongServices(ctx context.Context) (*ServicesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, KongBaseURL+"/services", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var servicesResponse ServicesResponse
	if err := json.Unmarshal(body, &servicesResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &servicesResponse, nil
}

func createKongRoute(ctx context.Context, route Route) (*Route, error) {
	url := fmt.Sprintf("%s/services/%s/routes", KongBaseURL, route.Service.Name)
	path := fmt.Sprintf("/%s", preparePath(route.Path))
	if strings.Contains(path, "S+)") {
		path = "~" + path
	}

	type request struct {
		Name          string   `json:"name"`
		Paths         []string `json:"paths"`
		Methods       []string `json:"methods"`
		Tags          []string `json:"tags"`
		RegexPriority int      `json:"regex_priority"`
	}

	payload := request{
		Name:          route.Name,
		Paths:         []string{path},
		Methods:       []string{route.Method},
		Tags:          route.Tags,
		RegexPriority: route.RegexPriority,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status code: %s - %s", resp.Status, body)
	}

	var createdRouteResponse struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &createdRouteResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	route.ID = createdRouteResponse.ID
	return &route, nil
}

func addRequestTransformerPlugin(ctx context.Context, route *Route) error {
	type Replace struct {
		URI string `json:"uri"`
	}
	type Config struct {
		Replace    Replace `json:"replace"`
		HTTPMethod string  `json:"http_method"`
	}

	type Route struct {
		ID string `json:"id"`
	}

	type request struct {
		Name   string `json:"name"`
		Config Config `json:"config"`
		Route  Route  `json:"route"`
	}
	payload := request{
		Name: "request-transformer",
		Config: Config{
			Replace: Replace{
				URI: fmt.Sprintf("/%s", preparePluginPath(route.Path)),
			},
			HTTPMethod: route.Method,
		},
		Route: Route{
			ID: route.ID,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling request transformer: %v", err)
	}

	url := KongBaseURL + "/plugins"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("status code: %s - %s", resp.Status, body)
	}

	return nil
}

func addAuthorizePlugin(ctx context.Context, route *Route) error {
	type Config struct {
		RequiredPermissions []string `json:"required_permissions"`
		AuthorizeURL        string   `json:"authorize_url"`
	}

	type Route struct {
		ID string `json:"id"`
	}

	type request struct {
		Name   string `json:"name"`
		Config Config `json:"config"`
		Route  Route  `json:"route"`
	}

	payload := request{
		Name: "ps-authorize",
		Config: Config{
			AuthorizeURL: AuthorizeURL,
			RequiredPermissions: func() []string {
				if len(route.Permissions) > 0 {
					return route.Permissions
				}
				return []string{"NONE"}
			}(),
		},
		Route: Route{
			ID: route.ID,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling request transformer: %v", err)
	}

	url := KongBaseURL + "/plugins"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("status code: %s - %s", resp.Status, body)
	}

	return nil
}

func deleteKongRoutes(ctx context.Context, routes []Route) {
	var wg sync.WaitGroup
	for _, route := range routes {
		wg.Add(1)
		go func(route Route) {
			defer wg.Done()
			url := KongBaseURL + "/routes/" + route.ID
			req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
			if err != nil {
				fmt.Printf("Failed to create request for route %s: %v\n", route.ID, err)
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("Failed to delete route %s: %v\n", route.ID, err)
				return
			}

			defer func(Body io.ReadCloser) {
				if err = Body.Close(); err != nil {

				}
			}(resp.Body)

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Failed to read response body for route %s: %v\n", route.ID, err)
				return
			}
			if resp.StatusCode != http.StatusNoContent {
				fmt.Printf("Failed to delete route %s, Status Code: %d, Body: %s\n", route.ID, resp.StatusCode, body)
			} else {
				fmt.Printf("Route deleted: %s\n", route.ID)
			}
		}(route)
	}
	wg.Wait()
}
