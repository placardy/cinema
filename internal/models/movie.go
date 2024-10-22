package models

import (
	"github.com/google/uuid"
)

type Movie struct {
	ID          uuid.UUID
	Title       string
	Description string
	ReleaseDate string
	Rating      float64
}
