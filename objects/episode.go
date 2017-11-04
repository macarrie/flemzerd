package objects

import "time"

type Episode struct {
	AbsoluteNumber int
	Number         int
	Season         int
	Name           string
	Date           time.Time
	Id             int
	Overview       string
}
