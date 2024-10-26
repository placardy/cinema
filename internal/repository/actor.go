package repository

import (
	"cinema/internal/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	query := `
        INSERT INTO actors (id, name, gender, date_of_birth)
        VALUES ($1, $2, $3, $4) RETURNING id`

	err := a.db.QueryRow(query, id, actor.Name, actor.Gender, actor.DateOfBirth).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to add actor: %w", err)
	}
	return id, nil
}

// Получить актера по id
func (a *actor) GetActor(id uuid.UUID) (*models.Actor, error) {
	query := `SELECT id, name, gender, date_of_birth FROM actors WHERE id = $1`
	var actor models.Actor
	err := a.db.QueryRow(query, id).Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.DateOfBirth)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Актёр не найден
		}
		return nil, fmt.Errorf("failed to get actor: %w", err)
	}
	return &actor, nil
}

// Обновить актера
func (a *actor) UpdateActor(id uuid.UUID, actor models.UpdateActor) error {
	query := `
        UPDATE actors SET
            name = COALESCE($1, name),
            gender = COALESCE($2, gender),
            date_of_birth = COALESCE($3, date_of_birth)
        WHERE id = $4`

	_, err := a.db.Exec(query, actor.Name, actor.Gender, actor.DateOfBirth, id)
	if err != nil {
		return fmt.Errorf("failed to update actor: %w", err)
	}
	return nil
}

// Удалить актера
func (a *actor) DeleteActor(id uuid.UUID) error {
	query := `DELETE FROM actors WHERE id = $1`
	_, err := a.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete actor: %w", err)
	}
	return nil
}

// Получить всех актеров
func (a *actor) GetAllActors(limit, offset int) ([]*models.Actor, error) {
	var actors []*models.Actor
	query := `SELECT id, name, gender, date_of_birth FROM actors LIMIT $1 OFFSET $2`

	rows, err := a.db.Query(query, limit, offset) // Передаем limit и offset
	if err != nil {
		return nil, fmt.Errorf("failed to get actors: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var actor models.Actor
		err := rows.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("failed to scan actor: %w", err)
		}
		actors = append(actors, &actor)
	}

	// Проверяем наличие ошибок при итерации по строкам
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating rows: %w", err)
	}

	return actors, nil
}

// Получить актеров с фильмами с пагинацией
func (r *actor) GetActorsWithMovies(limit, offset int) ([]models.Actor, error) {
	query := `
    WITH paginated_actors AS (
        SELECT 
            a.id AS actor_id,
            a.name AS actor_name,
            a.gender AS actor_gender,
            a.date_of_birth AS actor_birth_date
        FROM 
            actors a
        ORDER BY 
            a.name
        LIMIT $1 OFFSET $2
    )
    SELECT 
        pa.actor_id,
        pa.actor_name,
        pa.actor_gender,
        pa.actor_birth_date,
        m.id AS movie_id,
        m.title AS movie_title,
        m.description AS movie_description,
        m.release_date AS movie_release_date,
        m.rating AS movie_rating
    FROM 
        paginated_actors pa
    LEFT JOIN 
        movie_actors ma ON pa.actor_id = ma.actor_id
    LEFT JOIN 
        movies m ON ma.movie_id = m.id;`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actorsMap = make(map[uuid.UUID]*models.Actor)

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
			return nil, err
		}

		// Создаем нового актера, если его еще нет в мапе
		if _, exists := actorsMap[actorID]; !exists {
			actorsMap[actorID] = &models.Actor{
				ID:          actorID,
				Name:        actorName,
				Gender:      actorGender,
				DateOfBirth: actorBirthDate,
				Movies:      []models.Movie{},
			}
		}

		// Добавляем фильм только если он существует
		if movieID != nil {
			actorsMap[actorID].Movies = append(actorsMap[actorID].Movies, models.Movie{
				ID:          *movieID,
				Title:       *movieTitle,
				Description: *movieDesc,
				ReleaseDate: *movieRelease,
				Rating:      *movieRating,
			})
		}
	}

	var actors []models.Actor
	for _, actor := range actorsMap {
		actors = append(actors, *actor)
	}

	return actors, nil
}
