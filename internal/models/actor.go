package models

import (
	"time"

	"github.com/google/uuid"
)

type Actor struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Gender      string    `json:"gender"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Movies      []Movie   `json:"movies"`
}

type CreateActor struct {
	Name        string    `json:"name" binding:"required"`
	Gender      string    `json:"gender" binding:"required"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`
}

type UpdateActor struct {
	Name        *string    `json:"name"`
	Gender      *string    `json:"gender"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}
