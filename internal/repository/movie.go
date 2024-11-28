package repository

import (
	"database/sql"
	"fmt"
	"log"
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

func (m *movie) BeginTransaction() (*sql.Tx, error) {
	tx, err := m.db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction in BeginTransaction: %v", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

// Функция для проверки, существует ли фильм по id
func (m *movie) CheckMovieExists(movieID uuid.UUID) (bool, error) {
	query := sq.
		Select("EXISTS (SELECT 1 FROM movies WHERE id = $1)").
		PlaceholderFormat(sq.Dollar).
		From("movies").
		Where(sq.Eq{"id": movieID})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Printf("Error building query for checking movie existence in CheckMovieExists: %v", err)
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var exists bool
	err = m.db.QueryRow(sqlQuery, args...).Scan(&exists)
	if err != nil {
		log.Printf("Error checking movie existence in CheckMovieExists: %v", err)
		return false, fmt.Errorf("error checking movie existence: %w", err)
	}

	return exists, nil
}

func (r *movie) CheckActorsExist(actorIDs []uuid.UUID) (bool, error) {
	if len(actorIDs) == 0 {
		return true, nil // Пустой список валиден
	}

	query := sq.
		Select("COUNT(*)").
		From("actors").
		Where(sq.Eq{"id": actorIDs}).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Printf("Error building query in CheckActorsExist: %v", err)
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var count int
	err = r.db.QueryRow(sqlQuery, args...).Scan(&count)
	if err != nil {
		log.Printf("Error checking actors existence in CheckActorsExist: %v", err)
		return false, fmt.Errorf("failed to check actors existence: %w", err)
	}

	// Проверяем, что количество найденных записей совпадает с количеством переданных ID
	return count == len(actorIDs), nil
}

func (r *movie) AddMovieActorRelations(tx *sql.Tx, movieID uuid.UUID, actorIDs []uuid.UUID) error {
	if len(actorIDs) == 0 {
		return nil // Если нет актеров для добавления, ничего не делаем
	}

	// Строим запрос на добавление всех связей актера с фильмом
	queryBuilder := sq.Insert("movie_actors").
		Columns("movie_id", "actor_id").
		Suffix("ON CONFLICT (movie_id, actor_id) DO NOTHING") // Для предотвращения дублирования

	// Добавляем каждую пару movie_id, actor_id
	for _, actorID := range actorIDs {
		queryBuilder = queryBuilder.Values(movieID, actorID)
	}

	// Генерация SQL-запроса
	sqlQuery, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Printf("Error building query in AddMovieActorRelations: %v", err)
		return fmt.Errorf("failed to build query: %w", err)
	}

	// Выполнение запроса
	_, err = tx.Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Error adding movie-actor relations in AddMovieActorRelations: %v", err)
		return fmt.Errorf("failed to add movie-actor relations: %w", err)
	}

	return nil
}

// Удаление связей по movieID
func (m *movie) RemoveMovieActorRelations(tx *sql.Tx, movieID uuid.UUID) error {
	query := sq.Delete("movie_actors").Where(sq.Eq{"movie_id": movieID})

	sqlQuery, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Printf("Error building delete query in RemoveMovieActorRelations: %v", err)
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = tx.Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Error deleting old relations in RemoveMovieActorRelations: %v", err)
		return fmt.Errorf("failed to delete old relations: %w", err)
	}
	return nil
}

func (r *movie) RemoveSelectedMovieActorRelations(tx *sql.Tx, movieID uuid.UUID, actorIDs []uuid.UUID) error {
	if len(actorIDs) == 0 {
		return nil // Если нет актеров, которых нужно удалить, ничего не делаем
	}

	// Формируем запрос на удаление
	query := sq.Delete("movie_actors").
		Where(sq.Eq{"movie_id": movieID}).
		Where(sq.Eq{"actor_id": actorIDs}).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Printf("Error building delete query in RemoveSelectedMovieActorRelations: %v", err)
		return fmt.Errorf("failed to build query: %w", err)
	}

	// Выполнение запроса
	_, err = tx.Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Error removing movie-actor relations in RemoveSelectedMovieActorRelations: %v", err)
		return fmt.Errorf("failed to remove movie-actor relations: %w", err)
	}

	return nil
}

