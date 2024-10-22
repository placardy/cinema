package models

import (
	"github.com/google/uuid"
)

type Actor struct {
	ID        uuid.UUID
	Name      string
	Gender    string
	BirthDate string
}
