package interfaces

import "lld/lld-2/bookmyshow/models"

// for updating and fetching from database

type MovieRepo interface {
	SaveMovie(mv models.Movie) error
	FindMovieByName(name string) models.Movie
}

type UserRepo interface {
	SaveUser(u models.User) error
	FindUserById(id int) models.User
}

type TheaterRepo interface {
	SaveTheater(t *models.Theater) error
	FindTheaterByName(name string) models.Theater
	SaveScreen(s models.Screen, theaterName string) error
	// FindScreenById(id int, theaterName string) (*models.Screen, error)
}

type ShowRepo interface {
	SaveShow(s models.Show) error
	FindShowById(id int) models.Show
	FindShowsByMovieId(movieId string) []models.Show
	FindShowsByCity(city string) []models.Show
}

type ShowSeatRepo interface {
	SaveShowSeats(s []models.ShowSeat) error
	FindShowSeatsByShowId(showId int) []models.ShowSeat
	UpdateSeatStatus(showId int, seatId int, status models.SeatStatus) error
	BulkUpdateSeatStatus(showId int, seatIds []int, status models.SeatStatus) error
	LockSeats(showId int, seatIds []int) error
	ReleaseExpiredLocks() error // should be triggered by a background job to release locks after a certain timeout
}

// type ScreenRepo interface {
// 	SaveScreen(s models.Screen) error
// 	FindScreenById(id int) models.Screen
// }
