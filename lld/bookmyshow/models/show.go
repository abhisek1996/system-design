package models

import "time"

type Show struct {
	Id        int
	Movie     Movie
	Theater   Theater
	Screen    Screen
	StartTime time.Time
	EndTime   time.Time
}
