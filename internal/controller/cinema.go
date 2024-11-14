package controller

import (
	"cinema/internal/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// Добавление актера
func (c *Cinema) AddActor(ctx *gin.Context) {
	var newActor models.CreateActor
	if err := ctx.ShouldBindJSON(&newActor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"}) // status 400
		return
	}
	actorID, err := c.actor.AddActor(newActor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // status 500
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"actor_id": actorID}) // status 200
}

// Получение актера по ID
func (c *Cinema) GetActor(ctx *gin.Context) {
	idParam := ctx.Query("id")
	actorID, err := uuid.Parse(idParam)
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

// Обновление актера
func (c *Cinema) UpdateActor(ctx *gin.Context) {
	idParam := ctx.Query("id")
	actorID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid actor ID"})
		return
	}

	// Декодируем JSON из тела запроса в структуру UpdateActor
	var updateActor models.UpdateActor
	if err := ctx.ShouldBindJSON(&updateActor); err != nil {
		// Если тело запроса недействительно, возвращаем статус 400
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Вызываем метод обновления актера
	if err := c.actor.UpdateActor(actorID, updateActor); err != nil {
		// Если произошла ошибка при обновлении актера, возвращаем статус 500
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Успешный ответ, статус 204 без содержимого
	ctx.Status(http.StatusNoContent)
}

// Удаление актера
func (c *Cinema) DeleteActor(ctx *gin.Context) {
	idParam := ctx.Query("id")
	actorID, err := uuid.Parse(idParam)
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

// Добавить фильм
func (c *Cinema) AddMovie(ctx *gin.Context) {
	var newMovie models.CreateMovie

	if err := ctx.ShouldBindJSON(&newMovie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	id, err := c.movie.AddMovie(newMovie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

// Добавить связь между актером и фильмом
func (c *Cinema) AddMovieActorRelation(ctx *gin.Context) {
	var relation struct {
		ActorID uuid.UUID `json:"actor_id"`
		MovieID uuid.UUID `json:"movie_id"`
	}

	if err := ctx.ShouldBindJSON(&relation); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := c.movie.AddMovieActorRelation(relation.ActorID, relation.MovieID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Удалить связь между актером и фильмом
func (c *Cinema) RemoveMovieActorRelation(ctx *gin.Context) {
	var relation struct {
		ActorID uuid.UUID `json:"actor_id"`
		MovieID uuid.UUID `json:"movie_id"`
	}

	if err := ctx.ShouldBindJSON(&relation); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := c.movie.RemoveMovieActorRelation(relation.ActorID, relation.MovieID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Получить фильм по ID
func (c *Cinema) GetMovieByID(ctx *gin.Context) {
	// Получаем параметр id из URL
	idStr := ctx.DefaultQuery("id", "")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	// Получаем фильм по ID
	movie, err := c.movie.GetMovieByID(id)
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
	sortBy := ctx.DefaultQuery("sortBy", "")
	order := ctx.DefaultQuery("order", "")

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
func (c *Cinema) SearchMoviesByTitle(ctx *gin.Context) {
	// Получаем фрагмент названия фильма из параметров запроса
	titleFragment := ctx.DefaultQuery("title", "")

	// Получаем лимит и смещение из запроса
	limit, offset, err := parseLimitOffset(ctx)
	if err != nil {
		// Если есть ошибка с лимитом или смещением, возвращаем ошибку
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем фильмы, которые соответствуют фрагменту названия
	movies, err := c.movie.SearchMoviesByTitle(titleFragment, limit, offset)
	if err != nil {
		// Если произошла ошибка при поиске фильмов, возвращаем ошибку
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем найденные фильмы в формате JSON
	ctx.JSON(http.StatusOK, movies)
}

// Поиск фильмов по актеру
func (c *Cinema) SearchMoviesByActorName(ctx *gin.Context) {
	actorNameFragment := ctx.DefaultQuery("actor_name", "")

	limit, offset, err := parseLimitOffset(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movies, err := c.movie.SearchMoviesByActorName(actorNameFragment, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, movies)
}

// Обновить фильм
func (c *Cinema) UpdateMovie(ctx *gin.Context) {
	idStr := ctx.DefaultQuery("id", "")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	// Декодируем тело запроса в структуру updatedMovie
	var updatedMovie models.UpdateMovie
	if err := ctx.ShouldBindJSON(&updatedMovie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := c.movie.UpdateMovie(id, updatedMovie); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем статус 204 (No Content), так как обновление прошло успешно
	ctx.Status(http.StatusNoContent)
}

// Удалить фильм по ID
func (c *Cinema) DeleteMovie(ctx *gin.Context) {
	// Получаем ID фильма из параметров запроса
	idStr := ctx.DefaultQuery("id", "")
	id, err := uuid.Parse(idStr)
	if err != nil {
		// Если ID фильма некорректен, возвращаем ошибку 400
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	// Удаляем фильм
	if err := c.movie.DeleteMovie(id); err != nil {
		// Если произошла ошибка при удалении фильма, возвращаем ошибку 500
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем статус 204 (No Content), так как удаление прошло успешно
	ctx.Status(http.StatusNoContent)
}

// Получить все фильмы
func (c *Cinema) GetAllMovies(ctx *gin.Context) {
	// Получаем лимит и смещение из параметров запроса
	limit, offset, err := parseLimitOffset(ctx)
	if err != nil {
		// Если есть ошибка с лимитом или смещением, возвращаем ошибку 400
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем все фильмы с пагинацией
	movies, err := c.movie.GetAllMovies(limit, offset)
	if err != nil {
		// Если произошла ошибка при получении фильмов, возвращаем ошибку 500
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем фильмы в формате JSON
	ctx.JSON(http.StatusOK, movies)
}
