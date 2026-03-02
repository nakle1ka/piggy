package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nakle1ka/piggy/internal/middleware"
)

func getUserId(c *gin.Context) (uuid.UUID, bool) {
	key, exists := c.Get(middleware.UserIdKey)
	if !exists {
		return uuid.Nil, false
	}

	return key.(uuid.UUID), true
}

func getPiggyId(c *gin.Context) (uuid.UUID, bool) {
	id := c.Param("id")
	piggyId, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, false
	}
	return piggyId, true
}
