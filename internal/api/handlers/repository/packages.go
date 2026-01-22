package repository

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
)

// GetUserPackages returns all packages for the authenticated user
// GET /api/repositories/packages?package_type=npm
func (h *Handler) GetUserPackages(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		pkghttp.UnauthorizedResponse(c, "User not found in context")
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Invalid user ID format", err)
		return
	}

	packageType := c.Query("package_type")

	accessToken, err := h.getAccessToken(userUUID)
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	packages, err := h.repositoryService.GetUserPackages(c.Request.Context(), accessToken, packageType)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch packages", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "Packages fetched successfully", gin.H{
		"packages":      packages,
		"package_count": len(packages),
		"package_type":  packageType,
	})
}

// GetOrgPackages returns all packages for a specific organization
// GET /api/repositories/:org/packages?package_type=container
func (h *Handler) GetOrgPackages(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		pkghttp.UnauthorizedResponse(c, "User not found in context")
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Invalid user ID format", err)
		return
	}

	org := c.Param("org")
	if org == "" {
		pkghttp.BadRequestResponse(c, "Organization name is required")
		return
	}

	packageType := c.Query("package_type")

	accessToken, err := h.getAccessToken(userUUID)
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	packages, err := h.repositoryService.GetOrgPackages(c.Request.Context(), accessToken, org, packageType)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch packages", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "Packages fetched successfully", gin.H{
		"organization":  org,
		"packages":      packages,
		"package_count": len(packages),
		"package_type":  packageType,
	})
}
