package middlewares

import (
	"errors"
	"log"
	"net/http"

	"beautyessentials.com/internal/api/responses"
	"github.com/gin-gonic/gin"
)

// ErrorType represents different types of errors
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	// ErrorTypeInternal represents internal server errors
	ErrorTypeInternal ErrorType = "INTERNAL_ERROR"
	// ErrorTypeNotFound represents not found errors
	ErrorTypeNotFound ErrorType = "NOT_FOUND"
	// ErrorTypeUnauthorized represents unauthorized errors
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	// ErrorTypeForbidden represents forbidden errors
	ErrorTypeForbidden ErrorType = "FORBIDDEN"
)

// AppError represents an application error
type AppError struct {
	Type        ErrorType
	Message     string
	Description string
	Data        interface{}
}

// Error implements the error interface
func (e AppError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(message string, description string, data ...interface{}) AppError {
	var errorData interface{}
	if len(data) > 0 {
		errorData = data[0]
	}
	return AppError{
		Type:        ErrorTypeValidation,
		Message:     message,
		Description: description,
		Data:        errorData,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(message string, description string) AppError {
	return AppError{
		Type:        ErrorTypeInternal,
		Message:     message,
		Description: description,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, description string) AppError {
	return AppError{
		Type:        ErrorTypeNotFound,
		Message:     message,
		Description: description,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string, description string) AppError {
	return AppError{
		Type:        ErrorTypeUnauthorized,
		Message:     message,
		Description: description,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string, description string) AppError {
	return AppError{
		Type:        ErrorTypeForbidden,
		Message:     message,
		Description: description,
	}
}

// ErrorHandler is a middleware that handles errors
func ErrorHandler(respHelper *responses.ResponseHelper) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last().Err
			
			// Log the error
			log.Printf("Error: %v", err)
			
			// Handle different types of errors
			var appErr AppError
			if errors.As(err, &appErr) {
				// Handle application errors
				switch appErr.Type {
				case ErrorTypeValidation:
					respHelper.SendError(c, appErr.Message, appErr.Description, http.StatusUnprocessableEntity, appErr.Data)
				case ErrorTypeInternal:
					respHelper.SendError(c, appErr.Message, appErr.Description, http.StatusInternalServerError)
				case ErrorTypeNotFound:
					respHelper.SendError(c, appErr.Message, appErr.Description, http.StatusNotFound)
				case ErrorTypeUnauthorized:
					respHelper.SendError(c, appErr.Message, appErr.Description, http.StatusUnauthorized)
				case ErrorTypeForbidden:
					respHelper.SendError(c, appErr.Message, appErr.Description, http.StatusForbidden)
				default:
					respHelper.SendError(c, "An unexpected error occurred", err.Error(), http.StatusInternalServerError)
				}
			} else {
				// Handle generic errors
				respHelper.SendError(c, "An unexpected error occurred", err.Error(), http.StatusInternalServerError)
			}
		}
	}
}