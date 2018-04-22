package objects

import "github.com/jinzhu/gorm"

type DownloadingItem struct {
	gorm.Model
	Downloading         bool
	FailedTorrents      []Torrent
	CurrentTorrent      Torrent
	CurrentDownloaderId string
	DownloadFailed      bool
}
