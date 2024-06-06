package messagebroker

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
)

type DriverInterface interface {
	Close()
	Produce(name string, message interface{}, delaySeconds int64) error
	RegisterConsumer(name string, callback func(message []byte) error) error
}

type Queue struct {
	Log    logger.Logger
	Config config.Config
	Driver DriverInterface
}

func NewQueue(log logger.Logger, config config.Config) *Queue {
	return &Queue{
		Log:    log,
		Config: config,
	}
}

func RegisterAllQueues(events ...port.EventQueue) {
	for _, event := range events {
		event.Register()
	}
}
