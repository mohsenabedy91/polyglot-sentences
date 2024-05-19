package serviceerror

type Error interface {
	GetMessage() ErrorMessage
	GetAttributes() map[string]interface{}
	String() string
}

type ServiceError struct {
	message    ErrorMessage
	attributes map[string]interface{}
}

func NewServiceError(msg ErrorMessage, attrs ...map[string]interface{}) *ServiceError {

	var attributes map[string]interface{}

	if len(attrs) > 0 {
		attributes = attrs[0]
	}

	return &ServiceError{
		message:    msg,
		attributes: attributes,
	}
}

func (r ServiceError) GetMessage() ErrorMessage {
	return r.message
}

func (r ServiceError) GetAttributes() map[string]interface{} {
	return r.attributes
}

func (r ServiceError) String() string {
	return string(r.message)
}
