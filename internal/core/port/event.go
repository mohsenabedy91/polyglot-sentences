package port

type EventQueue interface {
	Name() string
	Publish(msg interface{})
	Consume(message []byte) error
	Register()
}
