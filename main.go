package main

import (
	"fmt"
	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	flag "github.com/ogier/pflag"
	"strconv"
	"time"

	provider "github.com/macarrie/flemzerd/providers"
	"github.com/macarrie/flemzerd/providers/tvdb"

	indexer "github.com/macarrie/flemzerd/indexers"
	"github.com/macarrie/flemzerd/indexers/torznab"

	notifier "github.com/macarrie/flemzerd/notifiers"
	"github.com/macarrie/flemzerd/notifiers/pushbullet"

	downloader "github.com/macarrie/flemzerd/downloaders"
	"github.com/macarrie/flemzerd/downloaders/transmission"
)

var config configuration.Configuration

func initProviders(config configuration.Configuration) {
	log.Debug("Initializing Providers")

	var newProviders []provider.Provider
	for providerType, providerElt := range config.Providers {
		switch providerType {
		case "tvdb":
			newProviders = append(newProviders, tvdb.New(providerElt["apikey"], providerElt["username"], providerElt["userkey"]))
		default:
			log.WithFields(log.Fields{
				"providerType": providerType,
			}).Warning("Unknown provider type")
		}

		if len(newProviders) != 0 {
			for _, newProvider := range newProviders {
				newProvider.Init()
				provider.AddProvider(newProvider)
				log.WithFields(log.Fields{
					"provider": providerType,
				}).Info("Provider added to list of providers")
			}
			newProviders = []provider.Provider{}
		}
	}
}

func initIndexers(config configuration.Configuration) {
	log.Debug("Initializing Indexers")

	var newIndexers []indexer.Indexer
	for indexerType, indexerList := range config.Indexers {
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

func initDownloaders(config configuration.Configuration) {
	log.Debug("Initializing Downloaders")

	var newDownloaders []downloader.Downloader
	for name, downloaderObject := range config.Downloaders {
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

func initNotifiers(config configuration.Configuration) {
	log.Debug("Initializing Notifiers")
	//fmt.Printf("%#v\n", config)
	for name, notifierObject := range config.Notifiers {
		//fmt.Printf("Name: %v, Notifier: %#v\n", name, notifierObject)
		switch name {
		case "pushbullet":
			pushbulletNotifier := pushbullet.New(map[string]string{"AccessToken": notifierObject["accesstoken"]})
			notifier.AddNotifier(pushbulletNotifier)

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

func NotifyRecentEpisode(show provider.Show, episode provider.Episode) {
	for _, episodeId := range notifier.Retention {
		airDate, err := time.Parse("2006-01-02", episode.Date)
		if err != nil {
			continue
		}

		if airDate.Before(time.Now().AddDate(0, 0, -14)) {
			notifier.RemoveFromRetention(episodeId)
		}
	}

	notificationTitle := fmt.Sprintf("%v: New episode aired (S%03dE%03d)", show.Name, episode.Season, episode.Number)
	notificationContent := fmt.Sprintf("New episode aired on %v\n%v Season %03d Episode %03d: %v", episode.Date, show.Name, episode.Season, episode.Number, episode.Name)

	err := notifier.NotifyRecentEpisode(episode.Id, notificationTitle, notificationContent)
	if err != nil {
		log.Warning("Failed to send all notifications")
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

	config, err := configuration.Load()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot load configuration file")
	}

	err = configuration.Check(config)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Errors in configuration file")
	}

	initNotifiers(config)
	initProviders(config)
	initIndexers(config)
	initDownloaders(config)

	// Load configuration objects
	var showObjects []provider.Show

	for _, show := range config.Shows {
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

	torrentErr := downloader.AddTorrent("http://localhost:9117/dl/kickasstorrent/?jackett_apikey=cigs498n8oqmtwqygegelo9hdgjd28ag&path=Wm5YUXk4bUFPdkkyQTJ2aFhSaEZlVnNIMmczc3VjdXdqZjR5dGd2R3hNVXR0ckhQdm9vWExWbnIyOXp5QXk2SGRDN3VIU00rN2s2cmN2YloraC9ZMVpRUldXRWoxeVcvbi9JYjVjRTN2N0JOc0g2c01RTHVYUjZIMnFQbXh5UTIxamVCR3UxOVpMZ0xwbXRMeWtCS3E0S2hGZDI1eTdpL1dicGJBaDhXbEZCc0toKzVwR3VsczRHUlg3NEU2UFdzOHVPcmxLOHg4eUtxQ2wzd2dXdU9MNkIydGc2RVpvcXdrMlBaaE1kKzQ2RHR6U0RYVDRHOVlZazNRbFVsNXpTaC9pZVorNitkcTJxODg5L2w2MkNOajEzSXRXbVpFMmNkNk9oeUF3RUdhajhONVJjV2RTODcrTThOUGFPQk5kVVdGQXg3ZEE5SmNxNitiejF6d1hHcGpkMzJqcWVqTWFFR0lFZ0xpZ3gvRXdrPQ2&file=Brooklyn+Nine-Nine+S05E02+The+Big+House+2+1080p+AMZN+WEB-DL+DD%2B5.1+H+264-ViSUM.torrent")
	if torrentErr != nil {
		log.Debug("Add torrent error: ", torrentErr)
	} else {
		log.Debug("Torrent added")
	}

	log.Debug("Starting polling loop")
	loopTicker := time.NewTicker(15 * time.Second)
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
				NotifyRecentEpisode(show, recentEpisode)

				torrentList, err := indexer.GetTorrentForEpisode(show.Name, recentEpisode.Season, recentEpisode.Number)
				if err != nil {
					log.Warning(err)
				}
				log.Debug("Torrents found: ", len(torrentList))
				for _, torrent := range torrentList {
					log.Debug(fmt.Sprintf("Torrent (%04v) %v", torrent.Attributes["seeders"], torrent.Title))
				}
			}
		}

		log.Debug("========== Polling loop end ==========\n")
		<-loopTicker.C
	}
}
