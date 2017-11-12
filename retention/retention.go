package retention

import (
	"time"

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

func Load() {
	InitStruct()
}

func Save() {
	// TODO: Write retention data into a file
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
			Episode: e,
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
