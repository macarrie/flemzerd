package healthcheck

import (
	"testing"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"

	log "github.com/macarrie/flemzerd/logging"
)

func init() {
	log.Setup(true)

	db.DbPath = "/tmp/flemzerd.db"
	db.Load()
	db.ResetDb()

	// go test makes a cd into package directory when testing. We must go up by one level to load our testdata
	configuration.UseFile("../testdata/test_config.toml")
	configuration.Load()
}

func TestPoll(t *testing.T) {
	CheckHealth()
}

func TestRunAndStop(t *testing.T) {
	configuration.Load()

	configuration.Config.System.CheckInterval = 1

	go Run()
	time.Sleep(10 * time.Second)
	Stop()
}
