package objects

import "time"

type Movie struct {
	Id            int
	Title         string
	OriginalTitle string
	Overview      string
	Poster        string
	Date          time.Time
}
