package repository

import (
	"cinema/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateActor(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActor(db)

	// Test input
	actor := models.CreateActor{
		Name:        "John Doe",
		Gender:      "Male",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Mock expectation
	expectedID := uuid.New() // Mock the ID that will be returned by the query
	mock.ExpectQuery(`INSERT INTO actors`).
		WithArgs(sqlmock.AnyArg(), actor.Name, actor.Gender, actor.DateOfBirth). // Use AnyArg for the UUID
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	// Execute
	resultID, err := repo.CreateActor(actor)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedID, resultID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetActor(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActor(db)

	actorID := uuid.New()
	actor := models.Actor{
		ID:          actorID,
		Name:        "John Doe",
		Gender:      "Male",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	mock.ExpectQuery(`SELECT id, name, gender, date_of_birth FROM actors`).
		WithArgs(actorID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "gender", "date_of_birth"}).
			AddRow(actor.ID, actor.Name, actor.Gender, actor.DateOfBirth))

	// Execute
	result, err := repo.GetActor(actorID)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, actor.ID, result["id"])
	assert.Equal(t, actor.Name, result["name"])
	assert.Equal(t, actor.Gender, result["gender"])
	assert.Equal(t, actor.DateOfBirth, result["date_of_birth"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteActor(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActor(db)

	actorID := uuid.New()

	mock.ExpectExec(`DELETE FROM actors`).
		WithArgs(actorID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute
	err = repo.DeleteActor(actorID)

	// Verify
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateActor(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActor(db)

	actorID := uuid.New()
	name := "John Doe Updated"
	gender := "Male"
	dateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	updateData := models.UpdateActor{
		Name:        &name,
		Gender:      &gender,
		DateOfBirth: &dateOfBirth,
	}

	mock.ExpectExec(`UPDATE actors`).
		WithArgs(updateData.Name, updateData.Gender, updateData.DateOfBirth, actorID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute
	err = repo.UpdateActor(actorID, updateData)

	// Verify
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllActors(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActor(db)

	actors := []models.Actor{
		{
			ID:          uuid.New(),
			Name:        "Actor One",
			Gender:      "Male",
			DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          uuid.New(),
			Name:        "Actor Two",
			Gender:      "Female",
			DateOfBirth: time.Date(1985, 5, 5, 0, 0, 0, 0, time.UTC),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "gender", "date_of_birth"}).
		AddRow(actors[0].ID, actors[0].Name, actors[0].Gender, actors[0].DateOfBirth).
		AddRow(actors[1].ID, actors[1].Name, actors[1].Gender, actors[1].DateOfBirth)

	mock.ExpectQuery(`SELECT id, name, gender, date_of_birth FROM actors LIMIT .* OFFSET .*`).
		WillReturnRows(rows)

	// Execute
	result, err := repo.GetAllActors(2, 0)

	// Verify
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, actors[0].Name, result[0]["name"])
	assert.Equal(t, actors[1].Name, result[1]["name"])
	assert.NoError(t, mock.ExpectationsWereMet())
}
