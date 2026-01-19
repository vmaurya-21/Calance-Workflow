package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
)

// GetProfile returns the current authenticated user
// GET /api/auth/me
func (h *Handler) GetProfile(c *gin.Context) {
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

	user, err := h.userRepository.FindByID(userUUID)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch user", err)
		return
	}
	if user == nil {
		pkghttp.NotFoundResponse(c, "User not found")
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "User fetched successfully", user.ToResponse())
}

// Logout logs out the user
// POST /api/auth/logout
func (h *Handler) Logout(c *gin.Context) {
	// JWT is stateless, so logout is handled client-side
	pkghttp.SuccessResponse(c, http.StatusOK, "Logged out successfully", nil)
}
