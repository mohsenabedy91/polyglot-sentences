package port

type Event interface {
	Name() string
	Publish(message interface{})
	Consume(message []byte) error
	Register()
}
