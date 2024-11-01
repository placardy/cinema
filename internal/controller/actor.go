package controller

import (
	"cinema/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type serviceActor interface {
	AddActor(actor models.CreateActor) (uuid.UUID, error)
	GetActor(id uuid.UUID) (*models.Actor, error)
	UpdateActor(id uuid.UUID, actor models.UpdateActor) error
	DeleteActor(id uuid.UUID) error
	GetAllActors(limit, offset int) ([]*models.Actor, error)
	GetActorsWithMovies(limit, offset int) ([]*models.Actor, error)
}

type actor struct {
	service serviceActor
}

func NewActor(service serviceActor) *actor {
	return &actor{service: service}
}

// Добавление актера
func (c *actor) AddActor(w http.ResponseWriter, r *http.Request) {
	var newActor models.CreateActor
	if err := json.NewDecoder(r.Body).Decode(&newActor); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	actorID, err := c.service.AddActor(newActor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]uuid.UUID{"actor_id": actorID})
}

// Получение актера по ID
func (c *actor) GetActor(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	actorID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}
	actor, err := c.service.GetActor(actorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if actor == nil {
		http.Error(w, "Actor not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(actor)
}

// Обновление актера
func (c *actor) UpdateActor(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	actorID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}

	var updateActor models.UpdateActor
	if err := json.NewDecoder(r.Body).Decode(&updateActor); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := c.service.UpdateActor(actorID, updateActor); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Удаление актера
func (c *actor) DeleteActor(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	actorID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}
	if err := c.service.DeleteActor(actorID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Получение всех актеров с пагинацией
func (c *actor) GetAllActors(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actors, err := c.service.GetAllActors(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(actors)
}

// Получение актеров с фильмами
func (c *actor) GetActorsWithMovies(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actors, err := c.service.GetActorsWithMovies(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(actors)
}

func parseLimitOffset(r *http.Request) (int, int, error) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		return 0, 0, fmt.Errorf("invalid or missing limit")
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		return 0, 0, fmt.Errorf("invalid or missing offset")
	}
	return limit, offset, nil
}
