package objects

import (
	"github.com/jinzhu/gorm"
)

type Notification struct {
	gorm.Model
	Type      uint8
	Movie     Movie
	MovieID   int
	TvShow    TvShow
	TvShowID  int
	Episode   Episode
	EpisodeID int
	Read      bool
	Title     string
	Content   string
}
