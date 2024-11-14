package main

import (
	"cinema/internal/controller"
	"cinema/internal/postgres"
	"cinema/internal/repository"
	"cinema/internal/routes"
	"cinema/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func Run() error {
	r := gin.Default()
	// Подключение к базе данных
	db, err := postgres.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Инициализация слоев репозиториев и сервисов
	movieStore := repository.NewMovie(db)
	actorStore := repository.NewActor(db)

	movieService := service.NewMovie(movieStore)
	actorService := service.NewActor(actorStore)

	// Создание контроллера для работы с фильмами и актерами
	cinemaController := controller.NewCinema(movieService, actorService)

	// Настройка маршрутов
	routes.SetupRoutes(r, cinemaController)

	r.Run(":8080")
	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal("Error running application:", err)
	}
}