func (m *movie) AddMovie(tx *sql.Tx, movie models.CreateMovie) (uuid.UUID, error) {
	id := uuid.New()
	query := sq.
		Insert("movies").
		Columns("id", "title", "description", "release_date", "rating").
		Values(id, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Printf("Error building query for adding movie in AddMovie: %v", err)
		return uuid.Nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = tx.QueryRow(sqlQuery, args...).Scan(&id)
	if err != nil {
		log.Printf("Error adding movie in AddMovie: %v", err)
		return uuid.Nil, fmt.Errorf("failed to add movie: %w", err)
	}

	return id, nil
}

func (m *movie) GetMovieByID(id uuid.UUID) (map[string]interface{}, error) {
	query := sq.
		Select("id", "title", "description", "release_date", "rating").
		From("movies").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Printf("Error building query in GetMovieByID: %v", err)
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
		log.Printf("Error scanning row in GetMovieByID: %v", err)
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
		log.Printf("Error building query in GetMoviesByActorID: %v", err)
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := m.db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Error executing query in GetMoviesByActorID: %v", err)
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
			log.Printf("Error scanning row in GetMoviesByActorID: %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
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

func (m *movie) GetMoviesWithFilters(sortBy string, order string, limit, offset int) ([]map[string]interface{}, error) {
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
		log.Printf("Error building query in GetMoviesWithFilters: %v", err)
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	// Выполняем запрос
	rows, err := m.db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Error executing query in GetMoviesWithFilters: %v", err)
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
			log.Printf("Error scanning row in GetMoviesWithFilters: %v", err)
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

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows in GetMoviesWithFilters: %v", err)
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return rawMovies, nil
}

func (m *movie) SearchMoviesByTitleAndActor(filterTitle, filterActor string, limit, offset int) ([]map[string]interface{}, error) {
	query := sq.
		Select("DISTINCT m.id", "m.title", "m.description", "m.release_date", "m.rating").
		From("movies m").
		LeftJoin("movie_actors ma ON m.id = ma.movie_id").
		LeftJoin("actors a ON ma.actor_id = a.id").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	// Добавляем условия только при наличии фильтров
	if filterTitle != "" {
		query = query.Where(sq.ILike{"m.title": fmt.Sprintf("%%%s%%", filterTitle)})
	}
	if filterActor != "" {
		query = query.Where(sq.ILike{"a.name": fmt.Sprintf("%%%s%%", filterActor)})
	}

	// Конвертация в SQL
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Printf("Error building query in SearchMoviesByTitleAndActor: %v", err)
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	// Выполнение запроса
	rows, err := m.db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Error executing query in SearchMoviesByTitleAndActor: %v", err)
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	defer rows.Close()

	// Обработка результатов
	var rawMovies []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var title, description string
		var releaseDate time.Time
		var rating float64
		err := rows.Scan(&id, &title, &description, &releaseDate, &rating)
		if err != nil {
			log.Printf("Error scanning row in SearchMoviesByTitleAndActor: %v", err)
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}

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
func (m *movie) UpdateMovie(tx *sql.Tx, id uuid.UUID, movie models.UpdateMovie) error {
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
		log.Printf("Error building query in UpdateMovie: %v", err)
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = tx.Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Error executing query in UpdateMovie: %v", err)
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
		log.Printf("Error building query in DeleteMovie: %v", err)
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = m.db.Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("Error executing query in DeleteMovie: %v", err)
		return fmt.Errorf("failed to delete movie: %w", err)
	}
	return nil
}
