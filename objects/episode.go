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
	Name              string
	Date              time.Time
	Overview          string
	Notified          bool
	Downloaded        bool
	DownloadingItem   DownloadingItem
	DownloadingItemID int
}
