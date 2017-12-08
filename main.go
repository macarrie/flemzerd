package main

import (
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/retention"
	flag "github.com/ogier/pflag"

	provider "github.com/macarrie/flemzerd/providers"
	"github.com/macarrie/flemzerd/providers/impl/tvdb"

	indexer "github.com/macarrie/flemzerd/indexers"
	"github.com/macarrie/flemzerd/indexers/impl/torznab"

	notifier "github.com/macarrie/flemzerd/notifiers"
	"github.com/macarrie/flemzerd/notifiers/impl/desktop"
	"github.com/macarrie/flemzerd/notifiers/impl/pushbullet"

	downloader "github.com/macarrie/flemzerd/downloaders"
	"github.com/macarrie/flemzerd/downloaders/impl/transmission"

	. "github.com/macarrie/flemzerd/objects"
)

func initProviders() {
	log.Debug("Initializing Providers")

	var newProviders []provider.Provider
	for providerType, providerElt := range configuration.Config.Providers {
		switch providerType {
		case "tvdb":
			np, _ := impl.New(providerElt["apikey"])
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

	initNotifiers()
	initProviders()
	initIndexers()
	initDownloaders()

	retentionErr := retention.Load()
	if retentionErr != nil {
		log.WithFields(log.Fields{
			"error": retentionErr,
		}).Warning("Could not load retention data. Starting daemon with empty retention")
	}

	//	 Load configuration objects
	var showObjects []TvShow

	for _, show := range configuration.Config.Shows {
		showName := show
		show, err := provider.FindShow(show)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"show":  showName,
			}).Warning("Unable to get show informations")
		} else {
			showObjects = append(showObjects, show)
		}
	}
	if len(showObjects) == 0 {
		log.Fatal("Impossible to get show informations for shows defined in configuration. Shutting down")
	}

	loopTicker := time.NewTicker(time.Duration(configuration.Config.System.EpisodeCheckInterval) * time.Minute)
	go func() {
		log.Debug("Starting polling loop")
		for {
			log.Debug("========== Polling loop start ==========")

			for _, show := range showObjects {
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

					// Send only download for first 10 torrents. If the first 10 don't work, the download is probably fucked, or another problem is happening and trying more torrents won't change anything
					toDownload := downloader.FillToDownloadTorrentList(recentEpisode, torrentList)
					if len(toDownload) == 0 {
						retention.ChangeDownloadingState(recentEpisode, false)
						continue
					}
					go downloader.Download(show, recentEpisode, toDownload)
				}
			}

			log.Debug("========== Polling loop end ==========\n")
			<-loopTicker.C
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	// Listen for channels
	for {
		<-signalChannel
		log.Info("Shutting down...")

		err := retention.Save()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to save retention data")
		}

		os.Exit(0)
	}
}
