package objects

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Torrent struct {
	gorm.Model
	CurrentTorrentID uint
	TorrentId        string
	Name             string
	Link             string
	DownloadDir      string
	Seeders          int
	PercentDone      float64
	TotalSize        int64
	ETA              time.Time
	RateDownload     int64
	RateUpload       int64
	Status           int
}
