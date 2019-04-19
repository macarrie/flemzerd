package stats

import (
	"expvar"

	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"

	. "github.com/macarrie/flemzerd/objects"
)

type StatsFields struct {
	Movies struct {
		Tracked     *expvar.Int
		Downloading *expvar.Int
		Downloaded  *expvar.Int
		Removed     *expvar.Int
	}
	Shows struct {
		Tracked     *expvar.Int
		Downloading *expvar.Int
		Downloaded  *expvar.Int
		Removed     *expvar.Int
	}
	Notifications struct {
		Read   *expvar.Int
		Unread *expvar.Int
		Count  *expvar.Int
	}
}

var Stats StatsFields

func LoadStats() {
	initMoviesStats(&Stats)
	initShowsStats(&Stats)
	initEpisodesStats(&Stats)
	initNotificationsStats(&Stats)
}

func initMoviesStats(stats *StatsFields) {
	trackedMovies, err := db.GetTrackedMovies()
	if err != nil {
		log.Error("Error while getting downloading movies from db: ", err)
	}

	downloadingMovies, err := db.GetDownloadingMovies()
	if err != nil {
		log.Error("Error while getting downloading movies from db: ", err)
	}

	downloadedMovies, err := db.GetDownloadedMovies()
	if err != nil {
		log.Error("Error while getting downloaded movies from db: ", err)
	}

	stats.Movies.Tracked.Set(int64(len(trackedMovies)))
	stats.Movies.Tracked.Set(int64(len(downloadingMovies)))
	stats.Movies.Tracked.Set(int64(len(downloadedMovies)))
}

func initShowsStats(stats *StatsFields) {
	trackedShows, err := db.GetTrackedTvShows()
	if err != nil {
		log.Error("Error while gettings tracked shows from db: ", err)
	}
	stats.Shows.Tracked.Set(int64(len(trackedShows)))
}

func initEpisodesStats(stats *StatsFields) {
	downloadingEpisodes, err := db.GetDownloadingEpisodes()
	if err != nil {
		log.Error("Error while getting downloading episodes from db: ", err)
	}
	downloadedEpisodes, err := db.GetDownloadedEpisodes()
	if err != nil {
		log.Error("Error while getting downloaded episodes from db: ", err)
	}

	stats.Shows.Tracked.Set(int64(len(downloadingEpisodes)))
	stats.Shows.Tracked.Set(int64(len(downloadedEpisodes)))
}

func initNotificationsStats(stats *StatsFields) {
	var readNotifications []Notification
	db.Client.Where(&Notification{Read: true}).Order("created_at DESC").Find(&readNotifications)

	var unreadNotifications []Notification
	// Use plain SQL string instead of struct because GORM does not perform where on zero values when querying with structs
	db.Client.Where("read = false").Find(&unreadNotifications)

	stats.Notifications.Read.Set(int64(len(readNotifications)))
	stats.Notifications.Unread.Set(int64(len(unreadNotifications)))
	stats.Notifications.Count.Set(int64(len(readNotifications)) + int64(len(unreadNotifications)))
}
