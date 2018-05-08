package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
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

	. "github.com/macarrie/flemzerd/objects"
)

func initConfiguration(debug bool) {
	daemon.SdNotify(false, "READY=0")
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

	server.Stop()
	if configuration.Config.Interface.Enabled {
		// Start HTTP server
		go server.Start(configuration.Config.Interface.Port, debug)
	}
	daemon.SdNotify(false, "READY=1")
}

func initProviders() {
	log.Debug("Initializing Providers")
	provider.Reset()

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

		default:
			log.WithFields(log.Fields{
				"notifierType": name,
			}).Warning("Unknown notifier type")
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

func downloadChainFunc() {
	if configuration.Config.System.TrackShows {
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
				reqEpisode := Episode{}
				req := db.Client.Where(Episode{
					Name:   recentEpisode.Name,
					Season: recentEpisode.Season,
					Number: recentEpisode.Number,
				}).Find(&reqEpisode)
				if req.RecordNotFound() {
					recentEpisode.TvShow = show
					db.Client.Create(&recentEpisode)
				} else {
					recentEpisode = reqEpisode
				}

				err := notifier.NotifyRecentEpisode(&recentEpisode)
				if err != nil {
					log.Warning(err)
				}

				if recentEpisode.Downloaded {
					log.WithFields(log.Fields{
						"show":   show.Name,
						"number": recentEpisode.Number,
						"season": recentEpisode.Season,
						"name":   recentEpisode.Name,
					}).Debug("Episode already downloaded, nothing to do")
					continue
				}

				if recentEpisode.DownloadingItem.Downloading {
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

				toDownload := downloader.FillEpisodeToDownloadTorrentList(&recentEpisode, torrentList)
				if len(toDownload) == 0 {
					downloader.MarkEpisodeFailedDownload(&show, &recentEpisode)
					continue
				}
				notifier.NotifyEpisodeDownloadStart(&recentEpisode)
				go downloader.DownloadEpisode(show, recentEpisode, toDownload)
			}
		}
	}

	if configuration.Config.System.TrackMovies {
		for _, movie := range provider.Movies {
			if movie.Date.After(time.Now()) {
				log.WithFields(log.Fields{
					"movie":        movie.Title,
					"release_date": movie.Date,
				}).Debug("Movie not yet released, ignoring")
				continue
			}

			err := notifier.NotifyNewMovie(&movie)
			if err != nil {
				log.Warning(err)
			}

			if movie.Downloaded {
				log.WithFields(log.Fields{
					"movie": movie.Title,
				}).Debug("Movie already downloaded, nothing to do")
				continue
			}

			if movie.DownloadingItem.Downloading {
				log.WithFields(log.Fields{
					"movie": movie.Title,
				}).Debug("Movie already being downloaded, nothing to do")
				continue
			}

			torrentList, err := indexer.GetTorrentForMovie(movie)
			if err != nil {
				log.Warning(err)
				continue
			}
			log.WithFields(log.Fields{
				"movie": movie.Title,
				"nb":    len(torrentList),
			}).Debug("Torrents found")

			toDownload := downloader.FillMovieToDownloadTorrentList(&movie, torrentList)
			if len(toDownload) == 0 {
				log.Error("Download list empty")
				downloader.MarkMovieFailedDownload(&movie)
				continue
			}
			notifier.NotifyMovieDownloadStart(&movie)
			go downloader.DownloadMovie(movie, toDownload)
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

	dbErr := db.Load()
	if dbErr != nil {
		log.WithFields(log.Fields{
			"error": dbErr,
		}).Warning("Could not load retention data. Starting daemon with empty retention")
	}

	initConfiguration(*debugMode)

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
				//Even if not able to download, retrieve media info for UI if enabled
				if configuration.Config.System.TrackShows {
					provider.GetTVShowsInfoFromConfig()
				}
				if configuration.Config.System.TrackMovies {
					provider.GetMoviesInfoFromConfig()
				}
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
				if configuration.Config.System.TrackShows {
					provider.GetTVShowsInfoFromConfig()
				}
				if configuration.Config.System.TrackMovies {
					provider.GetMoviesInfoFromConfig()
				}

				if recovery {
					downloader.RecoverFromRetention()
					recovery = false
				}

				downloadChainFunc()
			}

			log.Debug("========== Polling loop end ==========\n")
			<-loopTicker.C
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	for {
		switch sig := <-signalChannel; sig {
		case syscall.SIGINT, syscall.SIGTERM:
			log.Info("Shutting down...")

			server.Stop()

			log.Info("Closing DB connection")
			db.Client.Close()

			os.Exit(0)
		case syscall.SIGUSR1:
			log.Info("Signal received: reloading configuration")
			initConfiguration(*debugMode)
		}
	}
}
