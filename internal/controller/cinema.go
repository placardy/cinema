package controller

import (
	"cinema/internal/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type serviceMovie interface {
	// связи
	AddMovieActorRelations(movieID uuid.UUID, actorIDs []uuid.UUID) error
	RemoveSelectedMovieActorRelations(movieID uuid.UUID, actorIDs []uuid.UUID) error
	UpdateMovieActorRelations(movieID uuid.UUID, actorIDs []uuid.UUID) error
	// фильмы
	AddMovie(movie models.CreateMovie) (uuid.UUID, error)
	GetMovieByID(id uuid.UUID) (*models.Movie, error)
	GetMoviesByActorID(actorID uuid.UUID, limit, offset int) ([]*models.Movie, error)
	GetMoviesWithFilters(sortBy string, order string, limit, offset int) ([]*models.Movie, error)
	SearchMoviesByTitleAndActor(filterTitle, filterActor string, limit, offset int) ([]*models.Movie, error)
	UpdateMovie(id uuid.UUID, movie models.UpdateMovie) error
	DeleteMovie(id uuid.UUID) error
}

type serviceActor interface {
	AddActor(actor models.CreateActor) (uuid.UUID, error)
	GetActor(id uuid.UUID) (*models.Actor, error)
	UpdateActor(id uuid.UUID, actor models.UpdateActor) error
	DeleteActor(id uuid.UUID) error
	GetAllActors(limit, offset int) ([]*models.Actor, error)
	GetActorsWithMovies(limit, offset int) ([]*models.ActorWithMovies, error)
}

type Cinema struct {
	movie serviceMovie
	actor serviceActor
}

func NewCinema(movie serviceMovie, actor serviceActor) *Cinema {
	return &Cinema{movie: movie, actor: actor}
}

func parseLimitOffset(ctx *gin.Context) (int, int, error) {
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil || limit <= 0 {
		return 0, 0, fmt.Errorf("invalid or missing limit")
	}
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil || offset < 0 {
		return 0, 0, fmt.Errorf("invalid or missing offset")
	}
	return limit, offset, nil
}

