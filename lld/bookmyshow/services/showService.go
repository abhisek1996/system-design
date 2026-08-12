package services

import "lld/lld-2/bookmyshow/interfaces"

type ShowService struct {
	showRepo interfaces.ShowRepo
}

// type ShowService interface {
// 	AddShow(movieName string, theaterName string, screenName string, showTime string) error
// 	GetShow(id int) (models.Show, error)
// 	GetShowsByMovie(movieId int) ([]models.Show, error)
// 	GetShowsByCity(city string) ([]models.Show, error)
// }

func (ss *ShowService) AddShow(movieId string, theaterId string, screenName string, startTime string, endTime string) error {
	// create show object and save to repo
	return nil
}

func (ss *ShowService) GetShow(id int) (string, error) {
	return ss.showRepo.FindShowById(id).Movie.Name, nil
}

func (ss *ShowService) GetShowsByMovie(movieId string) ([]string, error) {
	shows := ss.showRepo.FindShowsByMovieId(movieId)
	var showNames []string
	for _, show := range shows {
		showNames = append(showNames, show.Movie.Name)
	}
	return showNames, nil
}

func (ss *ShowService) GetShowsByCity(city string) ([]string, error) {
	shows := ss.showRepo.FindShowsByCity(city)
	var showNames []string
	for _, show := range shows {
		showNames = append(showNames, show.Movie.Name)
	}
	return showNames, nil
}
