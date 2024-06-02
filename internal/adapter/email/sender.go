package email

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
)

type Sender interface {
	Send(to string, name string, subject string, body string) error
}

func NewSender(log logger.Logger, cfg config.SendGrid) Sender {
	return newSendGrid(log, cfg)
}
