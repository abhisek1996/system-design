package services

import (
	"lld/lld-2/bookmyshow/interfaces"
	"lld/lld-2/bookmyshow/models"
)

type TheaterService struct {
	theaterRepo interfaces.TheaterRepo
}

// type TheaterService interface {
// 	AddTheater(name string, city string) error
// 	GetTheater(id int) (models.Theater, error)
// 	AddScreen(theaterName string, screenNumber int) error
// 	GetScreen(id int) (models.Screen, error)
// }

func (ts *TheaterService) AddTheater(name string, city string) error {
	theater := models.Theater{
		Name: name,
		City: city,
	}
	ts.theaterRepo.SaveTheater(&theater)
	return nil
}

func (ts *TheaterService) GetTheater(name string) (string, error) {

	return ts.theaterRepo.FindTheaterByName(name).Name, nil
}

func (ts *TheaterService) AddScreen(theaterName string, screenNumber int, seats []models.Seat) error {
	screen := models.Screen{
		Id:    screenNumber,
		Name:  "test",
		Seats: seats,
	}
	return ts.theaterRepo.SaveScreen(screen, theaterName)
}

// func (ts *TheaterService) GetScreen(id int, theaterName string) (models.Screen, error) {
// 	return *ts.theaterRepo.FindScreenById(id, theaterName), nil
// }
