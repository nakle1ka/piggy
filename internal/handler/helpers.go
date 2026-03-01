package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nakle1ka/piggy/internal/middleware"
)

func getUserId(c *gin.Context) (uuid.UUID, bool) {
	key, exists := c.Get(middleware.UserIdKey)
	if !exists {
		c.Status(http.StatusUnauthorized)
		return uuid.Nil, false
	}

	userId, err := uuid.Parse(key.(string))
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return uuid.Nil, false
	}

	return userId, true
}

func getPiggyId(c *gin.Context) (uuid.UUID, bool) {
	id := c.Param("id")
	piggyId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid piggy id",
		})
		return uuid.Nil, false
	}
	return piggyId, true
}
