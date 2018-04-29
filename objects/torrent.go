package objects

import "github.com/jinzhu/gorm"

type Torrent struct {
	gorm.Model
	FailedTorrentID  uint
	CurrentTorrentID uint
	TorrentId        string
	Name             string
	Link             string
	DownloadDir      string
	Seeders          int
}
