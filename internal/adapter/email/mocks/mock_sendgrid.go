package mocks

import (
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/mock"
)

type MockSendGrid struct {
	mock.Mock
}

func (r *MockSendGrid) Send(to string, name string, subject string, body string) error {
	args := r.Called(to, name, subject, body)
	return args.Error(0)
}

type MockClient struct {
	mock.Mock
}

func (r *MockClient) Send(email *mail.SGMailV3) (*rest.Response, error) {
	args := r.Called(email)
	return args.Get(0).(*rest.Response), args.Error(1)
}
