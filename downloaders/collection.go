package downloader

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/mediacenters"
	"github.com/macarrie/flemzerd/notifiers"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/stats"

	"github.com/macarrie/flemzerd/downloadable"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

var downloadersCollection []Downloader

type ContextStorage struct {
	Context context.Context
	Cancel  func()
}

type DownloadRoutineStruct map[uint]ContextStorage

var EpisodeDownloadRoutines DownloadRoutineStruct
var MovieDownloadRoutines DownloadRoutineStruct

var downloadRoutinesMutex sync.Mutex

func init() {
	EpisodeDownloadRoutines = make(DownloadRoutineStruct)
	MovieDownloadRoutines = make(DownloadRoutineStruct)
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

func GetTorrentStatus(t *Torrent) error {
	return downloadersCollection[0].GetTorrentStatus(t)
}

func HandleTorrentDownload(ctx context.Context, d downloadable.Downloadable, torrent Torrent) (err error, aborted bool) {
	downloadingItem := d.GetDownloadingItem()
	downloadRoutinesStruct := getDownloadRoutinesStruct(d)

	// If current torrent is set, we are recovering a download process and not adding torrents
	if downloadingItem.CurrentTorrent.ID == 0 {
		torrentId, err := AddTorrent(torrent)
		if err != nil {
			RemoveTorrent(torrent)
			downloadingItem.FailedTorrents = append(downloadingItem.FailedTorrents, torrent)
			downloadingItem.CurrentTorrent = Torrent{}
			d.SetDownloadingItem(downloadingItem)
			db.SaveDownloadable(&d)

			return errors.Wrap(err, "Couldn't add torrent in downloader. Skipping to next torrent in list"), false
		}

		downloadingItem.CurrentTorrent = torrent
		downloadingItem.CurrentDownloaderId = torrentId
		d.SetDownloadingItem(downloadingItem)
		db.SaveDownloadable(&d)
	}

	downloadingItem.CurrentTorrent = torrent
	db.Client.Save(&downloadingItem)
	_ = StartTorrent(torrent)

	downloadErr, downloadAborted := WaitForDownload(ctx, torrent)
	if downloadAborted {
		return nil, downloadAborted
	}
	if downloadErr != nil {
		if err := RemoveTorrent(torrent); err != nil {
			log.WithFields(log.Fields{
				"torrent": torrent.Name,
				"error":   err,
			}).Error("Could not remove torrent on download error")
		}
		downloadingItem.FailedTorrents = append(downloadingItem.FailedTorrents, torrent)
		d.SetDownloadingItem(downloadingItem)
		db.SaveDownloadable(&d)

		log.WithFields(log.Fields{
			"error":   downloadErr,
			"torrent": torrent.Name,
		}).Debug("Error during torrent download. Finish current torrent download")

		return errors.Wrap(downloadErr, "error during download"), false
	}

	// If function has not returned yet, download ended with no errors !
	d.GetLog().Info("Item successfully downloaded")
	notifier.NotifyDownloadedItem(d)

	downloadingItem.Pending = false
	downloadingItem.Downloading = false
	downloadingItem.Downloaded = true
	d.SetDownloadingItem(downloadingItem)
	db.SaveDownloadable(&d)

	err = MoveItemToLibrary(d)
	if err != nil {
		d.GetLog().WithFields(log.Fields{
			"temporary_path": downloadingItem.CurrentTorrent.DownloadDir,
			"library_path":   configuration.Config.Library.MoviePath,
			"error":          err,
		}).Error("Could not move item from temporary download path to library folder")
	} else {
		mediacenter.RefreshLibrary()
	}

	if err := RemoveTorrent(torrent); err != nil {
		log.WithFields(log.Fields{
			"torrent": torrent.Name,
			"error":   err,
		}).Error("Could not remove torrent after download success ")
	}

	downloadRoutinesMutex.Lock()
	ctxStore, ok := downloadRoutinesStruct[d.GetId()]
	if ok {
		ctxStore.Cancel()
	}
	downloadRoutinesMutex.Unlock()

	return nil, false
}

func WaitForDownload(ctx context.Context, t Torrent) (err error, aborted bool) {
	downloadLoopTicker := time.NewTicker(20 * time.Second)
	for {
		log.WithFields(log.Fields{
			"torrent": t.Name,
		}).Debug("Checking torrent download progress")

		err := GetTorrentStatus(&t)
		if err != nil {
			return errors.Wrap(err, "error when getting download status"), false
		}
		db.Client.Save(&t)
		status := t.Status

		switch status {
		case TORRENT_STOPPED:
			return errors.New("Torrent stopped in download client"), false
		case TORRENT_SEEDING:
			// Download complete ! Return with no error
			return nil, false
		}
		select {
		case <-downloadLoopTicker.C:
			continue
		case <-ctx.Done():
			return nil, true
		}
	}
}

func getDownloadRoutinesStruct(d downloadable.Downloadable) DownloadRoutineStruct {
	switch d.(type) {
	case *Movie:
		return MovieDownloadRoutines
	case *Episode:
		return EpisodeDownloadRoutines
	default:
		return EpisodeDownloadRoutines
	}
}

func Download(d downloadable.Downloadable) error {
	routinesStruct := getDownloadRoutinesStruct(d)

	downloadingItem := d.GetDownloadingItem()
	var recovery bool
	if downloadingItem.CurrentTorrent.ID != 0 {
		recovery = true
	} else {
		recovery = false
	}

	if downloadingItem.Downloaded || downloadingItem.Downloading {
		return errors.New("Item is currently downloading or already downloaded. Skipping")
	}

	d.GetLog().WithFields(log.Fields{
		"recovery": recovery,
	}).Info("Starting download process")

	downloadingItem.Pending = false
	downloadingItem.DownloadFailed = false
	downloadingItem.TorrentsNotFound = false
	downloadingItem.FailedTorrents = []Torrent{}
	downloadingItem.Downloading = true
	d.SetDownloadingItem(downloadingItem)
	db.SaveDownloadable(&d)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		delete(routinesStruct, d.GetId())
	}()

	downloadRoutinesMutex.Lock()
	routinesStruct[d.GetId()] = ContextStorage{
		Context: ctx,
		Cancel:  cancel,
	}
	downloadRoutinesMutex.Unlock()

	if recovery {
		AddTorrentMapping(downloadingItem.CurrentTorrent.TorrentId, downloadingItem.CurrentDownloaderId)
	}

	for _, torrent := range downloadingItem.TorrentList {
		torrent.DownloadDir = fmt.Sprintf("%s/%s/", configuration.Config.Library.CustomTmpPath, xid.New())

		if db.TorrentHasFailed(downloadingItem, torrent) {
			continue
		}

		d.SetDownloadingItem(downloadingItem)
		db.SaveDownloadable(&d)

		var torrentDownloadError error
		var downloadAborted bool

		if recovery {
			torrentDownloadError, downloadAborted = HandleTorrentDownload(ctx, d, torrent)
		} else {
			torrentDownloadError, downloadAborted = HandleTorrentDownload(ctx, d, torrent)
		}

		if downloadAborted {
			d.GetLog().Info("Download manually aborted. Cleaning up current download artifacts.")

			currentDownloadPath := downloadingItem.CurrentTorrent.DownloadDir
			if err := RemoveTorrent(downloadingItem.CurrentTorrent); err != nil {
				log.WithFields(log.Fields{
					"torrent": downloadingItem.CurrentTorrent.Name,
					"error":   err,
				}).Error("Could not remove torrent from downloader when aborting download")
			}
			if err := os.Remove(currentDownloadPath); err != nil {
				log.WithFields(log.Fields{
					"path":  currentDownloadPath,
					"error": err,
				}).Error("Could not remove downloaded data for media when aborting download")
			}

			db.Client.Unscoped().Delete(&downloadingItem)
			downloadingItem = DownloadingItem{}
			d.SetDownloadingItem(downloadingItem)
			db.SaveDownloadable(&d)

			return nil
		}
		if torrentDownloadError != nil {
			log.WithFields(log.Fields{
				"err":     torrentDownloadError,
				"torrent": torrent.Name,
			}).Warning("Couldn't dowonload torrent. Skipping to next torrent in list")
			// downloadingItem.CurrentTorrent = nil
		} else {
			return nil
		}

		// If function has not returned yet, it means the download failed
		if len(downloadingItem.FailedTorrents) > configuration.Config.System.TorrentDownloadAttemptsLimit {
			MarkDownloadAsFailed(d)
			return errors.New("Download failed, no torrents could be downloaded")
		}

		continue
	}

	return errors.New("No torrents in current torrent list could be downloaded")
}

