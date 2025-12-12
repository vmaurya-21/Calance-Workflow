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
	userRepository    *repositories.UserRepository
}

// NewRepositoryController creates a new repository controller
func NewRepositoryController(
	repositoryService *services.GitHubRepositoryService,
	tokenRepository *repositories.TokenRepository,
	userRepository *repositories.UserRepository,
) *RepositoryController {
	return &RepositoryController{
		repositoryService: repositoryService,
		tokenRepository:   tokenRepository,
		userRepository:    userRepository,
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

// CreateTag creates and pushes a tag for a specific commit
// POST /api/repositories/tags
func (rc *RepositoryController) CreateTag(c *gin.Context) {
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

	// Parse request body
	var tagRequest services.CreateTagRequest
	if err := c.ShouldBindJSON(&tagRequest); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	// Determine the owner (use provided owner or fetch from authenticated user)
	var owner string
	if tagRequest.Owner != "" {
		owner = tagRequest.Owner
	} else {
		// Fetch user from database to get GitHub username
		user, err := rc.userRepository.FindByID(userUUID)
		if err != nil {
			utils.InternalServerErrorResponse(c, "Failed to fetch user", err)
			return
		}
		if user == nil {
			utils.UnauthorizedResponse(c, "User not found")
			return
		}
		owner = user.Username
	}

	// Validate tag name
	if tagRequest.TagName == "" {
		utils.BadRequestResponse(c, "Tag name is required")
		return
	}

	// Validate commit SHA format (40 character hex string)
	if len(tagRequest.CommitSHA) != 40 {
		utils.BadRequestResponse(c, "Invalid commit SHA format. Must be 40 character hash")
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

	log.Printf("DEBUG - Creating tag '%s' for commit %s in repository: %s/%s",
		tagRequest.TagName, tagRequest.CommitSHA, owner, tagRequest.Repo)

	// Create the tag
	reference, err := rc.repositoryService.CreateTag(c.Request.Context(), token.AccessToken, owner, tagRequest)
	if err != nil {
		log.Printf("ERROR - Failed to create tag: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to create tag", err)
		return
	}

	log.Printf("DEBUG - Successfully created tag '%s' at ref %s", tagRequest.TagName, reference.Ref)

	utils.SuccessResponse(c, http.StatusCreated, "Tag created and pushed successfully", gin.H{
		"owner":      owner,
		"repository": tagRequest.Repo,
		"tag_name":   tagRequest.TagName,
		"commit_sha": tagRequest.CommitSHA,
		"ref":        reference.Ref,
		"object_sha": reference.Object.SHA,
		"url":        reference.URL,
	})
}
