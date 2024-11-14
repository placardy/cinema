package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware проверяет, имеет ли пользователь одну из разрешенных ролей
func RoleMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(RoleKey)
		if !exists || role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found"})
			c.Abort()
			return
		}

		// Проверка роли пользователя
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		// Если роль не подходит, возвращаем ошибку
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		c.Abort()
	}
}