func AbortDownload(d downloadable.Downloadable) {
	d.GetLog().Info("Aborting item download")
	downloadingItem := d.GetDownloadingItem()
	downloadRoutinesStruct := getDownloadRoutinesStruct(d)

	if !downloadingItem.Pending {
		downloadRoutinesMutex.Lock()
		ctxStore, ok := downloadRoutinesStruct[d.GetId()]
		if ok {
			ctxStore.Cancel()
		} else {
			d.GetLog().Warning("Could not find download routine to abort")
		}
		downloadRoutinesMutex.Unlock()
	}

	downloadingItem.Pending = false
	downloadingItem.Downloading = false
	downloadingItem.Downloaded = false
	downloadingItem.DownloadFailed = false
	d.SetDownloadingItem(downloadingItem)
	db.SaveDownloadable(&d)

	switch d.(type) {
	case *Movie:
		stats.Stats.Movies.Downloading -= 1
	case *Episode:
		stats.Stats.Episodes.Downloading -= 1
	}
}

func MarkDownloadAsFailed(d downloadable.Downloadable) {
	d.GetLog().Error("Download failed, no torrents could be downloaded")

	notifier.NotifyFailedDownload(d)

	downloadingItem := d.GetDownloadingItem()
	downloadingItem.DownloadFailed = true
	downloadingItem.Downloading = false
	d.SetDownloadingItem(downloadingItem)
	db.SaveDownloadable(&d)
}

