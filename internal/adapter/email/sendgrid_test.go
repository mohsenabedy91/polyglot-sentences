package email_test

import (
	"github.com/bxcodec/faker/v4"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/email"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/sendgrid/rest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSend_Successful(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	mockClient := new(email.MockClient)

	conf := config.SendGrid{
		Name:    "Test",
		Address: "test@google.com",
		Key:     "SG.fake_key",
	}

	sender := email.NewSender(mockLogger, conf)
	sender.SetClient(mockClient)

	mockResponse := &rest.Response{
		StatusCode: 202,
		Body:       "Accepted",
	}
	mockClient.On("Send", mock.Anything).Return(mockResponse, nil)

	mockLogger.On("Info", logger.SendGrid, logger.SendGridSendEmail, "Email sent successfully", mock.Anything).Return()

	to := faker.Email()
	name := faker.Name()
	subject := faker.Word()
	body := faker.Sentence()

	err := sender.Send(to, name, subject, body)
	require.NoError(t, err)

	mockLogger.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestSend_Failure(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	mockClient := new(email.MockClient)

	conf := config.SendGrid{
		Name:    "Test",
		Address: "test@google.com",
		Key:     "SG.fake_key",
	}

	sender := email.NewSender(mockLogger, conf)
	sender.SetClient(mockClient)

	mockError := serviceerror.New(serviceerror.FailedSendEmail)

	mockClient.On("Send", mock.Anything).Return(&rest.Response{}, mockError)

	mockLogger.On("Error", logger.SendGrid, logger.SendGridSendEmail, mockError.Error(), mock.Anything).Return()

	to := faker.Email()
	name := faker.Name()
	subject := faker.Word()
	body := faker.Sentence()

	err := sender.Send(to, name, subject, body)
	require.Error(t, err)
	require.Equal(t, mockError, err)

	mockLogger.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestSend_LoggingSuccess(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	mockClient := new(email.MockClient)

	conf := config.SendGrid{
		Name:    "Test",
		Address: "test@google.com",
		Key:     "SG.fake_key",
	}

	sender := email.NewSender(mockLogger, conf)
	sender.SetClient(mockClient)

	mockResponse := &rest.Response{
		StatusCode: 202,
		Body:       "Accepted",
	}
	mockClient.On("Send", mock.Anything).Return(mockResponse, nil)

	to := faker.Email()
	name := faker.Name()
	subject := faker.Word()
	body := faker.Sentence()

	extra := map[logger.ExtraKey]interface{}{
		"To":      to,
		"Address": conf.Address,
		"Body":    body,
		"Response": &rest.Response{
			StatusCode: 202,
			Body:       "Accepted",
		},
	}
	mockLogger.On("Info", logger.SendGrid, logger.SendGridSendEmail, "Email sent successfully", extra).Return()

	err := sender.Send(to, name, subject, body)
	require.NoError(t, err)

	mockLogger.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestSend_LoggingFailure(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	mockClient := new(email.MockClient)

	conf := config.SendGrid{
		Name:    "Test",
		Address: "test@google.com",
		Key:     "SG.fake_key",
	}

	sender := email.NewSender(mockLogger, conf)
	sender.SetClient(mockClient)

	mockError := serviceerror.New(serviceerror.FailedSendEmail)

	mockResponse := &rest.Response{}
	mockClient.On("Send", mock.Anything).Return(mockResponse, mockError)

	to := faker.Email()
	name := faker.Name()
	subject := faker.Word()
	body := faker.Sentence()

	extra := map[logger.ExtraKey]interface{}{
		"To":       to,
		"Address":  conf.Address,
		"Body":     body,
		"Response": mockResponse,
	}
	mockLogger.On("Error", logger.SendGrid, logger.SendGridSendEmail, mockError.Error(), extra).Return()

	err := sender.Send(to, name, subject, body)
	require.Error(t, err)
	require.Equal(t, mockError, err)

	mockLogger.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}
