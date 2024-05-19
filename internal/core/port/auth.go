package port

import "github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"

type AuthService interface {
	GenerateToken(userUUID string) (*string, serviceerror.Error)
}
