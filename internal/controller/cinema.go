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
func (c *cinema) AddActor(w http.ResponseWriter, r *http.Request) {
	var newActor models.CreateActor
	if err := json.NewDecoder(r.Body).Decode(&newActor); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	actorID, err := c.actor.AddActor(newActor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]uuid.UUID{"actor_id": actorID})
}

// Получение актера по ID
func (c *cinema) GetActor(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	actorID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}
	actor, err := c.actor.GetActor(actorID)
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
func (c *cinema) UpdateActor(w http.ResponseWriter, r *http.Request) {
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
	if err := c.actor.UpdateActor(actorID, updateActor); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Удаление актера
func (c *cinema) DeleteActor(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	actorID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}
	if err := c.actor.DeleteActor(actorID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Получение всех актеров с пагинацией
func (c *cinema) GetAllActors(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actors, err := c.actor.GetAllActors(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(actors)
}

// Получение актеров с фильмами
func (c *cinema) GetActorsWithMovies(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actors, err := c.actor.GetActorsWithMovies(limit, offset)
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
func (c *cinema) AddMovie(w http.ResponseWriter, r *http.Request) {
	var newMovie models.CreateMovie
	if err := json.NewDecoder(r.Body).Decode(&newMovie); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := c.movie.AddMovie(newMovie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uuid.UUID{"id": id})
}

// Добавить связь между актером и фильмом
func (c *cinema) AddMovieActorRelation(w http.ResponseWriter, r *http.Request) {
	var relation struct {
		ActorID uuid.UUID `json:"actor_id"`
		MovieID uuid.UUID `json:"movie_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&relation); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.movie.AddMovieActorRelation(relation.ActorID, relation.MovieID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Удалить связь между актером и фильмом
func (c *cinema) RemoveMovieActorRelation(w http.ResponseWriter, r *http.Request) {
	var relation struct {
		ActorID uuid.UUID `json:"actor_id"`
		MovieID uuid.UUID `json:"movie_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&relation); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.movie.RemoveMovieActorRelation(relation.ActorID, relation.MovieID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Получить фильм по ID
func (c *cinema) GetMovieByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	movie, err := c.movie.GetMovieByID(id)
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
func (c *cinema) GetMoviesByActorID(w http.ResponseWriter, r *http.Request) {
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
	movies, err := c.movie.GetMoviesByActorID(actorID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// Получить фильмы с фильтрацией
func (c *cinema) GetMoviesWithFilters(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	order := r.URL.Query().Get("order")
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movies, err := c.movie.GetMoviesWithFilters(sortBy, order, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// Поиск фильмов по названию
func (c *cinema) SearchMoviesByTitle(w http.ResponseWriter, r *http.Request) {
	titleFragment := r.URL.Query().Get("title")
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movies, err := c.movie.SearchMoviesByTitle(titleFragment, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// Поиск фильмов по актеру
func (c *cinema) SearchMoviesByActorName(w http.ResponseWriter, r *http.Request) {
	actorNameFragment := r.URL.Query().Get("actor_name")
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movies, err := c.movie.SearchMoviesByActorName(actorNameFragment, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// Обновить фильм
func (c *cinema) UpdateMovie(w http.ResponseWriter, r *http.Request) {
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

	if err := c.movie.UpdateMovie(id, updatedMovie); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Удалить фильм по ID
func (c *cinema) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	if err := c.movie.DeleteMovie(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Получить все фильмы
func (c *cinema) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movies, err := c.movie.GetAllMovies(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
