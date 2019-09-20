package healthcheck

import (
	"time"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	"github.com/macarrie/flemzerd/downloaders"
	"github.com/macarrie/flemzerd/helpers/modules"
	"github.com/macarrie/flemzerd/indexers"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/mediacenters"
	"github.com/macarrie/flemzerd/notifiers"
	"github.com/macarrie/flemzerd/providers"
	"golang.org/x/sys/unix"
)

var RunTicker *time.Ticker
var CanDownload bool

func Run() {
	dbErr := db.Load()
	if dbErr != nil {
		log.WithFields(log.Fields{
			"error": dbErr,
		}).Warning("Could not connect to database. Starting daemon without any previous data")
	}

	RunTicker = time.NewTicker(time.Duration(configuration.Config.System.HealthcheckInterval) * time.Minute)
	go func() {
		log.Debug("Starting healthcheck loop")
		for {
			CheckHealth()
			<-RunTicker.C
		}
	}()
}

func Stop() {
	log.Info("Closing DB connection")
}

func CheckHealth() {
	CanDownload = true

	log.Debug("========== Healthcheck loop start ==========")

	if _, err := notifier.Status(); err != nil {
		log.Error("No notifier alive. No notifications will be sent until next polling.")
	}

	if mods, err := provider.Status(); err != nil && !modules_helper.AtLeastOneAlive(mods) {
		log.Error("No provider alive. Impossible to retrieve media informations, stopping download chain until next polling.")
		CanDownload = false
	}

	if mods, err := indexer.Status(); err != nil && !modules_helper.AtLeastOneAlive(mods) {
		log.Error("No indexer alive. Impossible to retrieve torrents for media, stopping download chain until next polling.")
		CanDownload = false
	}

	if _, err := downloader.Status(); err != nil {
		log.Error("No downloader alive. Impossible to download media, stopping download chain until next polling.")
		CanDownload = false
	}

	if _, err := mediacenter.Status(); err != nil {
		log.Error("Mediacenter not alive. Post download library refresh may not be done correctly")
	}

	if err := unix.Access(configuration.Config.Library.CustomTmpPath, unix.W_OK); err != nil {
		log.WithFields(log.Fields{
			"path": configuration.Config.Library.CustomTmpPath,
		}).Error("Cannot write into tmp path. Media will not be able to be downloaded.")
		CanDownload = false
	}
	log.Debug("========== Healthcheck loop end ==========\n")
}
