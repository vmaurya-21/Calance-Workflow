package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	c.JSON(statusCode, response)
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, errors interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: "Validation failed",
		Data:    errors,
	})
}

// UnauthorizedResponse sends an unauthorized error response
func UnauthorizedResponse(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Message: message,
		Error:   "Unauthorized",
	})
}

// NotFoundResponse sends a not found error response
func NotFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: message,
		Error:   "Not found",
	})
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusInternalServerError, message, err)
}
