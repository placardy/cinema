package service

import (
	"cinema/internal/models"

	"github.com/google/uuid"
)

type storeActor interface {
	AddActor(actor models.CreateActor) (uuid.UUID, error)
	GetActor(id uuid.UUID) (*models.Actor, error)
	UpdateActor(id uuid.UUID, actor models.UpdateActor) error
	DeleteActor(id uuid.UUID) error
	GetAllActors(limit, offset int) ([]*models.Actor, error)
	GetActorsWithMovies(limit, offset int) ([]*models.Actor, error)
}

type actor struct {
	store storeActor
}

func NewActor(store storeActor) actor {
	return actor{store: store}
}

// Добавление актера
func (a actor) AddActor(actor models.CreateActor) (uuid.UUID, error) {
	return a.store.AddActor(actor)
}

// Получение актера по ID
func (a actor) GetActor(id uuid.UUID) (*models.Actor, error) {
	return a.store.GetActor(id)
}

// Обновление актера по ID
func (a actor) UpdateActor(id uuid.UUID, actor models.UpdateActor) error {
	return a.store.UpdateActor(id, actor)
}

// Удаление актера
func (a actor) DeleteActor(id uuid.UUID) error {
	return a.store.DeleteActor(id)
}

// Получение актеров с пагинацией
func (a actor) GetAllActors(limit, offset int) ([]*models.Actor, error) {
	return a.store.GetAllActors(limit, offset)
}

// Получение актеров с фильмами с пагинацией
func (a actor) GetActorsWithMovies(limit, offset int) ([]*models.Actor, error) {
	return a.store.GetActorsWithMovies(limit, offset)
}
