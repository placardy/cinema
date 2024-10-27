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
	// birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	// actorID, _ := actorRepo.AddActor(models.CreateActor{"TestActor", "male", birthDate})

	// // добавить фильм
	// release := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	// movieID, _ := movieRepo.AddMovie(models.CreateMovie{"TestMovie", "testdisc", release, 5})

	// добавить связь актер-фильм
	// movieRepo.AddActorToMovieRelation(actorID, movieID)

	// Поиск фильмов по имени актера
	// movies, _ := movieRepo.SearchMoviesByActorName("ar")
	// for _, movie := range movies {
	// 	fmt.Println(movie.Title)
	// }

	// Поиск фильмов по названию
	movies, err := movieRepo.GetMovies("ratidgfdgng", "DESC", 4, 2)
	if err != nil {
		fmt.Println(err)
	}
	for _, movie := range movies {
		fmt.Println("movie:", movie.Title, "rating:", movie.Rating)
	}

	// Получить актеров и фильмы
	actors, err := actorRepo.GetActorsWithMovies(7, 1)
	if err != nil {
		fmt.Println(err)
	}
	for _, actor := range actors {
		fmt.Println("actor:", actor.Name, "id", actor.ID)
	}
	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal("Error running application:", err)
	}
}
