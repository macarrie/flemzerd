package downloader

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	media_helper "github.com/macarrie/flemzerd/helpers/media"
	log "github.com/macarrie/flemzerd/logging"
	mediacenter "github.com/macarrie/flemzerd/mediacenters"
	notifier "github.com/macarrie/flemzerd/notifiers"
	. "github.com/macarrie/flemzerd/objects"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

var downloadersCollection []Downloader
var EpisodeDownloadRoutines map[uint](chan bool)
var MovieDownloadRoutines map[uint](chan bool)

func init() {
	EpisodeDownloadRoutines = make(map[uint](chan bool))
	MovieDownloadRoutines = make(map[uint](chan bool))
}

func AddDownloader(d Downloader) {
	downloadersCollection = append(downloadersCollection, d)
}

func Status() ([]Module, error) {
	var modList []Module
	var errorList *multierror.Error

	for _, downloader := range downloadersCollection {
		mod, downloaderAliveError := downloader.Status()
		if downloaderAliveError != nil {
			log.WithFields(log.Fields{
				"error": downloaderAliveError,
			}).Warning("Downloader is not alive")
			errorList = multierror.Append(errorList, downloaderAliveError)
		}
		modList = append(modList, mod)
	}

	return modList, errorList.ErrorOrNil()
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
		return "", errors.Wrap(err, "cannot add torrent in downloader")
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

func EpisodeHandleTorrentDownload(e *Episode, recovery bool, stopChannel chan bool) (err error, aborted bool) {
	torrent := e.DownloadingItem.CurrentTorrent
	if !recovery {
		torrentId, err := AddTorrent(torrent)
		if err != nil {
			RemoveTorrent(torrent)
			e.DownloadingItem.FailedTorrents = append(e.DownloadingItem.FailedTorrents, torrent)
			e.DownloadingItem.CurrentTorrent = Torrent{}
			db.Client.Save(&e)
			return errors.Wrap(err, "Couldn't add torrent in downloader. Skipping to next torrent in list"), false
		}

		e.DownloadingItem.CurrentTorrent = torrent
		e.DownloadingItem.CurrentDownloaderId = torrentId
		db.Client.Save(&e)
	}

	StartTorrent(torrent)

	downloadErr, downloadAborted := WaitForDownload(torrent, stopChannel)
	if downloadAborted {
		return nil, downloadAborted
	}
	if downloadErr != nil {
		RemoveTorrent(torrent)
		e.DownloadingItem.FailedTorrents = append(e.DownloadingItem.FailedTorrents, torrent)
		db.Client.Save(&e)

		log.WithFields(log.Fields{
			"error":   downloadErr,
			"torrent": torrent.Name,
		}).Debug("Error during torrent download. Finish current torrent download")

		return errors.Wrap(downloadErr, "error during download"), false
	}

	// If function has not returned yet, download ended with no errors !
	log.WithFields(log.Fields{
		"show":   media_helper.GetShowTitle(e.TvShow),
		"season": e.Season,
		"number": e.Number,
		"name":   e.Name,
	}).Info("Episode successfully downloaded")
	notifier.NotifyDownloadedEpisode(e)

	e.DownloadingItem.Pending = false
	e.DownloadingItem.Downloading = false
	e.DownloadingItem.Downloaded = true
	db.Client.Save(e)

	err = MoveEpisodeToLibrary(e)
	if err != nil {
		log.WithFields(log.Fields{
			"show":           media_helper.GetShowTitle(e.TvShow),
			"episode":        e.Name,
			"season":         e.Season,
			"number":         e.Number,
			"temporary_path": e.DownloadingItem.CurrentTorrent.DownloadDir,
			"library_path":   configuration.Config.Library.ShowPath,
			"error":          err,
		}).Error("Could not move episode from temporay download path to library folder")
	} else {
		mediacenter.RefreshLibrary()
	}

	RemoveTorrent(torrent)

	return nil, false
}

func MovieHandleTorrentDownload(m *Movie, recovery bool, stopChannel chan bool) (err error, aborted bool) {
	torrent := m.DownloadingItem.CurrentTorrent
	if !recovery {
		torrentId, err := AddTorrent(torrent)
		if err != nil {
			RemoveTorrent(torrent)
			m.DownloadingItem.FailedTorrents = append(m.DownloadingItem.FailedTorrents, torrent)
			m.DownloadingItem.CurrentTorrent = Torrent{}
			db.Client.Save(&m)
			return errors.Wrap(err, "Couldn't add torrent in downloader. Skipping to next torrent in list"), false
		}

		m.DownloadingItem.CurrentTorrent = torrent
		m.DownloadingItem.CurrentDownloaderId = torrentId
		db.Client.Save(&m)
	}

	StartTorrent(torrent)

	downloadErr, downloadAborted := WaitForDownload(torrent, stopChannel)
	if downloadAborted {
		return nil, downloadAborted
	}
	if downloadErr != nil {
		RemoveTorrent(torrent)
		m.DownloadingItem.FailedTorrents = append(m.DownloadingItem.FailedTorrents, torrent)
		db.Client.Save(&m)

		log.WithFields(log.Fields{
			"error":   downloadErr,
			"torrent": torrent.Name,
		}).Debug("Error during torrent download. Finish current torrent download")

		return errors.Wrap(downloadErr, "error during download"), false
	}

	// If function has not returned yet, download ended with no errors !
	log.WithFields(log.Fields{
		"name": media_helper.GetMovieTitle(*m),
	}).Info("Movie successfully downloaded")
	notifier.NotifyDownloadedMovie(m)

	m.DownloadingItem.Pending = false
	m.DownloadingItem.Downloading = false
	m.DownloadingItem.Downloaded = true
	db.Client.Save(m)

	err = MoveMovieToLibrary(m)
	if err != nil {
		log.WithFields(log.Fields{
			"movie":          media_helper.GetMovieTitle(*m),
			"temporary_path": m.DownloadingItem.CurrentTorrent.DownloadDir,
			"library_path":   configuration.Config.Library.MoviePath,
			"error":          err,
		}).Error("Could not move movie from temporary download path to library folder")
	} else {
		mediacenter.RefreshLibrary()
	}

	RemoveTorrent(torrent)

	return nil, false
}

func WaitForDownload(t Torrent, stopChannel chan bool) (err error, aborted bool) {
	downloadLoopTicker := time.NewTicker(1 * time.Minute)
	for {
		log.WithFields(log.Fields{
			"torrent": t.Name,
		}).Debug("Checking torrent download progress")

		// Check if download has been manually stopped
		select {
		case <-stopChannel:
			return nil, true
		default:
			log.WithFields(log.Fields{
				"torrent": t.Name,
			}).Debug("Download in progress")
		}

		status, err := GetTorrentStatus(t)
		if err != nil {
			return errors.Wrap(err, "error when getting download status"), false
		}

		switch status {
		case TORRENT_STOPPED:
			return errors.New("Torrent stopped in download client"), false
		case TORRENT_SEEDING:
			// Download complete ! Return with no error
			return nil, false
		}
		<-downloadLoopTicker.C
	}
}

func DownloadEpisode(e Episode, torrentList []Torrent, stopChannel chan bool, recovery bool) error {
	recoveryDone := false

	if e.DownloadingItem.Downloaded || e.DownloadingItem.Downloading {
		return errors.New("Episode downloading or already downloaded. Skipping")
	}

	log.WithFields(log.Fields{
		"show":   media_helper.GetShowTitle(e.TvShow),
		"season": e.Season,
		"number": e.Number,
		"name":   e.Name,
	}).Info("Starting download process")

	e.DownloadingItem.Pending = false
	e.DownloadingItem.Downloading = true
	db.Client.Save(&e)

	if recovery {
		AddTorrentMapping(e.DownloadingItem.CurrentTorrent.TorrentId, e.DownloadingItem.CurrentDownloaderId)
	}

	for _, torrent := range torrentList {
		torrent.DownloadDir = fmt.Sprintf("%s/%s/", configuration.Config.Library.CustomTmpPath, xid.New())

		if db.TorrentHasFailed(e.DownloadingItem, torrent) {
			continue
		}

		e.DownloadingItem.CurrentTorrent = torrent
		db.Client.Save(&e)

		var torrentDownload error
		var downloadAborted bool
		if recoveryDone {
			torrentDownload, downloadAborted = EpisodeHandleTorrentDownload(&e, false, stopChannel)
		} else {
			recoveryDone = true
			torrentDownload, downloadAborted = EpisodeHandleTorrentDownload(&e, true, stopChannel)
		}
		if downloadAborted {
			log.WithFields(log.Fields{
				"show":   media_helper.GetShowTitle(e.TvShow),
				"season": e.Season,
				"number": e.Number,
				"name":   e.Name,
			}).Info("Download manually aborted. Cleaning up current download artifacts.")

			currentDownloadPath := e.DownloadingItem.CurrentTorrent.DownloadDir
			RemoveTorrent(e.DownloadingItem.CurrentTorrent)
			os.Remove(currentDownloadPath)

			e.DownloadingItem.AbortPending = false
			e.DownloadingItem.FailedTorrents = []Torrent{}
			e.DownloadingItem.CurrentTorrent = Torrent{}
			e.DownloadingItem.CurrentDownloaderId = ""
			db.Client.Save(&e)
			db.Client.Delete(&e)

			return nil
		}
		if torrentDownload != nil {
			log.WithFields(log.Fields{
				"err":     torrentDownload,
				"torrent": torrent.Name,
			}).Warning("Couldn't download torrent. Skipping to next torrent in list")
		} else {
			return nil
		}

		continue
	}

	// If function has not returned yet, it means the download failed
	if len(e.DownloadingItem.FailedTorrents) > configuration.Config.System.TorrentDownloadAttemptsLimit {
		MarkEpisodeFailedDownload(&e)
		return errors.New("Download failed, no torrents could be downloaded")
	}

	return errors.New("No torrents in current torrent list could be downloaded")
}

func DownloadMovie(m Movie, torrentList []Torrent, stopChannel chan bool, recovery bool) error {
	recoveryDone := true
	if m.DownloadingItem.Downloaded || m.DownloadingItem.Downloading {
		return errors.New("Movie downloading or already downloaded. Skipping")
	}

	log.WithFields(log.Fields{
		"name": media_helper.GetMovieTitle(m),
	}).Info("Starting download process")

	m.DownloadingItem.Pending = false
	m.DownloadingItem.Downloading = true
	db.Client.Save(&m)

	if recovery {
		AddTorrentMapping(m.DownloadingItem.CurrentTorrent.TorrentId, m.DownloadingItem.CurrentDownloaderId)
	}

	for _, torrent := range torrentList {
		torrent.DownloadDir = fmt.Sprintf("%s/%s/", configuration.Config.Library.CustomTmpPath, xid.New())

		if db.TorrentHasFailed(m.DownloadingItem, torrent) {
			continue
		}

		m.DownloadingItem.CurrentTorrent = torrent
		db.Client.Save(&m)

		var torrentDownload error
		var downloadAborted bool
		if recoveryDone {
			torrentDownload, downloadAborted = MovieHandleTorrentDownload(&m, false, stopChannel)
		} else {
			recoveryDone = true
			torrentDownload, downloadAborted = MovieHandleTorrentDownload(&m, true, stopChannel)
		}
		if downloadAborted {
			log.WithFields(log.Fields{
				"movie": media_helper.GetMovieTitle(m),
			}).Info("Download manually aborted. Cleaning up current download artifacts.")

			currentDownloadPath := m.DownloadingItem.CurrentTorrent.DownloadDir
			RemoveTorrent(m.DownloadingItem.CurrentTorrent)
			os.Remove(currentDownloadPath)

			m.DownloadingItem.AbortPending = false
			m.DownloadingItem.FailedTorrents = []Torrent{}
			m.DownloadingItem.CurrentTorrent = Torrent{}
			m.DownloadingItem.CurrentDownloaderId = ""

			db.Client.Save(&m)

			return nil
		}
		if torrentDownload != nil {
			log.WithFields(log.Fields{
				"err":     torrentDownload,
				"torrent": torrent.Name,
			}).Warning("Couldn't download torrent. Skipping to next torrent in list")
		} else {
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

func AbortEpisodeDownload(e *Episode) {
	log.WithFields(log.Fields{
		"id":      e.ID,
		"show":    media_helper.GetShowTitle(e.TvShow),
		"episode": e.Name,
		"season":  e.Season,
		"number":  e.Number,
	}).Info("Aborting episode download")

	if !e.DownloadingItem.Pending {
		close(EpisodeDownloadRoutines[e.ID])
		e.DownloadingItem.AbortPending = true
	} else {
		e.DownloadingItem.AbortPending = false
	}

	e.DownloadingItem.Pending = false
	e.DownloadingItem.Downloading = false
	e.DownloadingItem.Downloaded = false
	e.DownloadingItem.DownloadFailed = false
	db.Client.Save(e)
}

func AbortMovieDownload(m *Movie) {
	log.WithFields(log.Fields{
		"id":    m.ID,
		"title": media_helper.GetMovieTitle(*m),
	}).Info("Aborting movie download")

	if !m.DownloadingItem.Pending {
		close(MovieDownloadRoutines[m.ID])
		m.DownloadingItem.AbortPending = true
	} else {
		m.DownloadingItem.AbortPending = false
	}

	m.DownloadingItem.Pending = false
	m.DownloadingItem.Downloading = false
	m.DownloadingItem.Downloaded = false
	m.DownloadingItem.DownloadFailed = false
	db.Client.Save(m)
}

func MarkEpisodeFailedDownload(e *Episode) {
	log.WithFields(log.Fields{
		"show":   media_helper.GetShowTitle(e.TvShow),
		"season": e.Season,
		"number": e.Number,
		"name":   e.Name,
	}).Error("Download failed, no torrents could be downloaded")

	notifier.NotifyFailedEpisode(e)

	e.DownloadingItem.DownloadFailed = true
	e.DownloadingItem.Downloading = false
	db.Client.Save(&e)
}

func MarkMovieFailedDownload(m *Movie) {
	log.WithFields(log.Fields{
		"movie": media_helper.GetMovieTitle(*m),
	}).Error("Download failed, no torrents could be downloaded")

	notifier.NotifyFailedMovie(m)

	m.DownloadingItem.DownloadFailed = true
	m.DownloadingItem.Downloading = false
	db.Client.Save(&m)
}

func sanitizeStringForFilename(src string) string {
	reg, _ := regexp.Compile("[^a-z0-9]+")
	return reg.ReplaceAllString(strings.ToLower(src), "_")
}

func MoveEpisodeToLibrary(episode *Episode) error {
	log.WithFields(log.Fields{
		"show":           media_helper.GetShowTitle(episode.TvShow),
		"episode":        episode.Name,
		"season":         episode.Season,
		"number":         episode.Number,
		"temporary_path": episode.DownloadingItem.CurrentTorrent.DownloadDir,
		"library_path":   configuration.Config.Library.ShowPath,
	}).Debug("Moving episode to library")

	destinationPath := fmt.Sprintf("%s/%s/Season %d/s%02de%02d", configuration.Config.Library.ShowPath, sanitizeStringForFilename(media_helper.GetShowTitle(episode.TvShow)), episode.Season, episode.Season, episode.Number)
	err := os.MkdirAll(destinationPath, 0755)
	if err != nil {
		return errors.Wrap(err, "Could not create library folder for episode")
	}

	target := fmt.Sprintf("%s/%s", destinationPath, episode.DownloadingItem.CurrentTorrent.Name)
	err = os.Rename(episode.DownloadingItem.CurrentTorrent.DownloadDir, target)
	if err != nil {
		return errors.Wrap(err, "Could not move episode to library")
	}

	err = os.Remove(episode.DownloadingItem.CurrentTorrent.DownloadDir)
	if err != nil {
		log.WithFields(log.Fields{
			"path":  episode.DownloadingItem.CurrentTorrent.DownloadDir,
			"error": err,
		}).Warning("Could not remove temporary folder for download")
	}

	episode.DownloadingItem.CurrentTorrent.DownloadDir = destinationPath
	db.Client.Save(episode)

	return nil
}

func MoveMovieToLibrary(movie *Movie) error {
	log.WithFields(log.Fields{
		"movie":          media_helper.GetMovieTitle(*movie),
		"temporary_path": movie.DownloadingItem.CurrentTorrent.DownloadDir,
		"library_path":   configuration.Config.Library.MoviePath,
	}).Debug("Moving movie to library")

	destinationPath := fmt.Sprintf("%s/%s", configuration.Config.Library.MoviePath, sanitizeStringForFilename(media_helper.GetMovieTitle(*movie)))
	err := os.MkdirAll(destinationPath, 0755)
	if err != nil {
		return errors.Wrap(err, "Could not create library folder for movie")
	}

	target := fmt.Sprintf("%s/%s", destinationPath, movie.DownloadingItem.CurrentTorrent.Name)
	err = os.Rename(movie.DownloadingItem.CurrentTorrent.DownloadDir, target)
	if err != nil {
		return errors.Wrap(err, "Could not move movie to library")
	}

	err = os.Remove(movie.DownloadingItem.CurrentTorrent.DownloadDir)
	if err != nil {
		log.WithFields(log.Fields{
			"path":  movie.DownloadingItem.CurrentTorrent.DownloadDir,
			"error": err,
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
