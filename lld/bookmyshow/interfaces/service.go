package interfaces

import "lld/lld-2/bookmyshow/models"

type MovieService interface {
	AddMovie(mv models.Movie) error
	GetMovie(name string) (models.Movie, error)
}

type UserService interface {
	RegisterUser(name string, email string) error
	GetUser(id int) (models.User, error)
}

type TheaterService interface {
	AddTheater(name string, city string) error
	GetTheater(id int) (models.Theater, error)
	AddScreen(theaterName string, screenNumber int) error
	GetScreen(id int) (models.Screen, error)
}

type ShowService interface {
	AddShow(movieName string, theaterName string, screenName string, showTime string, showEndTime string) error
	GetShow(id int) (models.Show, error)
	GetShowsByMovie(movieId int) ([]models.Show, error)
	GetShowsByCity(city string) ([]models.Show, error)
}

type BookingService interface {
	InitiateBooking(show models.Show, seats []int, userId int) error
	ConfirmBooking(bookingId int) error
	CancelBooking(bookingId int) error
}

type PaymentService interface {
	ProcessPayment(userId int, amount float64) error
}

type NotificationService interface {
	SendNotification(userId int, message string) error
}
