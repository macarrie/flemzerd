// Package db provides access to local database and gives some shortcuts for getting certain types of objects
// The ultimate goal is to create an abstraction layer for data access so an underlying database structure change does not require to change the entire codebase
package db

import (
	"os"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/stats"
	"golang.org/x/sys/unix"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/macarrie/flemzerd/downloadable"
)

// Client represents gorm database connection
var Client *gorm.DB

// DbPath: Path of database file
var DbPath = DB_PATH

// Current session data
var Session SessionData

var DbFatal func(code int) = os.Exit

// InitDb initializes and migrates database tables
func InitDb() {
	Client.AutoMigrate(&SessionData{}, &TvShow{}, &TvSeason{}, &Episode{}, &Movie{}, &MediaIds{}, &Torrent{}, &DownloadingItem{}, &Notification{})
}

// Reset DB tables to an empty state. Mainly used in test suite.
func ResetDb() {
	Client.DropTable(&SessionData{})
	Client.DropTable(&TvShow{})
	Client.DropTable(&Episode{})
	Client.DropTable(&Movie{})
	Client.DropTable(&MediaIds{})
	Client.DropTable(&Torrent{})
	Client.DropTable(&DownloadingItem{})
	Client.DropTable(&Notification{})
	InitDb()
}

// Load initializes database connexion
func Load() error {
	dbObj, err := gorm.Open("sqlite3", DbPath)
	if err != nil {
		log.Fatal("Cannot open db: ", err)
	}

	err = unix.Access(DbPath, unix.W_OK)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"path":  DbPath,
		}).Error("Cannot write into database")
		DbFatal(1)
	}

	Client = dbObj.Set("gorm:auto_preload", true)
	Client.DB().SetMaxOpenConns(1)

	InitDb()

	Client.First(&Session)

	log.Info("Connection to Db successful")
	return nil
}

// Returns tracked movies, and an error. Tracked movies are movies retrieved from watchlists that are not already downloaded or currently downloading
// Most recent tracked movies are returned first
// A non nil error is returned if a problem was encoutered when getting movies from the database.
func GetTrackedMovies() ([]Movie, error) {
	var movies []Movie
	var retList []Movie
	// TODO: Handle error
	Client.Order("created_at DESC").Find(&movies)

	for _, m := range movies {
		if !m.DownloadingItem.Pending && !m.DownloadingItem.Downloading && !m.DownloadingItem.Downloaded {
			retList = append(retList, m)
		}
	}

	stats.Stats.Movies.Tracked = len(retList)
	return retList, nil
}

// Returns tracked shows, and an error. Tracked shows are shows added from watchlists that have been retrieved from the local database?
// Returned shows are ordered by status and title
// A non nil error is returned if a problem was encoutered when getting shows from the database.
func GetTrackedTvShows() ([]TvShow, error) {
	var tvShows []TvShow

	// TODO: Handle error
	Client.Order("status").Order("title").Find(&tvShows)

	stats.Stats.Shows.Tracked = len(tvShows)
	return tvShows, nil
}

// Gets downloading (or download pending) episodes from database
// Deleted episodes are returned too.
// Returned episodes are ordered by id (descending)
func GetDownloadingEpisodes() ([]Episode, error) {
	var episodes []Episode
	var retList []Episode
	Client.Unscoped().Where("downloading_item_id <> 0").Order("id DESC").Find(&episodes)

	for index, e := range episodes {
		if e.DownloadingItemID == 0 {
			continue
		}

		Client.Model(&e).Association("DownloadingItem").Find(&episodes[index].DownloadingItem)
		if (e.DownloadingItem.Downloading || e.DownloadingItem.Pending) && !e.DownloadingItem.Downloaded {
			Client.Model(&e).Association("TvShow").Find(&episodes[index].TvShow)
			retList = append(retList, e)
		}
	}

	stats.Stats.Episodes.Downloading = len(retList)
	return retList, nil
}

