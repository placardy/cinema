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
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float64   `json:"rating"`
	ActorIDs    []uuid.UUID `json:"actor_ids"`
}

type MovieWithActors struct {
	Movie
	Actors []Actor `json:"actors"`
}

type UpdateMovie struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	ReleaseDate *time.Time `json:"release_date"`
	Rating      *float64   `json:"rating"`
	ActorIDs    *[]uuid.UUID `json:"actor_ids"`
}

// Validation struct
func (cm CreateMovie) Validate() []ValidationError {
	var errs []ValidationError

	// Title min=1 max=150
	if len(cm.Title) < 1 || len(cm.Title) > 150 {
		errs = append(errs, ValidationError{
			Field:   "title",
			Message: "Title must be between 1 and 150 characters",
		})
	}

	// Description max=1000
	if len(cm.Description) > 1000 {
		errs = append(errs, ValidationError{
			Field:   "description",
			Message: "Description must not exceed 1000 characters",
		})
	}

	// ReleaseDate required and ReleaseDate <= now
	if cm.ReleaseDate.IsZero() {
		errs = append(errs, ValidationError{
			Field:   "release_date",
			Message: "Release date is required",
		})
	} else if cm.ReleaseDate.After(time.Now()) {
		errs = append(errs, ValidationError{
			Field:   "release_date",
			Message: "Release date cannot be in the future",
		})
	}

	// Ratin min=0 max=10
	if cm.Rating < 0 || cm.Rating > 10 {
		errs = append(errs, ValidationError{
			Field:   "rating",
			Message: "Rating must be between 0 and 10",
		})
	}

	// ActorIDs max_count = 100
	if len(cm.ActorIDs) > 100 {
		errs = append(errs, ValidationError{
			Field:   "actor_ids",
			Message: "Too many actors in the list, maximum allowed is 100",
		})
	}

	return errs
}

func (um UpdateMovie) Validate() []ValidationError {
	var errs []ValidationError

	// Title min=1 max=150
	if um.Title != nil {
		if len(*um.Title) < 1 || len(*um.Title) > 150 {
			errs = append(errs, ValidationError{
				Field:   "title",
				Message: "Title must be between 1 and 150 characters",
			})
		}
	}

	// Description max=1000
	if um.Description != nil {
		if len(*um.Description) > 1000 {
			errs = append(errs, ValidationError{
				Field:   "description",
				Message: "Description must not exceed 1000 characters",
			})
		}
	}

	// ReleaseDate <= now
	if um.ReleaseDate != nil {
		if um.ReleaseDate.After(time.Now()) {
			errs = append(errs, ValidationError{
				Field:   "release_date",
				Message: "Release date cannot be in the future",
			})
		}
	}

	// Ratin min=0 max=10
	if um.Rating != nil {
		if *um.Rating < 0 || *um.Rating > 10 {
			errs = append(errs, ValidationError{
				Field:   "rating",
				Message: "Rating must be between 0 and 10",
			})
		}
	}

	// ActorIDs max_count = 100
	// if len(*um.ActorIDs) > 100 {
	// 	errs = append(errs, ValidationError{
	// 		Field:   "actor_ids",
	// 		Message: "Too many actors in the list, maximum allowed is 100",
	// 	})
	// }

	return errs
}
