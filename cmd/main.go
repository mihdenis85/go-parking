package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amend-parking-backend/internal/api"
	"github.com/amend-parking-backend/internal/config"
	"github.com/amend-parking-backend/internal/database"
	"github.com/amend-parking-backend/internal/repository"
	"github.com/amend-parking-backend/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/amend-parking-backend/docs"
)

// @title           Parking Service API
// @version         1.0
// @description     API для управления парковочными местами
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API Key для аутентификации

// @host      localhost:8000
// @BasePath  /

func main() {
	config.LoadConfig()

	if err := database.InitializeDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	log.Println("Application startup")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	repo := repository.NewRepository()
	svc := service.NewService(repo)

	api.SetupRoutes(router, svc)

	srv := &http.Server{
		Addr:    ":" + config.Settings.ServerPort,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %s", config.Settings.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Application shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
