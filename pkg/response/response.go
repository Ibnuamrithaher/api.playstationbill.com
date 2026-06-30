package response

import (
	"github.com/gin-gonic/gin"
)

// Response defines the standard JSON response format for the application
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SendSuccess sends a success JSON response with the given status code, message, and data
func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendError sends an error JSON response with the given status code, message, and error details
func SendError(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
