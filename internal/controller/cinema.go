package controller

type cinema struct {
	movie serviceMovie
	actor serviceActor
}

func NewCinema(movie serviceMovie, actor serviceActor) cinema {
	return cinema{movie: movie, actor: actor}
}
