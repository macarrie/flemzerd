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
	Client = dbObj.Set("gorm:auto_preload", true)

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

func GetTrackedMovies() ([]Movie, error) {
	var movies []Movie
	var retList []Movie
	Client.Order("created_at DESC").Find(&movies)

	for _, m := range movies {
		if !m.DownloadingItem.Downloading && !m.Downloaded {
			retList = append(retList, m)
		}
	}

	return retList, nil
}

func GetTrackedTvShows() ([]TvShow, error) {
	var tvShows []TvShow

	Client.Order("status").Order("name").Find(&tvShows)

	return tvShows, nil
}

func GetDownloadingEpisodes() ([]Episode, error) {
	var episodes []Episode
	var retList []Episode
	Client.Find(&episodes)

	for _, e := range episodes {
		if e.DownloadingItem.Downloading && !e.Downloaded {
			retList = append(retList, e)
		}
	}
	return retList, nil
}

func GetDownloadingMovies() ([]Movie, error) {
	var movies []Movie
	var retList []Movie
	Client.Find(&movies)

	for _, m := range movies {
		if m.DownloadingItem.Downloading && !m.Downloaded {
			retList = append(retList, m)
		}
	}

	return retList, nil
}

func GetDownloadedEpisodes() ([]Episode, error) {
	var episodes []Episode
	Client.Where(&Episode{
		Downloaded: true,
	}).Order("id DESC").Find(&episodes)
	return episodes, nil
}

func GetDownloadedMovies() ([]Movie, error) {
	var movies []Movie
	Client.Where(&Movie{
		Downloaded: true,
	}).Order("id DESC").Find(&movies)
	return movies, nil
}

func SaveTraktToken(token string) {
	Session.TraktToken = token
	Client.Save(&Session)
}

func LoadTraktToken() string {
	return Session.TraktToken
}
