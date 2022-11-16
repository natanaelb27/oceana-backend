package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GinGetUUID(c *gin.Context, key string) (b uuid.UUID) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(uuid.UUID)
	}
	return
}
