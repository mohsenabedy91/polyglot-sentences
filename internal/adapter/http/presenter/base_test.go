package presenter_test

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/mocks"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper function to create a gin context for testing
func createTestGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestNewResponse(t *testing.T) {
	ctx, _ := createTestGinContext()
	mockTranslator := new(translation.MockTranslator)
	response := presenter.NewResponse(ctx, mockTranslator, handler.StatusCodeMapping)

	require.NotNil(t, response)
	require.Equal(t, response.GetStatusCodeMapping(), handler.StatusCodeMapping)
	require.Empty(t, response.GetResponse())
}

func TestResponse_InvalidRequest(t *testing.T) {
	mockTranslator := new(translation.MockTranslator)
	mockTranslator.On("Lang", "errors.invalidRequestBody", mock.Anything, (*string)(nil)).Return("Invalid request body")

	ctx, _ := createTestGinContext()
	response := presenter.NewResponse(ctx, mockTranslator, handler.StatusCodeMapping)

	err := errors.New("EOF")
	response.InvalidRequest(err).Echo()

	require.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
	require.Equal(t, map[string]interface{}{"error": "Invalid request body"}, response.GetResponse())
	mockTranslator.AssertExpectations(t)
}

func TestResponse_InvalidRequest_Validation(t *testing.T) {
	ctx, _ := createTestGinContext()

	mockTranslator := new(translation.MockTranslator)
	mockTranslator.On("Lang", "errors.invalidRequestBody", mock.Anything, (*string)(nil)).Return("Invalid request body")

	response := presenter.NewResponse(ctx, mockTranslator, handler.StatusCodeMapping)

	err := errors.New("EOF")
	response.Validation(err).Echo()

	require.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
	require.Equal(t, map[string]interface{}{"error": "Invalid request body"}, response.GetResponse())
	mockTranslator.AssertExpectations(t)
}

func TestResponse_Validation(t *testing.T) {
	ctx, _ := createTestGinContext()

	mockTranslator := new(translation.MockTranslator)
	response := presenter.NewResponse(ctx, mockTranslator)

	err := errors.New("validation error")
	response.Validation(err).Echo(http.StatusUnprocessableEntity)

	require.Equal(t, http.StatusUnprocessableEntity, ctx.Writer.Status())
	require.Equal(t, map[string]interface{}{
		"validationErrors": []presenter.ValidationError{
			{
				Field:   "unknown",
				Message: "validation error",
			},
		},
	}, response.GetResponse())
}

func TestResponse_Payload(t *testing.T) {
	ctx, _ := createTestGinContext()
	mockTranslator := new(translation.MockTranslator)
	response := presenter.NewResponse(ctx, mockTranslator, handler.StatusCodeMapping)

	data := map[string]interface{}{"key": "value"}
	response.Payload(data)

	require.Equal(t, http.StatusOK, ctx.Writer.Status())
	require.Equal(t, map[string]interface{}{"data": data}, response.GetResponse())
}

func TestResponse_Meta(t *testing.T) {
	ctx, _ := createTestGinContext()
	mockTranslator := new(translation.MockTranslator)
	response := presenter.NewResponse(ctx, mockTranslator, handler.StatusCodeMapping)

	meta := map[string]interface{}{"meta_key": "meta_value"}
	response.Meta(meta)

	require.Equal(t, map[string]interface{}{"meta": meta}, response.GetResponse())
}

func TestResponse_Error(t *testing.T) {
	ctx, _ := createTestGinContext()

	mockTranslator := new(translation.MockTranslator)
	mockTranslator.On("Lang", "errors.serverError", mock.Anything, (*string)(nil)).Return("An error occurred")

	response := presenter.NewResponse(ctx, mockTranslator, handler.StatusCodeMapping)

	err := serviceerror.New("errors.serverError")
	response.Error(err).Echo()

	require.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	require.Equal(t, map[string]interface{}{"error": "An error occurred"}, response.GetResponse())
	mockTranslator.AssertExpectations(t)
}

func TestResponse_ErrorMsg(t *testing.T) {
	ctx, _ := createTestGinContext()
	mockTranslator := new(translation.MockTranslator)
	response := presenter.NewResponse(ctx, mockTranslator, handler.StatusCodeMapping)

	err := errors.New("error message")
	response.ErrorMsg(err)

	require.Equal(t, map[string]interface{}{"error": "error message"}, response.GetResponse())
}

func TestResponse_Message(t *testing.T) {
	ctx, _ := createTestGinContext()

	mockTranslator := new(translation.MockTranslator)
	mockTranslator.On("Lang", "A message", mock.Anything, (*string)(nil)).Return("A message")

	response := presenter.NewResponse(ctx, mockTranslator, handler.StatusCodeMapping)

	response.Message("A message")

	require.Equal(t, map[string]interface{}{"message": "A message"}, response.GetResponse())
	mockTranslator.AssertExpectations(t)
}

func TestResponse_Echo(t *testing.T) {
	ctx, w := createTestGinContext()
	mockTranslator := new(translation.MockTranslator)
	response := presenter.NewResponse(ctx, mockTranslator, handler.StatusCodeMapping)

	response.Payload(gin.H{"key": "value"})
	response.Echo(http.StatusOK)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"data":{"key":"value"}}`, w.Body.String())
}

func TestResponse_Echo_without_status_code(t *testing.T) {
	ctx, w := createTestGinContext()

	mockTranslator := new(translation.MockTranslator)
	mockTranslator.On("Lang", "errors.test", mock.Anything, (*string)(nil)).Return("error message")

	response := presenter.NewResponse(ctx, mockTranslator)

	err := errors.New("errors.test")
	response.Error(err).Echo()

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.Equal(t, map[string]interface{}{"error": "error message"}, response.GetResponse())

	mockTranslator.AssertExpectations(t)
}

func TestResponse_InvalidRequest_InvalidError(t *testing.T) {
	ctx, _ := createTestGinContext()
	mockTranslator := new(translation.MockTranslator)
	response := presenter.NewResponse(ctx, mockTranslator, handler.StatusCodeMapping)

	err := errors.New("some other error")
	response.InvalidRequest(err)

	require.Nil(t, response.GetServiceError())
	require.Empty(t, response.GetResponse())
}

func TestTranslate(t *testing.T) {
	mockTranslator := new(translation.MockTranslator)
	mockTranslator.On("Lang", "attributes.Email", mock.Anything, (*string)(nil)).Return("Email")
	mockTranslator.On("Lang", "validation.required", mock.Anything, (*string)(nil)).Return("The Email field is required.")

	mockFieldError := new(mocks.MockFieldError)
	mockFieldError.On("Field").Return("Email")
	mockFieldError.On("Tag").Return("required")
	mockFieldError.On("Param").Return("")

	errs := validator.ValidationErrors{
		mockFieldError,
	}

	expected := []presenter.ValidationError{
		{
			Field:   "Email",
			Message: "The Email field is required.",
		},
	}

	result := presenter.Translate(mockTranslator, errs)
	require.Equal(t, expected, result)

	mockTranslator.AssertExpectations(t)
	mockFieldError.AssertExpectations(t)
}
