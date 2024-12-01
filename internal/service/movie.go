package service

import (
	"cinema/internal/models"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type storeMovie interface {
	BeginTransaction() (*sql.Tx, error)
	CheckMovieExists(movieID uuid.UUID) (bool, error)
	CheckActorsExist(actorIDs []uuid.UUID) (bool, error)
	AddMovieActorRelations(tx *sql.Tx, movieID uuid.UUID, actorIDs []uuid.UUID) error
	RemoveMovieActorRelations(tx *sql.Tx, movieID uuid.UUID) error
	RemoveSelectedMovieActorRelations(tx *sql.Tx, movieID uuid.UUID, actorIDs []uuid.UUID) error
	CreateMovie(tx *sql.Tx, movie models.CreateMovie) (uuid.UUID, error)
	GetMovieByID(id uuid.UUID) (map[string]interface{}, error)
	GetMoviesByActorID(actorID uuid.UUID, limit, offset int) ([]map[string]interface{}, error)
	GetMoviesWithFilters(sortBy string, order string, limit, offset int) ([]map[string]interface{}, error)
	SearchMoviesByTitleAndActor(titleFragment, actorNameFragment string, limit, offset int) ([]map[string]interface{}, error)
	UpdateMovie(tx *sql.Tx, id uuid.UUID, movie models.UpdateMovie) error
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

func (s *movie) ValidateActorIDs(actorIDs []uuid.UUID) error {
	exists, err := s.store.CheckActorsExist(actorIDs)
	if err != nil {
		return fmt.Errorf("failed to validate actor IDs: %w", err)
	}
	if !exists {
		return fmt.Errorf("one or more actors in the list do not exist")
	}
	return nil
}

// Валидация movieID (проверка на существование)
func (m *movie) ValidateMovieID(movieID uuid.UUID) error {
	exists, err := m.store.CheckMovieExists(movieID)
	if err != nil {
		return fmt.Errorf("failed to check movie existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("movie with ID %s not found", movieID)
	}
	return nil
}

func (s *movie) AddMovieActorRelations(movieID uuid.UUID, actorIDs []uuid.UUID) error {
	// Validate movie and actors
	if err := s.ValidateMovieID(movieID); err != nil {
		return err
	}
	if err := s.ValidateActorIDs(actorIDs); err != nil {
		return err
	}

	// Transaction for adding relations
	err := withTransactionError(s.store.BeginTransaction, func(tx *sql.Tx) error {
		err := s.store.AddMovieActorRelations(tx, movieID, actorIDs)
		if err != nil {
			log.Printf("Error adding movie-actor relations for movie ID %v: %v", movieID, err)
			return fmt.Errorf("failed to add movie-actor relations: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Printf("Transaction failed while adding movie-actor relations for movie ID %v: %v", movieID, err)
		return err
	}

	return nil

}

func (s *movie) UpdateMovieActorRelations(movieID uuid.UUID, actorIDs []uuid.UUID) error {
	// Validate movie and actors
	if err := s.ValidateMovieID(movieID); err != nil {
		return err
	}
	if err := s.ValidateActorIDs(actorIDs); err != nil {
		return err
	}

	// Transaction for updating relations (remove old, add new)
	err := withTransactionError(s.store.BeginTransaction, func(tx *sql.Tx) error {
		// Remove old relations
		if err := s.store.RemoveMovieActorRelations(tx, movieID); err != nil {
			log.Printf("Error removing old movie-actor relations for movie ID %v: %v", movieID, err)
			return fmt.Errorf("failed to remove old relations: %w", err)
		}
		// Add new relations
		if err := s.store.AddMovieActorRelations(tx, movieID, actorIDs); err != nil {
			log.Printf("Error adding new movie-actor relations for movie ID %v: %v", movieID, err)
			return fmt.Errorf("failed to add new relations: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Printf("Transaction failed while updating movie-actor relations for movie ID %v: %v", movieID, err)
		return err
	}
	return nil
}

func (s *movie) RemoveSelectedMovieActorRelations(movieID uuid.UUID, actorIDs []uuid.UUID) error {
	// Validate movie and actors
	if err := s.ValidateMovieID(movieID); err != nil {
		return err
	}
	if err := s.ValidateActorIDs(actorIDs); err != nil {
		return err
	}

	// Transaction for removing specific relations
	err := withTransactionError(s.store.BeginTransaction, func(tx *sql.Tx) error {
		err := s.store.RemoveSelectedMovieActorRelations(tx, movieID, actorIDs)
		if err != nil {
			log.Printf("Error removing selected movie-actor relations for movie ID %v: %v", movieID, err)
			return fmt.Errorf("failed to remove movie-actor relations: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Printf("Transaction failed while removing selected movie-actor relations for movie ID %v: %v", movieID, err)
		return err
	}
	return nil
}

func (s *movie) CreateMovie(movie models.CreateMovie) (uuid.UUID, error) {
	// Validate actors before proceeding
	if err := s.ValidateActorIDs(movie.ActorIDs); err != nil {
		return uuid.Nil, err
	}

	// Transaction for adding a new movie and its relations
	movieID, err := withTransactionUUID(s.store.BeginTransaction, func(tx *sql.Tx) (uuid.UUID, error) {
		// Add the movie
		movieID, err := s.store.CreateMovie(tx, movie)

		if err != nil {
			log.Printf("[CreateMovie] Failed to add movie %v: %v", movie.Title, err)
			return uuid.Nil, fmt.Errorf("failed to add movie: %w", err)
		}
		// Add actor relations
		err = s.store.AddMovieActorRelations(tx, movieID, movie.ActorIDs)
		if err != nil {
			log.Printf("[CreateMovie] Failed to add movie-actor relations for movie ID %v: %v", movieID, err)
			return uuid.Nil, fmt.Errorf("failed to add movie-actor relations: %w", err)
		}
		return movieID, nil
	})

	if err != nil {
		log.Printf("[CreateMovie] Transaction failed for movie %v: %v", movie.Title, err)
		return uuid.Nil, err
	}

	return movieID, nil
}

// Получение фильма по ID
func (s *movie) GetMovieByID(movieID uuid.UUID) (*models.Movie, error) {
	// Вызываем репозиторий, чтобы получить сырые данные
	rawData, err := s.store.GetMovieByID(movieID)
	if err != nil {
		log.Printf("[GetMovieByID] Failed to retrieve movie with ID %v: %v", movieID, err)
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
		log.Printf("[GetMoviesByActorID] Failed to retrieve movies for actor ID %v: %v", actorID, err)
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
		log.Printf("[GetMoviesWithFilters] Failed to fetch movies with sortBy=%s, order=%s: %v", sortBy, order, err)
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
		log.Printf("[SearchMoviesByTitleAndActor] Failed to search movies by titleFragment=%q, actorNameFragment=%q: %v", titleFragment, actorNameFragment, err)
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

func (s *movie) UpdateMovie(movieID uuid.UUID, movie models.UpdateMovie) error {
	// Validate movie ID
	if err := s.ValidateMovieID(movieID); err != nil {
		return err
	}

	// Transaction for updating movie and its relations
	err := withTransactionError(s.store.BeginTransaction, func(tx *sql.Tx) error {
		// Update the movie details
		err := s.store.UpdateMovie(tx, movieID, movie)
		if err != nil {
			log.Printf("[UpdateMovie] Failed to update movie details for ID %v: %v", movieID, err)
			return fmt.Errorf("[UpdateMovie] failed to update movie: %w", err)
		}

		// Update actor relations if provided
		if movie.ActorIDs != nil {
			err = s.store.RemoveMovieActorRelations(tx, movieID)
			if err != nil {
				log.Printf("[UpdateMovie] Failed to remove old relations for movie ID %v: %v", movieID, err)
				return fmt.Errorf("[UpdateMovie] failed to remove old relations: %w", err)
			}

			err = s.store.AddMovieActorRelations(tx, movieID, *movie.ActorIDs)
			if err != nil {
				log.Printf("[UpdateMovie] Failed to add new relations for movie ID %v: %v", movieID, err)
				return fmt.Errorf("[UpdateMovie] failed to add new relations: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("[UpdateMovie] Transaction failed for movie ID %v: %v", movieID, err)
		return err
	}
	return nil
}

// Удаление фильма по ID
func (m *movie) DeleteMovie(movieID uuid.UUID) error {
	return m.store.DeleteMovie(movieID)
}
