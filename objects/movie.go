package objects

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Movie struct {
	gorm.Model
	MediaIds          MediaIds
	MediaIdsID        int
	Title             string
	OriginalTitle     string
	Overview          string
	Poster            string
	Date              time.Time
	Notified          bool
	Downloaded        bool
	DownloadingItem   DownloadingItem
	DownloadingItemID int
}
