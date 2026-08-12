package repository

import "lld/lld-2/bookmyshow/models"

type ShowRepo struct {
	showList []models.Show
}

func NewShowRepo() *ShowRepo {
	return &ShowRepo{
		showList: make([]models.Show, 0),
	}
}

func (s *ShowRepo) SaveShow(show models.Show) error {
	s.showList = append(s.showList, show)
	return nil
}

func (s *ShowRepo) FindShowById(id int) models.Show {
	return s.showList[id]
}

func (s *ShowRepo) FindShowsByMovieId(movieId string) []models.Show {
	// code to find shows by movie id from database
	var shows []models.Show
	for _, show := range s.showList {
		if show.Movie.Name == movieId {
			shows = append(shows, show)
		}
	}
	return shows
}

func (s *ShowRepo) FindShowsByCity(city string) []models.Show {
	// code to find shows by city from database
	var shows []models.Show
	for _, show := range s.showList {
		if show.Theater.City == city {
			shows = append(shows, show)
		}
	}
	return shows
}

func (s *ShowRepo) CreateShow(movieId string, theaterId string, screenName string, startTime string, endTime string) error {
	// generate show seats and save all show seats.
	// add a show
	return nil
}
