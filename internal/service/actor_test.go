package service

import (
	"cinema/internal/models"
	"cinema/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateActor(t *testing.T) {
	// Создание контроллера для мока
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание мока
	mockStore := mocks.NewMockstoreActor(ctrl)

	// Данные для теста
	actor := models.CreateActor{
		Name:        "Leonardo DiCaprio",
		Gender:      "Male",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	expectedID := uuid.New()

	// Настройка мока
	mockStore.EXPECT().CreateActor(actor).Return(expectedID, nil)

	// Создаём сервис с использованием мока
	actorService := NewActor(mockStore)

	// Вызов тестируемого метода
	resultID, err := actorService.CreateActor(actor)

	// Проверки
	assert.NoError(t, err)
	assert.Equal(t, expectedID, resultID)
}
