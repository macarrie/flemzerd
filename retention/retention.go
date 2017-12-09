package retention

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

type DownloadingEpisode struct {
	Episode         Episode
	Downloading     bool
	FailedTorrents  map[string]*Torrent
	CurrentDownload struct {
		Torrent      Torrent
		DownloaderId string
	}
}

type RetentionData struct {
	NotifiedEpisodes    map[int]Episode
	DownloadedEpisodes  map[int]Episode
	DownloadingEpisodes map[int]*DownloadingEpisode
	FailedEpisodes      map[int]Episode
}

var retentionData RetentionData

func InitStruct() {
	retentionData.NotifiedEpisodes = make(map[int]Episode)
	retentionData.DownloadedEpisodes = make(map[int]Episode)
	retentionData.DownloadingEpisodes = make(map[int]*DownloadingEpisode)
	retentionData.FailedEpisodes = make(map[int]Episode)
}

func Load() error {
	InitStruct()

	file, err := os.Open(RETENTION_FILE_PATH)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if !json.Valid(data) {
		return errors.New("Retention file contains invalid data, aborting retention load")
	}

	unmarshalError := json.Unmarshal(data, &retentionData)
	if unmarshalError != nil {
		return unmarshalError
	}

	log.Info("Retention data successfully loaded")
	return nil
}

func Save() error {
	file, err := os.OpenFile(RETENTION_FILE_PATH, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(retentionData)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	log.Debug("Retention data saved")

	return nil
}

func HasBeenNotified(e Episode) bool {
	_, ok := retentionData.NotifiedEpisodes[e.Id]
	return ok
}

func HasBeenDownloaded(e Episode) bool {
	_, ok := retentionData.DownloadedEpisodes[e.Id]
	return ok
}

func IsInDownloadProcess(e Episode) bool {
	_, ok := retentionData.DownloadingEpisodes[e.Id]
	return ok
}

func IsDownloading(e Episode) bool {
	if !IsInDownloadProcess(e) {
		return false
	}

	ep, _ := retentionData.DownloadingEpisodes[e.Id]
	return ep.Downloading
}

func IsInFailedTorrents(e Episode, t Torrent) bool {
	if !IsDownloading(e) {
		return false
	}

	_, ok := retentionData.DownloadingEpisodes[e.Id].FailedTorrents[t.Id]
	return ok
}

func GetFailedTorrentsCount(e Episode) int {
	if !IsDownloading(e) {
		return 0
	}

	return len(retentionData.DownloadingEpisodes[e.Id].FailedTorrents)
}

func HasDownloadFailed(e Episode) bool {
	_, present := retentionData.FailedEpisodes[e.Id]

	return present
}

func AddNotifiedEpisode(e Episode) {
	if !HasBeenNotified(e) {
		retentionData.NotifiedEpisodes[e.Id] = e
	}
}

func AddDownloadedEpisode(e Episode) {
	if !HasBeenDownloaded(e) {
		retentionData.DownloadedEpisodes[e.Id] = e
	}

	delete(retentionData.DownloadingEpisodes, e.Id)
}

func AddDownloadingEpisode(e Episode) {
	if !IsDownloading(e) {
		retentionData.DownloadingEpisodes[e.Id] = &DownloadingEpisode{
			Episode:        e,
			Downloading:    true,
			FailedTorrents: make(map[string]*Torrent),
		}
	}
}

func AddFailedEpisode(e Episode) {
	if !HasDownloadFailed(e) {
		retentionData.FailedEpisodes[e.Id] = e
	}
}

func AddFailedTorrent(e Episode, t Torrent) {
	if !IsDownloading(e) {
		AddDownloadingEpisode(e)
	}
	retentionData.DownloadingEpisodes[e.Id].FailedTorrents[t.Id] = &t
}

func SetCurrentDownload(e Episode, t Torrent, downloaderId string) {
	if !IsDownloading(e) {
		return
	}
	retentionData.DownloadingEpisodes[e.Id].CurrentDownload.Torrent = t
	retentionData.DownloadingEpisodes[e.Id].CurrentDownload.DownloaderId = downloaderId
}

func GetDownloadingEpisodes() ([]Episode, error) {
	var episodesList []Episode
	for _, val := range retentionData.DownloadingEpisodes {
		episodesList = append(episodesList, val.Episode)
	}
	return episodesList, nil
}

func GetCurrentDownloadFromRetention(e Episode) (Torrent, string, error) {
	if !IsDownloading(e) {
		return Torrent{}, "", errors.New("Episode is not getting downloaded, no torrents to fetch")
	}
	currentDownload := retentionData.DownloadingEpisodes[e.Id].CurrentDownload
	return currentDownload.Torrent, currentDownload.DownloaderId, nil
}

func ChangeDownloadingState(e Episode, state bool) {
	if !IsInDownloadProcess(e) {
		return
	}

	downloadingEpisode := retentionData.DownloadingEpisodes[e.Id]
	downloadingEpisode.Downloading = state
}

func CleanOldNotifiedEpisodes() {
	for _, ep := range retentionData.NotifiedEpisodes {
		// Remove notified episodes older than 14 days from retention to avoid
		if ep.Date.Before(time.Now().AddDate(0, 0, -14)) {
			RemoveNotifiedEpisode(ep)
		}
	}
}

func RemoveNotifiedEpisode(e Episode) {
	delete(retentionData.NotifiedEpisodes, e.Id)
}

func RemoveDownloadedEpisode(e Episode) {
	delete(retentionData.DownloadedEpisodes, e.Id)
}

func RemoveDownloadingEpisode(e Episode) {
	delete(retentionData.DownloadingEpisodes, e.Id)
}

func RemoveFailedEpisode(e Episode) {
	delete(retentionData.FailedEpisodes, e.Id)
}
