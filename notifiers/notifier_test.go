package notifier

type MockNotifier struct{}

var mockNotificationCounter int
var notifierInitialized bool

func (n MockNotifier) Setup(authCredentials map[string]string) {
	return
}

func (n MockNotifier) Init() error {
	notifierInitialized = true
	return nil
}

func (n MockNotifier) Send(title, content string) error {
	mockNotificationCounter += 1
	return nil
}
