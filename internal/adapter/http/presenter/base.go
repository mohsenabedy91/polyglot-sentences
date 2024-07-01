package presenter

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
)

type Response struct {
	ctx               *gin.Context
	response          map[string]interface{}
	statusCodeMapping map[serviceerror.ErrorMessage]int
	translation       *translation.Translation
	serviceError      serviceerror.Error
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
		serviceErr := serviceerror.New(serviceerror.InvalidRequestBody)
		r.serviceError = serviceErr
		var errorResponse = Error{
			Error: r.translation.Lang(serviceErr.Error(), nil, nil),
		}
		r.response["error"] = errorResponse.Error
	}

	return r
}

func (r *Response) Validation(err error) *Response {
	r.InvalidRequest(err)
	if r.serviceError != nil {
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
	var serviceErr serviceerror.Error
	if ok := errors.As(err, &serviceErr); !ok {
		statusErr := status.Convert(err)
		msg, _ := statusErr.WithDetails()
		serviceErr = serviceerror.New(serviceerror.ErrorMessage(msg.Message()))
	}
	r.serviceError = serviceErr

	var errorResponse = Error{
		Error: r.translation.Lang(serviceErr.Error(), serviceErr.GetAttributes(), nil),
	}
	r.response["error"] = errorResponse.Error
	return r
}

func (r *Response) ErrorMsg(err error) *Response {
	r.response["error"] = err.Error()
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
		if r.serviceError != nil {
			if val, ok := r.statusCodeMapping[r.serviceError.GetErrorMessage()]; !ok {
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

func translate(trans *translation.Translation, err error) (validationErrors []ValidationError) {

	if ok := errors.As(err, &validator.ValidationErrors{}); !ok {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "unknown",
			Message: err.Error(),
		})
		return validationErrors
	}

	for _, validationErr := range err.(validator.ValidationErrors) {

		attribute := trans.Lang(fmt.Sprintf("attributes.%s", validationErr.Field()), nil, nil)

		validationError := ValidationError{
			Field: validationErr.Field(),
			Message: func() string {
				return trans.Lang(
					fmt.Sprintf("validation.%s", validationErr.Tag()),
					map[string]interface{}{
						"attribute":         attribute,
						validationErr.Tag(): validationErr.Param(),
					},
					nil,
				)
			}(),
		}

		validationErrors = append(validationErrors, validationError)
	}

	return validationErrors
}
