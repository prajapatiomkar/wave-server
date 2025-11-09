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
	"github.com/prajapatiomkar/wave-server/internal/websocket"
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

	// Initialize repositories
	userRepo := repositories.NewUserRepository(config.GetDB())
	messageRepo := repositories.NewMessageRepository(config.GetDB())

	// Initialize services
	authService := services.NewAuthService(userRepo)
	messageService := services.NewMessageService(messageRepo, userRepo)

	// Initialize WebSocket hub
	hub := websocket.NewHub(messageService)
	go hub.Run()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	wsHandler := handlers.NewWebSocketHandler(hub)
	messageHandler := handlers.NewMessageHandler(messageService)

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowWebSockets:  true,
	}
	router.Use(cors.New(corsConfig))

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

		// WebSocket route (handles auth internally)
		api.GET("/ws", wsHandler.HandleConnection)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", authHandler.GetMe)
			protected.GET("/messages/:room_id", messageHandler.GetMessageHistory)
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
