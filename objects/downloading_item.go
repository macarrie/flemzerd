package objects

import "github.com/jinzhu/gorm"

type DownloadingItem struct {
	gorm.Model
	Pending             bool
	Downloading         bool
	Downloaded          bool      `json:"downloaded"`
	FailedTorrents      []Torrent `gorm:"foreignkey:FailedTorrentID"`
	CurrentTorrent      Torrent   `gorm:"foreignkey:CurrentTorrentID"`
	CurrentDownloaderId string
	DownloadFailed      bool
	TorrentsNotFound    bool
}
