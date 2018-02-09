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

func IsAlive() error {
	var downloaderAliveError error
	for _, downloader := range downloadersCollection {
		downloaderAliveError = downloader.IsAlive()
		if downloaderAliveError == nil {
			return nil
		} else {
			log.WithFields(log.Fields{
				"error": downloaderAliveError,
			}).Warning("Downloader not alive")
		}
	}
	return downloaderAliveError
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

func HandleTorrentDownload(e Episode, torrent Torrent, recovery bool) error {
	if !recovery {
		torrentId, err := AddTorrent(torrent)
		if err != nil {
			return errors.New("Couldn't add torrent in downloader. Skipping to next torrent in list")
		}
		retention.SetCurrentDownload(e, torrent, torrentId)
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
			retention.AddFailedTorrent(e, torrent)

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

		torrentDownload := HandleTorrentDownload(e, torrent, false)
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
			if !retention.HasBeenDownloaded(e) {
				notifier.NotifyDownloadedEpisode(show, e)
			}
		}
		continue
	}

	// If function has not returned yet, it means the download failed
	if retention.GetFailedTorrentsCount(e) > configuration.Config.System.TorrentDownloadAttemptsLimit {
		MarkFailedDownload(show, e)
		return errors.New("Download failed, no torrents could be downloaded")
	}

	return nil
}

func MarkFailedDownload(show TvShow, e Episode) {
	retention.AddFailedEpisode(e)
	log.WithFields(log.Fields{
		"show":   show.Name,
		"season": e.Season,
		"number": e.Number,
		"name":   e.Name,
	}).Error("Download failed, no torrents could be downloaded")

	notifier.NotifyFailedEpisode(show, e)
	retention.ChangeDownloadingState(e, false)
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

func RecoverFromRetention() {
	downloadingEpisodesFromRetention, err := retention.GetDownloadingEpisodes()
	if err != nil {
		log.Error(err)
		return
	}

	if len(downloadingEpisodesFromRetention) != 0 {
		log.Debug("Launching watch threads for downloading episodes found in retention")
	}
	for _, ep := range downloadingEpisodesFromRetention {
		torrent, downloaderId, err := retention.GetCurrentDownloadFromRetention(ep)
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

		go HandleTorrentDownload(ep, torrent, true)
	}
}
