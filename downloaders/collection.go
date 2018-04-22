package downloader

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/notifiers"
	. "github.com/macarrie/flemzerd/objects"
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

func EpisodeHandleTorrentDownload(e *Episode, recovery bool) error {
	torrent := e.DownloadingItem.CurrentTorrent
	if !recovery {
		torrentId, err := AddTorrent(torrent)
		if err != nil {
			return errors.New("Couldn't add torrent in downloader. Skipping to next torrent in list")
		}
		e.DownloadingItem.CurrentTorrent = torrent
		e.DownloadingItem.CurrentDownloaderId = torrentId
		db.Client.Save(&e)
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
			e.DownloadingItem.FailedTorrents = append(e.DownloadingItem.FailedTorrents, torrent)
			db.Client.Save(&e)

			log.WithFields(log.Fields{
				"error":   downloadErr,
				"torrent": torrent.Name,
			}).Debug("Error during torrent download. Finish current torrent download")

			return retryErr
		}
	}

	// If function has not returned yet, download ended with no errors !
	e.Downloaded = true
	db.Client.Save(&e)

	RemoveTorrent(torrent)

	return nil
}

func MovieHandleTorrentDownload(m *Movie, recovery bool) error {
	torrent := m.DownloadingItem.CurrentTorrent
	if !recovery {
		torrentId, err := AddTorrent(torrent)
		if err != nil {
			return errors.New("Couldn't add torrent in downloader. Skipping to next torrent in list")
		}

		m.DownloadingItem.CurrentTorrent = torrent
		m.DownloadingItem.CurrentDownloaderId = torrentId
		db.Client.Save(&m)
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
			m.DownloadingItem.FailedTorrents = append(m.DownloadingItem.FailedTorrents, torrent)
			db.Client.Save(&m)

			log.WithFields(log.Fields{
				"error":   downloadErr,
				"torrent": torrent.Name,
			}).Debug("Error during torrent download. Finish current torrent download")

			return retryErr
		}
	}

	// If function has not returned yet, download ended with no errors !
	m.Downloaded = true
	db.Client.Save(&m)

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

func DownloadEpisode(show *TvShow, e *Episode, torrentList []Torrent) error {
	if e.Downloaded || e.DownloadingItem.Downloading {
		return errors.New("Episode downloading or already downloaded. Skipping")
	}

	log.WithFields(log.Fields{
		"show":   show.Name,
		"season": e.Season,
		"number": e.Number,
		"name":   e.Name,
	}).Info("Starting download process")

	e.DownloadingItem.Downloading = true
	db.Client.Save(&e)

	for _, torrent := range torrentList {
		torrent.DownloadDir = fmt.Sprintf("%s/%s/Season %d/", configuration.Config.Library.ShowPath, show.Name, e.Season)

		if db.TorrentHasFailed(e.DownloadingItem, torrent) {
			continue
		}

		e.DownloadingItem.CurrentTorrent = torrent
		db.Client.Save(&e)

		torrentDownload := EpisodeHandleTorrentDownload(e, false)
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

			e.Downloaded = true
			db.Client.Save(&e)

			return nil
		}
		continue
	}

	// If function has not returned yet, it means the download failed
	if len(e.DownloadingItem.FailedTorrents) > configuration.Config.System.TorrentDownloadAttemptsLimit {
		MarkEpisodeFailedDownload(show, e)
		return errors.New("Download failed, no torrents could be downloaded")
	}

	return nil
}

func DownloadMovie(m *Movie, torrentList []Torrent) error {
	if m.Downloaded || m.DownloadingItem.Downloading {
		return errors.New("Movie downloading or already downloaded. Skipping")
	}

	log.WithFields(log.Fields{
		"name": m.Title,
	}).Info("Starting download process")

	m.DownloadingItem.Downloading = true
	m.DownloadingItem.Downloading = true
	db.Client.Save(&m)

	for _, torrent := range torrentList {
		torrent.DownloadDir = fmt.Sprintf("%s/%s/", configuration.Config.Library.MoviePath, m.Title)

		if db.TorrentHasFailed(m.DownloadingItem, torrent) {
			continue
		}

		m.DownloadingItem.CurrentTorrent = torrent
		db.Client.Save(&m)

		torrentDownload := MovieHandleTorrentDownload(m, false)
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

			m.Downloaded = true
			db.Client.Save(&m)

			return nil
		}
		continue
	}

	// If function has not returned yet, it means the download failed
	if len(m.DownloadingItem.FailedTorrents) > configuration.Config.System.TorrentDownloadAttemptsLimit {
		MarkMovieFailedDownload(m)
		return errors.New("Download failed, no torrents could be downloaded")
	}

	return nil
}

func MarkEpisodeFailedDownload(show *TvShow, e *Episode) {
	log.WithFields(log.Fields{
		"show":   show.Name,
		"season": e.Season,
		"number": e.Number,
		"name":   e.Name,
	}).Error("Download failed, no torrents could be downloaded")

	notifier.NotifyFailedEpisode(show, e)

	e.DownloadingItem.DownloadFailed = true
	e.DownloadingItem.Downloading = false
	db.Client.Save(&e)
}

func MarkMovieFailedDownload(m *Movie) {
	log.WithFields(log.Fields{
		"movie": m.Title,
	}).Error("Download failed, no torrents could be downloaded")

	notifier.NotifyFailedMovie(m)

	m.DownloadingItem.DownloadFailed = true
	m.DownloadingItem.Downloading = false
	db.Client.Save(&m)
}

func FillEpisodeToDownloadTorrentList(e *Episode, list []Torrent) []Torrent {
	var torrentList []Torrent
	for _, torrent := range list {
		if !db.TorrentHasFailed(e.DownloadingItem, torrent) {
			torrentList = append(torrentList, torrent)
		}
	}

	if len(torrentList) < 10 {
		return torrentList
	} else {
		return torrentList[:10]
	}
}

func FillMovieToDownloadTorrentList(m *Movie, list []Torrent) []Torrent {
	var torrentList []Torrent
	for _, torrent := range list {
		if !db.TorrentHasFailed(m.DownloadingItem, torrent) {
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
		if !ep.DownloadingItem.Downloading {
			log.WithFields(log.Fields{
				"episode": ep.Name,
				"season":  ep.Season,
				"number":  ep.Number,
			}).Debug("Episode is not getting downloaded, no torrents to fetch")
			continue
		}
		AddTorrentMapping(ep.DownloadingItem.CurrentTorrent.TorrentId, ep.DownloadingItem.CurrentDownloaderId)

		log.WithFields(log.Fields{
			"episode": ep.Name,
			"season":  ep.Season,
			"number":  ep.Number,
		}).Debug("Launched download processing recovery")

		go EpisodeHandleTorrentDownload(&ep, true)
	}

	if len(downloadingMoviesFromRetention) != 0 {
		log.Debug("Launching watch threads for downloading movies found in retention")
	}
	for _, m := range downloadingMoviesFromRetention {
		if !m.DownloadingItem.Downloading {
			log.WithFields(log.Fields{
				"movie": m.Title,
			}).Debug("Movie is not getting downloaded, no torrents to fetch")
			continue
		}
		AddTorrentMapping(m.DownloadingItem.CurrentTorrent.TorrentId, m.DownloadingItem.CurrentDownloaderId)

		log.WithFields(log.Fields{
			"name": m.Title,
		}).Debug("Launched download processing recovery")

		go MovieHandleTorrentDownload(&m, true)
	}
}
