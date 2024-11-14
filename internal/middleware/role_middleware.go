package middleware

import (
	"net/http"
)

// RoleMiddleware проверяет, имеет ли пользователь одну из разрешенных ролей
func RoleMiddleware(allowedRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(RoleKey).(string)
			if !ok || role == "" {
				http.Error(w, "Role not found", http.StatusUnauthorized)
				return
			}

			// Проверка роли пользователя
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					// Если роль разрешена, передаем управление следующему обработчику
					next.ServeHTTP(w, r)
					return
				}
			}

			// Если роль не подходит, возвращаем ошибку
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
