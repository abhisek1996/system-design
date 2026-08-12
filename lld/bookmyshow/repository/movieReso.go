package repository

import "lld/lld-2/bookmyshow/models"

type MovieRepo struct {
	MovieList map[string]models.Movie
}

func NewMovieRepo() *MovieRepo {
	return &MovieRepo{
		MovieList: make(map[string]models.Movie),
	}
}

func (m *MovieRepo) SaveMovie(mv models.Movie) error {
	m.MovieList[mv.Name] = mv
	return nil
}

func (m *MovieRepo) FindMovieByName(name string) models.Movie {
	return m.MovieList[name]
}
