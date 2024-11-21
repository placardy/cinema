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
}

type ActorWithMovies struct {
	Actor
	Movies []Movie `json:"movies"`
}

type CreateActor struct {
	Name        string    `json:"name"`
	Gender      string    `json:"gender"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

type UpdateActor struct {
	Name        *string    `json:"name"`
	Gender      *string    `json:"gender"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}

// Validate для CreateActor
func (ca CreateActor) Validate() []ValidationError {
	var errs []ValidationError

	// Name min=1 max=100
	if len(ca.Name) < 1 || len(ca.Name) > 100 {
		errs = append(errs, ValidationError{
			Field:   "name",
			Message: "Name must be between 1 and 100 characters",
		})
	}

	// Gender must be "male", "female", or "other"
	validGenders := map[string]bool{"male": true, "female": true, "other": true}
	if !validGenders[ca.Gender] {
		errs = append(errs, ValidationError{
			Field:   "gender",
			Message: "Gender must be 'male', 'female', or 'other'",
		})
	}

	// DateOfBirth must not be in the future
	if ca.DateOfBirth.IsZero() {
		errs = append(errs, ValidationError{
			Field:   "date_of_birth",
			Message: "Date of birth is required",
		})
	} else if ca.DateOfBirth.After(time.Now()) {
		errs = append(errs, ValidationError{
			Field:   "date_of_birth",
			Message: "Date of birth cannot be in the future",
		})
	}

	return errs
}

// Validate для UpdateActor
func (ua UpdateActor) Validate() []ValidationError {
	var errs []ValidationError

	// Name min=1 max=100
	if ua.Name != nil {
		if len(*ua.Name) < 1 || len(*ua.Name) > 100 {
			errs = append(errs, ValidationError{
				Field:   "name",
				Message: "Name must be between 1 and 100 characters",
			})
		}
	}

	// Gender must be "male", "female", or "other"
	if ua.Gender != nil {
		validGenders := map[string]bool{"male": true, "female": true, "other": true}
		if !validGenders[*ua.Gender] {
			errs = append(errs, ValidationError{
				Field:   "gender",
				Message: "Gender must be 'male', 'female', or 'other'",
			})
		}
	}

	// DateOfBirth must not be in the future
	if ua.DateOfBirth != nil {
		if ua.DateOfBirth.After(time.Now()) {
			errs = append(errs, ValidationError{
				Field:   "date_of_birth",
				Message: "Date of birth cannot be in the future",
			})
		}
	}

	return errs
}
