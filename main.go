package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/retention"
	"github.com/macarrie/flemzerd/server"
	flag "github.com/ogier/pflag"

	provider "github.com/macarrie/flemzerd/providers"
	"github.com/macarrie/flemzerd/providers/impl/tmdb"
	"github.com/macarrie/flemzerd/providers/impl/tvdb"

	indexer "github.com/macarrie/flemzerd/indexers"
	"github.com/macarrie/flemzerd/indexers/impl/torznab"

	notifier "github.com/macarrie/flemzerd/notifiers"
	"github.com/macarrie/flemzerd/notifiers/impl/desktop"
	"github.com/macarrie/flemzerd/notifiers/impl/pushbullet"

	downloader "github.com/macarrie/flemzerd/downloaders"
	"github.com/macarrie/flemzerd/downloaders/impl/transmission"

	watchlist "github.com/macarrie/flemzerd/watchlists"
	"github.com/macarrie/flemzerd/watchlists/impl/manual"
	"github.com/macarrie/flemzerd/watchlists/impl/trakt"
)

func initProviders() {
	log.Debug("Initializing Providers")

	var newProviders []provider.Provider
	for providerType, providerElt := range configuration.Config.Providers {
		switch providerType {
		case "tmdb":
			np, _ := tmdb.New(providerElt["apikey"])
			newProviders = append(newProviders, np)
		case "tvdb":
			np, _ := tvdb.New(providerElt["apikey"])
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
					"provider": providerType,
				}).Info("Provider added to list of providers")
			}
			newProviders = []provider.Provider{}
		}
	}
}

func initIndexers() {
	log.Debug("Initializing Indexers")

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
					"downloader": name,
				}).Info("Downloader added to list of downloaders")
			}
			newDownloaders = []downloader.Downloader{}
		}
	}
}

func initNotifiers() {
	log.Debug("Initializing Notifiers")

	for name, notifierObject := range configuration.Config.Notifiers {
		switch name {
		case "pushbullet":
			pushbulletNotifier := pushbullet.New(map[string]string{"AccessToken": notifierObject["accesstoken"]})
			notifier.AddNotifier(pushbulletNotifier)

			log.WithFields(log.Fields{
				"notifier": name,
			}).Info("Notifier added to list of notifiers")

		case "desktop":
			desktopNotifier := desktop.New()
			notifier.AddNotifier(desktopNotifier)

			log.WithFields(log.Fields{
				"notifier": name,
			}).Info("Notifier added to list of notifiers")

		default:
			log.WithFields(log.Fields{
				"notifierType": name,
			}).Warning("Unknown notifier type")
		}
	}
}

func initWatchlists() {
	log.Debug("Initializing Watchlists")

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
					"watchlist": watchlistType,
				}).Info("Watchlist added to list of watchlists")
			}
			newWatchlists = []watchlist.Watchlist{}
		}
	}
}

func downloadChainFunc() {
	for _, show := range provider.TVShows {
		recentEpisodes, err := provider.FindRecentlyAiredEpisodesForShow(show)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"show":  show.Name,
			}).Warning("No recent episodes found")
			continue
		}

		for _, recentEpisode := range recentEpisodes {
			err := notifier.NotifyRecentEpisode(show, recentEpisode)
			if err != nil {
				log.Warning(err)
			}

			if retention.HasBeenDownloaded(recentEpisode) {
				log.WithFields(log.Fields{
					"show":   show.Name,
					"number": recentEpisode.Number,
					"season": recentEpisode.Season,
					"name":   recentEpisode.Name,
				}).Debug("Episode already downloaded, nothing to do")
				continue
			}

			if retention.IsDownloading(recentEpisode) {
				log.WithFields(log.Fields{
					"show":   show.Name,
					"number": recentEpisode.Number,
					"season": recentEpisode.Season,
					"name":   recentEpisode.Name,
				}).Debug("Episode already being downloaded, nothing to do")
				continue
			}

			torrentList, err := indexer.GetTorrentForEpisode(show.Name, recentEpisode.Season, recentEpisode.Number)
			if err != nil {
				log.Warning(err)
				continue
			}
			log.Debug("Torrents found: ", len(torrentList))

			toDownload := downloader.FillToDownloadTorrentList(recentEpisode, torrentList)
			if len(toDownload) == 0 {
				downloader.MarkFailedDownload(show, recentEpisode)
				continue
			}
			go downloader.Download(show, recentEpisode, toDownload)
		}
	}
}

func main() {
	debugMode := flag.BoolP("debug", "d", false, "Start in debug mode")
	configFilePath := flag.StringP("config", "c", "", "Configuration file path to use")

	flag.Parse()

	if *debugMode {
		log.Setup(true)
	} else {
		log.Setup(false)
	}

	if *configFilePath != "" {
		log.Info("Loading provided configuration file")
		configuration.UseFile(*configFilePath)
	}

	err := configuration.Load()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot load configuration file")
	}

	configuration.Check()

	retentionErr := retention.Load()
	if retentionErr != nil {
		log.WithFields(log.Fields{
			"error": retentionErr,
		}).Warning("Could not load retention data. Starting daemon with empty retention")
	}

	initNotifiers()
	initProviders()
	initIndexers()
	initDownloaders()
	initWatchlists()

	if configuration.Config.Interface.Enabled {
		// Start HTTP server
		go server.Start(configuration.Config.Interface.Port, *debugMode)
	}

	//	 Load configuration objects
	var recovery bool = true

	loopTicker := time.NewTicker(time.Duration(configuration.Config.System.EpisodeCheckInterval) * time.Minute)
	go func() {
		log.Debug("Starting polling loop")
		for {
			var executeDownloadChain bool = true
			log.Debug("========== Polling loop start ==========")

			if _, err := notifier.Status(); err != nil {
				log.Error("No notifier alive. No notifications will be sent until next polling.")
			}

			if _, err := provider.Status(); err != nil {
				log.Error("No provider alive. Impossible to retrieve TVShow informations, stopping download chain until next polling.")
				executeDownloadChain = false
			} else {
				//Even if not able to download, retrieve media info for UI
				provider.GetTVShowsInfoFromConfig()
			}

			if _, err := indexer.Status(); err != nil {
				log.Error("No indexer alive. Impossible to retrieve torrents for TVShows, stopping download chain until next polling.")
				executeDownloadChain = false
			}

			if _, err := downloader.Status(); err != nil {
				log.Error("No downloader alive. Impossible to download TVShow, stopping download chain until next polling.")
				executeDownloadChain = false
			}

			if executeDownloadChain {
				if len(provider.TVShows) == 0 {
					provider.GetTVShowsInfoFromConfig()
				}

				if recovery {
					downloader.RecoverFromRetention()
					recovery = false
				}

				downloadChainFunc()

				err := retention.Save()
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Error("Failed to save retention data")
				}
			}
			log.Debug("========== Polling loop end ==========\n")
			<-loopTicker.C
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	for {
		<-signalChannel
		log.Info("Shutting down...")

		server.Stop()

		err := retention.Save()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to save retention data")
		}

		os.Exit(0)
	}
}
