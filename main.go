package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/coreos/go-systemd/daemon"
	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	flag "github.com/ogier/pflag"

	"github.com/macarrie/flemzerd/providers"
	"github.com/macarrie/flemzerd/providers/impl/tmdb"
	"github.com/macarrie/flemzerd/providers/impl/tvdb"

	"github.com/macarrie/flemzerd/indexers"
	"github.com/macarrie/flemzerd/indexers/impl/torznab"

	"github.com/macarrie/flemzerd/notifiers"
	"github.com/macarrie/flemzerd/notifiers/impl/desktop"
	"github.com/macarrie/flemzerd/notifiers/impl/eventlog"
	"github.com/macarrie/flemzerd/notifiers/impl/kodi"
	"github.com/macarrie/flemzerd/notifiers/impl/pushbullet"
	"github.com/macarrie/flemzerd/notifiers/impl/telegram"

	"github.com/macarrie/flemzerd/downloaders"
	"github.com/macarrie/flemzerd/downloaders/impl/transmission"

	"github.com/macarrie/flemzerd/watchlists"
	"github.com/macarrie/flemzerd/watchlists/impl/manual"
	"github.com/macarrie/flemzerd/watchlists/impl/trakt"

	"github.com/macarrie/flemzerd/mediacenters"
	"github.com/macarrie/flemzerd/mediacenters/impl/kodi"

	"github.com/macarrie/flemzerd/healthcheck"
	"github.com/macarrie/flemzerd/scheduler"
	"github.com/macarrie/flemzerd/server"
)

func initConfiguration() {
	err := configuration.Load()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot load configuration file")
	}

	configuration.Check()

	initNotifiers()
	initProviders()
	initIndexers()
	initDownloaders()
	initWatchlists()
	initMediaCenters()
	initStats()
}

func initProviders() {
	log.Debug("Initializing Providers")
	provider.Reset()

	var newProviders []provider.Provider
	for providerType, _ := range configuration.Config.Providers {
		switch providerType {
		case "tmdb":
			np, _ := tmdb.New(configuration.TMDB_API_KEY, configuration.Config.Providers[providerType]["order"])
			newProviders = append(newProviders, np)
		case "tvdb":
			np, _ := tvdb.New(configuration.TVDB_API_KEY, configuration.Config.Providers[providerType]["order"])
			newProviders = append(newProviders, np)
		default:
			log.WithFields(log.Fields{
				"providerType": providerType,
			}).Warning("Unknown provider type")
		}

		if len(newProviders) != 0 {
			for _, newProvider := range newProviders {
				provider.AddProvider(newProvider)
				log.WithFields(log.Fields{
					"provider": newProvider.GetName(),
				}).Info("Provider added to list of providers")
			}
			newProviders = []provider.Provider{}
		}
	}
}

func initIndexers() {
	log.Debug("Initializing Indexers")
	indexer.Reset()

	var newIndexers []indexer.Indexer
	for indexerType, indexerList := range configuration.Config.Indexers {
		switch indexerType {
		case "torznab":
			for _, indexer := range indexerList {
				newIndexers = append(newIndexers, torznab.New(indexer["name"], indexer["url"], indexer["apikey"]))
			}
		default:
			log.WithFields(log.Fields{
				"indexerType": indexerType,
			}).Warning("Unknown indexer type")
		}

		if len(newIndexers) != 0 {
			for _, newIndexer := range newIndexers {
				indexer.AddIndexer(newIndexer)
				log.WithFields(log.Fields{
					"indexer": newIndexer.GetName(),
				}).Info("Indexer added to list of indexers")
			}
			newIndexers = []indexer.Indexer{}
		}
	}
}

func initDownloaders() {
	log.Debug("Initializing Downloaders")
	downloader.Reset()

	var newDownloaders []downloader.Downloader
	for name, downloaderObject := range configuration.Config.Downloaders {
		switch name {
		case "transmission":
			address := downloaderObject["address"]
			port, _ := strconv.Atoi(downloaderObject["port"])
			user, authNeeded := downloaderObject["user"]
			password := downloaderObject["password"]
			if !authNeeded {
				user = ""
				password = ""
			}

			transmissionDownloader := transmission.New(address, port, user, password)
			newDownloaders = append(newDownloaders, transmissionDownloader)
		default:
			log.WithFields(log.Fields{
				"downloaderType": name,
			}).Warning("Unknown downloader type")
		}

		if len(newDownloaders) != 0 {
			for _, newDownloader := range newDownloaders {
				newDownloader.Init()
				downloader.AddDownloader(newDownloader)
				log.WithFields(log.Fields{
					"downloader": newDownloader.GetName(),
				}).Info("Downloader added to list of downloaders")
			}
			newDownloaders = []downloader.Downloader{}
		}
	}
}

