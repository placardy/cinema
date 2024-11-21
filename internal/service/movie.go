package service

import (
	"cinema/internal/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type storeMovie interface {
	BeginTransaction() (*sql.Tx, error)
	CheckActorExists(actorID uuid.UUID) (bool, error)
	CheckMovieExists(movieID uuid.UUID) (bool, error)
	AddMovieActorRelations(tx *sql.Tx, movieID uuid.UUID, actorIDs []uuid.UUID) error
	RemoveMovieActorRelations(tx *sql.Tx, movieID uuid.UUID) error
	RemoveSelectedMovieActorRelations(tx *sql.Tx, movieID uuid.UUID, actorIDs []uuid.UUID) error
	AddMovie(tx *sql.Tx, movie models.CreateMovie) (uuid.UUID, error)
	GetMovieByID(id uuid.UUID) (map[string]interface{}, error)
	GetMoviesByActorID(actorID uuid.UUID, limit, offset int) ([]map[string]interface{}, error)
	GetMoviesWithFilters(sortBy string, order string, limit, offset int) ([]map[string]interface{}, error)
	SearchMoviesByTitleAndActor(titleFragment, actorNameFragment string, limit, offset int) ([]map[string]interface{}, error)
	UpdateMovie(id uuid.UUID, movie models.UpdateMovie) error
	DeleteMovie(id uuid.UUID) error
}

type movie struct {
	store storeMovie
}

func NewMovie(store storeMovie) *movie {
	return &movie{store: store}
}

// Обертка транзакция возвращающая err
func withTransactionError(beginTx func() (*sql.Tx, error), action func(tx *sql.Tx) error) error {
	tx, err := beginTx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = action(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Обертка транзакция возвращающая uuid err
func withTransactionUUID(beginTx func() (*sql.Tx, error), action func(tx *sql.Tx) (uuid.UUID, error)) (uuid.UUID, error) {
	var result uuid.UUID

	tx, err := beginTx()
	if err != nil {
		return result, fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err = action(tx)
	if err != nil {
		tx.Rollback()
		return result, err
	}

	err = tx.Commit()
	if err != nil {
		return result, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

// Валидация actorIDs (проверка на существование)
func (m *movie) ValidateActorID(actorID uuid.UUID) error {
	exists, err := m.store.CheckActorExists(actorID) // Проверяем одного актера
	if err != nil {
		return fmt.Errorf("failed to check actor existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("actor with ID %s not found", actorID)
	}
	return nil
}

// Валидация movieID (проверка на существование)
func (m *movie) ValidateMovieID(movieID uuid.UUID) error {
	exists, err := m.store.CheckActorExists(movieID)
	if err != nil {
		return fmt.Errorf("failed to check actor existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("movie with ID %s not found", movieID)
	}
	return nil
}

// Добавить связи фильм-актёр по movieID
func (m *movie) AddMovieActorRelations(movieID uuid.UUID, actorIDs []uuid.UUID) error {
	// 1. Проверяем существование фильма
	err := m.ValidateMovieID(movieID)
	if err != nil {
		return err
	}

	// 2. Проверяем существование актеров
	for _, actorID := range actorIDs {
		err := m.ValidateActorID(actorID)
		if err != nil {
			return err
		}
	}

	// Используем обертку для транзакций
	return withTransactionError(m.store.BeginTransaction, func(tx *sql.Tx) error {
		// 3. Добавление связей между фильмом и актерами
		err := m.store.AddMovieActorRelations(tx, movieID, actorIDs)
		if err != nil {
			return fmt.Errorf("failed to add actor-movie relations: %w", err)
		}
		return nil
	})
}

func (m *movie) UpdateMovieActorsRealations(movieID uuid.UUID, actorIDs []uuid.UUID) error {
	// Проверка существования фильма
	err := m.ValidateMovieID(movieID)
	if err != nil {
		return err
	}

	for _, actorID := range actorIDs {
		err := m.ValidateActorID(actorID)
		if err != nil {
			return err
		}
	}

	// Транзакционные операции: удаление и добавление связей
	// 2. Начинаем транзакцию
	return withTransactionError(m.store.BeginTransaction, func(tx *sql.Tx) error {
		// 3. Удаление связей между фильмом и актерами
		err = m.store.RemoveMovieActorRelations(tx, movieID)
		if err != nil {
			return fmt.Errorf("failed to remove old relations: %w", err)
		}

		err = m.store.AddMovieActorRelations(tx, movieID, actorIDs)
		if err != nil {
			return fmt.Errorf("failed to add new relations: %w", err)
		}

		return nil
	})

}

// Удалить связи фильм-актёры
func (m *movie) RemoveSelectedMovieActorRelations(movieID uuid.UUID, actorIDs []uuid.UUID) error {
	err := m.ValidateMovieID(movieID)
	if err != nil {
		return err
	}

	for _, actorID := range actorIDs {
		err := m.ValidateActorID(actorID)
		if err != nil {
			return err
		}
	}

	// 2. Начинаем транзакцию
	return withTransactionError(m.store.BeginTransaction, func(tx *sql.Tx) error {
		// 3. Удаление связей между фильмом и актерами
		err := m.store.RemoveSelectedMovieActorRelations(tx, movieID, actorIDs)
		if err != nil {
			return fmt.Errorf("failed to remove actor-movie relations: %w", err)
		}
		return nil
	})
}

// Добавление фильма
func (m *movie) AddMovie(movie models.CreateMovie) (uuid.UUID, error) {
	// 1. Проверяем существование актеров
	for _, actorID := range movie.ActorIDs {
		err := m.ValidateActorID(actorID)
		if err != nil {
			return uuid.Nil, err
		}
	}

	// Используем для начала транзакции
	return withTransactionUUID(m.store.BeginTransaction, func(tx *sql.Tx) (uuid.UUID, error) {
		// Добавляем фильм
		movieID, err := m.store.AddMovie(tx, movie)
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to add movie: %w", err)
		}

		// Добавляем связи между фильмом и актерами
		if err := m.store.AddMovieActorRelations(tx, movieID, movie.ActorIDs); err != nil {
			return uuid.Nil, fmt.Errorf("failed to add actor-movie relations: %w", err)
		}

		return movieID, nil
	})
}

// Получение фильма по ID
func (s *movie) GetMovieByID(movieID uuid.UUID) (*models.Movie, error) {
	err := s.ValidateMovieID(movieID)
	if err != nil {
		return nil, err
	}
	// Вызываем репозиторий, чтобы получить сырые данные
	rawData, err := s.store.GetMovieByID(movieID)
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

// Поиск фильмов по названию и актеру
func (m *movie) SearchMoviesByTitleAndActor(titleFragment string, actorNameFragment string, limit, offset int) ([]*models.Movie, error) {
	// Получаем сырые данные из репозитория
	rawData, err := m.store.SearchMoviesByTitleAndActor(titleFragment, actorNameFragment, limit, offset)
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
func (m *movie) DeleteMovie(movieID uuid.UUID) error {
	err := m.ValidateMovieID(movieID)
	if err != nil {
		return err
	}
	return m.store.DeleteMovie(movieID)
}
