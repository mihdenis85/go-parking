package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/amend-parking-backend/internal/config"
)

const XAPIKeyHeader = "X-API-Key"

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(XAPIKeyHeader)
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"detail": "Invalid API Key. Check 'X-API-Key' header.",
			})
			c.Abort()
			return
		}

		if apiKey != config.Settings.ParkingServiceAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"detail": "Invalid API Key. Check 'X-API-Key' header.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
