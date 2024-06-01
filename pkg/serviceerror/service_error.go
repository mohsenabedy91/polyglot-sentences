package serviceerror

import (
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
)

type Error interface {
	GetErrorMessage() ErrorMessage
	GetAttributes() map[string]interface{}
	Error() string
}

type ServiceError struct {
	message    ErrorMessage
	attributes map[string]interface{}
	trans      *translation.Translation
}

func New(msg ErrorMessage, attrs ...map[string]interface{}) *ServiceError {

	var attributes map[string]interface{}

	if len(attrs) > 0 {
		attributes = attrs[0]
	}

	return &ServiceError{
		message:    msg,
		attributes: attributes,
	}
}

func (r ServiceError) GetErrorMessage() ErrorMessage {
	return r.message
}

func (r ServiceError) GetAttributes() map[string]interface{} {
	return r.attributes
}

func (r ServiceError) Error() string {
	return string(r.message)
}

func NewServerError() *ServiceError {
	return New(ServerError)
}
