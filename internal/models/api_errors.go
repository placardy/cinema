package models

type APIError struct {
	Code    string      `json:"code"`              // Код ошибки (например, VALIDATION_ERROR, INVALID_JSON, etc.)
	Message string      `json:"message"`           // Общее описание ошибки
	Details interface{} `json:"details,omitempty"` // Дополнительные данные об ошибке (может быть nil или содержать ValidationError, массив и т.д.)
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
