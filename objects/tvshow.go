package objects

import "time"

type TvShow struct {
	Banner     string
	Poster     string
	FirstAired time.Time
	Id         int
	Overview   string
	Name       string
	Status     string
}
