package main

import (
	"log"
	"os"

	"shorter-backend/config"
	"shorter-backend/handlers"
	"shorter-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connections
	config.ConnectDatabase()
	config.ConnectRedis()

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",  // Next.js frontend
		"http://localhost:3001",  // Alternative port
		"https://your-domain.com", // Production domain
	}
	
	// Allow all origins in development
	if os.Getenv("GIN_MODE") != "release" {
		corsConfig.AllowAllOrigins = true
	}
	
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "URL Shortener API is running",
		})
	})

	// Advanced health check
	r.GET("/health/detailed", handlers.GetDetailedHealth)

	// Admin routes (consider adding authentication in production)
	admin := r.Group("/admin")
	{
		admin.GET("/stats", handlers.GetSystemStats)
		admin.GET("/activity", handlers.GetRecentActivity)
	}

	// API routes
	api := r.Group("/api")
	api.Use(middleware.RateLimitMiddleware(middleware.GeneralLimiter))
	{
		// URL shortening (with stricter rate limit)
		api.POST("/shorten", middleware.RateLimitMiddleware(middleware.CreateURLLimiter), handlers.ShortenURL)
		
		// Get all URLs with pagination
		api.GET("/urls", handlers.GetAllURLs)
		
		// Get statistics for a specific short URL
		api.GET("/stats/:code", handlers.GetURLStats)
		
		// QR Code generation (with specific rate limit)
		qr := api.Group("/qr", middleware.RateLimitMiddleware(middleware.QRCodeLimiter))
		{
			qr.GET("/:code/image", handlers.GenerateQRCode)
			qr.GET("/:code", handlers.GetQRCodeHTML)
		}
	}

	// Redirect routes (without /api prefix for clean short URLs)
	r.GET("/:code", handlers.RedirectURL)

	// Static file serving for frontend (if built)
	r.Static("/static", "./static")

	// Catch-all for frontend routing (SPA)
	r.NoRoute(func(c *gin.Context) {
		// Serve index.html for client-side routing
		c.File("./static/index.html")
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
} 