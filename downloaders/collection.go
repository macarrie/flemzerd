package downloader

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/notifiers"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/retention"
)

var downloadersCollection []Downloader

func AddDownloader(d Downloader) {
	downloadersCollection = append(downloadersCollection, d)
}

func Status() ([]Module, error) {
	var modList []Module
	var aggregatedErrorMessage bytes.Buffer

	for _, downloader := range downloadersCollection {
		mod, downloaderAliveError := downloader.Status()
		if downloaderAliveError != nil {
			log.WithFields(log.Fields{
				"error": downloaderAliveError,
			}).Warning("Downloader is not alive")
			aggregatedErrorMessage.WriteString(downloaderAliveError.Error())
			aggregatedErrorMessage.WriteString("\n")
		}
		modList = append(modList, mod)
	}

	var retError error
	if aggregatedErrorMessage.Len() == 0 {
		retError = nil
	} else {
		retError = errors.New(aggregatedErrorMessage.String())
	}
	return modList, retError
}

func Reset() {
	downloadersCollection = []Downloader{}
}

func AddTorrent(t Torrent) (string, error) {
	if len(downloadersCollection) == 0 {
		return "", errors.New("Cannot add torrents, no downloaders are configured")
	}

	id, err := downloadersCollection[0].AddTorrent(t)
	if err != nil {
		return "", err
	}

	return id, nil
}

func AddTorrentMapping(flemzerId string, downloaderId string) {
	downloadersCollection[0].AddTorrentMapping(flemzerId, downloaderId)
}

func StartTorrent(t Torrent) error {
	return nil
}

func RemoveTorrent(t Torrent) error {
	if len(downloadersCollection) == 0 {
		return errors.New("Cannot remove torrents, no downloaders are configured")
	}

	return downloadersCollection[0].RemoveTorrent(t)
}

func GetTorrentStatus(t Torrent) (int, error) {
	return downloadersCollection[0].GetTorrentStatus(t)
}

func EpisodeHandleTorrentDownload(e Episode, torrent Torrent, recovery bool) error {
	if !recovery {
		torrentId, err := AddTorrent(torrent)
		if err != nil {
			return errors.New("Couldn't add torrent in downloader. Skipping to next torrent in list")
		}
		retention.SetCurrentEpisodeDownload(e, torrent, torrentId)
	}

	StartTorrent(torrent)

	retryCount := 0

	// Try twice to download a torrent before marking it as rubbish
	downloadErr := WaitForDownload(torrent)
	if downloadErr != nil {
		log.WithFields(log.Fields{
			"error":   downloadErr,
			"torrent": torrent.Name,
		}).Debug("Error during torrent download. Retrying download")

		RemoveTorrent(torrent)
		AddTorrent(torrent)
		retryCount++
		retryErr := WaitForDownload(torrent)
		if retryErr != nil {
			RemoveTorrent(torrent)
			retention.AddFailedEpisodeTorrent(e, torrent)

			log.WithFields(log.Fields{
				"error":   downloadErr,
				"torrent": torrent.Name,
			}).Debug("Error during torrent download. Finish current torrent download")

			return retryErr
		}
	}

	// If function has not returned yet, download ended with no errors !
	retention.AddDownloadedEpisode(e)
	RemoveTorrent(torrent)
	return nil
}

func MovieHandleTorrentDownload(m Movie, torrent Torrent, recovery bool) error {
	if !recovery {
		torrentId, err := AddTorrent(torrent)
		if err != nil {
			return errors.New("Couldn't add torrent in downloader. Skipping to next torrent in list")
		}
		retention.SetCurrentMovieDownload(m, torrent, torrentId)
	}

	StartTorrent(torrent)

	retryCount := 0

	// Try twice to download a torrent before marking it as rubbish
	downloadErr := WaitForDownload(torrent)
	if downloadErr != nil {
		log.WithFields(log.Fields{
			"error":   downloadErr,
			"torrent": torrent.Name,
		}).Debug("Error during torrent download. Retrying download")

		RemoveTorrent(torrent)
		AddTorrent(torrent)
		retryCount++
		retryErr := WaitForDownload(torrent)
		if retryErr != nil {
			RemoveTorrent(torrent)
			retention.AddFailedMovieTorrent(m, torrent)

			log.WithFields(log.Fields{
				"error":   downloadErr,
				"torrent": torrent.Name,
			}).Debug("Error during torrent download. Finish current torrent download")

			return retryErr
		}
	}

	// If function has not returned yet, download ended with no errors !
	retention.AddDownloadedMovie(m)
	RemoveTorrent(torrent)
	return nil
}

func WaitForDownload(t Torrent) error {
	downloadLoopTicker := time.NewTicker(1 * time.Minute)
	for {
		log.WithFields(log.Fields{
			"torrent": t.Name,
		}).Debug("Checking torrent download progress")

		status, err := GetTorrentStatus(t)
		if err != nil {
			return err
		}

		switch status {
		case TORRENT_STOPPED:
			return errors.New("Torrent stopped in download client")
		case TORRENT_SEEDING:
			// Download complete ! Return with no error
			return nil
		}
		<-downloadLoopTicker.C
	}
}

