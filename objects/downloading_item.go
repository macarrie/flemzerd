package objects

import "github.com/jinzhu/gorm"

type DownloadingItem struct {
	gorm.Model
	Downloading         bool
	FailedTorrents      []Torrent `gorm:"foreignkey:FailedTorrentID"`
	CurrentTorrent      Torrent   `gorm:"foreignkey:CurrentTorrentID"`
	CurrentDownloaderId string
	DownloadFailed      bool
}
