package controller

import (
	_ "cinema/docs"
	"cinema/internal/models"
	"cinema/internal/utils"
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
	CreateMovie(movie models.CreateMovie) (uuid.UUID, error)
	GetMovieByID(id uuid.UUID) (*models.Movie, error)
	GetMoviesByActorID(actorID uuid.UUID, limit, offset int) ([]*models.Movie, error)
	GetMoviesWithFilters(sortBy string, order string, limit, offset int) ([]*models.Movie, error)
	SearchMoviesByTitleAndActor(filterTitle, filterActor string, limit, offset int) ([]*models.Movie, error)
	UpdateMovie(id uuid.UUID, movie models.UpdateMovie) error
	DeleteMovie(id uuid.UUID) error
}

type serviceActor interface {
	CreateActor(actor models.CreateActor) (uuid.UUID, error)
	GetActor(id uuid.UUID) (*models.Actor, error)
	GetAllActors(limit, offset int) ([]*models.Actor, error)
	GetActorsWithMovies(limit, offset int) ([]*models.ActorWithMovies, error)
	UpdateActor(id uuid.UUID, actor models.UpdateActor) error
	DeleteActor(id uuid.UUID) error
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
		return 0, 0, fmt.Errorf("invalid or missing limit: %s", ctx.Query("limit"))
	}
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil || offset < 0 {
		return 0, 0, fmt.Errorf("invalid or missing limit: %s", ctx.Query("limit"))
	}
	return limit, offset, nil
}

// AddMovieActorRelations godoc
// @Summary      Add actors to a movie
// @Description  Add a list of actors to a movie by movie ID
// @Tags         movie-actors
// @Accept       json
// @Produce      json
// @Param        movie_id  path     string  true  "ID of the movie"
// @Param        actor_ids body     []uuid.UUID  true  "List of actor IDs to be added to the movie"
// @Success      204       {string}  string  "No Content"
// @Failure      400       {object}  models.APIError "Invalid request body"
// @Failure      500       {object}  models.APIError "Internal Server Error"
// @Router       /api/movies/{movie_id}/actors [post]
func (c *Cinema) AddMovieActorRelations(ctx *gin.Context) {
	movieIDStr := ctx.Param("movie_id") // Получаем movie_id из параметров пути
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		// Если преобразование не удалось, возвращаем ошибку 400
		utils.BadRequestResponse(ctx, "Invalid movie ID format")
		return
	}

	var actorIDs []uuid.UUID
	if err := ctx.ShouldBindJSON(&actorIDs); err != nil {
		utils.InvalidJSONResponse(ctx)
		return
	}

	if err := c.movie.AddMovieActorRelations(movieID, actorIDs); err != nil {
		utils.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}

// UpdateMovieActorRelations godoc
// @Summary      Update actors for a movie
// @Description  Update the list of actors for a movie by movie ID
// @Tags         movie-actors
// @Accept       json
// @Produce      json
// @Param        movie_id  path     string  true  "ID of the movie"
// @Param        actor_ids body     []uuid.UUID  true  "List of actor IDs to be updated for the movie"
// @Success      204       {string}  string  "No Content"
// @Failure      400       {object}  models.APIError "Invalid request body"
// @Failure      500       {object}  models.APIError "Internal Server Error"
// @Router       /api/movies/{movie_id}/actors [put]
func (c *Cinema) UpdateMovieActorRelations(ctx *gin.Context) {
	movieIDStr := ctx.Param("movie_id") // Получаем movie_id из параметров пути
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		// Если преобразование не удалось, возвращаем ошибку 400
		utils.BadRequestResponse(ctx, "Invalid movie ID format")
		return
	}

	var actorIDs []uuid.UUID
	if err := ctx.ShouldBindJSON(&actorIDs); err != nil {
		utils.InvalidJSONResponse(ctx)
		return
	}

	if err := c.movie.UpdateMovieActorRelations(movieID, actorIDs); err != nil {
		utils.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}

// RemoveSelectedMovieActorRelations godoc
// @Summary      Remove actors from a movie
// @Description  Remove a list of actors from a movie by movie ID
// @Tags         movie-actors
// @Accept       json
// @Produce      json
// @Param        movie_id  path     string  true  "ID of the movie"
// @Param        actor_ids body     []uuid.UUID  true  "List of actor IDs to be removed from the movie"
// @Success      204       {string}  string  "No Content"
// @Failure      400       {object}  models.APIError "Invalid request body"
// @Failure      500       {object}  models.APIError "Internal Server Error"
// @Router       /api/movies/{movie_id}/actors [delete]
func (c *Cinema) RemoveSelectedMovieActorRelations(ctx *gin.Context) {
	movieIDStr := ctx.Param("movie_id") // Получаем movie_id из параметров пути
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		// Если преобразование не удалось, возвращаем ошибку 400
		utils.BadRequestResponse(ctx, "Invalid movie ID format")
		return
	}

	var actorIDs []uuid.UUID
	if err := ctx.ShouldBindJSON(&actorIDs); err != nil {
		utils.InvalidJSONResponse(ctx)
		return
	}

	if err := c.movie.RemoveSelectedMovieActorRelations(movieID, actorIDs); err != nil {
		utils.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}

