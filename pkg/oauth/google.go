package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
)

type GoogleUserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	AvatarURL     string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

type GoogleService interface {
	UserGoogleInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error)
}

type OAuth struct {
	log    logger.Logger
	config config.Oauth
}

func New(log logger.Logger, config config.Oauth) *OAuth {
	return &OAuth{config: config, log: log}
}

func (r *OAuth) UserGoogleInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	var oauthCfg = oauth2.Config{
		ClientID:     r.config.Google.ClientId,
		ClientSecret: r.config.Google.ClientSecret,
		RedirectURL:  r.config.Google.CallbackURL,
		Endpoint:     google.Endpoint,
	}

	client := oauthCfg.Client(
		ctx,
		&oauth2.Token{
			AccessToken: accessToken,
		},
	)

	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		r.log.Error(logger.Google, logger.ExternalService, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	if response.StatusCode == http.StatusOK {
		var userInfo *GoogleUserInfo

		body, _ := io.ReadAll(response.Body)
		if err = json.Unmarshal(body, &userInfo); err != nil {

			r.log.Error(logger.Google, logger.ExternalService, fmt.Sprintf("Error unmarshalling message, error: %v", err), nil)
			return nil, serviceerror.NewServerError()
		}

		return userInfo, nil
	}

	r.log.Error(logger.Google, logger.ExternalService, fmt.Sprintf("error: Unable to fetch user info. Status Code: %s", response.Status), nil)
	return nil, serviceerror.NewServerError()
}
