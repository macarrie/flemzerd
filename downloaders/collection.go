package downloader

import (
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
	if len(downloadersCollection) == 0 {
		return errors.New("Cannot remove torrents, no downloaders are configured")
	}

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

func Download(show TvShow, e Episode, torrentList []Torrent) error {
	if retention.HasBeenDownloaded(e) || retention.IsDownloading(e) {
		return errors.New("Episode downloading or already downloaded. Skipping")
	}

	log.WithFields(log.Fields{
		"show":   show.Name,
		"season": e.Season,
		"number": e.Number,
		"name":   e.Name,
	}).Info("Starting download process")

	retention.AddDownloadingEpisode(e)
	retention.ChangeDownloadingState(e, true)

	for _, torrent := range torrentList {
		torrent.DownloadDir = fmt.Sprintf("%s/%s/Season %d/", configuration.Config.Library.ShowPath, show.Name, e.Season)

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
			return nil
		}
	}

	// If function has not returned yet, it means the download failed
	if retention.GetFailedTorrentsCount(e) > configuration.Config.System.TorrentDownloadAttemptsLimit {
		retention.AddFailedEpisode(e)
		log.WithFields(log.Fields{
			"show":   show.Name,
			"season": e.Season,
			"number": e.Number,
			"name":   e.Name,
		}).Error("Download failed, no torrents could be downloaded")

		notifier.NotifyFailedEpisode(show, e)
		retention.ChangeDownloadingState(e, false)

		return errors.New("Download failed, no torrents could be downloaded")
	}

	return nil
}

func FillToDownloadTorrentList(e Episode, list []Torrent) []Torrent {
	var torrentList []Torrent
	for _, torrent := range list {
		if !retention.IsInFailedTorrents(e, torrent) {
			torrentList = append(torrentList, torrent)
		}
	}

	if len(torrentList) < 10 {
		return torrentList
	} else {
		return torrentList[:10]
	}
}
