package repository

import (
	"database/sql"
	"fmt"

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
func (m *movie) AddMovie(movie models.CreateMovie) (uuid.UUID, error) {
	id := uuid.New()
	query := `
        INSERT INTO movies (id, title, description, release_date, rating)
        VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := m.db.QueryRow(query, id, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to add movie: %w", err)
	}
	return id, nil
}

// Создание связи между фильмом и актером
func (m *movie) AddActorToMovieRelation(actorID, movieID uuid.UUID) error {
	query := `
        INSERT INTO actor_movies (actor_id, movie_id)
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING` // избегание дубликата

	_, err := m.db.Exec(query, actorID, movieID)
	if err != nil {
		return fmt.Errorf("failed to add actor to movie: %w", err)
	}
	return nil
}

// Удаление связи между фильмом и актером
func (m *movie) RemoveActorFromMovieRelation(actorID, movieID uuid.UUID) error {
	query := `
        DELETE FROM actor_movies 
        WHERE actor_id = $1 AND movie_id = $2`

	_, err := m.db.Exec(query, actorID, movieID)
	if err != nil {
		return fmt.Errorf("failed to remove actor from movie: %w", err)
	}
	return nil
}

// Получить фильм по id
func (m *movie) GetMovie(id uuid.UUID) (*models.Movie, error) {
	query := `
		SELECT id, title, description, release_date, rating 
		FROM movies 
		WHERE id = $1`
	var movie models.Movie
	err := m.db.QueryRow(query, id).Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Movie not found
		}
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}
	return &movie, nil
}

// Получить фильмы по actor id
func (m *movie) GetMoviesByActorID(actorID uuid.UUID, limit, offset int) ([]*models.Movie, error) {
	query := `
		SELECT m.id, m.title, m.description, m.release_date, m.rating
		FROM movies m
		JOIN movie_actors ma ON m.id = ma.movie_id
		WHERE ma.actor_id = $1 
		LIMIT $2 OFFSET $3`
	rows, err := m.db.Query(query, actorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movies []*models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	return movies, nil
}

// Получить фильмы с фильтрацией
func (m *movie) GetMovies(sortBy string, order string, limit, offset int) ([]*models.Movie, error) {
	//валидация
	validSortColumns := map[string]struct{}{
		"title":        {},
		"release_date": {},
		"rating":       {},
	}
	validOrder := map[string]struct{}{
		"ASC":  {},
		"DESC": {},
	}
	if _, ok := validSortColumns[sortBy]; !ok {
		sortBy = "rating"
	}
	if _, ok := validOrder[order]; !ok {
		order = "DESC"
	}
	query := fmt.Sprintf(`SELECT id, title, description, release_date, rating 
	FROM movies 
	ORDER BY %s %s 
	LIMIT $1 OFFSET $2`, sortBy, order)

	rows, err := m.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*models.Movie
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
func (m *movie) SearchMoviesByTitle(titleFragment string, limit, offset int) ([]*models.Movie, error) {
	var movies []*models.Movie
	query := `
		SELECT id, title, description, release_date, rating 
		FROM movies 
		WHERE title ILIKE '%' || $1 || '%'
		LIMIT $2 OFFSET $3`

	rows, err := m.db.Query(query, titleFragment, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating); err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}
		movies = append(movies, &movie)
	}
	return movies, nil
}

// Поиск фильмов по актеру
func (m *movie) SearchMoviesByActorName(actorNameFragment string, limit, offset int) ([]*models.Movie, error) {
	var movies []*models.Movie
	query := `
        SELECT DISTINCT m.id, m.title, m.description, m.release_date, m.rating 
        FROM movies m
        JOIN movie_actors ma ON m.id = ma.movie_id
        JOIN actors a ON ma.actor_id = a.id
        WHERE a.name ILIKE '%' || $1 || '%'
		LIMIT $2 OFFSET $3`

	rows, err := m.db.Query(query, actorNameFragment, limit, offset)
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

// Обновить фильм
func (m *movie) UpdateMovie(id uuid.UUID, movie models.UpdateMovie) error {
	query := `
        UPDATE movies SET
            title = COALESCE($1, title),
            description = COALESCE($2, description),
            release_date = COALESCE($3, release_date),
            rating = COALESCE($4, rating)
        WHERE id = $5`

	_, err := m.db.Exec(query, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating, id)
	if err != nil {
		return fmt.Errorf("failed to update movie: %w", err)
	}
	return nil
}

// Удалить фильм по id
func (m *movie) DeleteMovie(id uuid.UUID) error {
	query := `DELETE FROM movies WHERE id = $1`
	_, err := m.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}
	return nil
}

// Получение списка фильмов
func (m *movie) GetAllMovies(limit, offset int) ([]*models.Movie, error) {
	var movies []*models.Movie
	query := `SELECT id, title, description, release_date, rating FROM movies LIMIT $1 OFFSET $2`

	rows, err := m.db.Query(query, limit, offset) // Передаем limit и offset
	if err != nil {
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}
		movies = append(movies, &movie)
	}

	// Проверяем наличие ошибок при итерации по строкам
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating rows: %w", err)
	}

	return movies, nil
}
