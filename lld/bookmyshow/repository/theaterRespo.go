package repository

import (
	"fmt"
	"lld/lld-2/bookmyshow/models"
)

type TheaterRepo struct {
	TheaterList map[string]*models.Theater
}

func NewTheaterRepo() *TheaterRepo {
	return &TheaterRepo{
		TheaterList: make(map[string]*models.Theater),
	}
}

func (tr *TheaterRepo) SaveTheater(t *models.Theater) error {
	tr.TheaterList[t.Name] = t
	return nil
}

func (tr *TheaterRepo) FindTheaterByName(name string) models.Theater {
	return *tr.TheaterList[name]
}

func (tr *TheaterRepo) SaveScreen(s models.Screen, theaterName string) error {
	theater := tr.TheaterList[theaterName]
	if theater == nil {
		return fmt.Errorf("theater not found")
	}
	theater.Screens = append(theater.Screens, s)
	tr.TheaterList[theater.Name] = theater
	return nil
}

// func (tr *TheaterRepo) FindScreenById(id int, theaterName string) (*models.Screen, error) {
// 	theater := tr.TheaterList[theaterName]
// 	if theater == nil {
// 		return nil, fmt.Errorf("theater not found")
// 	}

// 	for _, screen := range theater.Screens {
// 		if screen.Id == id {
// 			return &screen, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("screen not found")
// }
