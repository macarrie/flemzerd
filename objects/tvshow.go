package objects

import "time"

const (
	TVSHOW_RETURNING = 1
	TVSHOW_PLANNED   = 2
	TVSHOW_ENDED     = 3
	TVSHOW_UNKNOWN   = 4
)

type TvShow struct {
	Banner     string
	Poster     string
	FirstAired time.Time
	Id         int
	Overview   string
	Name       string
	Status     int
}
