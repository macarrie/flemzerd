package notifier

type Notifier interface {
	IsAlive() error
	Send(title, content string) error
}
