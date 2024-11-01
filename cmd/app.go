package main

import (
	"cinema/internal/controller"
	"cinema/internal/postgres"
	"cinema/internal/repository"
	"cinema/internal/service"
	"log"

	_ "github.com/lib/pq"
)

func Run() error {
	db, err := postgres.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	movieRepo := repository.NewMovie(db)
	actorRepo := repository.NewActor(db)
	movieService := service.NewMovie(movieRepo)
	actorService := service.NewActor(actorRepo)
	cinemaCtrl := controller.NewCinema(movieService, actorService)
	movies := cinemaCtrl.GetAllActors()

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal("Error running application:", err)
	}
}
