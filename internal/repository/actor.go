package repository

import (
	"cinema/internal/models"
	"database/sql"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type actor struct {
	db *sql.DB
}

func NewActor(db *sql.DB) *actor {
	return &actor{db: db}
}

// Добавить актера
func (a *actor) AddActor(actor models.CreateActor) (uuid.UUID, error) {
	id := uuid.New()
	query := sq.
		Insert("actors").
		Columns("id", "name", "gender", "date_of_birth").
		Values(id, actor.Name, actor.Gender, actor.DateOfBirth).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Printf("[AddActor] Error building query: %v", err)
		return uuid.Nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = a.db.QueryRow(sqlQuery, args...).Scan(&id)
	if err != nil {
		log.Printf("[AddActor] Error executing query: %v", err)
		return uuid.Nil, fmt.Errorf("failed to add actor: %w", err)
	}
	return id, nil
}

func (a *actor) GetActor(id uuid.UUID) (map[string]interface{}, error) {
	query := sq.
		Select("id", "name", "gender", "date_of_birth").
		From("actors").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Printf("[GetActor] Error building query: %v", err)
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rawData := make(map[string]interface{})
	var actorID uuid.UUID
	var name, gender string
	var dateOfBirth time.Time

	err = a.db.QueryRow(sqlQuery, args...).Scan(&actorID, &name, &gender, &dateOfBirth)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Актёр не найден
		}
		log.Printf("[GetActor] Error executing query: %v", err)
		return nil, fmt.Errorf("failed to get actor: %w", err)
	}

	rawData["id"] = actorID
	rawData["name"] = name
	rawData["gender"] = gender
	rawData["date_of_birth"] = dateOfBirth

	return rawData, nil
}

func (a *actor) GetAllActors(limit, offset int) ([]map[string]interface{}, error) {
	query := sq.
		Select("id", "name", "gender", "date_of_birth").
		From("actors").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Printf("[GetAllActors] Error building query: %v", err)
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := a.db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("[GetAllActors] Error executing query: %v", err)
		return nil, fmt.Errorf("failed to get actors: %w", err)
	}
	defer rows.Close()

	var rawActors []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var name, gender string
		var dateOfBirth time.Time

		err := rows.Scan(&id, &name, &gender, &dateOfBirth)
		if err != nil {
			log.Printf("[GetAllActors] Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan actor: %w", err)
		}

		rawActor := map[string]interface{}{
			"id":            id,
			"name":          name,
			"gender":        gender,
			"date_of_birth": dateOfBirth,
		}
		rawActors = append(rawActors, rawActor)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[GetAllActors] Error iterating rows: %v", err)
		return nil, fmt.Errorf("error occurred while iterating rows: %w", err)
	}

	return rawActors, nil
}

func (a *actor) GetActorsWithMovies(limit, offset int) ([]map[string]interface{}, error) {
	// Создаем подзапрос для пагинации актеров
	actorsQuery := sq.
		Select("id AS actor_id", "name AS actor_name", "gender AS actor_gender", "date_of_birth AS actor_birth_date").
		From("actors").
		OrderBy("name ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	// Основной запрос с соединением фильмов
	query := sq.
		Select("pa.actor_id", "pa.actor_name", "pa.actor_gender", "pa.actor_birth_date", "m.id AS movie_id",
			"m.title AS movie_title", "m.description AS movie_description", "m.release_date AS movie_release_date", "m.rating AS movie_rating").
		FromSelect(actorsQuery, "pa").
		LeftJoin("movie_actors ma ON pa.actor_id = ma.actor_id").
		LeftJoin("movies m ON ma.movie_id = m.id").
		OrderBy("pa.actor_name ASC", "pa.actor_id").
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Printf("[GetActorsWithMovies] Error building query: %v", err)
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := a.db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("[GetActorsWithMovies] Error executing query: %v", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Срез для хранения актеров в порядке их получения
	// var result []map[string]interface{}
	actorsMap := make(map[uuid.UUID]map[string]interface{})

	for rows.Next() {
		var (
			actorID        uuid.UUID
			actorName      string
			actorGender    string
			actorBirthDate time.Time
			movieID        *uuid.UUID
			movieTitle     *string
			movieDesc      *string
			movieRelease   *time.Time
			movieRating    *float64
		)

		if err := rows.Scan(
			&actorID,
			&actorName,
			&actorGender,
			&actorBirthDate,
			&movieID,
			&movieTitle,
			&movieDesc,
			&movieRelease,
			&movieRating,
		); err != nil {
			log.Printf("[GetActorsWithMovies] Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		// Проверяем, есть ли актёр в словаре
		if actor, exists := actorsMap[actorID]; exists {
			// Если актёр уже есть, добавляем информацию о фильме
			if movieID != nil {
				actor["movies"] = append(actor["movies"].([]map[string]interface{}), map[string]interface{}{
					"id":           *movieID,
					"title":        *movieTitle,
					"description":  *movieDesc,
					"release_date": *movieRelease,
					"rating":       *movieRating,
				})
			}
		} else {
			// Если актёра ещё нет, создаём его запись
			newActor := map[string]interface{}{
				"id":            actorID,
				"name":          actorName,
				"gender":        actorGender,
				"date_of_birth": actorBirthDate,
				"movies":        []map[string]interface{}{},
			}

			// Добавляем информацию о фильме (если она есть)
			if movieID != nil {
				newActor["movies"] = append(newActor["movies"].([]map[string]interface{}), map[string]interface{}{
					"id":           *movieID,
					"title":        *movieTitle,
					"description":  *movieDesc,
					"release_date": *movieRelease,
					"rating":       *movieRating,
				})
			}

			// Сохраняем актёра в словарь
			actorsMap[actorID] = newActor
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("[GetActorsWithMovies] Error iterating over rows: %v", err)
		return nil, fmt.Errorf("error occurred while iterating rows: %w", err)
	}

	// Преобразуем карту в срез для возврата
	var result []map[string]interface{}
	for _, actor := range actorsMap {
		result = append(result, actor)
	}

	return result, nil
}

func (a *actor) UpdateActor(id uuid.UUID, actor models.UpdateActor) error {
	query := sq.
		Update("actors").
		Set("name", sq.Expr("COALESCE(?, name)", actor.Name)).
		Set("gender", sq.Expr("COALESCE(?, gender)", actor.Gender)).
		Set("date_of_birth", sq.Expr("COALESCE(?, date_of_birth)", actor.DateOfBirth)).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Printf("[UpdateActor] Error building query: %v", err)
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = a.db.Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("[UpdateActor] Error executing query: %v", err)
		return fmt.Errorf("failed to update actor: %w", err)
	}
	return nil
}

func (a *actor) DeleteActor(id uuid.UUID) error {
	query := sq.
		Delete("actors").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Printf("[DeleteActor] Error building query: %v", err)
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = a.db.Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("[DeleteActor] Error executing query: %v", err)
		return fmt.Errorf("failed to delete actor: %w", err)
	}
	return nil
}
