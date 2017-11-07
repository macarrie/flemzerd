package retention

import (
	"time"

	. "github.com/macarrie/flemzerd/objects"
)

type RetentionData struct {
	NotifiedEpisodes   []Episode
	DownloadedEpisodes []Episode
}

var retentionData RetentionData

func HasBeenNotified(e Episode) bool {
	for _, episode := range retentionData.NotifiedEpisodes {
		if e.Id == episode.Id {
			return true
		}
	}

	return false
}

func Load() {
	// TODO: Read retention data from file
}

func Save() {
	// TODO: Write retention data into a file
}

func AddNotifiedEpisode(e Episode) {
	if !HasBeenNotified(e) {
		retentionData.NotifiedEpisodes = append(retentionData.NotifiedEpisodes, e)
	}
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
	var notifiedEpisodes []Episode

	for _, episode := range retentionData.NotifiedEpisodes {
		if episode.Id != e.Id {
			notifiedEpisodes = append(notifiedEpisodes, episode)
		}
	}

	retentionData.NotifiedEpisodes = notifiedEpisodes
}
