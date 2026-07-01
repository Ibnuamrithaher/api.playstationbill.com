package response

import (
	"github.com/gin-gonic/gin"
)

// PaginationMeta defines the structure for pagination metadata
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Response defines the standard JSON response format for the application
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Status  int         `json:"status"`
	Meta    interface{} `json:"meta,omitempty"`
}

// SendSuccess sends a success JSON response with the given status code, message, and data
func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Status:  statusCode,
	})
}

// SendSuccessWithMeta sends a success JSON response with pagination metadata
func SendSuccessWithMeta(c *gin.Context, statusCode int, message string, data interface{}, meta PaginationMeta) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Status:  statusCode,
		Meta:    meta,
	})
}

// SendError sends an error JSON response with the given status code, message, and error details
func SendError(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Errors:  errors,
		Status:  statusCode,
	})
}
