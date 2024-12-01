package utils

import (
	"cinema/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Метод для создания ошибки с валидацией
func ValidationErrorResponse(ctx *gin.Context, validationErrors []models.ValidationError) {
	ctx.JSON(http.StatusBadRequest, models.APIError{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed for one or more fields",
		Details: validationErrors,
	})
}

// Метод для создания ошибки с некорректным форматом запроса
func InvalidJSONResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, models.APIError{
		Code:    "INVALID_JSON",
		Message: "Invalid JSON format",
		Details: nil,
	})
}

// Метод для создания ошибки при внутренних ошибках
func InternalServerErrorResponse(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, models.APIError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: message,
		Details: nil,
	})
}

// Метод для ошибки 404 - не найдено
func NotFoundResponse(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, models.APIError{
		Code:    "NOT_FOUND",
		Message: message,
		Details: nil,
	})
}

// Метод для ошибки 400 - неверный ID (например, UUID)
func BadRequestResponse(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, models.APIError{
		Code:    "BAD_REQUEST",
		Message: message,
		Details: nil,
	})
}
