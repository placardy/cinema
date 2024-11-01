package controller

type movie struct {
	service serviceMovie
}

func NewMovie(service serviceMovie) *movie {
	return &movie{service: service}
}
