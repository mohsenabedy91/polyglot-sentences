package oauth_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/oauth"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockHTTPClientProvider struct {
	mock.Mock
}

func (r *MockHTTPClientProvider) GetClient(ctx context.Context, token *oauth2.Token) *http.Client {
	args := r.Called(ctx, token)
	return args.Get(0).(*http.Client)
}

type RoundTripper struct {
	handler http.Handler
	err     error
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if r.err != nil {
		return nil, r.err
	}

	recorder := httptest.NewRecorder()
	r.handler.ServeHTTP(recorder, req)
	resp := recorder.Result()

	if resp.Body != nil {
		resp.Body = &errorReadCloser{ReadCloser: resp.Body}
	}

	return resp, nil
}

type errorReadCloser struct {
	io.ReadCloser
}

func (r *errorReadCloser) Close() error {
	return errors.New("close error")
}

func TestUserGoogleInfo(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		expectedError  error
		expectedEmail  string
		closeBodyError bool
		httpClientErr  error
	}{
		{
			name: "Successful user info retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(oauth.GoogleUserInfo{
					Email: "user@example.com",
				})
			},
			expectedError: nil,
			expectedEmail: "user@example.com",
		},
		{
			name: "Unsuccessful user info retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedError: serviceerror.NewServerError(),
		},
		{
			name: "Error during HTTP request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectedError: serviceerror.NewServerError(),
		},
		{
			name: "Invalid JSON response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("invalid json"))
			},
			expectedError: serviceerror.NewServerError(),
		},
		{
			name: "Error closing response body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(oauth.GoogleUserInfo{
					Email: "user@example.com",
				})
			},
			expectedError:  nil,
			expectedEmail:  "user@example.com",
			closeBodyError: true,
		},
		{
			name:          "HTTP client error",
			handler:       nil,
			expectedError: serviceerror.NewServerError(),
			httpClientErr: errors.New("http client error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockLogger := new(logger.MockLogger)
			mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			mockClientProvider := new(MockHTTPClientProvider)
			if test.httpClientErr != nil {
				mockClientProvider.On("GetClient", mock.Anything, mock.Anything).Return(&http.Client{
					Transport: &RoundTripper{handler: test.handler, err: test.httpClientErr},
				})
			} else {
				mockClientProvider.On("GetClient", mock.Anything, mock.Anything).Return(&http.Client{
					Transport: &RoundTripper{handler: test.handler},
				})
			}

			conf := config.Oauth{
				Google: config.Google{
					ClientId:     "client_id",
					ClientSecret: "client_secret",
					CallbackURL:  "callback_url",
				},
			}

			oauthService := oauth.New(mockLogger, conf, mockClientProvider)
			ctx := context.Background()
			accessToken := "access_token"

			if test.closeBodyError {
				mockClientProvider.On("GetClient", mock.Anything, mock.Anything).Return(&http.Client{
					Transport: &RoundTripper{
						handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							recorder := httptest.NewRecorder()
							test.handler(recorder, r)
							resp := recorder.Result()
							resp.Body = &errorReadCloser{ReadCloser: resp.Body}
							w.WriteHeader(resp.StatusCode)

							write, err := w.Write(recorder.Body.Bytes())
							require.NoError(t, err)
							require.Greater(t, write, 0)
						}),
					},
				})
			}

			userInfo, err := oauthService.UserGoogleInfo(ctx, accessToken)
			if test.expectedError != nil {
				require.Error(t, err)
				require.Nil(t, userInfo)
			} else {
				require.NoError(t, err)
				require.NotNil(t, userInfo)
				require.Equal(t, test.expectedEmail, userInfo.Email)
			}

			mockLogger.AssertExpectations(t)
			mockClientProvider.AssertExpectations(t)
		})
	}
}

func TestClientProvider_GetClient(t *testing.T) {
	conf := config.Oauth{
		Google: config.Google{
			ClientId:     "mock-client-id",
			ClientSecret: "mock-client-secret",
			CallbackURL:  "http://localhost/callback",
		},
	}

	clientProvider := oauth.ClientProvider{
		Conf: conf,
	}

	ctx := context.Background()
	token := &oauth2.Token{
		AccessToken: "mock-access-token",
	}

	client := clientProvider.GetClient(ctx, token)

	require.NotNil(t, client, "Expected non-nil HTTP client")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	tokenSource := oauth2.StaticTokenSource(token)
	oauthClient := oauth2.NewClient(ctx, tokenSource)
	oauthReq, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	oauthResp, err := oauthClient.Do(oauthReq)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, oauthResp.StatusCode)
}
