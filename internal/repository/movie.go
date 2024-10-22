package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"

	"cinema/internal/models"
)

type movie struct {
	db *sql.DB
}

func NewMovie(db *sql.DB) *movie {
	return &movie{db: db}
}

// Добавить фильм
func (m *movie) CreateMovie(title, description, release_date string, rating float64) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO movies (id, title, description, release_date, rating) VALUES ($1, $2, $3, $4, $5)`
	_, err := m.db.Exec(query, id, title, description, release_date, rating)
	if err != nil {
		log.Printf("Failed to insert movie: %v", err)
		return uuid.Nil, err
	}
	return id, nil
}

// Создание связи между фильмом и актером
func (r *movie) CreateMovieActorRelation(movieID, actorID uuid.UUID) error {
	query := `INSERT INTO movie_actors (movie_id, actor_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, movieID, actorID)
	if err != nil {
		return err
	}
	return nil
}

// Получить фильм по id
func (a *actor) GetMovieByID(id uuid.UUID) (*models.Movie, error) {
	var movie models.Movie
	query := `SELECT * FROM movies WHERE id = $1`
	err := a.db.QueryRow(query, id).Scan(&movie.Title, &movie.Description, &movie.Rating, &movie.ReleaseDate)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

// Обновить фильм
func (m *movie) UpdateMovie(id uuid.UUID, title string, description string, releaseDate string, rating float64) error {
	query := `UPDATE movies SET title = $1, description = $2, release_date = $3, rating = $4 WHERE id = $5`
	_, err := m.db.Exec(query, title, description, releaseDate, rating, id)
	if err != nil {
		return err
	}
	return nil
}

// Частично обновить фильм
func (m *movie) PartialUpdateMovie(id uuid.UUID, fields map[string]interface{}) error {
	query := `UPDATE movies SET `
	args := []interface{}{}
	i := 1
	for key, value := range fields {
		if i > 1 {
			query += ", "
		}
		query += fmt.Sprintf("%s=$%d", key, i)
		args = append(args, value)
		i++
	}
	query += fmt.Sprintf(" WHERE id=$%d", i)
	args = append(args, id)

	_, err := m.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to partially update movie: %v", err)
	}
	return nil
}

// Удалить фильм по id
func (m *movie) DeleteMovie(id uuid.UUID) error {
	query := `DELETE FROM movies WHERE id = $1`
	_, err := m.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// Получить фильмы с фильтрацией
func (m *movie) GetMovies(sortBy string, order string) ([]*models.Movie, error) {
	var movies []*models.Movie
	query := `SELECT * FROM movies`
	switch sortBy {
	case "title":
		query += " ORDER BY title"
	case "release_date":
		query += " ORDER BY release_date"
	case "rating":
		query += " ORDER BY rating"
	default:
		query += " ORDER BY rating"
	}

	if order == "asc" {
		query += " ASC"
	} else {
		query += " DESC"
	}

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	return movies, nil
}

// Поиск фильмов по названию
func (m *movie) SearchMoviesByTitle(titleFragment string) ([]*models.Movie, error) {
	var movies []*models.Movie
	query := `SELECT id, title, description, release_date, rating FROM movies WHERE title ILIKE '%' || $1 || '%'`

	rows, err := m.db.Query(query, titleFragment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	return movies, nil
}

// Поиск фильмов по актеру
func (m *movie) SearchMoviesByActorName(actorNameFragment string) ([]*models.Movie, error) {
	var movies []*models.Movie
	query := `
        SELECT DISTINCT m.id, m.title, m.description, m.release_date, m.rating 
        FROM movies m
        JOIN movie_actors ma ON m.id = ma.movie_id
        JOIN actors a ON ma.actor_id = a.id
        WHERE a.name ILIKE '%' || $1 || '%'`

	rows, err := m.db.Query(query, actorNameFragment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	return movies, nil
}
