package repository

import (
	"database/sql"
	"log"

	"github.com/google/uuid"

	"cinema/internal/models"
)

type actor struct {
	db *sql.DB
}

func NewActor(db *sql.DB) *actor {
	return &actor{db: db}
}

// Добавить актера
func (a *actor) CreateActor(name, gender, birthDate string) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO actors (id, name, gender, birth_date) VALUES ($1, $2, $3, $4)`
	_, err := a.db.Exec(query, id, name, gender, birthDate)
	if err != nil {
		log.Printf("Failed to insert actor: %v", err)
		return uuid.Nil, err
	}
	return id, nil
}

// Получить актера по id
func (a *actor) GetActorByID(id uuid.UUID) (*models.Actor, error) {
	var actor models.Actor
	query := `SELECT * FROM actors WHERE id = $1`
	err := a.db.QueryRow(query, id).Scan(&actor.Name, &actor.Gender, &actor.BirthDate)
	if err != nil {
		return nil, err
	}
	return &actor, nil
}

// Обновить актера по id
func (a *actor) UpdateActor(id uuid.UUID, name, gender, birthDate string) error {
	query := `UPDATE actors SET name = $1, gender = $2, birthdate = $3 WHERE id = $4`
	_, err := a.db.Exec(query, name, gender, birthDate, id)
	if err != nil {
		return err
	}
	return nil
}

// Удалить актера
func (a *actor) DeleteActor(id uuid.UUID) error {
	query := `DELETE FROM actors WHERE id = $1`
	_, err := a.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
