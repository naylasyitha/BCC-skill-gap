package middleware

import (
	"net/http"
	"os"
	"project-bcc/internal/usecase"
	"project-bcc/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authRepo usecase.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Sesi tidak ditemukan, silakan login terlebih dahulu",
			})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.ValidateToken(token, os.Getenv("JWT_SECRET"))
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "expired") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "Sesi Anda telah berakhir, silakan login kembali",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token kredensial tidak valid",
			})
			return
		}

		user, err := authRepo.FindByID(c.Request.Context(), claims.UserID)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Pengguna tidak ditemukan",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak, hanya admin yang diizinkan",
			})
			return
		}
		c.Next()
	}
}
