package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vmaurya-21/Calance-Workflow/internal/config"
	"github.com/vmaurya-21/Calance-Workflow/internal/controllers"
	"github.com/vmaurya-21/Calance-Workflow/internal/middleware"

	// "github.com/vmaurya-21/Calance-Workflow/internal/repositories"
	"github.com/vmaurya-21/Calance-Workflow/internal/services"
	// "gorm.io/gorm"
)

// SetupRouter configures all routes for the application
func SetupRouter(db interface{}, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Create router
	r := gin.Default()

	// Apply CORS middleware
	r.Use(middleware.CORSMiddleware(cfg.Frontend.AllowedOrigins))

	// TODO: Initialize repositories
	// userRepo := repositories.NewUserRepository(db)

	// Initialize services
	githubOAuthService := services.NewGitHubOAuthService(cfg)

	// Initialize controllers
	authController := controllers.NewAuthController(githubOAuthService, nil, cfg)

	// Health check route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "healthy",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.GET("/github", authController.GitHubLogin)
			auth.GET("/github/callback", authController.GitHubCallback)

			// TODO: Protected auth routes (database required)
			// auth.GET("/me", middleware.AuthMiddleware(), authController.GetCurrentUser)
			// auth.POST("/logout", middleware.AuthMiddleware(), authController.Logout)

			// Organizations endpoint (requires valid JWT)
			auth.GET("/organizations", authController.GetUserOrganizations)

			// Repository endpoints (requires valid JWT)
			auth.GET("/repositories", authController.GetUserRepositories)
			auth.GET("/organizations/:org/repositories", authController.GetOrganizationRepositories)
		}

		// Protected API routes example
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Add your protected routes here
			// Example:
			// protected.GET("/users", userController.List)
			// protected.GET("/profile", userController.GetProfile)
		}
	}

	return r
}
