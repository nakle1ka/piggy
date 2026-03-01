package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/piggy/internal/pkg/auth"
)

const UserIdKey = "user-id"

func JWTAuth(tm auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		accessToken := strings.TrimPrefix(header, "Bearer ")

		claims, err := tm.ValidateToken(accessToken)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(UserIdKey, claims.Subject)
		c.Next()
	}
}
