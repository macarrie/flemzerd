build:
	go build

test:
	go test ./...

fmt:
	gofmt -l -w .

install:
	go get "github.com/ogier/pflag"
	go get "github.com/spf13/viper"
	go get "github.com/sirupsen/logrus"
