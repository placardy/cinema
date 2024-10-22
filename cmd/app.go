package main

import (
	"cinema/internal/postgres"
	"cinema/internal/repository"
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

	movieRepo := repository.NewMovie(db)
	actorRepo := repository.NewActor(db)

	// добавить актера
	actorID, _ := actorRepo.CreateActor("TestActor", "male", "1970-05-05")
	// добавить фильм
	movieID, _ := movieRepo.CreateMovie("TestMovie", "testdisc", "2024-05-05", 5)
	// добавить связь актер-фильм
	movieRepo.CreateMovieActorRelation(actorID, movieID)

	// Поиск фильмов по имени актера
	movies, _ := movieRepo.SearchMoviesByActorName("ar")
	for _, movie := range movies {
		fmt.Println(movie.Title)
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal("Error running application:", err)
	}
}
