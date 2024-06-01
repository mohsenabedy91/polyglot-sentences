package messagebroker

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
)

type DriverInterface interface {
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

type JobQueue interface {
	RegisterQueue()
	Name() string
	Consume(message []byte) (err error)
}

func RegisterAllQueues(jobs ...JobQueue) {
	for _, r := range jobs {
		r.RegisterQueue()
	}
}
