package models

import "time"

type Booking struct {
	Id      int
	User    User
	Show    Show
	Seats   []Seat
	Amount  float64
	Date    time.Time
	Payment *Payment
}
