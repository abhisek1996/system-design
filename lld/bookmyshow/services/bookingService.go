package services

import "lld/lld-2/bookmyshow/interfaces"

type BookingService struct {
	// mu           mutex.Lock
	ShowseatRepo interfaces.ShowSeatRepo
}

func NewBookingService(showseatRepo interfaces.ShowSeatRepo) *BookingService {
	return &BookingService{
		ShowseatRepo: showseatRepo,
	}
}

func (bs *BookingService) InitiateBooking(showId int, seats []int, userId int) error {
	// bs.mu.Lock()
	// defer bs.mu.Unlock()

	// 1. Lock the seats

	// 2. Create a booking record with status "pending"
	
	for _, seat := range seats {
		showSeat := bs.ShowseatRepo.LockSeats(showId, seat)
	}

	// 3. Return booking id to user for confirmation




	return nil
}

func (bs *BookingService) ConfirmBooking(bookingId int) error {
	// implementation
	return nil
}

func (bs *BookingService) CancelBooking(bookingId int) error {
	// implementation
	return nil
}