func initNotifiers() {
	log.Debug("Initializing Notifiers")
	notifier.Reset()

	// Always add event log notifier
	notifier.AddNotifier(eventlog.New())

	for name, notifierObject := range configuration.Config.Notifiers {
		switch name {
		case "pushbullet":
			pushbulletNotifier := pushbullet.New(map[string]string{"AccessToken": notifierObject["accesstoken"]})
			notifier.AddNotifier(pushbulletNotifier)

			log.WithFields(log.Fields{
				"notifier": pushbulletNotifier.GetName(),
			}).Info("Notifier added to list of notifiers")

		case "desktop":
			desktopNotifier := desktop.New()
			notifier.AddNotifier(desktopNotifier)

			log.WithFields(log.Fields{
				"notifier": desktopNotifier.GetName(),
			}).Info("Notifier added to list of notifiers")

		case "kodi":
			kodiNotifier, err := kodi_notifier.New()
			if err != nil {
				log.WithFields(log.Fields{
					"notifier": "kodi",
					"error":    err,
				}).Error("Cannot connect to mediacenter for kodi notifications")
				break
			}
			notifier.AddNotifier(kodiNotifier)

			log.WithFields(log.Fields{
				"notifier": kodiNotifier.GetName(),
			}).Info("Notifier added to list of notifiers")

		case "telegram":
			telegramNotifier, err := telegram.New()
			if err != nil {
				log.WithFields(log.Fields{
					"notifier": "telegram",
					"error":    err,
				}).Error("Cannot connect to telegram for notifications")
				break
			}
			notifier.AddNotifier(telegramNotifier)

			log.WithFields(log.Fields{
				"notifier": telegramNotifier.GetName(),
			}).Info("Notifier added to list of notifiers")

		default:
			log.WithFields(log.Fields{
				"notifierType": name,
			}).Error("Unknown notifier type")
		}
	}
}

func initWatchlists() {
	log.Debug("Initializing Watchlists")
	watchlist.Reset()

	var newWatchlists []watchlist.Watchlist
	for watchlistType, _ := range configuration.Config.Watchlists {
		switch watchlistType {
		case "trakt":
			w, _ := trakt.New()
			newWatchlists = append(newWatchlists, w)
		case "manual":
			w, _ := manual.New()
			newWatchlists = append(newWatchlists, w)
		default:
			log.WithFields(log.Fields{
				"watchlistType": watchlistType,
			}).Warning("Unknown watchlist type")
		}

		if len(newWatchlists) != 0 {
			for _, newWatchlist := range newWatchlists {
				watchlist.AddWatchlist(newWatchlist)
				log.WithFields(log.Fields{
					"watchlist": newWatchlist.GetName(),
				}).Info("Watchlist added to list of watchlists")
			}
			newWatchlists = []watchlist.Watchlist{}
		}
	}
}

func initMediaCenters() {
	log.Debug("Initializing MediaCenters")
	mediacenter.Reset()

	var newMC []mediacenter.MediaCenter
	for mcType, _ := range configuration.Config.MediaCenters {
		switch mcType {
		case "kodi":
			mc, err := kodi.New()
			if err != nil {
				log.WithFields(log.Fields{
					"mediacenter": "kodi",
					"error":       err,
				}).Warning("Cannot connect to mediacenter")
			}
			newMC = append(newMC, mc)
		default:
			log.WithFields(log.Fields{
				"MediaCenterType": mcType,
			}).Warning("Unknown media center")
		}

		if len(newMC) != 0 {
			for _, mc := range newMC {
				mediacenter.AddMediaCenter(mc)
				log.WithFields(log.Fields{
					"mediacenter": mc.GetName(),
				}).Info("Mediacenter added to list of media centers")
			}
			newMC = []mediacenter.MediaCenter{}
		}
	}
}

func initStats() {
	log.Debug("Initializing Stats")

	_, _ = db.GetTrackedMovies()
	_, _ = db.GetDownloadingMovies()
	_, _ = db.GetDownloadedMovies()
	_, _ = db.GetRemovedMovies()

	_, _ = db.GetTrackedTvShows()
	_, _ = db.GetDownloadingEpisodes()
	_, _ = db.GetDownloadedEpisodes()
	_, _ = db.GetRemovedTvShows()

	_, _ = db.GetReadNotifications()
	_, _ = db.GetUnreadNotifications()
}

func main() {
	debugMode := flag.BoolP("debug", "d", false, "Start in debug mode")
	versionFlag := flag.BoolP("version", "v", false, "Display version number")
	configFilePath := flag.StringP("config", "c", "", "Configuration file path to use")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("flemzerd version %s", configuration.Version)
		os.Exit(0)
	}

	if *debugMode {
		log.Setup(true)
	} else {
		log.Setup(false)
	}

	log.Info("Starting flemzerd (version ", configuration.Version, ")")

	if *configFilePath != "" {
		log.Info("Loading provided configuration file")
		configuration.UseFile(*configFilePath)
	}

	dbErr := db.Load()
	if dbErr != nil {
		log.WithFields(log.Fields{
			"error": dbErr,
		}).Warning("Could not connect to database. Starting daemon without any previous data")
	}

	initConfiguration()

	// Perform initial healthcheck before starting check routines
	healthcheck.CheckHealth()

	scheduler.Run()
	healthcheck.Run()
	server.Stop()
	if configuration.Config.Interface.Enabled {
		go server.Start(configuration.Config.Interface.Port)
	}
	daemon.SdNotify(false, "READY=1")

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	for {
		switch sig := <-signalChannel; sig {
		case syscall.SIGINT, syscall.SIGTERM:
			log.Info("Shutting down...")
			server.Stop()
			scheduler.Stop()
			healthcheck.Stop()
			os.Exit(0)
		case syscall.SIGUSR1:
			log.Info("Signal received: reloading configuration")
			daemon.SdNotify(false, "READY=0")

			server.Stop()
			scheduler.Stop()
			healthcheck.Stop()

			initConfiguration()

			scheduler.Run()
			healthcheck.Run()
			if configuration.Config.Interface.Enabled {
				go server.Start(configuration.Config.Interface.Port)
			}

			daemon.SdNotify(false, "READY=1")
		}
	}
}
