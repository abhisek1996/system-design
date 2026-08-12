package models

type SeatStatus string

const (
	Available SeatStatus = "AVAILABLE"
	Booked    SeatStatus = "BOOKED"
	Locked    SeatStatus = "LOCKED"
)

// this has to me a separate entity as the same seat can be available for one show and booked for another show
type ShowSeat struct {
	Id     string
	ShowId int
	Seat   Seat
	Status SeatStatus
}
