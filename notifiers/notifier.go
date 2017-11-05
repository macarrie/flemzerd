package notifier

type Notifier interface {
	Send(title, content string) error
}
