package routes

import (
	"cinema/internal/controller"
	"cinema/internal/middleware"

	_ "cinema/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, cinemaController *controller.Cinema) {

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Маршруты для актеров
	actorGroup := router.Group("/api/actors")
	{
		actorGroup.GET("/:actor_id", cinemaController.GetActor)                  // Получить актера по ID
		actorGroup.GET("/:actor_id/movies", cinemaController.GetMoviesByActorID) // Получить фильмы по ID актера
		actorGroup.GET("/", cinemaController.GetAllActors)                       // Получить всех актеров
		actorGroup.GET("/with-movies", cinemaController.GetActorsWithMovies)     // Актеры с фильмами
	}

	// Маршруты для фильмов
	movieGroup := router.Group("/api/movies")
	{
		movieGroup.GET("/:movie_id", cinemaController.GetMovieByID)             // Получить фильм по ID
		movieGroup.GET("/", cinemaController.GetMoviesWithFilters)              // Фильтрация фильмов
		movieGroup.GET("/search", cinemaController.SearchMoviesByTitleAndActor) // Поиск по актеру и названию
	}

	// Маршруты для управления связями
	// Создание, удаление и обновление связей между фильмами и актерами
	relationGroup := router.Group("/api/movies/:movie_id/actors")
	{
		relationGroup.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware([]string{"admin"}))

		relationGroup.POST("/", cinemaController.AddMovieActorRelations)              // Добавить актеров в фильм
		relationGroup.DELETE("/", cinemaController.RemoveSelectedMovieActorRelations) // Удалить актеров из фильма
		relationGroup.PUT("/", cinemaController.UpdateMovieActorRelations)            // Обновить актеров в фильме
	}

	// Административные маршруты
	adminGroup := router.Group("/api")
	{
		adminGroup.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware([]string{"admin"}))

		// Фильмы
		adminGroup.POST("/movies", cinemaController.CreateMovie)             // Добавить фильм (admin)
		adminGroup.PUT("/movies/:movie_id", cinemaController.UpdateMovie)    // Обновить фильм (admin) // DONT WORK не работает на уровне сервиса нужно доделать!
		adminGroup.DELETE("/movies/:movie_id", cinemaController.DeleteMovie) // Удалить фильм (admin)

		// Актеры
		adminGroup.POST("/actors", cinemaController.CreateActor)             // Добавить актера (admin)
		adminGroup.PUT("/actors/:actor_id", cinemaController.UpdateActor)    // Обновить актера (admin)
		adminGroup.DELETE("/actors/:actor_id", cinemaController.DeleteActor) // Удалить актера (admin)
	}
}
