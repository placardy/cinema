package repository

import (
	"cinema/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBeginTransaction(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMovie(db)

	mock.ExpectBegin()

	tx, err := repo.BeginTransaction()
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckMovieExists(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMovie(db)

	movieID := uuid.New()

	mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM movies WHERE id = \$1\)`).
		WithArgs(movieID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true)) // Фильм существует

	exists, err := repo.CheckMovieExists(movieID)
	assert.NoError(t, err)
	assert.True(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckActorsExist(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMovie(db)

	actorIDs := []uuid.UUID{uuid.New(), uuid.New()}

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM actors WHERE id IN \(\$1,\$2\)`).
		WithArgs(actorIDs[0], actorIDs[1]).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	exists, err := repo.CheckActorsExist(actorIDs)
	assert.NoError(t, err)
	assert.True(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestCreateMovie(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMovie(db)

	movie := models.CreateMovie{
		Title:       "Inception",
		Description: "A mind-bending thriller.",
		ReleaseDate: time.Date(2010, 7, 16, 0, 0, 0, 0, time.UTC),
		Rating:      8.8,
		ActorIDs:    []uuid.UUID{},
	}

	movieID := uuid.New()

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO movies`).
		WithArgs(sqlmock.AnyArg(), movie.Title, movie.Description, movie.ReleaseDate, movie.Rating).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(movieID))

	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.NoError(t, err)

	resultID, err := repo.CreateMovie(tx, movie)

	assert.NoError(t, err)
	assert.Equal(t, movieID, resultID)

	err = tx.Commit()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
