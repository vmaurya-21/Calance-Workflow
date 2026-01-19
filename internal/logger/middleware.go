package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GinMiddleware returns a gin middleware for logging HTTP requests
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Build log event
		event := Logger.Info()

		// Add error if exists
		if len(c.Errors) > 0 {
			// Convert gin.Errors to []error
			errs := make([]error, len(c.Errors))
			for i, e := range c.Errors {
				errs[i] = e.Err
			}
			event = Logger.Error().Errs("errors", errs)
		} else if statusCode >= 500 {
			event = Logger.Error()
		} else if statusCode >= 400 {
			event = Logger.Warn()
		}

		// Log the request
		event.
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", raw).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("HTTP request")
	}
}
