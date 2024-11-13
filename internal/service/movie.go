package service

import (
	"cinema/internal/models"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type storeMovie interface {
	AddMovie(movie models.CreateMovie) (uuid.UUID, error)
	AddMovieActorRelation(actorID, movieID uuid.UUID) error
	RemoveMovieActorRelation(actorID, movieID uuid.UUID) error
	GetMovieByID(id uuid.UUID) (map[string]interface{}, error)
	GetMoviesByActorID(actorID uuid.UUID, limit, offset int) ([]map[string]interface{}, error)
	GetMoviesWithFilters(sortBy string, order string, limit, offset int) ([]map[string]interface{}, error)
	SearchMoviesByTitle(titleFragment string, limit, offset int) ([]map[string]interface{}, error)
	SearchMoviesByActorName(actorNameFragment string, limit, offset int) ([]map[string]interface{}, error)
	UpdateMovie(id uuid.UUID, movie models.UpdateMovie) error
	DeleteMovie(id uuid.UUID) error
	GetAllMovies(limit, offset int) ([]map[string]interface{}, error)
}

type movie struct {
	store storeMovie
}

func NewMovie(store storeMovie) *movie {
	return &movie{store: store}
}

// Добавление фильма
func (s *movie) AddMovie(rawData map[string]interface{}) (uuid.UUID, error) {
	// Сначала маппим rawData в структуру CreateMovie
	createMovie := models.CreateMovie{
		Title:       rawData["title"].(string),
		Description: rawData["description"].(string),
		ReleaseDate: rawData["release_date"].(time.Time),
		Rating:      rawData["rating"].(float64),
	}

	// Валидация данных перед добавлением
	if err := createMovie.Validate(); err != nil {
		return uuid.Nil, fmt.Errorf("validation failed: %w", err)
	}

	// Передаем данные в репозиторий для добавления
	return s.store.AddMovie(createMovie)
}

// Добавление отношения между актером и фильмом
func (m *movie) AddMovieActorRelation(actorID, movieID uuid.UUID) error {
	return m.store.AddMovieActorRelation(actorID, movieID)
}

// Удаление отношения между актером и фильмом
func (m *movie) RemoveMovieActorRelation(actorID, movieID uuid.UUID) error {
	return m.store.RemoveMovieActorRelation(actorID, movieID)
}

// Получение фильма по ID
func (s *movie) GetMovieByID(id uuid.UUID) (*models.Movie, error) {
	// Вызываем репозиторий, чтобы получить сырые данные
	rawData, err := s.store.GetMovieByID(id)
	if err != nil {
		return nil, err
	}
	if rawData == nil {
		return nil, nil // Фильм не найден
	}

	// Маппим сырые данные в структуру Movie
	movie := &models.Movie{
		ID:          rawData["id"].(uuid.UUID),
		Title:       rawData["title"].(string),
		Description: rawData["description"].(string),
		ReleaseDate: rawData["release_date"].(time.Time),
		Rating:      rawData["rating"].(float64),
	}

	return movie, nil
}

// Получение фильмов по ID актера с пагинацией
func (m *movie) GetMoviesByActorID(actorID uuid.UUID, limit, offset int) ([]*models.Movie, error) {
	rawData, err := m.store.GetMoviesByActorID(actorID, limit, offset)
	if err != nil {
		return nil, err
	}
	var movies []*models.Movie
	for _, data := range rawData {
		movie := &models.Movie{
			ID:          data["id"].(uuid.UUID),
			Title:       data["title"].(string),
			Description: data["description"].(string),
			ReleaseDate: data["release_date"].(time.Time),
			Rating:      data["rating"].(float64),
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

// Получение фильмов с сортировкой и пагинацией
func (m *movie) GetMoviesWithFilters(sortBy string, order string, limit, offset int) ([]*models.Movie, error) {
	// Получаем сырые данные из репозитория
	rawData, err := m.store.GetMoviesWithFilters(sortBy, order, limit, offset)
	if err != nil {
		return nil, err
	}

	// Маппим сырые данные в структуры моделей
	var movies []*models.Movie
	for _, data := range rawData {
		movie := &models.Movie{
			ID:          data["id"].(uuid.UUID),
			Title:       data["title"].(string),
			Description: data["description"].(string),
			ReleaseDate: data["release_date"].(time.Time),
			Rating:      data["rating"].(float64),
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// Поиск фильмов по названию
func (m *movie) SearchMoviesByTitle(titleFragment string, limit, offset int) ([]*models.Movie, error) {
	// Получаем сырые данные из репозитория
	rawData, err := m.store.SearchMoviesByTitle(titleFragment, limit, offset)
	if err != nil {
		return nil, err
	}

	// Маппим сырые данные в структуры моделей
	var movies []*models.Movie
	for _, data := range rawData {
		movie := &models.Movie{
			ID:          data["id"].(uuid.UUID),
			Title:       data["title"].(string),
			Description: data["description"].(string),
			ReleaseDate: data["release_date"].(time.Time),
			Rating:      data["rating"].(float64),
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// Поиск фильмов по имени актера
func (m *movie) SearchMoviesByActorName(actorNameFragment string, limit, offset int) ([]*models.Movie, error) {
	// Получаем сырые данные из репозитория
	rawData, err := m.store.SearchMoviesByActorName(actorNameFragment, limit, offset)
	if err != nil {
		return nil, err
	}

	// Маппим сырые данные в структуры моделей
	var movies []*models.Movie
	for _, data := range rawData {
		movie := &models.Movie{
			ID:          data["id"].(uuid.UUID),
			Title:       data["title"].(string),
			Description: data["description"].(string),
			ReleaseDate: data["release_date"].(time.Time),
			Rating:      data["rating"].(float64),
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// Обновление информации о фильме
func (m *movie) UpdateMovie(id uuid.UUID, movie models.UpdateMovie) error {
	return m.store.UpdateMovie(id, movie)
}

// Удаление фильма по ID
func (m *movie) DeleteMovie(id uuid.UUID) error {
	return m.store.DeleteMovie(id)
}

// Получение всех фильмов с пагинацией
func (m *movie) GetAllMovies(limit, offset int) ([]*models.Movie, error) {
	// Получаем сырые данные из репозитория
	rawMovies, err := m.store.GetAllMovies(limit, offset)
	if err != nil {
		return nil, err
	}

	var movies []*models.Movie
	// Маппим сырые данные в объекты
	for _, rawMovie := range rawMovies {
		movie := &models.Movie{
			ID:          rawMovie["id"].(uuid.UUID),
			Title:       rawMovie["title"].(string),
			Description: rawMovie["description"].(string),
			ReleaseDate: rawMovie["release_date"].(time.Time),
			Rating:      rawMovie["rating"].(float64),
		}
		movies = append(movies, movie)
	}

	return movies, nil
}
