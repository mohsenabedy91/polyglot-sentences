package messagebroker

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
)

type Queue struct {
	Log    logger.Logger
	Config config.Config
	Driver port.Driver
}

func NewQueue(log logger.Logger, config config.Config) *Queue {
	return &Queue{
		Log:    log,
		Config: config,
	}
}

func RegisterEvents(events ...port.Event) {
	for _, event := range events {
		event.Register()
	}
}
