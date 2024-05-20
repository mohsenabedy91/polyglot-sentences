package authservice

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"time"
)

type JWTService struct {
	log logger.Logger
	cfg config.Jwt
}

func New(log logger.Logger, cfg config.Jwt) *JWTService {
	return &JWTService{
		log: log,
		cfg: cfg,
	}
}

func (r JWTService) GenerateToken(userUUID string) (*string, error) {
	mapClaims := jwt.MapClaims{}
	mapClaims[config.AuthUserUUIDKey] = userUUID

	// please never change "exp" key
	mapClaims["exp"] = int(time.Now().Add(r.cfg.AccessTokenExpireDay * (24 * time.Hour)).Unix())

	jwtString, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		mapClaims,
	).SignedString([]byte(r.cfg.AccessTokenSecret))

	if err != nil {
		r.log.Error(logger.JWT, logger.JWTGenerate, err.Error(), nil)
		return nil, serviceerror.NewServiceError(serviceerror.ServerError)
	}

	return &jwtString, nil
}
