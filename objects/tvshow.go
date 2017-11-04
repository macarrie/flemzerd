package objects

import "time"

type TvShow struct {
	Aliases    []string
	Banner     string
	FirstAired time.Time
	Id         int
	Network    string
	Overview   string
	Name       string
	Status     string
}
