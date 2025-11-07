package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prajapatiomkar/wave-server/config"
	"github.com/prajapatiomkar/wave-server/internal/handlers"
	"github.com/prajapatiomkar/wave-server/internal/middleware"
	"github.com/prajapatiomkar/wave-server/internal/repositories"
	"github.com/prajapatiomkar/wave-server/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Connect to database
	config.ConnectDatabase()

	// Set Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowWebSockets:  true,
	}
	router.Use(cors.New(corsConfig))

	// Initialize repositories
	userRepo := repositories.NewUserRepository(config.GetDB())

	// Initialize services
	authService := services.NewAuthService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Chat API is running"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", authHandler.GetMe)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
