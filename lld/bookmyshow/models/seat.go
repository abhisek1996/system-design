package models

type SeatType string

const (
	SILVER   SeatType = "SILVER"
	GOLD     SeatType = "GOLD"
	PLATINUM SeatType = "PLATINUM"
)

type Seat struct {
	Id       int
	Row      int
	Col      int
	SeatType SeatType
}
