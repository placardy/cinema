package main

import (
	//"cinema/internal/controller"
	"cinema/internal/postgres"
	"cinema/internal/repository"
	"cinema/internal/service"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Run() error {
	db, err := postgres.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// movieStore := repository.NewMovie(db)
	actorStore := repository.NewActor(db)

	// movieService := service.NewMovie(movieStore)
	actorService := service.NewActor(actorStore)
	actors, nil := actorService.GetActorsWithMovies(5, 1)
	for _, actor := range actors {
		fmt.Println(actor.Name)
	}
	// cinemaController := controller.NewCinema(movieService, actorService)
	//cinemaController.AddMovie() GIT TEST
	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal("Error running application:", err)
	}
}
