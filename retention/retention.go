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
	Episode        Episode
	FailedTorrents map[string]Torrent
}

type RetentionData struct {
	NotifiedEpisodes    map[int]Episode
	DownloadedEpisodes  map[int]Episode
	DownloadingEpisodes map[int]DownloadingEpisode
}

var retentionData RetentionData

func HasBeenNotified(e Episode) bool {
	_, ok := retentionData.NotifiedEpisodes[e.Id]
	return ok
}

func HasBeenDownloaded(e Episode) bool {
	_, ok := retentionData.DownloadedEpisodes[e.Id]
	return ok
}

func IsDownloading(e Episode) bool {
	_, ok := retentionData.DownloadingEpisodes[e.Id]
	return ok
}

func IsInFailedTorrents(e Episode, t Torrent) bool {
	if !IsDownloading(e) {
		return false
	}

	_, ok := retentionData.DownloadingEpisodes[e.Id].FailedTorrents[t.Id]
	return ok
}

func InitStruct() {
	retentionData.NotifiedEpisodes = make(map[int]Episode)
	retentionData.DownloadedEpisodes = make(map[int]Episode)
	retentionData.DownloadingEpisodes = make(map[int]DownloadingEpisode)
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

	return nil
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
		retentionData.DownloadingEpisodes[e.Id] = DownloadingEpisode{
			Episode:        e,
			FailedTorrents: make(map[string]Torrent),
		}
	}
}

func AddFailedEpisode(e Episode) {
	// TODO
}

func AddFailedTorrent(e Episode, t Torrent) {
	if !IsDownloading(e) {
		AddDownloadingEpisode(e)
	}

	retentionData.DownloadingEpisodes[e.Id].FailedTorrents[t.Id] = t
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
