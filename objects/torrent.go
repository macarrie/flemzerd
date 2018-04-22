package objects

import "github.com/jinzhu/gorm"

type Torrent struct {
	gorm.Model
	DownloadingItemID int
	TorrentId         string
	Name              string
	Link              string
	DownloadDir       string
	Seeders           int
}
