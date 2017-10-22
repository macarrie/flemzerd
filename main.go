// TODO
// [] Reduire les logs en mode normal
// [] Faire une interface Downloaders et lancer un telechargement dessus (démon transmission pour commencer)
// [] Retention pour les telechargements
// [] Telecharger le fichier torrent dans /tmp/ (pas besoin ?)
// [] Virer les fichiers torrents a la fin du telechargement (ou telechargement echoué)
// [] Gerer le cas ou il y a une authentification sur Transmission
// [] Vérifier l'espace disque libre avant d'ajouter un torrent
// [] Configuration check/Indexers/Torznab
// [] Configuration check/Downloaders/Transmission
// [] Faire une boucle infinie pour regarder regulierement les nouveaux episodes et notifier si nouveaux resultats (poll_interval configurable)
//      [] Charger poll_interval dans la conf
//      [] Utiliser poll_interval comme intervalle de boucle
// [] Gérer les timeouts dans les requetes
// [] Pouvoir configurer un dossier destination
// [] Faire des hook en fin de téléchargement (maj kodi, création d'un dossier dans un dossier destination prédéfini et copier le fichier dedans)
// [] Tests unitaires (chiant)
// [] Doc (aussi chiant)

// TODO List de riche
// [] Retention fichier des derniers episodes notifiés. Si le demon tombe, pas de nouvelle notif pour rien. Lire au démarrage et écrire a l'extinction
// [] Script systemd pour transformer l'exec en démon et le gerer avec systemctl
// [] Ajouter un autre notifier
// [] Interface web pour voir vite fait ce qu'il se passe: nouveaux episodes, shows surveillés et telechargements en cours

package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"time"
    "strconv"
	"flemzerd/configuration"
	log "flemzerd/logging"

	provider "flemzerd/providers"
	"flemzerd/providers/tvdb"

	indexer "flemzerd/indexers"
	"flemzerd/indexers/torznab"

	notifier "flemzerd/notifiers"
	"flemzerd/notifiers/pushbullet"

	downloader "flemzerd/downloaders"
	"flemzerd/downloaders/transmission"
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
                "show": showName,
            }).Warning("Unable to get show informations")
		} else {
			showObjects = append(showObjects, show)
		}
	}
    if len(showObjects) == 0 {
        log.Fatal("Impossible to get show informations for shows defined in configuration. Shutting down")
    }

    downloader.AddTorrent("test")

	log.Debug("Starting polling loop")
	loopTicker := time.NewTicker(15 * time.Second)
	for {
		log.Debug("========== Polling loop start ==========")

		for _, show := range showObjects {
			recentEpisodes, err := provider.FindRecentlyAiredEpisodesForShow(show)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
                    "show": show.Name,
				}).Warning("No recent episodes found")
				continue
			}

			//log.Debug("Recent episodes: ", recentEpisod)
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