func DownloadEpisode(show TvShow, e Episode, torrentList []Torrent) error {
	if retention.EpisodeHasBeenDownloaded(e) || retention.EpisodeIsDownloading(e) {
		return errors.New(msg)
	}

	log.WithFields(log.Fields{
		"show":   show.Name,
		"season": e.Season,
		"number": e.Number,
		"name":   e.Name,
	}).Info("Starting download process")

	retention.AddDownloadingEpisode(e)
	retention.ChangeEpisodeDownloadingState(e, true)

	for _, torrent := range torrentList {
		torrent.DownloadDir = fmt.Sprintf("%s/%s/Season %d/", configuration.Config.Library.ShowPath, show.Name, e.Season)

		if retention.EpisodeIsInFailedTorrents(e, torrent) {
			continue
		}

		torrentDownload := EpisodeHandleTorrentDownload(e, torrent, false)
		if torrentDownload != nil {
			log.WithFields(log.Fields{
				"err":     torrentDownload,
				"torrent": torrent.Name,
			}).Warning("Couldn't download torrent. Skipping to next torrent in list")
		} else {
			log.WithFields(log.Fields{
				"show":   show.Name,
				"season": e.Season,
				"number": e.Number,
				"name":   e.Name,
			}).Info("Episode successfully downloaded")
			notifier.NotifyDownloadedEpisode(show, e)
			return nil
		}
		continue
	}

	// If function has not returned yet, it means the download failed
	if retention.EpisodeGetFailedTorrentsCount(e) > configuration.Config.System.TorrentDownloadAttemptsLimit {
		MarkEpisodeFailedDownload(show, e)
		return errors.New("Download failed, no torrents could be downloaded")
	}

	return nil
}

func DownloadMovie(m Movie, torrentList []Torrent) error {
	if retention.MovieHasBeenDownloaded(m) || retention.MovieIsDownloading(m) {
		return errors.New("Movie downloading or already downloaded. Skipping")
	}

	log.WithFields(log.Fields{
		"name": m.Title,
	}).Info("Starting download process")

	retention.AddDownloadingMovie(m)
	retention.ChangeMovieDownloadingState(m, true)

	for _, torrent := range torrentList {
		torrent.DownloadDir = fmt.Sprintf("%s/%s/", configuration.Config.Library.MoviePath, m.Title)

		if retention.MovieIsInFailedTorrents(m, torrent) {
			continue
		}

		torrentDownload := MovieHandleTorrentDownload(m, torrent, false)
		if torrentDownload != nil {
			log.WithFields(log.Fields{
				"err":     torrentDownload,
				"torrent": torrent.Name,
			}).Warning("Couldn't download torrent. Skipping to next torrent in list")
		} else {
			log.WithFields(log.Fields{
				"name": m.Title,
			}).Info("Movie successfully downloaded")
			notifier.NotifyDownloadedMovie(m)
			return nil
		}
		continue
	}

	// If function has not returned yet, it means the download failed
	if retention.MovieGetFailedTorrentsCount(m) > configuration.Config.System.TorrentDownloadAttemptsLimit {
		MarkMovieFailedDownload(m)
		return errors.New("Download failed, no torrents could be downloaded")
	}

	return nil
}

func MarkEpisodeFailedDownload(show TvShow, e Episode) {
	retention.AddFailedEpisode(e)
	log.WithFields(log.Fields{
		"show":   show.Name,
		"season": e.Season,
		"number": e.Number,
		"name":   e.Name,
	}).Error("Download failed, no torrents could be downloaded")

	notifier.NotifyFailedEpisode(show, e)
	retention.ChangeEpisodeDownloadingState(e, false)
}

func MarkMovieFailedDownload(m Movie) {
	retention.AddFailedMovie(m)
	log.WithFields(log.Fields{
		"movie": m.Title,
	}).Error("Download failed, no torrents could be downloaded")

	notifier.NotifyFailedMovie(m)
	retention.ChangeMovieDownloadingState(m, false)
}

func FillEpisodeToDownloadTorrentList(e Episode, list []Torrent) []Torrent {
	var torrentList []Torrent
	for _, torrent := range list {
		if !retention.EpisodeIsInFailedTorrents(e, torrent) {
			torrentList = append(torrentList, torrent)
		}
	}

	if len(torrentList) < 10 {
		return torrentList
	} else {
		return torrentList[:10]
	}
}

func FillMovieToDownloadTorrentList(m Movie, list []Torrent) []Torrent {
	var torrentList []Torrent
	for _, torrent := range list {
		if !retention.MovieIsInFailedTorrents(m, torrent) {
			torrentList = append(torrentList, torrent)
		}
	}

	if len(torrentList) < 10 {
		return torrentList
	} else {
		return torrentList[:10]
	}
}

func RecoverFromRetention() {
	downloadingEpisodesFromRetention, err := retention.GetDownloadingEpisodes()
	if err != nil {
		log.Error(err)
		return
	}
	downloadingMoviesFromRetention, err := retention.GetDownloadingMovies()
	if err != nil {
		log.Error(err)
		return
	}

	if len(downloadingEpisodesFromRetention) != 0 {
		log.Debug("Launching watch threads for downloading episodes found in retention")
	}
	for _, ep := range downloadingEpisodesFromRetention {
		torrent, downloaderId, err := retention.GetCurrentDownloadEpisodeFromRetention(ep)
		if err != nil {
			log.Error(nil)
			continue
		}
		AddTorrentMapping(torrent.Id, downloaderId)

		log.WithFields(log.Fields{
			"episode": ep.Name,
			"season":  ep.Season,
			"number":  ep.Number,
		}).Debug("Launched download processing recovery")

		go EpisodeHandleTorrentDownload(ep, torrent, true)
	}

	if len(downloadingMoviesFromRetention) != 0 {
		log.Debug("Launching watch threads for downloading movies found in retention")
	}
	for _, m := range downloadingMoviesFromRetention {
		torrent, downloaderId, err := retention.GetCurrentDownloadMovieFromRetention(m)
		if err != nil {
			log.Error(nil)
			continue
		}
		AddTorrentMapping(torrent.Id, downloaderId)

		log.WithFields(log.Fields{
			"name": m.Title,
		}).Debug("Launched download processing recovery")

		go MovieHandleTorrentDownload(m, torrent, true)
	}
}