// Gets downloading movies from database
// Returned episodes are ordered by id (descending)
func GetDownloadingMovies() ([]Movie, error) {
	var movies []Movie
	var retList []Movie
	Client.Where("downloading_item_id <> 0").Order("id DESC").Find(&movies)

	for _, m := range movies {
		if (m.DownloadingItem.Downloading || m.DownloadingItem.Pending) && !m.DownloadingItem.Downloaded {
			retList = append(retList, m)
		}
	}

	stats.Stats.Movies.Downloading = len(retList)
	return retList, nil
}

// Gets downloaded episodes from database
func GetDownloadedEpisodes() ([]Episode, error) {
	var episodes []Episode
	var retList []Episode
	Client.Unscoped().Find(&episodes)

	for _, e := range episodes {
		if e.DownloadingItem.Downloaded {
			retList = append(retList, e)
		}
	}

	stats.Stats.Episodes.Downloaded = len(retList)
	return retList, nil
}

// Gets downloaded movies from database
func GetDownloadedMovies() ([]Movie, error) {
	var movies []Movie
	var retList []Movie
	Client.Find(&movies)

	for _, m := range movies {
		if m.DownloadingItem.Downloaded {
			retList = append(retList, m)
		}
	}

	stats.Stats.Movies.Downloaded = len(retList)
	return retList, nil
}

// Gets removed movies from database
func GetRemovedMovies() ([]Movie, error) {
	var movies []Movie
	var retList []Movie

	if err := Client.Unscoped().Order("created_at DESC").Find(&movies).Error; err != nil {
		log.Error("Error while getting removed movies from db: ", err)
	}

	for _, m := range movies {
		if m.DeletedAt != nil {
			retList = append(retList, m)
		}
	}

	stats.Stats.Movies.Removed = len(retList)
	return retList, nil
}

// Gets removed movies from database
func GetRemovedTvShows() ([]TvShow, error) {
	var tvShows []TvShow
	var retList []TvShow

	if err := Client.Unscoped().Order("status").Order("title").Find(&tvShows).Error; err != nil {
		log.Error("Error while getting removed shows from db: ", err)
	}

	for _, show := range tvShows {
		if show.DeletedAt != nil {
			retList = append(retList, show)
		}
	}
	stats.Stats.Shows.Removed = len(retList)
	return retList, nil
}

func GetNotifications() ([]Notification, error) {
	var notifs []Notification
	Client.Order("created_at DESC").Find(&notifs)

	return notifs, nil
}

func GetReadNotifications() ([]Notification, error) {
	var notifs []Notification
	Client.Where(&Notification{Read: true}).Order("created_at DESC").Find(&notifs)

	stats.Stats.Notifications.Read = len(notifs)
	return notifs, nil
}
func GetUnreadNotifications() ([]Notification, error) {
	var notifs []Notification
	// Use plain SQL string instead of struct because GORM does not perform where on zero values when querying with structs
	Client.Where("read = false").Find(&notifs)

	stats.Stats.Notifications.Unread = len(notifs)
	return notifs, nil
}

// Saves given token as Trakt token in database
func SaveTraktToken(token string) {
	Session.TraktToken = token
	Client.Save(&Session)
}

// Saves given token as Trakt token in database
func SaveTelegramChatID(chat_id int64) {
	Session.TelegramChatID = chat_id
	Client.Save(&Session)
}

func SaveDownloadable(d *downloadable.Downloadable) {
	switch (*d).(type) {
	case *Movie:
		movie, ok := (*d).(*Movie)
		if ok {
			if errors := Client.Save(&movie).GetErrors(); len(errors) > 0 {
				log.WithFields(log.Fields{
					"errors": errors,
					"type":   "movie",
				}).Error("Could not save downloadable")
				return
			}

			return
		}
	case *Episode:
		episode, ok := (*d).(*Episode)
		if ok {
			if errors := Client.Save(&episode).GetErrors(); len(errors) > 0 {
				log.WithFields(log.Fields{
					"errors": errors,
					"type":   "episode",
				}).Error("Could not save downloadable")
				return
			}

			return
		}
	}
}
