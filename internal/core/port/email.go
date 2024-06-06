package port

type EmailSender interface {
	Send(to string, name string, subject string, body string) error
}
