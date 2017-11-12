package downloader

import (
	"errors"
	"time"

	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/notifiers"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/retention"
)

var downloadersCollection []Downloader

func AddDownloader(d Downloader) {
	downloadersCollection = append(downloadersCollection, d)
}

func AddTorrent(t Torrent) error {
	if len(downloadersCollection) == 0 {
		return errors.New("Cannot add torrents, no downloaders are configured")
	}

	err := downloadersCollection[0].AddTorrent(t)
	return err
}

func StartTorrent(t Torrent) error {
	return nil
}

func RemoveTorrent(t Torrent) error {
	return downloadersCollection[0].RemoveTorrent(t)
}

func GetTorrentStatus(t Torrent) (int, error) {
	return downloadersCollection[0].GetTorrentStatus(t)
}

func WaitForDownload(t Torrent) error {
	downloadLoopTicker := time.NewTicker(10 * time.Second)
	for {
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

func Download(show TvShow, e Episode, torrentList []Torrent) {
	if retention.HasBeenDownloaded(e) || retention.IsDownloading(e) {
		return
	}

	log.WithFields(log.Fields{
		"show":   show.Name,
		"season": e.Season,
		"number": e.Number,
		"name":   e.Name,
	}).Info("Starting download process")

	retention.AddDownloadingEpisode(e)

	for _, torrent := range torrentList {

		if retention.IsInFailedTorrents(e, torrent) {
			continue
		}

		err := AddTorrent(torrent)
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err.Error(),
				"torrent": torrent.Name,
			}).Warning("Couldn't add torrent in downloader. Skipping to next torrent in list")
			continue
		}

		StartTorrent(torrent)

		retryCount := 0
		downloadSuccess := true

		// Try twice to download a torrent before marking it as rubbish
		downloadErr := WaitForDownload(torrent)
		if downloadErr != nil {
			RemoveTorrent(torrent)
			AddTorrent(torrent)
			retryCount++
			retryErr := WaitForDownload(torrent)
			if retryErr != nil {
				RemoveTorrent(torrent)
				retention.AddFailedTorrent(e, torrent)
				downloadSuccess = false
				log.WithFields(log.Fields{
					"err":     downloadErr.Error(),
					"torrent": torrent.Name,
				}).Warning("Couldn't download torrent. Skipping to next torrent in list")
				continue
			}
		}

		if downloadSuccess {
			retention.AddDownloadedEpisode(e)
			RemoveTorrent(torrent)
			notifier.NotifyDownloadedEpisode(show, e)
			log.WithFields(log.Fields{
				"show":   show.Name,
				"season": e.Season,
				"number": e.Number,
				"name":   e.Name,
			}).Info("Episode successfully downloaded")
			return
		}
	}

	// If function has not returned yet, it means the download failed
	retention.AddFailedEpisode(e)
	log.WithFields(log.Fields{
		"show":   show.Name,
		"season": e.Season,
		"number": e.Number,
		"name":   e.Name,
	}).Error("Download failed, no torrents could be downloaded")
	notifier.NotifyFailedEpisode(show, e)
}
