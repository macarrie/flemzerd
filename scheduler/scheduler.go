package scheduler

import (
	"time"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	"github.com/macarrie/flemzerd/healthcheck"
	log "github.com/macarrie/flemzerd/logging"

	downloader "github.com/macarrie/flemzerd/downloaders"
	indexer "github.com/macarrie/flemzerd/indexers"
	notifier "github.com/macarrie/flemzerd/notifiers"
	provider "github.com/macarrie/flemzerd/providers"

	"github.com/macarrie/flemzerd/downloadable"

	. "github.com/macarrie/flemzerd/objects"
)

var RunTicker *time.Ticker

func Run() {
	// 	 Load configuration objects
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
	if err := db.Client.Close(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Could not close database connection")
	}
}

func Download(d downloadable.Downloadable) {
	downloadingItem := d.GetDownloadingItem()

	if downloadingItem.Downloaded {
		d.GetLog().Debug("Item already downloaded, nothing to do")
		return
	}

	if downloadingItem.Downloading {
		d.GetLog().Debug("Item already being downloaded, nothing to do")
		return
	}

	if !healthcheck.CanDownload {
		// Check modules health again just in case
		healthcheck.CheckHealth()
		if !healthcheck.CanDownload {
			_ = notifier.NotifyFailedDownload(d)
			return
		}
	}

	downloadingItem.Pending = true
	d.SetDownloadingItem(downloadingItem)
	db.SaveDownloadable(&d)

	var torrentList []Torrent
	downloadingItem.TorrentList = downloader.FillTorrentList(downloadingItem.TorrentList)
	if len(downloadingItem.TorrentList) == 0 {
		var err error
		torrentList, err = indexer.GetTorrents(d)
		if err != nil {
			log.Warning(err)
			notifier.NotifyTorrentsNotFound(d)

			downloadingItem.Downloading = false
			downloadingItem.Pending = false
			downloadingItem.TorrentsNotFound = true
			d.SetDownloadingItem(downloadingItem)
			db.SaveDownloadable(&d)

			return
		}
	} else {
		d.GetLog().Debug("Torrents already retrieved. Performing download with current torrent list")
		torrentList = downloadingItem.TorrentList
	}

	toDownload := downloader.FillTorrentList(torrentList)
	if len(toDownload) == 0 {
		d.GetLog().Debug("No torrents found")

		if !downloadingItem.TorrentsNotFound {
			notifier.NotifyTorrentsNotFound(d)
		}

		downloadingItem.Downloading = false
		downloadingItem.Pending = false
		downloadingItem.TorrentsNotFound = true
		d.SetDownloadingItem(downloadingItem)
		db.SaveDownloadable(&d)

		return
	} else {
		d.GetLog().WithFields(log.Fields{
			"nb": len(toDownload),
		}).Debug("Torrents found")

		downloadingItem.TorrentsNotFound = false
		downloadingItem.TorrentList = toDownload
		d.SetDownloadingItem(downloadingItem)
		db.SaveDownloadable(&d)
	}

	notifier.NotifyDownloadStart(d)

	go downloader.Download(d)
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

		ep.GetLog().Debug("Launched download processing recovery")
		recoveryEpisode := ep
		go Download(&recoveryEpisode)

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

		m.GetLog().Debug("Launched download processing recovery")

		recoveryMovie := m
		go Download(&recoveryMovie)
	}
}

func poll(recoveryDone *bool) {
	log.Debug("========== Polling loop start ==========")

	if configuration.Config.System.TrackShows {
		provider.GetTVShowsInfoFromConfig()
	}
	if configuration.Config.System.TrackMovies {
		provider.GetMoviesInfoFromConfig()
	}

	if recoveryDone != nil && !*recoveryDone && healthcheck.CanDownload {
		RecoverDownloadingItems()
		*recoveryDone = true
	}

	if configuration.Config.System.TrackShows {
		for _, show := range provider.TVShows {
			recentEpisodes, err := provider.FindRecentlyAiredEpisodesForShow(show)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"show":  show.GetTitle(),
				}).Warning("No recent episodes found")
				continue
			}

			for _, recentEpisode := range recentEpisodes {
				err := notifier.NotifyRecentEpisode(&recentEpisode)
				if err != nil {
					log.Warning(err)
				}

				downloadDelayPassed := time.Now().After(recentEpisode.Date.Add(time.Duration(configuration.Config.System.ShowDownloadDelay) * time.Hour))
				if healthcheck.CanDownload && configuration.Config.System.AutomaticShowDownload && downloadDelayPassed {
					Download(&recentEpisode)
				}
			}
		}
	}

	if configuration.Config.System.TrackMovies {
		for _, movie := range provider.Movies {
			if movie.Date.After(time.Now()) {
				log.WithFields(log.Fields{
					"movie":        movie.GetTitle(),
					"release_date": movie.Date,
				}).Debug("Movie not yet released, ignoring")
				continue
			}

			err := notifier.NotifyNewMovie(&movie)
			if err != nil {
				log.Warning(err)
			}

			downloadDelayPassed := time.Now().After(movie.Date.Add(time.Duration(configuration.Config.System.MovieDownloadDelay) * time.Hour))
			if healthcheck.CanDownload && configuration.Config.System.AutomaticMovieDownload && downloadDelayPassed {
				Download(&movie)
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
