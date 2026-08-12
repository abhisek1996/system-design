package services

import (
	"lld/lld-2/bookmyshow/interfaces"
	"lld/lld-2/bookmyshow/models"
)

type MovieService struct {
	MovieRepo interfaces.MovieRepo
}

func (ms *MovieService) AddMovie(name string, genre string) error {
	movie := models.Movie{
		Id:    0, // generate unique id
		Name:  name,
		Genre: genre,
	}
	ms.MovieRepo.SaveMovie(movie)
	return nil
}

func (ms *MovieService) GetMovie(name string) (string, error) {
	movie := ms.MovieRepo.FindMovieByName(name)
	return movie.Name, nil
}
