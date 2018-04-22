package objects

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	TVSHOW_RETURNING = 1
	TVSHOW_PLANNED   = 2
	TVSHOW_ENDED     = 3
	TVSHOW_UNKNOWN   = 4
)

type TvShow struct {
	gorm.Model
	MediaIds   MediaIds
	MediaIdsID int
	Episodes   []Episode
	Banner     string
	Poster     string
	FirstAired time.Time
	Overview   string
	Name       string
	Status     int
}
