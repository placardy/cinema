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

type serviceMovie interface {
	AddMovie(movie models.CreateMovie) (uuid.UUID, error)
	AddMovieActorRelation(actorID, movieID uuid.UUID) error
	RemoveMovieActorRelation(actorID, movieID uuid.UUID) error
	GetMovieByID(id uuid.UUID) (*models.Movie, error)
	GetMoviesByActorID(actorID uuid.UUID, limit, offset int) ([]*models.Movie, error)
	GetMoviesWithFilters(sortBy string, order string, limit, offset int) ([]*models.Movie, error)
	SearchMoviesByTitle(titleFragment string, limit, offset int) ([]*models.Movie, error)
	SearchMoviesByActorName(actorNameFragment string, limit, offset int) ([]*models.Movie, error)
	UpdateMovie(id uuid.UUID, movie models.UpdateMovie) error
	DeleteMovie(id uuid.UUID) error
	GetAllMovies(limit, offset int) ([]*models.Movie, error)
}

type cinema struct {
	movie serviceMovie
	actor serviceActor
}

func NewCinema(movie serviceMovie, actor serviceActor) *cinema {
	return &cinema{movie: movie, actor: actor}
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

// Добавить фильм
func (m *movie) AddMovie(w http.ResponseWriter, r *http.Request) {
	var newMovie models.CreateMovie
	if err := json.NewDecoder(r.Body).Decode(&newMovie); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := m.service.AddMovie(newMovie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uuid.UUID{"id": id})
}

// Добавить связь между актером и фильмом
func (m *movie) AddMovieActorRelation(w http.ResponseWriter, r *http.Request) {
	var relation struct {
		ActorID uuid.UUID `json:"actor_id"`
		MovieID uuid.UUID `json:"movie_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&relation); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := m.service.AddMovieActorRelation(relation.ActorID, relation.MovieID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Удалить связь между актером и фильмом
func (m *movie) RemoveMovieActorRelation(w http.ResponseWriter, r *http.Request) {
	var relation struct {
		ActorID uuid.UUID `json:"actor_id"`
		MovieID uuid.UUID `json:"movie_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&relation); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := m.service.RemoveMovieActorRelation(relation.ActorID, relation.MovieID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Получить фильм по ID
func (m *movie) GetMovieByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	movie, err := m.service.GetMovieByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if movie == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

// Получить фильмы по актеру
func (m *movie) GetMoviesByActorID(w http.ResponseWriter, r *http.Request) {
	actorIDStr := r.URL.Query().Get("actor_id")
	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}

	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	movies, err := m.service.GetMoviesByActorID(actorID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// Получить фильмы с фильтрацией
func (m *movie) GetMoviesWithFilters(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	order := r.URL.Query().Get("order")
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movies, err := m.service.GetMoviesWithFilters(sortBy, order, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// Поиск фильмов по названию
func (m *movie) SearchMoviesByTitle(w http.ResponseWriter, r *http.Request) {
	titleFragment := r.URL.Query().Get("title")
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movies, err := m.service.SearchMoviesByTitle(titleFragment, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// Поиск фильмов по актеру
func (m *movie) SearchMoviesByActorName(w http.ResponseWriter, r *http.Request) {
	actorNameFragment := r.URL.Query().Get("actor_name")
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movies, err := m.service.SearchMoviesByActorName(actorNameFragment, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// Обновить фильм
func (m *movie) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	var updatedMovie models.UpdateMovie
	if err := json.NewDecoder(r.Body).Decode(&updatedMovie); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := m.service.UpdateMovie(id, updatedMovie); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Удалить фильм по ID
func (m *movie) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	if err := m.service.DeleteMovie(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Получить все фильмы
func (m *movie) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movies, err := m.service.GetAllMovies(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
