package repository

import (
	"database/sql"
	"fmt"
	"strings"
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
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

// Добавить фильм
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
		return uuid.Nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = tx.QueryRow(sqlQuery, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to add movie: %w", err)
	}

	return id, nil
}

// Функция для проверки, существует ли актер по id
func (m *movie) CheckActorExists(actorID uuid.UUID) (bool, error) {
	query := sq.Select("EXISTS (SELECT 1 FROM actors WHERE id = ?)").
		PlaceholderFormat(sq.Dollar).
		From("actors").
		Where(sq.Eq{"id": actorID})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build check actor query: %w", err)
	}

	var exists bool
	err = m.db.QueryRow(sqlQuery, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // актер не найден
		}
		return false, fmt.Errorf("error checking actor existence: %w", err)
	}
	return exists, nil
}

// Функция для проверки, существует ли фильм по id
func (m *movie) CheckMovieExists(movieID uuid.UUID) (bool, error) {
	query := sq.
		Select("EXISTS (SELECT 1 FROM movies WHERE id = ?)").
		PlaceholderFormat(sq.Dollar).
		From("movies").
		Where(sq.Eq{"id": movieID})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var exists bool
	err = m.db.QueryRow(sqlQuery, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking movie existence: %w", err)
	}

	return exists, nil
}

// Функция для добавления новых связей по MovieID
func (m *movie) AddMovieActorRelations(tx *sql.Tx, movieID uuid.UUID, actorIDs []uuid.UUID) error {
	insertQuery := sq.Insert("movie_actors").
		Columns("movie_id", "actor_id")

	// Добавляем все связи
	for _, actorID := range actorIDs {
		insertQuery = insertQuery.Values(movieID, actorID)
	}

	// Генерация SQL-запроса
	sqlQuery, args, err := insertQuery.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	// Выполняем запрос в контексте транзакции
	_, err = tx.Exec(sqlQuery, args...)
	if err != nil {
		// Если возникает ошибка уникальности, игнорируем её
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil // Пропускаем ошибку дублирования
		}
		return fmt.Errorf("failed to insert actor-movie relations: %w", err)
	}

	return nil
}

// Удаление связей по movieID
func (m *movie) RemoveMovieActorRelations(tx *sql.Tx, movieID uuid.UUID) error {
	query := sq.Delete("movie_actors").Where(sq.Eq{"movie_id": movieID})

	sqlQuery, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = tx.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to delete old relations: %w", err)
	}
	return nil
}

// RemoveMovieActorRelations удаляет связи между фильмом и актёрами
func (m *movie) RemoveSelectedMovieActorRelations(tx *sql.Tx, movieID uuid.UUID, actorIDs []uuid.UUID) error {
	// Если список actorIDs пуст, нет необходимости выполнять запрос
	if len(actorIDs) == 0 {
		return nil
	}

	deleteQuery := sq.Delete("movie_actors").
		Where(sq.Eq{"movie_id": movieID}).
		Where(sq.Eq{"actor_id": actorIDs}) // Удаляем только указанных актёров

	sqlQuery, args, err := deleteQuery.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = tx.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to remove actor-movie relations: %w", err)
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
		query = query.Where(sq.Like{"m.title": fmt.Sprintf("%%%s%%", filterTitle)})
	}
	if filterActor != "" {
		query = query.Where(sq.Like{"a.name": fmt.Sprintf("%%%s%%", filterActor)})
	}

	// Конвертация в SQL
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	// Выполнение запроса
	rows, err := m.db.Query(sqlQuery, args...)
	if err != nil {
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