func sanitizeStringForFilename(src string) string {
	reg, _ := regexp.Compile("[^a-z0-9]+")
	return reg.ReplaceAllString(strings.ToLower(src), "_")
}

func MoveItemToLibrary(d downloadable.Downloadable) error {
	downloadingItem := d.GetDownloadingItem()

	var libraryPath string
	var destinationPath string

	switch d.(type) {
	case *Movie:
		libraryPath = configuration.Config.Library.MoviePath
		destinationPath = fmt.Sprintf("%s/%s", configuration.Config.Library.MoviePath, sanitizeStringForFilename(d.GetTitle()))
	case *Episode:
		libraryPath = configuration.Config.Library.ShowPath
		episode := *d.(*Episode)
		destinationPath = fmt.Sprintf("%s/%s/season_%d/s%02de%02d", configuration.Config.Library.ShowPath, sanitizeStringForFilename(episode.TvShow.GetTitle()), episode.Season, episode.Season, episode.Number)
	}

	d.GetLog().WithFields(log.Fields{
		"temporary_path": downloadingItem.CurrentTorrent.DownloadDir,
		"library_path":   libraryPath,
	}).Debug("Moving item to library")

	err := os.MkdirAll(destinationPath, 0755)
	if err != nil {
		return errors.Wrap(err, "Could not create library folder for item")
	}

	target := fmt.Sprintf("%s/%s", destinationPath, downloadingItem.CurrentTorrent.Name)
	err = os.Rename(downloadingItem.CurrentTorrent.DownloadDir, target)
	if err != nil {
		return errors.Wrap(err, "Could not move item to library")
	}

	downloadingItem.CurrentTorrent.DownloadDir = destinationPath
	db.SaveDownloadable(&d)

	return nil
}

func FillTorrentList(d downloadable.Downloadable, list []Torrent) []Torrent {
	downloadingItem := d.GetDownloadingItem()
	var torrentList []Torrent
	for _, torrent := range list {
		if !db.TorrentHasFailed(downloadingItem, torrent) {
			torrentList = append(torrentList, torrent)
		}
	}

	// TODO: Change according to download torrent limit
	if len(torrentList) < 10 {
		return torrentList
	} else {
		return torrentList[:10]
	}
}

// GetDownloader returns the registered downloader with name "name". An non-nil error is returned if no registered downloader are found with the required name
func GetDownloader(name string) (Downloader, error) {
	for _, d := range downloadersCollection {
		mod, _ := d.Status()
		if mod.Name == name {
			return d, nil
		}
	}

	return nil, fmt.Errorf("Downloader %s not found in configuration", name)
}
