package email

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridClient interface {
	Send(email *mail.SGMailV3) (*rest.Response, error)
}

type SendGrid struct {
	log    logger.Logger
	conf   config.SendGrid
	client SendGridClient
}

func NewSender(log logger.Logger, conf config.SendGrid) *SendGrid {
	client := sendgrid.NewSendClient(conf.Key)
	return &SendGrid{
		log:    log,
		conf:   conf,
		client: client,
	}
}

func (r *SendGrid) SetClient(client SendGridClient) {
	r.client = client
}

func (r *SendGrid) Send(to string, name string, subject string, body string) error {
	message := mail.NewV3Mail()
	message.SetFrom(mail.NewEmail(r.conf.Name, r.conf.Address))
	message.Subject = subject

	personalization := mail.NewPersonalization()
	personalization.AddTos(mail.NewEmail(name, to))

	content := mail.NewContent("text/html", body)
	message.AddContent(content)

	message.AddPersonalizations(personalization)

	response, err := r.client.Send(message)

	extra := map[logger.ExtraKey]interface{}{
		"To":       to,
		"Address":  r.conf.Address,
		"Body":     body,
		"Response": response,
	}
	if err != nil {
		r.log.Error(logger.SendGrid, logger.SendGridSendEmail, err.Error(), extra)
		return serviceerror.New(serviceerror.FailedSendEmail)
	}
	r.log.Info(logger.SendGrid, logger.SendGridSendEmail, "Email sent successfully", extra)
	return nil
}
