package routes

import (
	"cinema/internal/controller"
	"cinema/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cinemaController *controller.Cinema) {
	// Акторы (общие маршруты для получения и добавления актеров)
	actorGroup := router.Group("/api/actors")
	{
		actorGroup.GET("/:id", cinemaController.GetActor)                    // Получить актера по ID
		actorGroup.GET("/", cinemaController.GetAllActors)                   // Получить всех актеров
		actorGroup.GET("/with_movies", cinemaController.GetActorsWithMovies) // Актеры с фильмами
	}

	// Фильмы (общие маршруты для получения и добавления фильмов)
	movieGroup := router.Group("/api/movies")
	{
		movieGroup.GET("/:id", cinemaController.GetMovieByID)                     // Получить фильм по ID
		movieGroup.GET("/", cinemaController.GetAllMovies)                        // Получить все фильмы
		movieGroup.GET("/filters", cinemaController.GetMoviesWithFilters)         // Фильтрация фильмов
		movieGroup.GET("/search/title", cinemaController.SearchMoviesByTitle)     // Поиск фильмов по названию
		movieGroup.GET("/search/actor", cinemaController.SearchMoviesByActorName) // Поиск фильмов по актеру
		movieGroup.GET("/actor_movies", cinemaController.GetMoviesByActorID)      // Фильмы по актору
	}

	// Связь актеров и фильмов
	relationGroup := router.Group("/api/movie_actor")
	{
		relationGroup.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware([]string{"admin"}))

		relationGroup.POST("/", cinemaController.AddMovieActorRelation)      // Добавить связь
		relationGroup.DELETE("/", cinemaController.RemoveMovieActorRelation) // Удалить связь
	}

	// Административные маршруты
	adminGroup := router.Group("/api")
	{
		adminGroup.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware([]string{"admin"}))

		adminGroup.POST("/movies", cinemaController.AddMovie) // Добавить фильм (admin)
		adminGroup.POST("/actors", cinemaController.AddActor) // Добавить актера (admin)

		adminGroup.PUT("/actors/:id", cinemaController.UpdateActor)    // Обновить актера (admin)
		adminGroup.DELETE("/actors/:id", cinemaController.DeleteActor) // Удалить актера (admin)
	}
}
