package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

var SecretKey = []byte("your-secret-key") // Секретный ключ для проверки подписи токена

const RoleKey = "role"

// JWTAuthMiddleware для проверки токена и установки роли в контекст
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем токен из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Проверяем, начинается ли заголовок с "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		// Парсим и проверяем JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверка подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Извлекаем роль из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Role not found in token", http.StatusUnauthorized)
			return
		}

		// Устанавливаем роль в контексте
		ctx := context.WithValue(r.Context(), RoleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
