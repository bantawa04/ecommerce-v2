package responses

import (
	"github.com/gin-gonic/gin"
)

// ResponseHelper provides methods for standardized API responses
type ResponseHelper struct{}

// NewResponseHelper creates a new ResponseHelper
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// SendPaginatedResponse sends a paginated response
func (h *ResponseHelper) SendPaginatedResponse(c *gin.Context, data interface{}, message string, page, totalPages, perPage, totalItems int) {
	SendPaginatedResponse(c, data, message, page, totalPages, perPage, totalItems)
}

// SendResponse sends a success response with data
func (h *ResponseHelper) SendResponse(c *gin.Context, data interface{}, message string, statusCode int) {
	SendResponse(c, data, message, statusCode)
}

// SendError sends an error response
func (h *ResponseHelper) SendError(c *gin.Context, message string, description string, code int, data ...interface{}) {
	SendError(c, message, description, code, data...)
}

// SendSuccess sends a simple success message
func (h *ResponseHelper) SendSuccess(c *gin.Context, message string, statusCode int) {
	SendSuccess(c, message, statusCode)
}

// OkResponse is a shorthand for sending a 200 OK response
func (h *ResponseHelper) OkResponse(c *gin.Context, data interface{}, message string) {
	SendResponse(c, data, message, HTTPOk)
}

// CreatedResponse is a shorthand for sending a 201 Created response
func (h *ResponseHelper) CreatedResponse(c *gin.Context, data interface{}, message string) {
	SendResponse(c, data, message, HTTPCreated)
}

// ValidationError is a shorthand for sending a 422 Unprocessable Entity response
func (h *ResponseHelper) ValidationError(c *gin.Context, message string, description string, data ...interface{}) {
	SendError(c, message, description, HTTPUnprocessableEntity, data...)
}