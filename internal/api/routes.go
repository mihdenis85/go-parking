package api

import (
	"github.com/amend-parking-backend/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, svc *service.Service) {
	handlers := NewHandlers(svc)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/docs/index.html")
	})

	parking := router.Group("/parking")
	parking.Use(APIKeyAuth())
	{
		parking.GET("/free-spaces-count", handlers.GetCountOfFreeSpaces)
		parking.GET("/occupied-spaces-list", handlers.GetOccupiedSpaces)
		parking.POST("/park-car", handlers.ParkCar)
		parking.POST("/free-up", handlers.FreeUpParkingSpace)
		parking.GET("/parking-space-logs", handlers.GetParkingSpaceLogs)
	}
}
