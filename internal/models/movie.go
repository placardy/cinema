package models

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float64   `json:"rating"`
}

type CreateMovie struct {
	Title       string    `json:"title" binding:"required,min=1,max=150"`
	Description string    `json:"description" binding:"max=1000"`
	ReleaseDate time.Time `json:"release_date" binding:"required"`
	Rating      float64   `json:"rating" binding:"required,min=0,max=10"`
}

type UpdateMovie struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	ReleaseDate *time.Time `json:"release_date"`
	Rating      *float64   `json:"rating"`
}
