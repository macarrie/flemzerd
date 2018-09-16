// Package db provides access to local database and gives some shortcuts for getting certain types of objects
// The ultimate goal is to create an abstraction layer for data access so an underlying database structure change does not require to change the entire codebase
package db

import (
	"os"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"golang.org/x/sys/unix"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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

// Returns true if torrent is in the DownloadingItem FailedTorrents list, else returns false
func TorrentHasFailed(d DownloadingItem, t Torrent) bool {
	if !d.Downloading {
		return false
	}

	for _, torrent := range d.FailedTorrents {
		if t.TorrentId == torrent.TorrentId {
			return true
		}
	}

	return false
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

	return retList, nil
}

// Returns tracked shows, and an error. Tracked shows are shows added from watchlists that have been retrieved from the local database?
// Returned shows are ordered by status and name
// A non nil error is returned if a problem was encoutered when getting shows from the database.
func GetTrackedTvShows() ([]TvShow, error) {
	var tvShows []TvShow

	// TODO: Handle error
	Client.Order("status").Order("name").Find(&tvShows)

	return tvShows, nil
}

// Gets downloading (or download pending) episodes from database
// Deleted episodes are returned too.
// Returned episodes are ordered by id (descending)
func GetDownloadingEpisodes() ([]Episode, error) {
	var episodes []Episode
	var retList []Episode
	Client.Unscoped().Find(&episodes).Order("id DESC")

	for _, e := range episodes {
		if (e.DownloadingItem.Downloading || e.DownloadingItem.Pending) && !e.DownloadingItem.Downloaded {
			retList = append(retList, e)
		}
	}
	return retList, nil
}

// Gets downloading movies from database
// Returned episodes are ordered by id (descending)
func GetDownloadingMovies() ([]Movie, error) {
	var movies []Movie
	var retList []Movie
	Client.Find(&movies).Order("id DESC")

	for _, m := range movies {
		if (m.DownloadingItem.Downloading || m.DownloadingItem.Pending) && !m.DownloadingItem.Downloaded {
			retList = append(retList, m)
		}
	}

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

	return retList, nil
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
