package routes

import (
	"cinema/internal/controller"
	"cinema/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cinemaController *controller.Cinema) {
	// Маршруты для актеров
	actorGroup := router.Group("/api/actors")
	{
		actorGroup.GET("/:id", cinemaController.GetActor)                    // Получить актера по ID
		actorGroup.GET("/", cinemaController.GetAllActors)                   // Получить всех актеров
		actorGroup.GET("/with_movies", cinemaController.GetActorsWithMovies) // Актеры с фильмами
	}

	// Маршруты для фильмов
	movieGroup := router.Group("/api/movies")
	{
		movieGroup.GET("/:id", cinemaController.GetMovieByID)                          // Получить фильм по ID
		movieGroup.GET("/by_actor", cinemaController.GetMoviesByActorID)               // Фильмы по актеру
		movieGroup.GET("/filters", cinemaController.GetMoviesWithFilters)              // Фильтрация фильмов
		movieGroup.GET("/search/search", cinemaController.SearchMoviesByTitleAndActor) // Поиск по актеру и названию
		movieGroup.GET("/search/filters", cinemaController.GetMoviesWithFilters)       // Получить фильмы с фильтрами
	}

	// Административные маршруты связей
	relationGroup := router.Group("/api/movie_actors")
	{
		relationGroup.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware([]string{"admin"}))

		relationGroup.POST("/add", cinemaController.AddMovieActorRelations)                 // Добавить связи
		relationGroup.DELETE("/delete", cinemaController.RemoveSelectedMovieActorRelations) // Удалить связи
		relationGroup.PUT("update", cinemaController.UpdateMovieActorsRelations)            // Обновить связи
	}

	// Административные маршруты
	adminGroup := router.Group("/api")
	{
		adminGroup.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware([]string{"admin"}))

		adminGroup.POST("/movies", cinemaController.AddMovie)    // Добавить фильм (admin)
		adminGroup.POST("/movies", cinemaController.UpdateMovie) // Обновить фильм (admin)
		adminGroup.POST("/movies", cinemaController.DeleteMovie) // Удалить фильм (admin)

		adminGroup.POST("/actors", cinemaController.AddActor)          // Добавить актера (admin)
		adminGroup.PUT("/actors/:id", cinemaController.UpdateActor)    // Обновить актера (admin)
		adminGroup.DELETE("/actors/:id", cinemaController.DeleteActor) // Удалить актера (admin)
	}
}
