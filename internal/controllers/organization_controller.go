package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/repositories"
	"github.com/vmaurya-21/Calance-Workflow/internal/services"
	"github.com/vmaurya-21/Calance-Workflow/internal/utils"
)

type OrganizationController struct {
	organizationService *services.GitHubOrganizationService
	tokenRepository     *repositories.TokenRepository
}

// NewOrganizationController creates a new organization controller
func NewOrganizationController(
	organizationService *services.GitHubOrganizationService,
	tokenRepository *repositories.TokenRepository,
) *OrganizationController {
	return &OrganizationController{
		organizationService: organizationService,
		tokenRepository:     tokenRepository,
	}
}

// GetUserOrganizations returns all GitHub organizations for the authenticated user
// GET /api/organizations
func (oc *OrganizationController) GetUserOrganizations(c *gin.Context) {
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

	// Fetch token from database
	token, err := oc.tokenRepository.FindByUserID(userUUID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch access token", err)
		return
	}
	if token == nil {
		utils.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	log.Printf("DEBUG - Fetching organizations for user: %s", userUUID)

	// Get organizations using the GitHub token from database
	organizations, err := oc.organizationService.GetUserOrganizations(c.Request.Context(), token.AccessToken)
	if err != nil {
		log.Printf("ERROR - Failed to fetch organizations: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to fetch organizations", err)
		return
	}

	log.Printf("DEBUG - Successfully fetched %d organizations for user %s", len(organizations), userUUID)

	utils.SuccessResponse(c, http.StatusOK, "Organizations fetched successfully", gin.H{
		"organizations":      organizations,
		"organization_count": len(organizations),
	})
}

// GetOrganizationRepositories returns all repositories for a specific organization
// GET /api/organizations/:org/repositories
func (oc *OrganizationController) GetOrganizationRepositories(c *gin.Context) {
	// Get organization name from URL parameter
	orgName := c.Param("org")
	if orgName == "" {
		utils.BadRequestResponse(c, "Organization name is required")
		return
	}

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

	// Fetch token from database
	token, err := oc.tokenRepository.FindByUserID(userUUID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch access token", err)
		return
	}
	if token == nil {
		utils.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	log.Printf("DEBUG - Fetching repositories for organization: %s", orgName)

	// Get repositories for the organization
	repositories, err := oc.organizationService.GetOrganizationRepositories(c.Request.Context(), token.AccessToken, orgName)
	if err != nil {
		log.Printf("ERROR - Failed to fetch repositories for organization %s: %v", orgName, err)
		utils.InternalServerErrorResponse(c, "Failed to fetch repositories", err)
		return
	}

	log.Printf("DEBUG - Successfully fetched %d repositories for organization %s", len(repositories), orgName)

	utils.SuccessResponse(c, http.StatusOK, "Repositories fetched successfully", gin.H{
		"organization":     orgName,
		"repositories":     repositories,
		"repository_count": len(repositories),
	})
}

// GetUserRepositories returns all repositories accessible to the authenticated user from their organizations
// GET /api/repositories
func (oc *OrganizationController) GetUserRepositories(c *gin.Context) {
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

	// Fetch token from database
	token, err := oc.tokenRepository.FindByUserID(userUUID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch access token", err)
		return
	}
	if token == nil {
		utils.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	log.Printf("DEBUG - Fetching all accessible repositories for user: %s", userUUID)

	// Get repositories from all organizations
	repositoriesByOrg, err := oc.organizationService.GetUserRepositories(c.Request.Context(), token.AccessToken)
	if err != nil {
		log.Printf("ERROR - Failed to fetch repositories: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to fetch repositories", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Repositories fetched successfully", gin.H{
		"repositories_by_org": repositoriesByOrg,
		"organization_count":  len(repositoriesByOrg),
	})
}
