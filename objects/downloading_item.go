package objects

import "github.com/jinzhu/gorm"

type DownloadingItem struct {
	gorm.Model
	Pending             bool
	Downloading         bool
	Downloaded          bool
	FailedTorrents      []Torrent `gorm:"foreignkey:FailedTorrentID"`
	CurrentTorrent      Torrent   `gorm:"foreignkey:CurrentTorrentID"`
	CurrentDownloaderId string
	DownloadFailed      bool
	AbortPending        bool
	TorrentsNotFound    bool
}