// CreateMovie godoc
// @Summary      Create a new movie
// @Description  Adds a new movie to the database
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        movie  body      models.CreateMovie  true  "Movie details"
// @Success      201    {object}  map[string]string       "Movie created successfully"  example={"id": "1234"}
// @Failure      400    {object}  models.APIError  "Invalid JSON format or validation errors"
// @Failure      500    {object}  models.APIError  "Internal server error"
// @Router       /api/movies [post]
func (c *Cinema) CreateMovie(ctx *gin.Context) {
	var newMovie models.CreateMovie

	if err := ctx.ShouldBindJSON(&newMovie); err != nil {
		utils.InvalidJSONResponse(ctx)
		return
	}

	validationErrors := newMovie.Validate()
	if len(validationErrors) > 0 {
		utils.ValidationErrorResponse(ctx, validationErrors)
		return
	}

	id, err := c.movie.CreateMovie(newMovie)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetMovieByID godoc
// @Summary      Get movie by ID
// @Description  Retrieves a movie by its unique identifier
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        movie_id   path     string  true  "Movie ID" example("f47ac10b-58cc-4372-a567-0e02b2c3d479")
// @Success      200  {object}  models.Movie  "The movie details"
// @Failure      400  {object}  models.APIError  "Invalid movie ID format"
// @Failure      404  {object}  models.APIError  "Movie not found"
// @Failure      500  {object}  models.APIError  "Internal server error"
// @Router       /api/movies/{id} [get]
func (c *Cinema) GetMovieByID(ctx *gin.Context) {
	// Получаем параметр id из URL
	movieIDStr := ctx.Param("id")
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid movie ID format")
		return
	}

	// Получаем фильм по ID
	movie, err := c.movie.GetMovieByID(movieID)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	if movie == nil {
		utils.NotFoundResponse(ctx, "Movie not found")
		return
	}

	// Если фильм найден, возвращаем его с кодом 200
	ctx.JSON(http.StatusOK, movie)
}

// GetMoviesByActorID godoc
// @Summary      Get movies by actor ID
// @Description  Retrieve a list of movies by the actor's ID with optional pagination
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        actor_id  path    string  true  "Actor ID"  Format: uuid
// @Param        limit     query   int     false "Limit the number of movies returned"
// @Param        offset    query   int     false "Offset for pagination"
// @Success      200       {array}  models.Movie  "List of movies"
// @Failure      400       {object}  models.APIError "Invalid actor ID format or bad request"
// @Failure      500       {object}  models.APIError "Internal Server Error"
// @Router       /api/actors/{actor_id}/movies [get]
func (c *Cinema) GetMoviesByActorID(ctx *gin.Context) {
	actorIDStr := ctx.Param("id")
	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		// Если ID актера некорректен, возвращаем ошибку 400
		utils.BadRequestResponse(ctx, "Invalid actor ID format")
		return
	}

	limit, offset, err := parseLimitOffset(ctx)
	if err != nil {
		utils.BadRequestResponse(ctx, err.Error())
		return
	}

	movies, err := c.movie.GetMoviesByActorID(actorID, limit, offset)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, movies)
}

// GetMoviesWithFilters godoc
// @Summary      Get movies with filters
// @Description  Retrieve a list of movies with optional filters for sorting and pagination.
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        sortBy  query   string  false "Field to sort by" Enums(title, release_date, rating) default(rating)
// @Param        order   query   string  false "Sorting order" Enums(ASC, DESC) default(DESC)
// @Param        limit   query   int     false "Limit the number of movies returned" default(10)
// @Param        offset  query   int     false "Offset for pagination" default(0)
// @Success      200     {array} models.Movie   "List of filtered movies"
// @Failure      400     {object} models.APIError "Invalid request parameters"
// @Failure      500     {object} models.APIError "Internal server error"
// @Router       /api/movies [get]
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

// SearchMoviesByTitleAndActor godoc
// @Summary      Search movies by title and actor name
// @Description  Search for movies by a partial title and actor's name with optional pagination.
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        title        query   string  false "Movie title fragment" example("Inception")
// @Param        actor_name   query   string  false "Actor name fragment" example("Leonardo")
// @Param        limit        query   int     false "Limit the number of movies returned" default(10)
// @Param        offset       query   int     false "Offset for pagination" default(0)
// @Success      200          {array} models.Movie "List of movies matching the search"
// @Failure      400          {object} models.APIError "Invalid search parameters"
// @Failure      500          {object} models.APIError "Internal server error"
// @Router       /api/movies/search [get]
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

