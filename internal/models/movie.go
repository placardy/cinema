package models

import (
	"fmt"
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

// Валидация данных фильма
func (m *CreateMovie) Validate() error {
	if m.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if len(m.Title) > 150 {
		return fmt.Errorf("title must be less than 150 characters")
	}
	if m.Description != "" && len(m.Description) > 1000 {
		return fmt.Errorf("description must be less than 1000 characters")
	}
	if m.ReleaseDate.IsZero() {
		return fmt.Errorf("release date is required")
	}
	if m.Rating < 0 || m.Rating > 10 {
		return fmt.Errorf("rating must be between 0 and 10")
	}
	return nil
}
