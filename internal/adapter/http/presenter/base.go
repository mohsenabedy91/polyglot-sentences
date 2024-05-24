package presenter

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"io"
	"net/http"
)

type Response struct {
	ctx               *gin.Context
	response          map[string]interface{}
	statusCodeMapping map[serviceerror.ErrorMessage]int
	translation       *translation.Translation
	error             serviceerror.Error
}

type Error struct {
	Error string `json:"error" example:"error message"`
}

type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"The email must be a valid email address."`
}

func NewResponse(
	ctx *gin.Context,
	trans *translation.Translation,
	statusCodeMappings ...map[serviceerror.ErrorMessage]int,
) *Response {

	var statusCodeMapping = make(map[serviceerror.ErrorMessage]int)

	if len(statusCodeMappings) > 0 {
		statusCodeMapping = statusCodeMappings[0]
	}

	return &Response{
		ctx:               ctx,
		response:          make(map[string]interface{}),
		statusCodeMapping: statusCodeMapping,
		translation:       trans,
	}
}

func (r *Response) InvalidRequest(err error) *Response {
	if err.Error() == io.EOF.Error() {
		serviceErr := serviceerror.NewServiceError(serviceerror.InvalidRequestBody)
		r.error = serviceErr
		var errorResponse = Error{
			Error: r.translation.Lang(serviceErr.String(), nil, nil),
		}
		r.response["error"] = errorResponse.Error
	}

	return r
}

func (r *Response) Validation(err error) *Response {
	r.InvalidRequest(err)
	if r.error != nil {
		return r
	}

	r.response["validationErrors"] = translate(r.translation, err)
	return r
}

func (r *Response) Payload(data interface{}) *Response {
	r.response["data"] = data
	return r
}

func (r *Response) Meta(data interface{}) *Response {
	r.response["meta"] = data
	return r
}

func (r *Response) Error(err error) *Response {
	serviceErr, ok := err.(serviceerror.Error)
	if !ok {
		serviceErr = serviceerror.NewServiceError(serviceerror.ServerError)
	}
	r.error = serviceErr

	var errorResponse = Error{
		Error: r.translation.Lang(serviceErr.String(), serviceErr.GetAttributes(), nil),
	}
	r.response["error"] = errorResponse.Error
	return r
}

func (r *Response) ErrorMsg(err error) *Response {
	serviceErr, ok := err.(serviceerror.Error)
	if !ok {
		serviceErr = serviceerror.NewServiceError(serviceerror.ServerError)
	}
	var errorResponse = Error{
		Error: r.translation.Lang(serviceErr.String(), serviceErr.GetAttributes(), nil),
	}
	r.response["error"] = errorResponse.Error
	return r
}

func (r *Response) Message(msg string) *Response {
	r.response["message"] = r.translation.Lang(msg, nil, nil)
	return r
}

func (r *Response) Echo(overrideStatusCodes ...int) {

	var statusCode int

	if len(overrideStatusCodes) > 0 {
		statusCode = overrideStatusCodes[0]
	} else {
		if r.error != nil {
			if val, ok := r.statusCodeMapping[r.error.GetMessage()]; !ok {
				statusCode = http.StatusInternalServerError
			} else {
				statusCode = val
			}
		}
	}

	if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
		r.ctx.JSON(statusCode, r.response)
		return
	}

	r.ctx.AbortWithStatusJSON(statusCode, r.response)
}

func translate(t *translation.Translation, err error) (validationError []ValidationError) {

	if ok := errors.As(err, &validator.ValidationErrors{}); !ok {
		validationError = append(validationError, ValidationError{
			Field:   "unknown",
			Message: err.Error(),
		})
		return validationError
	}

	for _, err := range err.(validator.ValidationErrors) {
		validationError = append(validationError, ValidationError{
			Field: err.Field(),
			Message: func() string {
				return t.Lang(
					fmt.Sprintf("validation.%s", err.Tag()),
					map[string]interface{}{
						"attribute": err.Field(),
						err.Tag():   err.Param(),
					},
					nil,
				)
			}(),
		})
	}

	return validationError
}