// Добавить связь между актером и фильмом
func (c *Cinema) AddMovieActorRelations(ctx *gin.Context) {
	var relation struct {
		MovieID  uuid.UUID   `json:"movie_id"`
		ActorIDs []uuid.UUID `json:"actor_ids"`
	}

	if err := ctx.ShouldBindJSON(&relation); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := c.movie.AddMovieActorRelations(relation.MovieID, relation.ActorIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Удалить связь между актером и фильмом
func (c *Cinema) RemoveSelectedMovieActorRelations(ctx *gin.Context) {
	var relation struct {
		MovieID  uuid.UUID   `json:"movie_id"`
		ActorIDs []uuid.UUID `json:"actor_ids"`
	}

	if err := ctx.ShouldBindJSON(&relation); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := c.movie.RemoveSelectedMovieActorRelations(relation.MovieID, relation.ActorIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Удалить связь между актером и фильмом
func (c *Cinema) UpdateMovieActorRelations(ctx *gin.Context) {
	var relation struct {
		MovieID  uuid.UUID   `json:"movie_id"`
		ActorIDs []uuid.UUID `json:"actor_ids"`
	}

	if err := ctx.ShouldBindJSON(&relation); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := c.movie.UpdateMovieActorRelations(relation.MovieID, relation.ActorIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Добавить фильм
func (c *Cinema) AddMovie(ctx *gin.Context) {
	var newMovie models.CreateMovie

	if err := ctx.ShouldBindJSON(&newMovie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	validationErrors := newMovie.Validate()
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	id, err := c.movie.AddMovie(newMovie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

// Обновить фильм
func (c *Cinema) UpdateMovie(ctx *gin.Context) {
	movieIDStr := ctx.Param("id")
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var updatedMovie models.UpdateMovie
	if err := ctx.ShouldBindJSON(&updatedMovie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	validationErrors := updatedMovie.Validate()
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	if err := c.movie.UpdateMovie(movieID, updatedMovie); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Добавление актера
func (c *Cinema) AddActor(ctx *gin.Context) {
	var newActor models.CreateActor
	if err := ctx.ShouldBindJSON(&newActor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"}) // status 400
		return
	}

	validationErrors := newActor.Validate()
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors}) // status 400
		return
	}

	actorID, err := c.actor.AddActor(newActor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // status 500
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"actor_id": actorID}) // status 200
}

// Обновление актера
func (c *Cinema) UpdateActor(ctx *gin.Context) {
	actorIDStr := ctx.Param("id")
	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid actor ID"})
		return
	}

	var updateActor models.UpdateActor
	if err := ctx.ShouldBindJSON(&updateActor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	validationErrors := updateActor.Validate()
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	if err := c.actor.UpdateActor(actorID, updateActor); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Удаление актера
func (c *Cinema) DeleteActor(ctx *gin.Context) {
	actorIDStr := ctx.Param("id")
	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid actor ID"})
		return
	}
	if err := c.actor.DeleteActor(actorID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// Получение актера по ID
func (c *Cinema) GetActor(ctx *gin.Context) {
	actorIDStr := ctx.Param("id")
	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid actor ID"})
		return
	}
	actor, err := c.actor.GetActor(actorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if actor == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Actor not found"})
		return
	}
	ctx.JSON(http.StatusOK, actor)
}

// Получение всех актеров с пагинацией
func (c *Cinema) GetAllActors(ctx *gin.Context) {
	limit, offset, err := parseLimitOffset(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	actors, err := c.actor.GetAllActors(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, actors)
}

// Получение актеров с фильмами
func (c *Cinema) GetActorsWithMovies(ctx *gin.Context) {
	limit, offset, err := parseLimitOffset(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actors, err := c.actor.GetActorsWithMovies(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, actors)
}

// Получить фильм по ID
func (c *Cinema) GetMovieByID(ctx *gin.Context) {
	// Получаем параметр id из URL
	movieIDStr := ctx.Param("id")
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	// Получаем фильм по ID
	movie, err := c.movie.GetMovieByID(movieID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if movie == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	// Если фильм найден, возвращаем его с кодом 200
	ctx.JSON(http.StatusOK, movie)
}

// Получить фильмы по актеру
func (c *Cinema) GetMoviesByActorID(ctx *gin.Context) {
	actorIDStr := ctx.DefaultQuery("actor_id", "")
	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		// Если ID актера некорректен, возвращаем ошибку 400
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid actor ID"})
		return
	}

	limit, offset, err := parseLimitOffset(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movies, err := c.movie.GetMoviesByActorID(actorID, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, movies)
}

// Получить фильмы с фильтрацией
func (c *Cinema) GetMoviesWithFilters(ctx *gin.Context) {
	sortBy := ctx.DefaultQuery("sortBy", "rating")
	order := ctx.DefaultQuery("order", "DESC")
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
	sortBy = strings.ToLower(ctx.DefaultQuery("sortBy", "rating"))
	order = strings.ToUpper(ctx.DefaultQuery("order", "DESC"))

	limit, offset, err := parseLimitOffset(ctx)
	if err != nil {
		// Если есть ошибка с лимитом или смещением, возвращаем ошибку
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movies, err := c.movie.GetMoviesWithFilters(sortBy, order, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем список фильмов в формате JSON
	ctx.JSON(http.StatusOK, movies)
}

// Поиск фильмов по названию
func (c *Cinema) SearchMoviesByTitleAndActor(ctx *gin.Context) {
	// Получаем фрагмент названия фильма из параметров запроса
	titleFragment := ctx.DefaultQuery("title", "")
	actorNameFragment := ctx.DefaultQuery("actor_name", "")

	// Получаем лимит и смещение из запроса
	limit, offset, err := parseLimitOffset(ctx)
	if err != nil {
		// Если есть ошибка с лимитом или смещением, возвращаем ошибку
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем фильмы, которые соответствуют фрагменту названия
	movies, err := c.movie.SearchMoviesByTitleAndActor(titleFragment, actorNameFragment, limit, offset)
	if err != nil {
		// Если произошла ошибка при поиске фильмов, возвращаем ошибку
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем найденные фильмы в формате JSON
	ctx.JSON(http.StatusOK, movies)
}

// Удалить фильм по ID
func (c *Cinema) DeleteMovie(ctx *gin.Context) {
	movieIDStr := ctx.Param("id")
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	if err := c.movie.DeleteMovie(movieID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
