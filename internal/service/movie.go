package service

import (
	"cinema/internal/models"

	"github.com/google/uuid"
)

type storeMovie interface {
	AddMovie(movie models.CreateMovie) (uuid.UUID, error)
	AddMovieActorRelation(actorID, movieID uuid.UUID) error
	RemoveMovieActorRelation(actorID, movieID uuid.UUID) error
	GetMovieByID(id uuid.UUID) (*models.Movie, error)
	GetMoviesByActorID(actorID uuid.UUID, limit, offset int) ([]*models.Movie, error)
	GetMoviesWithFilters(sortBy string, order string, limit, offset int) ([]*models.Movie, error)
	SearchMoviesByTitle(titleFragment string, limit, offset int) ([]*models.Movie, error)
	SearchMoviesByActorName(actorNameFragment string, limit, offset int) ([]*models.Movie, error)
	UpdateMovie(id uuid.UUID, movie models.UpdateMovie) error
	DeleteMovie(id uuid.UUID) error
	GetAllMovies(limit, offset int) ([]*models.Movie, error)
}

type movie struct {
	store storeMovie
}

func NewMovie(store storeMovie) movie {
	return movie{store: store}
}

// Добавление фильма
func (m movie) AddMovie(movie models.CreateMovie) (uuid.UUID, error) {
	return m.store.AddMovie(movie)
}

// Добавление отношения между актером и фильмом
func (m movie) AddMovieActorRelation(actorID, movieID uuid.UUID) error {
	return m.store.AddMovieActorRelation(actorID, movieID)
}

// Удаление отношения между актером и фильмом
func (m movie) RemoveMovieActorRelation(actorID, movieID uuid.UUID) error {
	return m.store.RemoveMovieActorRelation(actorID, movieID)
}

// Получение фильма по ID
func (m movie) GetMovieByID(id uuid.UUID) (*models.Movie, error) {
	return m.store.GetMovieByID(id)
}

// Получение фильмов по ID актера с пагинацией
func (m movie) GetMoviesByActorID(actorID uuid.UUID, limit, offset int) ([]*models.Movie, error) {
	return m.store.GetMoviesByActorID(actorID, limit, offset)
}

// Получение фильмов с сортировкой и пагинацией
func (m movie) GetMoviesWithFilters(sortBy string, order string, limit, offset int) ([]*models.Movie, error) {
	return m.store.GetMoviesWithFilters(sortBy, order, limit, offset)
}

// Поиск фильмов по заголовку
func (m movie) SearchMoviesByTitle(titleFragment string, limit, offset int) ([]*models.Movie, error) {
	return m.store.SearchMoviesByTitle(titleFragment, limit, offset)
}

// Поиск фильмов по имени актера
func (m movie) SearchMoviesByActorName(actorNameFragment string, limit, offset int) ([]*models.Movie, error) {
	return m.store.SearchMoviesByActorName(actorNameFragment, limit, offset)
}

// Обновление информации о фильме
func (m movie) UpdateMovie(id uuid.UUID, movie models.UpdateMovie) error {
	return m.store.UpdateMovie(id, movie)
}

// Удаление фильма по ID
func (m movie) DeleteMovie(id uuid.UUID) error {
	return m.store.DeleteMovie(id)
}

// Получение всех фильмов с пагинацией
func (m movie) GetAllMovies(limit, offset int) ([]*models.Movie, error) {
	return m.store.GetAllMovies(limit, offset)
}
