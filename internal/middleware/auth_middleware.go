package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/utils"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Extract token from header
		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			if err == utils.ErrExpiredToken {
				utils.UnauthorizedResponse(c, "Token has expired")
			} else {
				utils.UnauthorizedResponse(c, "Invalid token")
			}
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID.String())
		c.Set("github_id", claims.GitHubID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens but doesn't require them
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.Next()
			return
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID.String())
		c.Set("github_id", claims.GitHubID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// GetUserIDFromContext gets the user ID from gin context
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return uuid.Nil, false
	}

	return userID, true
}

// CORSMiddleware handles CORS
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == strings.TrimSpace(allowedOrigin) {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
