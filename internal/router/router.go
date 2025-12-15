package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vmaurya-21/Calance-Workflow/internal/config"
	"github.com/vmaurya-21/Calance-Workflow/internal/controllers"
	"github.com/vmaurya-21/Calance-Workflow/internal/middleware"
	"github.com/vmaurya-21/Calance-Workflow/internal/repositories"
	"github.com/vmaurya-21/Calance-Workflow/internal/services"
	"gorm.io/gorm"
)

// SetupRouter configures all routes for the application
func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Create router
	r := gin.Default()

	// Apply CORS middleware
	r.Use(middleware.CORSMiddleware(cfg.Frontend.AllowedOrigins))

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	tokenRepo := repositories.NewTokenRepository(db)

	// Initialize services
	githubOAuthService := services.NewGitHubOAuthService(cfg)
	githubRepositoryService := services.NewGitHubRepositoryService()
	githubOrganizationService := services.NewGitHubOrganizationService(githubRepositoryService)

	// Initialize controllers
	authController := controllers.NewAuthController(githubOAuthService, userRepo, tokenRepo, cfg)
	organizationController := controllers.NewOrganizationController(githubOrganizationService, tokenRepo)
	repositoryController := controllers.NewRepositoryController(githubRepositoryService, tokenRepo, userRepo)

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
		// Auth routes
		auth := api.Group("/auth")
		{
			// Public auth routes
			auth.GET("/github", authController.GitHubLogin)
			auth.GET("/github/callback", authController.GitHubCallback)

			// Protected auth routes
			auth.GET("/me", middleware.AuthMiddleware(), authController.GetCurrentUser)
			auth.POST("/logout", middleware.AuthMiddleware(), authController.Logout)

			// Organization endpoints (requires valid JWT)
			auth.GET("/organizations", middleware.AuthMiddleware(), organizationController.GetUserOrganizations)
			auth.GET("/organizations/:org/repositories", middleware.AuthMiddleware(), organizationController.GetOrganizationRepositories)

			// Repository endpoints (requires valid JWT)
			auth.GET("/repositories", middleware.AuthMiddleware(), organizationController.GetUserRepositories)
			auth.GET("/repositories/:owner/:repo/branches", middleware.AuthMiddleware(), repositoryController.GetRepositoryBranches)
			auth.GET("/repositories/:owner/:repo/branches/:branch/commits", middleware.AuthMiddleware(), repositoryController.GetBranchCommits)
		}

		// Organization routes (protected)
		organizations := api.Group("/organizations")
		organizations.Use(middleware.AuthMiddleware())
		{
			organizations.GET("", organizationController.GetUserOrganizations)
			organizations.GET("/:org/repositories", organizationController.GetOrganizationRepositories)
		}

		// Repository routes (protected)
		repositories := api.Group("/repositories")
		repositories.Use(middleware.AuthMiddleware())
		{
			repositories.GET("", organizationController.GetUserRepositories)
			repositories.GET("/:owner/:repo/branches", repositoryController.GetRepositoryBranches)
			repositories.GET("/:owner/:repo/branches/:branch/commits", repositoryController.GetBranchCommits)
			repositories.GET("/:owner/:repo/tags", repositoryController.GetRepositoryTags)
			repositories.POST("/tags", repositoryController.CreateTag)
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
