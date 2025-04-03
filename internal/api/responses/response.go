package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTP status codes
const (
	HTTPOk                  = http.StatusOK
	HTTPCreated             = http.StatusCreated
	HTTPUnprocessableEntity = http.StatusUnprocessableEntity
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Meta    Meta        `json:"meta"`
}

// Meta contains pagination metadata
type Meta struct {
	Page       int `json:"page"`
	TotalPages int `json:"totalPages"`
	PerPage    int `json:"perPage"`
	TotalItems int `json:"totalItems"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success     bool        `json:"success"`
	Message     string      `json:"message"`
	Description string      `json:"description,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

// SendPaginatedResponse sends a paginated response
func SendPaginatedResponse(c *gin.Context, data interface{}, message string, page, totalPages, perPage, totalItems int) {
	response := PaginatedResponse{
		Success: true,
		Data:    data,
		Message: message,
		Meta: Meta{
			Page:       page,
			TotalPages: totalPages,
			PerPage:    perPage,
			TotalItems: totalItems,
		},
	}

	c.JSON(HTTPOk, response)
}

// SendResponse sends a success response with data
func SendResponse(c *gin.Context, data interface{}, message string, statusCode int) {
	response := Response{
		Success: true,
		Data:    data,
		Message: message,
	}

	c.JSON(statusCode, response)
}

// SendError sends an error response
func SendError(c *gin.Context, message string, description string, code int, data ...interface{}) {
	response := ErrorResponse{
		Success:     false,
		Message:     message,
		Description: description,
	}

	if len(data) > 0 {
		response.Data = data[0]
	}

	c.JSON(code, response)
}

// SendSuccess sends a simple success message
func SendSuccess(c *gin.Context, message string, statusCode int) {
	response := Response{
		Success: true,
		Message: message,
	}

	c.JSON(statusCode, response)
}