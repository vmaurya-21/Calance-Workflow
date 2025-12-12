package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/repositories"
	"github.com/vmaurya-21/Calance-Workflow/internal/services"
	"github.com/vmaurya-21/Calance-Workflow/internal/utils"
)

type RepositoryController struct {
	repositoryService *services.GitHubRepositoryService
	tokenRepository   *repositories.TokenRepository
}

// NewRepositoryController creates a new repository controller
func NewRepositoryController(
	repositoryService *services.GitHubRepositoryService,
	tokenRepository *repositories.TokenRepository,
) *RepositoryController {
	return &RepositoryController{
		repositoryService: repositoryService,
		tokenRepository:   tokenRepository,
	}
}

// GetRepositoryBranches returns all branches for a specific repository
// GET /api/repositories/:owner/:repo/branches
func (rc *RepositoryController) GetRepositoryBranches(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not found in context")
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Invalid user ID format", err)
		return
	}

	// Get owner and repo from URL parameters
	owner := c.Param("owner")
	repo := c.Param("repo")

	if owner == "" || repo == "" {
		utils.BadRequestResponse(c, "Owner and repository name are required")
		return
	}

	// Fetch token from database
	token, err := rc.tokenRepository.FindByUserID(userUUID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch access token", err)
		return
	}
	if token == nil {
		utils.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	log.Printf("DEBUG - Fetching branches for repository: %s/%s", owner, repo)

	// Get branches for the repository
	branches, err := rc.repositoryService.GetRepositoryBranches(c.Request.Context(), token.AccessToken, owner, repo)
	if err != nil {
		log.Printf("ERROR - Failed to fetch branches: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to fetch branches", err)
		return
	}

	log.Printf("DEBUG - Successfully fetched %d branches for repository %s/%s", len(branches), owner, repo)

	utils.SuccessResponse(c, http.StatusOK, "Branches fetched successfully", gin.H{
		"owner":        owner,
		"repository":   repo,
		"branches":     branches,
		"branch_count": len(branches),
	})
}

// GetBranchCommits returns the latest commits for a specific branch
// GET /api/repositories/:owner/:repo/branches/:branch/commits
func (rc *RepositoryController) GetBranchCommits(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not found in context")
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Invalid user ID format", err)
		return
	}

	// Get owner, repo, and branch from URL parameters
	owner := c.Param("owner")
	repo := c.Param("repo")
	branch := c.Param("branch")

	if owner == "" || repo == "" || branch == "" {
		utils.BadRequestResponse(c, "Owner, repository name, and branch are required")
		return
	}

	// Get per_page query parameter, default to 30
	perPage := 30
	if perPageStr := c.Query("per_page"); perPageStr != "" {
		var parsed int
		if _, err := fmt.Sscanf(perPageStr, "%d", &parsed); err == nil {
			if parsed > 0 && parsed <= 100 {
				perPage = parsed
			}
		}
	}

	// Fetch token from database
	token, err := rc.tokenRepository.FindByUserID(userUUID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch access token", err)
		return
	}
	if token == nil {
		utils.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	log.Printf("DEBUG - Fetching %d commits for branch %s in repository: %s/%s", perPage, branch, owner, repo)

	// Get commits for the branch
	commits, err := rc.repositoryService.GetBranchCommits(c.Request.Context(), token.AccessToken, owner, repo, branch, perPage)
	if err != nil {
		log.Printf("ERROR - Failed to fetch commits: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to fetch commits", err)
		return
	}

	log.Printf("DEBUG - Successfully fetched %d commits for branch %s", len(commits), branch)

	utils.SuccessResponse(c, http.StatusOK, "Commits fetched successfully", gin.H{
		"owner":        owner,
		"repository":   repo,
		"branch":       branch,
		"commits":      commits,
		"commit_count": len(commits),
	})
}
