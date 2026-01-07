package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware - Check if user has specific role
func AuthMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user_id and role from header or query
		userRole := c.GetHeader("X-User-Role")
		
		if userRole == "" {
			userRole = c.Query("user_role")
		}

		if userRole == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"result": "User role tidak ditemukan",
				"status": 0,
			})
			c.Abort()
			return
		}

		// Check if role is allowed
		isAllowed := false
		for _, role := range requiredRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{
				"result": "Anda tidak memiliki akses untuk resource ini",
				"status": 0,
			})
			c.Abort()
			return
		}

		// Store role in context
		c.Set("user_role", userRole)
		c.Next()
	}
}

// AdminOnly - Middleware untuk admin only
func AdminOnly() gin.HandlerFunc {
	return AuthMiddleware("admin")
}

// SellerOnly - Middleware untuk seller only
func SellerOnly() gin.HandlerFunc {
	return AuthMiddleware("admin", "seller")
}

// CustomerOnly - Middleware untuk customer only
func CustomerOnly() gin.HandlerFunc {
	return AuthMiddleware("customer", "admin")
}

// GetUserRoleFromContext - Helper function to get user role from context
func GetUserRoleFromContext(c *gin.Context) string {
	role, exists := c.Get("user_role")
	if !exists {
		return ""
	}
	return role.(string)
}
