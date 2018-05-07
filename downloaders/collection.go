package downloader

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/notifiers"
	. "github.com/macarrie/flemzerd/objects"

	"github.com/rs/xid"
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
	e.DownloadingItem.Downloading = false
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
	m.DownloadingItem.Downloading = false
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

func DownloadEpisode(show TvShow, e Episode, torrentList []Torrent) error {
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
		torrent.DownloadDir = fmt.Sprintf("%s/%s/", DOWNLOAD_TMP_DIR, xid.New())

		if db.TorrentHasFailed(e.DownloadingItem, torrent) {
			continue
		}

		e.DownloadingItem.CurrentTorrent = torrent
		db.Client.Save(&e)

		torrentDownload := EpisodeHandleTorrentDownload(&e, false)
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
			notifier.NotifyDownloadedEpisode(show, &e)

			e.Downloaded = true
			e.DownloadingItem.Downloading = false
			err := MoveEpisodeToLibrary(show, &e)
			if err != nil {
				log.WithFields(log.Fields{
					"show":           show.Name,
					"episode":        e.Name,
					"season":         e.Season,
					"number":         e.Number,
					"temporary_path": e.DownloadingItem.CurrentTorrent.DownloadDir,
					"library_path":   configuration.Config.Library.ShowPath,
					"error":          err,
				}).Error("Could not move episode from temporay download path to library folder")
			}
			db.Client.Save(&e)

			return nil
		}
		continue
	}

	// If function has not returned yet, it means the download failed
	if len(e.DownloadingItem.FailedTorrents) > configuration.Config.System.TorrentDownloadAttemptsLimit {
		MarkEpisodeFailedDownload(&show, &e)
		return errors.New("Download failed, no torrents could be downloaded")
	}

	return errors.New("No torrents in current torrent list could be downloaded")
}

func DownloadMovie(m Movie, torrentList []Torrent) error {
	if m.Downloaded || m.DownloadingItem.Downloading {
		return errors.New("Movie downloading or already downloaded. Skipping")
	}

	log.WithFields(log.Fields{
		"name": m.Title,
	}).Info("Starting download process")

	m.DownloadingItem.Downloading = true
	db.Client.Save(&m)

	for _, torrent := range torrentList {
		torrent.DownloadDir = fmt.Sprintf("%s/%s/", DOWNLOAD_TMP_DIR, xid.New())

		if db.TorrentHasFailed(m.DownloadingItem, torrent) {
			continue
		}

		m.DownloadingItem.CurrentTorrent = torrent
		db.Client.Save(&m)

		torrentDownload := MovieHandleTorrentDownload(&m, false)
		if torrentDownload != nil {
			log.WithFields(log.Fields{
				"err":     torrentDownload,
				"torrent": torrent.Name,
			}).Warning("Couldn't download torrent. Skipping to next torrent in list")
		} else {
			log.WithFields(log.Fields{
				"name": m.Title,
			}).Info("Movie successfully downloaded")
			notifier.NotifyDownloadedMovie(&m)

			m.Downloaded = true
			m.DownloadingItem.Downloading = false
			err := MoveMovieToLibrary(&m)
			if err != nil {
				log.WithFields(log.Fields{
					"movie":          m.Title,
					"temporary_path": m.DownloadingItem.CurrentTorrent.DownloadDir,
					"library_path":   configuration.Config.Library.ShowPath,
					"error":          err,
				}).Error("Could not move movie from temporay download path to library folder")
			}
			db.Client.Save(&m)

			return nil
		}
		continue
	}

	// If function has not returned yet, it means the download failed
	if len(m.DownloadingItem.FailedTorrents) > configuration.Config.System.TorrentDownloadAttemptsLimit {
		MarkMovieFailedDownload(&m)
		return errors.New("Download failed, no torrents could be downloaded")
	}

	return errors.New("No torrents in current torrent list could be downloaded")
}

func MarkEpisodeFailedDownload(show *TvShow, e *Episode) {
	log.WithFields(log.Fields{
		"show":   show.Name,
		"season": e.Season,
		"number": e.Number,
		"name":   e.Name,
	}).Error("Download failed, no torrents could be downloaded")

	notifier.NotifyFailedEpisode(*show, e)

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

func MoveEpisodeToLibrary(show TvShow, episode *Episode) error {
	log.WithFields(log.Fields{
		"show":           show.Name,
		"episode":        episode.Name,
		"season":         episode.Season,
		"number":         episode.Number,
		"temporary_path": episode.DownloadingItem.CurrentTorrent.DownloadDir,
		"library_path":   configuration.Config.Library.ShowPath,
	}).Debug("Moving episode to library")

	destinationPath := fmt.Sprintf("%s/%s/Season %d/", configuration.Config.Library.ShowPath, show.Name, episode.Season)
	err := os.Rename(episode.DownloadingItem.CurrentTorrent.DownloadDir, destinationPath)
	if err != nil {
		return err
	}

	err = os.Remove(episode.DownloadingItem.CurrentTorrent.DownloadDir)
	if err != nil {
		log.WithFields(log.Fields{
			"path": episode.DownloadingItem.CurrentTorrent.DownloadDir,
		}).Warning("Could not remove temporary folder for download")
	}

	episode.DownloadingItem.CurrentTorrent.DownloadDir = destinationPath
	db.Client.Save(episode)

	return nil
}

func MoveMovieToLibrary(movie *Movie) error {
	log.WithFields(log.Fields{
		"movie":          movie.Title,
		"temporary_path": movie.DownloadingItem.CurrentTorrent.DownloadDir,
		"library_path":   configuration.Config.Library.ShowPath,
	}).Debug("Moving episode to library")

	destinationPath := fmt.Sprintf("%s/%s/", configuration.Config.Library.MoviePath, movie.Title)
	err := os.Rename(movie.DownloadingItem.CurrentTorrent.DownloadDir, destinationPath)
	if err != nil {
		return err
	}

	err = os.Remove(movie.DownloadingItem.CurrentTorrent.DownloadDir)
	if err != nil {
		log.WithFields(log.Fields{
			"path": movie.DownloadingItem.CurrentTorrent.DownloadDir,
		}).Warning("Could not remove temporary folder for download")
	}

	movie.DownloadingItem.CurrentTorrent.DownloadDir = destinationPath
	db.Client.Save(movie)

	return nil
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
		AddTorrentMapping(m.DownloadingItem.CurrentTorrent.TorrentId, m.DownloadingItem.CurrentDownloaderId)

		log.WithFields(log.Fields{
			"name": m.Title,
		}).Debug("Launched download processing recovery")

		go MovieHandleTorrentDownload(&m, true)
	}
}
