package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
)

// getUserIDFromContext extracts the user ID from the gin context
func getUserIDFromContext(c *gin.Context) int64 {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	return userID.(int64)
}

// parseTime parses a time string in RFC3339 format
func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}
