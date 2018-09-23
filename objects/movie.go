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
	CustomTitle       string
	Overview          string
	Poster            string
	Date              time.Time
	Notified          bool
	DownloadingItem   DownloadingItem
	DownloadingItemID int
	UseDefaultTitle   bool
}
