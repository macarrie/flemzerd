package objects

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Episode struct {
	gorm.Model
	MediaIds          MediaIds
	MediaIdsID        int
	TvShow            TvShow
	TvShowID          int
	AbsoluteNumber    int
	Number            int
	Season            int
	Title             string
	Date              time.Time
	Overview          string
	Notified          bool
	DownloadingItem   DownloadingItem
	DownloadingItemID int
}
