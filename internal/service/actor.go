package service

import (
	"cinema/internal/models"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type storeActor interface {
	AddActor(actor models.CreateActor) (uuid.UUID, error)
	GetActor(id uuid.UUID) (map[string]interface{}, error)
	UpdateActor(id uuid.UUID, actor models.UpdateActor) error
	DeleteActor(id uuid.UUID) error
	GetAllActors(limit, offset int) ([]map[string]interface{}, error)
	GetActorsWithMovies(limit, offset int) ([]map[string]interface{}, error)
}

type actor struct {
	store storeActor
}

func NewActor(store storeActor) *actor {
	return &actor{store: store}
}

// Добавление актера
func (a *actor) AddActor(actor models.CreateActor) (uuid.UUID, error) {
	return a.store.AddActor(actor)
}

// Получение актера по ID
func (a *actor) GetActor(id uuid.UUID) (*models.Actor, error) {
	rawData, err := a.store.GetActor(id)
	if err != nil {
		return nil, err
	}
	if rawData == nil {
		return nil, nil // Актёр не найден
	}

	actor := &models.Actor{
		ID:          rawData["id"].(uuid.UUID),
		Name:        rawData["name"].(string),
		Gender:      rawData["gender"].(string),
		DateOfBirth: rawData["date_of_birth"].(time.Time),
	}

	return actor, nil
}

// Обновление актера по ID
func (a *actor) UpdateActor(id uuid.UUID, actor models.UpdateActor) error {
	return a.store.UpdateActor(id, actor)
}

// Удаление актера
func (a *actor) DeleteActor(id uuid.UUID) error {
	return a.store.DeleteActor(id)
}

// Получение актеров с пагинацией
func (a *actor) GetAllActors(limit, offset int) ([]*models.Actor, error) {
	// Получаем сырые данные от репозитория
	rawActors, err := a.store.GetAllActors(limit, offset)
	if err != nil {
		return nil, err
	}

	// Маппинг сырых данных в структуру Actor
	var actors []*models.Actor
	for _, rawActor := range rawActors {
		actor := &models.Actor{
			ID:          rawActor["id"].(uuid.UUID),
			Name:        rawActor["name"].(string),
			Gender:      rawActor["gender"].(string),
			DateOfBirth: rawActor["date_of_birth"].(time.Time),
		}
		actors = append(actors, actor)
	}

	return actors, nil
}

func (s *actor) GetActorsWithMovies(limit, offset int) ([]*models.ActorWithMovies, error) {
	// Получаем сырые данные от репозитория
	data, err := s.store.GetActorsWithMovies(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get actors with movies: %w", err)
	}

	// Маппинг данных
	var actors []*models.ActorWithMovies
	// Создаем словарь для фильмов, чтобы избежать повторного маппинга
	movieMap := make(map[uuid.UUID]models.Movie)

	// Перебираем полученные данные
	for _, row := range data {
		// Маппинг актера
		actorID := row["id"].(uuid.UUID)
		actorName := row["name"].(string)
		actorGender := row["gender"].(string)
		actorBirthDate := row["date_of_birth"].(time.Time)

		// Маппинг фильмов
		moviesData := row["movies"].([]map[string]interface{})
		var movies []models.Movie
		for _, movieData := range moviesData {
			// Проверяем, есть ли фильм в словаре
			movieID := movieData["id"].(uuid.UUID)
			if movie, exists := movieMap[movieID]; exists {
				// Если фильм уже в словаре, добавляем его в список
				movies = append(movies, movie)
			} else {
				// Если фильма нет в словаре, маппируем его и добавляем в список и словарь
				newMovie := models.Movie{
					ID:          movieID,
					Title:       movieData["title"].(string),
					Description: movieData["description"].(string),
					ReleaseDate: movieData["release_date"].(time.Time),
					Rating:      movieData["rating"].(float64),
				}
				movies = append(movies, newMovie)
				movieMap[movieID] = newMovie // Сохраняем фильм в словарь
			}
		}

		// Создаем объект актера с фильмами
		actor := &models.ActorWithMovies{
			Actor: models.Actor{
				ID:          actorID,
				Name:        actorName,
				Gender:      actorGender,
				DateOfBirth: actorBirthDate,
			},
			Movies: movies,
		}

		// Добавляем актера в срез
		actors = append(actors, actor)
	}

	return actors, nil
}
