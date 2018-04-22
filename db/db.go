package db

import (
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var Client *gorm.DB
var DbPath = DB_PATH
var Session SessionData

func InitDb() {
	Client.AutoMigrate(&SessionData{}, &TvShow{}, &Episode{}, &Movie{}, &MediaIds{}, &Torrent{}, &DownloadingItem{})
}

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

func Load() error {
	dbObj, err := gorm.Open("sqlite3", DbPath)
	if err != nil {
		log.Fatal("Cannot open db: ", err)
	}
	Client = dbObj

	InitDb()

	Client.First(&Session)

	log.Info("Connection to Db successful")
	return nil
}

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

func GetDownloadingEpisodes() ([]Episode, error) {
	var episodes []Episode
	Client.Table("episodes").Joins("JOIN downloading_items ON episodes.downloading_item_id = downloading_items.id").Where("downloading = 1 AND downloaded = 0").Scan(&episodes)
	return episodes, nil
}

func GetDownloadingMovies() ([]Movie, error) {
	var movies []Movie
	Client.Table("movies").Joins("JOIN downloading_items ON movies.downloading_item_id = downloading_items.id").Where("downloading = 1 AND downloaded = 0").Scan(&movies)
	return movies, nil
}

func GetDownloadedEpisodes() ([]Episode, error) {
	var episodes []Episode
	Client.Table("episodes").Where("downloaded = 1").Scan(&episodes)
	return episodes, nil
}

func GetDownloadedMovies() ([]Movie, error) {
	var movies []Movie
	Client.Table("movies").Where("downloaded = 1").Scan(&movies)
	return movies, nil
}

func SaveTraktToken(token string) {
	Session.TraktToken = token
	Client.Save(&Session)
}

func LoadTraktToken() string {
	return Session.TraktToken
}
