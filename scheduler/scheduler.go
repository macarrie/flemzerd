package scheduler

import (
	"strconv"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"

	provider "github.com/macarrie/flemzerd/providers"
	"github.com/macarrie/flemzerd/providers/impl/tmdb"
	"github.com/macarrie/flemzerd/providers/impl/tvdb"

	indexer "github.com/macarrie/flemzerd/indexers"
	"github.com/macarrie/flemzerd/indexers/impl/torznab"

	notifier "github.com/macarrie/flemzerd/notifiers"
	"github.com/macarrie/flemzerd/notifiers/impl/desktop"
	kodi_notifier "github.com/macarrie/flemzerd/notifiers/impl/kodi"
	"github.com/macarrie/flemzerd/notifiers/impl/pushbullet"
	"github.com/macarrie/flemzerd/notifiers/impl/telegram"

	downloader "github.com/macarrie/flemzerd/downloaders"
	"github.com/macarrie/flemzerd/downloaders/impl/transmission"

	watchlist "github.com/macarrie/flemzerd/watchlists"
	"github.com/macarrie/flemzerd/watchlists/impl/manual"
	"github.com/macarrie/flemzerd/watchlists/impl/trakt"

	mediacenter "github.com/macarrie/flemzerd/mediacenters"
	"github.com/macarrie/flemzerd/mediacenters/impl/kodi"

	. "github.com/macarrie/flemzerd/objects"
)

func initConfiguration(debug bool) {
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
}

func initProviders() {
	log.Debug("Initializing Providers")
	provider.Reset()

	var newProviders []provider.Provider
	for providerType, _ := range configuration.Config.Providers {
		switch providerType {
		case "tmdb":
			np, _ := tmdb.New(configuration.TMDB_API_KEY)
			newProviders = append(newProviders, np)
		case "tvdb":
			np, _ := tvdb.New(configuration.TVDB_API_KEY)
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

		case "kodi":
			kodiNotifier, err := kodi_notifier.New()
			if err != nil {
				log.WithFields(log.Fields{
					"notifier": "kodi",
					"error":    err,
				}).Warning("Cannot connect to mediacenter for kodi notifications")
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
				}).Warning("Cannot connect to telegram for notifications")
			}
			notifier.AddNotifier(telegramNotifier)

			log.WithFields(log.Fields{
				"notifier": telegramNotifier.GetName(),
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

func Run(debug bool) {
	dbErr := db.Load()
	if dbErr != nil {
		log.WithFields(log.Fields{
			"error": dbErr,
		}).Warning("Could not load retention data. Starting daemon with empty retention")
	}

	initConfiguration(debug)

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

			if _, err := mediacenter.Status(); err != nil {
				log.Error("Mediacenter not alive. Post download library refresh may not be done correctly")
			}

			if recovery {
				downloader.RecoverFromRetention()
				recovery = false
			}

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

						if executeDownloadChain {
							DownloadEpisode(recentEpisode)
						}
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

					if executeDownloadChain {
						DownloadMovie(movie)
					}
				}
			}

			log.Debug("========== Polling loop end ==========\n")
			<-loopTicker.C
		}
	}()
}

func Stop() {
	log.Info("Closing DB connection")
	db.Client.Close()
}

func Reload(debug bool) {
	initConfiguration(debug)
}

func DownloadEpisode(episode Episode) {
	if episode.DownloadingItem.Downloaded {
		log.WithFields(log.Fields{
			"show":   episode.TvShow.Name,
			"number": episode.Number,
			"season": episode.Season,
			"name":   episode.Name,
		}).Debug("Episode already downloaded, nothing to do")
		return
	}

	if episode.DownloadingItem.Downloading {
		log.WithFields(log.Fields{
			"show":   episode.TvShow.Name,
			"number": episode.Number,
			"season": episode.Season,
			"name":   episode.Name,
		}).Debug("Episode already being downloaded, nothing to do")
		return
	}

	episode.DownloadingItem.Pending = true
	db.Client.Save(&episode)

	torrentList, err := indexer.GetTorrentForEpisode(episode.TvShow.Name, episode.Season, episode.Number)
	if err != nil {
		log.Warning(err)
		return
	}
	log.Debug("Torrents found: ", len(torrentList))

	toDownload := downloader.FillEpisodeToDownloadTorrentList(&episode, torrentList)
	if len(toDownload) == 0 {
		downloader.MarkEpisodeFailedDownload(&episode)
		return
	}
	notifier.NotifyEpisodeDownloadStart(&episode)
	downloader.EpisodeDownloadRoutines[episode.ID] = make(chan bool, 1)

	go downloader.DownloadEpisode(episode, toDownload, downloader.EpisodeDownloadRoutines[episode.ID])
}

func DownloadMovie(movie Movie) {
	if movie.DownloadingItem.Downloaded {
		log.WithFields(log.Fields{
			"movie": movie.Title,
		}).Debug("Movie already downloaded, nothing to do")
		return
	}

	if movie.DownloadingItem.Downloading {
		log.WithFields(log.Fields{
			"movie": movie.Title,
		}).Debug("Movie already being downloaded, nothing to do")
		return
	}

	movie.DownloadingItem.Pending = true
	db.Client.Save(&movie)

	torrentList, err := indexer.GetTorrentForMovie(movie)
	if err != nil {
		log.Warning(err)
		return
	}
	log.WithFields(log.Fields{
		"movie": movie.Title,
		"nb":    len(torrentList),
	}).Debug("Torrents found")

	toDownload := downloader.FillMovieToDownloadTorrentList(&movie, torrentList)
	if len(toDownload) == 0 {
		log.Error("Download list empty")
		downloader.MarkMovieFailedDownload(&movie)
		return
	}
	notifier.NotifyMovieDownloadStart(&movie)
	downloader.MovieDownloadRoutines[movie.ID] = make(chan bool, 1)

	go downloader.DownloadMovie(movie, toDownload, downloader.MovieDownloadRoutines[movie.ID])
}
