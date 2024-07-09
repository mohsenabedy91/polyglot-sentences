package port

type Driver interface {
	Close()
	Produce(name string, message interface{}, delaySeconds int64) error
	RegisterConsumer(name string, callback func(message []byte) error) error
}
