package email

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGrid struct {
	log logger.Logger
	cfg config.SendGrid
}

func NewSender(log logger.Logger, cfg config.SendGrid) *SendGrid {
	return &SendGrid{
		log: log,
		cfg: cfg,
	}
}

func (s *SendGrid) Send(to string, name string, subject string, body string) error {
	message := mail.NewV3Mail()
	message.SetFrom(mail.NewEmail(s.cfg.Name, s.cfg.Address))
	message.Subject = subject

	personalization := mail.NewPersonalization()
	personalization.AddTos(mail.NewEmail(name, to))

	content := mail.NewContent("text/html", body)
	message.AddContent(content)

	message.AddPersonalizations(personalization)

	client := sendgrid.NewSendClient(s.cfg.Key)
	response, err := client.Send(message)

	extra := map[logger.ExtraKey]interface{}{
		"To":       to,
		"Address":  s.cfg.Address,
		"Body":     body,
		"Response": response,
	}
	if err != nil {
		s.log.Error(logger.SendGrid, logger.SendGridSendEmail, err.Error(), extra)
		return err
	}
	s.log.Info(logger.SendGrid, logger.SendGridSendEmail, "Email sent successfully", extra)
	return nil
}
