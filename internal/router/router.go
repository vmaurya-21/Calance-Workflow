package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vmaurya-21/Calance-Workflow/internal/config"
	"gorm.io/gorm"

	// New modular handlers
	authHandler "github.com/vmaurya-21/Calance-Workflow/internal/api/handlers/auth"
	orgHandler "github.com/vmaurya-21/Calance-Workflow/internal/api/handlers/organization"
	repoHandler "github.com/vmaurya-21/Calance-Workflow/internal/api/handlers/repository"
	workflowHandler "github.com/vmaurya-21/Calance-Workflow/internal/api/handlers/workflow"

	// Domain services
	authDomain "github.com/vmaurya-21/Calance-Workflow/internal/domain/auth"
	orgDomain "github.com/vmaurya-21/Calance-Workflow/internal/domain/organization"
	repoDomain "github.com/vmaurya-21/Calance-Workflow/internal/domain/repository"
	workflowDomain "github.com/vmaurya-21/Calance-Workflow/internal/domain/workflow"

	// Infrastructure
	database "github.com/vmaurya-21/Calance-Workflow/internal/infrastructure/database/repositories"

	// Utilities
	"github.com/vmaurya-21/Calance-Workflow/internal/logger"
	"github.com/vmaurya-21/Calance-Workflow/internal/middleware"
)

// SetupRouter configures all routes for the application
func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Create router
	r := gin.Default()

	// Apply logging middleware
	r.Use(logger.GinMiddleware())

	// Apply CORS middleware
	r.Use(middleware.CORSMiddleware(cfg.Frontend.AllowedOrigins))

	// Initialize repositories
	userRepo := database.NewUserRepository(db)
	tokenRepo := database.NewTokenRepository(db)

	// Initialize domain services
	scopes := []string{"user:email", "read:user", "read:org", "repo", "workflow", "read:packages"}
	authService := authDomain.NewService(cfg.GitHub.ClientID, cfg.GitHub.ClientSecret, cfg.GitHub.RedirectURL, scopes)
	workflowService := workflowDomain.NewService()
	repositoryService := repoDomain.NewService()
	organizationService := orgDomain.NewService()

	// Initialize handlers
	authHandlers := authHandler.NewHandler(authService, userRepo, tokenRepo)
	organizationHandlers := orgHandler.NewHandler(organizationService, tokenRepo)
	repositoryHandlers := repoHandler.NewHandler(repositoryService, tokenRepo)
	workflowHandlers := workflowHandler.NewHandler(workflowService, tokenRepo)

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
			auth.GET("/github", authHandlers.Login)
			auth.GET("/github/callback", authHandlers.Callback)

			// Protected auth routes
			auth.GET("/me", middleware.AuthMiddleware(), authHandlers.GetProfile)
			auth.POST("/logout", middleware.AuthMiddleware(), authHandlers.Logout)

			// Organization endpoints (requires valid JWT)
			auth.GET("/organizations", middleware.AuthMiddleware(), organizationHandlers.List)
			auth.GET("/organizations/:org/repositories", middleware.AuthMiddleware(), organizationHandlers.GetRepositories)

			// Repository endpoints (requires valid JWT)
			auth.GET("/repositories", middleware.AuthMiddleware(), organizationHandlers.GetUserRepositories)
			auth.GET("/repositories/:owner/:repo/branches", middleware.AuthMiddleware(), repositoryHandlers.GetBranches)
			auth.GET("/repositories/:owner/:repo/branches/:branch/commits", middleware.AuthMiddleware(), repositoryHandlers.GetCommits)
		}

		// Organization routes (protected)
		organizations := api.Group("/organizations")
		organizations.Use(middleware.AuthMiddleware())
		{
			organizations.GET("", organizationHandlers.List)
			organizations.GET("/:org/repositories", organizationHandlers.GetRepositories)
		}

		// Repository routes (protected)
		repositories := api.Group("/repositories")
		repositories.Use(middleware.AuthMiddleware())
		{
			repositories.GET("", organizationHandlers.GetUserRepositories)
			repositories.GET("/:owner/:repo/branches", repositoryHandlers.GetBranches)
			repositories.GET("/:owner/:repo/branches/:branch/commits", repositoryHandlers.GetCommits)
			repositories.GET("/:owner/:repo/tags", repositoryHandlers.GetTags)
			repositories.POST("/tags", repositoryHandlers.CreateTag)

			// GitHub Actions workflow runs endpoints
			repositories.GET("/:owner/:repo/actions/runs", repositoryHandlers.GetWorkflowRuns)
			repositories.GET("/:owner/:repo/actions/runs/:run_id", repositoryHandlers.GetWorkflowRunDetail)

			// GitHub Actions job logs endpoint (returns all logs for a job)
			repositories.GET("/:owner/:repo/actions/jobs/:job_id/logs", repositoryHandlers.GetJobLogs)
		}

		// GitHub Packages routes (protected)
		packages := api.Group("/packages")
		packages.Use(middleware.AuthMiddleware())
		{
			packages.GET("/user", repositoryHandlers.GetUserPackages)
			packages.GET("/org/:org", repositoryHandlers.GetOrgPackages)
		}

		// Workflow routes (protected)
		workflows := api.Group("/workflows")
		workflows.Use(middleware.AuthMiddleware())
		{
			workflows.GET("/:owner/:repo", workflowHandlers.List)
			workflows.POST("/create", workflowHandlers.Create)
			workflows.POST("/preview", workflowHandlers.Preview)

			// Workflow edit endpoints
			workflows.GET("/:owner/:repo/file", workflowHandlers.GetWorkflowContent)
			workflows.PUT("/:owner/:repo/file", workflowHandlers.UpdateWorkflow)
		}
	}

	return r
}
