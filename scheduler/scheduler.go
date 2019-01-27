package scheduler

import (
	"strconv"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	"golang.org/x/sys/unix"

	provider "github.com/macarrie/flemzerd/providers"
	"github.com/macarrie/flemzerd/providers/impl/tmdb"
	"github.com/macarrie/flemzerd/providers/impl/tvdb"

	indexer "github.com/macarrie/flemzerd/indexers"
	"github.com/macarrie/flemzerd/indexers/impl/torznab"

	notifier "github.com/macarrie/flemzerd/notifiers"
	"github.com/macarrie/flemzerd/notifiers/impl/desktop"
	"github.com/macarrie/flemzerd/notifiers/impl/eventlog"
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

	media_helper "github.com/macarrie/flemzerd/helpers/media"
	. "github.com/macarrie/flemzerd/objects"
)

var RunTicker *time.Ticker

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
		}).Warning("Could not connect to database. Starting daemon without any previous data")
	}

	initConfiguration(debug)

	//	 Load configuration objects
	var recoveryDone bool = false

	RunTicker = time.NewTicker(time.Duration(configuration.Config.System.CheckInterval) * time.Minute)
	go func() {
		log.Debug("Starting polling loop")
		for {
			poll(&recoveryDone)
			<-RunTicker.C
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

func DownloadEpisode(episode Episode, recovery bool) {
	if recovery {
		episode.DownloadingItem.Downloading = false
	}

	if episode.DownloadingItem.Downloaded {
		log.WithFields(log.Fields{
			"show":   media_helper.GetShowTitle(episode.TvShow),
			"number": episode.Number,
			"season": episode.Season,
			"title":  episode.Title,
		}).Debug("Episode already downloaded, nothing to do")
		return
	}

	if episode.DownloadingItem.Downloading {
		log.WithFields(log.Fields{
			"show":   media_helper.GetShowTitle(episode.TvShow),
			"number": episode.Number,
			"season": episode.Season,
			"title":  episode.Title,
		}).Debug("Episode already being downloaded, nothing to do")
		return
	}

	episode.DownloadingItem.Pending = true
	db.Client.Save(&episode)

	torrentList, err := indexer.GetTorrentForEpisode(episode)
	if err != nil {
		log.Warning(err)
		return
	}

	if recovery && episode.DownloadingItem.CurrentTorrent.ID != 0 {
		torrentList = append([]Torrent{episode.DownloadingItem.CurrentTorrent}, torrentList...)
	}

	toDownload := downloader.FillEpisodeToDownloadTorrentList(&episode, torrentList)
	if len(toDownload) == 0 {
		log.WithFields(log.Fields{
			"show":    media_helper.GetShowTitle(episode.TvShow),
			"episode": episode.Title,
			"nb":      len(torrentList),
		}).Warning("No torrents found")

		// Dont send notification when no torrents have been found on previous torrent download
		if !episode.DownloadingItem.TorrentsNotFound {
			if err := notifier.SendNotification(Notification{
				Type:    NOTIFICATION_NO_TORRENTS,
				Episode: episode,
			}); err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Warning("Could not send 'no torrents found' notification")
			} else {
				episode.DownloadingItem.TorrentsNotFound = true
			}
		}

		episode.DownloadingItem.Downloading = false
		episode.DownloadingItem.Pending = false
		db.Client.Save(&episode)

		return
	} else {
		log.WithFields(log.Fields{
			"show":    episode.TvShow.OriginalTitle,
			"episode": episode.Title,
			"nb":      len(torrentList),
		}).Debug("Torrents found")

		episode.DownloadingItem.TorrentsNotFound = false
		db.Client.Save(&episode)
	}

	notifier.NotifyEpisodeDownloadStart(&episode)

	go downloader.DownloadEpisode(episode, toDownload, recovery)
}

func DownloadMovie(movie Movie, recovery bool) {
	if movie.DownloadingItem.Downloaded {
		log.WithFields(log.Fields{
			"movie": media_helper.GetMovieTitle(movie),
		}).Debug("Movie already downloaded, nothing to do")
		return
	}

	if movie.DownloadingItem.Downloading {
		log.WithFields(log.Fields{
			"movie": media_helper.GetMovieTitle(movie),
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
	if recovery && movie.DownloadingItem.CurrentTorrent.ID != 0 {
		torrentList = append([]Torrent{movie.DownloadingItem.CurrentTorrent}, torrentList...)
	}

	toDownload := downloader.FillMovieToDownloadTorrentList(&movie, torrentList)
	if len(toDownload) == 0 {
		log.WithFields(log.Fields{
			"movie": media_helper.GetMovieTitle(movie),
			"nb":    len(torrentList),
		}).Debug("Torrents found")

		if !movie.DownloadingItem.TorrentsNotFound {
			if err := notifier.SendNotification(Notification{
				Type:  NOTIFICATION_NO_TORRENTS,
				Movie: movie,
			}); err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Warning("Could not send 'no torrents found' notification")
			} else {
				movie.DownloadingItem.TorrentsNotFound = true
			}
		}

		movie.DownloadingItem.Downloading = false
		movie.DownloadingItem.Pending = false
		db.Client.Save(&movie)

		return
	} else {
		log.WithFields(log.Fields{
			"movie": movie.OriginalTitle,
			"nb":    len(torrentList),
		}).Debug("Torrents found")

		movie.DownloadingItem.TorrentsNotFound = false
		db.Client.Save(&movie)
	}

	notifier.NotifyMovieDownloadStart(&movie)

	go downloader.DownloadMovie(movie, toDownload, recovery)
}

func RecoverDownloadingItems() {
	downloadingEpisodesFromRetention, err := db.GetDownloadingEpisodes()
	if err != nil {
		log.Error(err)
		return
	}
	downloadingMoviesFromRetention, err := db.GetDownloadingMovies()
	if err != nil {
		log.Error(err)
		return
	}

	if len(downloadingEpisodesFromRetention) != 0 {
		log.Debug("Launching watch threads for downloading episodes found in retention")
	}
	for _, ep := range downloadingEpisodesFromRetention {
		ep.DeletedAt = nil
		ep.DownloadingItem.Pending = false
		ep.DownloadingItem.Downloading = false
		ep.DownloadingItem.Downloaded = false
		db.Client.Save(&ep)

		log.WithFields(log.Fields{
			"episode": ep.Title,
			"season":  ep.Season,
			"number":  ep.Number,
		}).Debug("Launched download processing recovery")
		recoveryEpisode := ep
		go DownloadEpisode(recoveryEpisode, true)

	}

	if len(downloadingMoviesFromRetention) != 0 {
		log.Debug("Launching watch threads for downloading movies found in retention")
	}
	for _, m := range downloadingMoviesFromRetention {
		m.DeletedAt = nil
		m.DownloadingItem.Pending = false
		m.DownloadingItem.Downloading = false
		m.DownloadingItem.Downloaded = false
		db.Client.Save(&m)

		log.WithFields(log.Fields{
			"title": m.OriginalTitle,
		}).Debug("Launched download processing recovery")

		recoveryMovie := m
		go DownloadMovie(recoveryMovie, true)
	}
}

func poll(recoveryDone *bool) {
	var executeDownloadChain bool = true
	log.Debug("========== Polling loop start ==========")

	if _, err := notifier.Status(); err != nil {
		log.Error("No notifier alive. No notifications will be sent until next polling.")
	}

	if _, err := provider.Status(); err != nil {
		log.Error("No provider alive. Impossible to retrieve media informations, stopping download chain until next polling.")
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
		log.Error("No indexer alive. Impossible to retrieve torrents for media, stopping download chain until next polling.")
		executeDownloadChain = false
	}

	if _, err := downloader.Status(); err != nil {
		log.Error("No downloader alive. Impossible to download media, stopping download chain until next polling.")
		executeDownloadChain = false
	}

	if _, err := mediacenter.Status(); err != nil {
		log.Error("Mediacenter not alive. Post download library refresh may not be done correctly")
	}

	if err := unix.Access(configuration.Config.Library.CustomTmpPath, unix.W_OK); err != nil {
		log.WithFields(log.Fields{
			"path": configuration.Config.Library.CustomTmpPath,
		}).Error("Cannot write into tmp path. Media will not be able to be downloaded.")
		executeDownloadChain = false
	}

	if recoveryDone != nil && !*recoveryDone && executeDownloadChain {
		RecoverDownloadingItems()
		*recoveryDone = true
	}

	if configuration.Config.System.TrackShows {
		for _, show := range provider.TVShows {
			recentEpisodes, err := provider.FindRecentlyAiredEpisodesForShow(show)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"show":  media_helper.GetShowTitle(show),
				}).Warning("No recent episodes found")
				continue
			}

			for _, recentEpisode := range recentEpisodes {
				err := notifier.NotifyRecentEpisode(&recentEpisode)
				if err != nil {
					log.Warning(err)
				}

				if executeDownloadChain && configuration.Config.System.AutomaticShowDownload {
					DownloadEpisode(recentEpisode, false)
				}
			}
		}
	}

	if configuration.Config.System.TrackMovies {
		for _, movie := range provider.Movies {
			if movie.Date.After(time.Now()) {
				log.WithFields(log.Fields{
					"movie":        media_helper.GetMovieTitle(movie),
					"release_date": movie.Date,
				}).Debug("Movie not yet released, ignoring")
				continue
			}

			err := notifier.NotifyNewMovie(&movie)
			if err != nil {
				log.Warning(err)
			}

			if executeDownloadChain && configuration.Config.System.AutomaticMovieDownload {
				DownloadMovie(movie, false)
			}
		}
	}

	log.Debug("========== Polling loop end ==========\n")
}

func ResetRunTicker() {
	RunTicker.Stop()
	RunTicker = time.NewTicker(time.Duration(configuration.Config.System.CheckInterval) * time.Minute)
	poll(nil)
}
