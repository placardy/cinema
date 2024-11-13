package repository

import (
	"database/sql"
	"fmt"
	"time"

	"cinema/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
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
	query := sq.
		Insert("movies").
		Columns("id", "title", "description", "release_date", "rating").
		Values(id, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = m.db.QueryRow(sqlQuery, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to add movie: %w", err)
	}

	return id, nil
}

// Создание связи между фильмом и актером
func (m *movie) AddMovieActorRelation(actorID, movieID uuid.UUID) error {
	query := sq.
		Insert("actor_movies").
		Columns("actor_id", "movie_id").
		Values(actorID, movieID).
		Suffix("ON CONFLICT DO NOTHING").
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = m.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to add actor to movie: %w", err)
	}
	return nil
}

// Удаление связи между фильмом и актером
func (m *movie) RemoveMovieActorRelation(actorID, movieID uuid.UUID) error {
	query := sq.
		Delete("actor_movies").
		Where(sq.Eq{"actor_id": actorID, "movie_id": movieID}).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	_, err = m.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to remove actor from movie: %w", err)
	}
	return nil
}

// Получить фильм по id
func (m *movie) GetMovieByID(id uuid.UUID) (map[string]interface{}, error) {
	query := sq.
		Select("id", "title", "description", "release_date", "rating").
		From("movies").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	var idRaw uuid.UUID
	var title, description string
	var releaseDate time.Time
	var rating float64
	err = m.db.QueryRow(sqlQuery, args...).Scan(&idRaw, &title, &description, &releaseDate, &rating)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Фильм не найден
		}
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}

	rawData := map[string]interface{}{
		"id":           idRaw,
		"title":        title,
		"description":  description,
		"release_date": releaseDate,
		"rating":       rating,
	}

	return rawData, nil
}

// Получить фильмы по actor id
func (m *movie) GetMoviesByActorID(actorID uuid.UUID, limit, offset int) ([]map[string]interface{}, error) {
	query := sq.
		Select("m.id", "m.title", "m.description", "m.release_date", "m.rating").
		From("movies m").
		Join("movie_actors ma ON m.id = ma.movie_id").
		Where(sq.Eq{"ma.actor_id": actorID}).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := m.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie by actorID: %w", err)
	}
	defer rows.Close()

	var rawData []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var title, description string
		var releaseDate time.Time
		var rating float64
		if err := rows.Scan(&id, &title, &description, &releaseDate, &rating); err != nil {
			return nil, err
		}
		movieData := map[string]interface{}{
			"id":           id,
			"title":        title,
			"description":  description,
			"release_date": releaseDate,
			"rating":       rating,
		}
		rawData = append(rawData, movieData)
	}
	return rawData, nil
}

// Получить фильмы с фильтрацией
func (m *movie) GetMoviesWithFilters(sortBy string, order string, limit, offset int) ([]map[string]interface{}, error) {
	// Валидация
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

	// Строим SQL запрос
	query := sq.
		Select("id", "title", "description", "release_date", "rating").
		From("movies").
		OrderBy(fmt.Sprintf("%s %s", sortBy, order)).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	// Преобразуем запрос в SQL
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	// Выполняем запрос
	rows, err := m.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get movies with filtration: %w", err)
	}
	defer rows.Close()

	var rawMovies []map[string]interface{}
	// Считываем результаты
	for rows.Next() {
		var id uuid.UUID
		var title, description string
		var releaseDate time.Time
		var rating float64
		err := rows.Scan(&id, &title, &description, &releaseDate, &rating)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Добавляем сырые данные в срез
		rawMovie := map[string]interface{}{
			"id":           id,
			"title":        title,
			"description":  description,
			"release_date": releaseDate,
			"rating":       rating,
		}
		rawMovies = append(rawMovies, rawMovie)
	}

	return rawMovies, nil
}

// Поиск фильмов по названию
func (m *movie) SearchMoviesByTitle(titleFragment string, limit, offset int) ([]map[string]interface{}, error) {
	query := sq.
		Select("id", "title", "description", "release_date", "rating").
		From("movies").
		Where(sq.Like{"title": fmt.Sprintf("%%%s%%", titleFragment)}).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := m.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies by title: %w", err)
	}
	defer rows.Close()

	var rawMovies []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var title, description string
		var releaseDate time.Time
		var rating float64
		err := rows.Scan(&id, &title, &description, &releaseDate, &rating)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}

		// Добавляем сырые данные в срез
		rawMovie := map[string]interface{}{
			"id":           id,
			"title":        title,
			"description":  description,
			"release_date": releaseDate,
			"rating":       rating,
		}
		rawMovies = append(rawMovies, rawMovie)
	}

	return rawMovies, nil
}

// Поиск фильмов по актеру
func (m *movie) SearchMoviesByActorName(actorNameFragment string, limit, offset int) ([]map[string]interface{}, error) {
	query := sq.
		Select("DISTINCT m.id", "m.title", "m.description", "m.release_date", "m.rating").
		From("movies m").
		Join("movie_actors ma ON m.id = ma.movie_id").
		Join("actors a ON ma.actor_id = a.id").
		Where(sq.Like{"a.name": fmt.Sprintf("%%%s%%", actorNameFragment)}).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := m.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies by actor: %w", err)
	}
	defer rows.Close()

	var rawMovies []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var title, description string
		var releaseDate time.Time
		var rating float64
		err := rows.Scan(&id, &title, &description, &releaseDate, &rating)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}

		// Добавляем сырые данные в срез
		rawMovie := map[string]interface{}{
			"id":           id,
			"title":        title,
			"description":  description,
			"release_date": releaseDate,
			"rating":       rating,
		}
		rawMovies = append(rawMovies, rawMovie)
	}

	return rawMovies, nil
}

// Обновить фильм
func (m *movie) UpdateMovie(id uuid.UUID, movie models.UpdateMovie) error {
	query := sq.
		Update("movies").
		Set("title", sq.Expr("COALESCE(?, title)", movie.Title)).
		Set("description", sq.Expr("COALESCE(?, description)", movie.Description)).
		Set("release_date", sq.Expr("COALESCE(?, release_date)", movie.ReleaseDate)).
		Set("rating", sq.Expr("COALESCE(?, rating)", movie.Rating)).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = m.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to update movie: %w", err)
	}
	return nil
}

// Удалить фильм по id
func (m *movie) DeleteMovie(id uuid.UUID) error {
	query := sq.
		Delete("movies").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = m.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}
	return nil
}

// Получение списка фильмов
func (m *movie) GetAllMovies(limit, offset int) ([]map[string]interface{}, error) {
	var movies []map[string]interface{}
	query := sq.
		Select("id", "title", "description", "release_date", "rating").
		From("movies").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := m.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}
	defer rows.Close()

	// Считываем сырые данные в map
	for rows.Next() {
		// Определяем переменные для каждого поля
		var id uuid.UUID
		var title, description string
		var releaseDate time.Time
		var rating float64

		// Сканируем строки в переменные
		if err := rows.Scan(&id, &title, &description, &releaseDate, &rating); err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}

		// Добавляем в map, чтобы передать в сервис
		movie := map[string]interface{}{
			"id":           id,
			"title":        title,
			"description":  description,
			"release_date": releaseDate,
			"rating":       rating,
		}

		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating rows: %w", err)
	}

	return movies, nil
}