// UpdateMovie godoc
// @Summary      Update movie details
// @Description  Update the details of a movie based on its ID
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        movie_id        path    string  true  "Movie ID" example("f47ac10b-58cc-4372-a567-0e02b2c3d479")
// @Param        movie     body    models.UpdateMovie true "Updated movie details"
// @Success      204       "Movie successfully updated"
// @Failure      400       {object} models.APIError "Invalid request body or parameters"
// @Failure      500       {object} models.APIError "Internal server error"
// @Router       /api/movies/{id} [put]
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

// DeleteMovie godoc
// @Summary      Delete movie
// @Description  Delete a movie by its ID
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        movie_id   path    string  true  "Movie ID" example("f47ac10b-58cc-4372-a567-0e02b2c3d479")
// @Success      204  "Movie successfully deleted"
// @Failure      400  {object} models.APIError "Invalid movie ID"
// @Failure      500  {object} models.APIError "Internal server error"
// @Router       /api/movies/{id} [delete]
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

// CreateActor godoc
// @Summary      Create a new actor
// @Description  Create a new actor with the provided data
// @Tags         Actors
// @Accept       json
// @Produce      json
// @Param        actor     body    models.CreateActor true "New actor details"
// @Success      201       {object} map[string]string "Actor ID" example({ "actor_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479" })
// @Failure      400       {object} models.APIError "Invalid request body or validation errors"
// @Failure      500       {object} models.APIError "Internal server error"
// @Router       /api/actors [post]
func (c *Cinema) CreateActor(ctx *gin.Context) {
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

	actorID, err := c.actor.CreateActor(newActor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // status 500
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"actor_id": actorID}) // status 200
}

// GetActor godoc
// @Summary      Get actor by ID
// @Description  Retrieve an actor's details by their unique ID
// @Tags         Actors
// @Accept       json
// @Produce      json
// @Param        actor_id    path    string  true  "Actor ID" example("f47ac10b-58cc-4372-a567-0e02b2c3d479")
// @Success      200   {object} models.Actor "Actor details"
// @Failure      400   {object} models.APIError "Invalid actor ID"
// @Failure      404   {object} models.APIError "Actor not found"
// @Failure      500   {object} models.APIError "Internal server error"
// @Router       /api/actors/{id} [get]
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

// GetAllActors godoc
// @Summary      Get all actors with pagination
// @Description  Retrieve a list of all actors with optional pagination
// @Tags         Actors
// @Accept       json
// @Produce      json
// @Param        limit   query   int     false "Limit the number of actors returned" default(10)
// @Param        offset  query   int     false "Offset for pagination" default(0)
// @Success      200     {array} models.Actor "List of actors"
// @Failure      400     {object} models.APIError "Invalid pagination parameters"
// @Failure      500     {object} models.APIError "Internal server error"
// @Router       /api/actors [get]
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

// GetActorsWithMovies godoc
// @Summary      Get actors with their movies
// @Description  Retrieve a list of actors with the movies they have appeared in
// @Tags         Actors
// @Accept       json
// @Produce      json
// @Param        limit   query   int     false "Limit the number of actors returned" default(10)
// @Param        offset  query   int     false "Offset for pagination" default(0)
// @Success      200     {array} models.ActorWithMovies "List of actors with movies"
// @Failure      400     {object} models.APIError "Invalid pagination parameters"
// @Failure      500     {object} models.APIError "Internal server error"
// @Router       /api/actors/with-movies [get]
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

// UpdateActor godoc
// @Summary      Update actor details
// @Description  Update an actor's details based on their ID
// @Tags         Actors
// @Accept       json
// @Produce      json
// @Param        actor_id      path    string  true  "Actor ID" example("f47ac10b-58cc-4372-a567-0e02b2c3d479")
// @Param        actor   body    models.UpdateActor true "Updated actor details"
// @Success      204     "Actor successfully updated"
// @Failure      400     {object} models.APIError "Invalid request body or parameters"
// @Failure      500     {object} models.APIError "Internal server error"
// @Router       /api/actors/{id} [put]
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

// DeleteActor godoc
// @Summary      Delete actor
// @Description  Delete an actor based on their ID
// @Tags         Actors
// @Accept       json
// @Produce      json
// @Param        actor_id   path    string  true  "Actor ID" example("f47ac10b-58cc-4372-a567-0e02b2c3d479")
// @Success      204  "Actor successfully deleted"
// @Failure      400  {object} models.APIError "Invalid actor ID"
// @Failure      500  {object} models.APIError "Internal server error"
// @Router       /api/actors/{id} [delete]
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
