package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
)

// CaseConverterMiddleware converts request body from camelCase to snake_case
// and response body from snake_case to camelCase
func CaseConverterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip for non-JSON requests
		if !strings.Contains(c.GetHeader("Content-Type"), "application/json") && c.Request.Method != "GET" {
			c.Next()
			return
		}

		// Read the request body
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}

		// Close the original body
		c.Request.Body.Close()

		// If the request body is not empty, convert it
		if len(requestBody) > 0 {
			// Convert request from camelCase to snake_case
			var requestMap map[string]interface{}
			if err := json.Unmarshal(requestBody, &requestMap); err == nil {
				convertedRequest := convertMapKeysToSnakeCase(requestMap)
				newRequestBody, err := json.Marshal(convertedRequest)
				if err == nil {
					c.Request.Body = io.NopCloser(bytes.NewBuffer(newRequestBody))
					c.Request.ContentLength = int64(len(newRequestBody))
				}
			} else {
				// If we can't unmarshal as map, try as array
				var requestArray []interface{}
				if err := json.Unmarshal(requestBody, &requestArray); err == nil {
					convertedRequest := convertArrayKeysToSnakeCase(requestArray)
					newRequestBody, err := json.Marshal(convertedRequest)
					if err == nil {
						c.Request.Body = io.NopCloser(bytes.NewBuffer(newRequestBody))
						c.Request.ContentLength = int64(len(newRequestBody))
					}
				} else {
					// If conversion fails, restore original body
					c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
				}
			}
		} else {
			// Restore empty body
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create a custom response writer to capture the response
		writer := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		// Process request
		c.Next()

		// Only proceed if content type is JSON
		contentType := writer.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			return
		}

		// Get response body
		responseBody := writer.body.Bytes()

		// If the response body is not empty, convert it
		if len(responseBody) > 0 {
			// Convert response from snake_case to camelCase
			var responseMap map[string]interface{}
			if err := json.Unmarshal(responseBody, &responseMap); err == nil {
				convertedResponse := convertMapKeysToCamelCase(responseMap)
				newResponseBody, err := json.Marshal(convertedResponse)
				if err == nil {
					// Reset headers and write the converted response
					writer.Header().Set("Content-Length", string(rune(len(newResponseBody))))
					writer.ResponseWriter.WriteHeader(writer.status)
					writer.ResponseWriter.Write(newResponseBody)
					return
				}
			} else {
				// If we can't unmarshal as map, try as array
				var responseArray []interface{}
				if err := json.Unmarshal(responseBody, &responseArray); err == nil {
					convertedResponse := convertArrayKeysToCamelCase(responseArray)
					newResponseBody, err := json.Marshal(convertedResponse)
					if err == nil {
						// Reset headers and write the converted response
						writer.Header().Set("Content-Length", string(rune(len(newResponseBody))))
						writer.ResponseWriter.WriteHeader(writer.status)
						writer.ResponseWriter.Write(newResponseBody)
						return
					}
				}
			}
		}

		// If conversion fails, write the original response
		writer.ResponseWriter.Write(responseBody)
	}
}

// responseBodyWriter is a custom response writer that captures the response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

// Write captures the response body
func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return len(b), nil
}

// WriteHeader captures the status code
func (r *responseBodyWriter) WriteHeader(statusCode int) {
	r.status = statusCode
}

// WriteString captures the response body
func (r *responseBodyWriter) WriteString(s string) (int, error) {
	r.body.WriteString(s)
	return len(s), nil
}

// convertMapKeysToSnakeCase converts all keys in a map to snake_case
func convertMapKeysToSnakeCase(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		switch v := v.(type) {
		case map[string]interface{}:
			result[strcase.ToSnake(k)] = convertMapKeysToSnakeCase(v)
		case []interface{}:
			result[strcase.ToSnake(k)] = convertArrayKeysToSnakeCase(v)
		default:
			result[strcase.ToSnake(k)] = v
		}
	}
	return result
}

// convertArrayKeysToSnakeCase converts all keys in an array of maps to snake_case
func convertArrayKeysToSnakeCase(a []interface{}) []interface{} {
	result := make([]interface{}, len(a))
	for i, v := range a {
		switch v := v.(type) {
		case map[string]interface{}:
			result[i] = convertMapKeysToSnakeCase(v)
		case []interface{}:
			result[i] = convertArrayKeysToSnakeCase(v)
		default:
			result[i] = v
		}
	}
	return result
}

// convertMapKeysToCamelCase converts all keys in a map to camelCase
func convertMapKeysToCamelCase(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		switch v := v.(type) {
		case map[string]interface{}:
			result[strcase.ToLowerCamel(k)] = convertMapKeysToCamelCase(v)
		case []interface{}:
			result[strcase.ToLowerCamel(k)] = convertArrayKeysToCamelCase(v)
		default:
			result[strcase.ToLowerCamel(k)] = v
		}
	}
	return result
}

// convertArrayKeysToCamelCase converts all keys in an array of maps to camelCase
func convertArrayKeysToCamelCase(a []interface{}) []interface{} {
	result := make([]interface{}, len(a))
	for i, v := range a {
		switch v := v.(type) {
		case map[string]interface{}:
			result[i] = convertMapKeysToCamelCase(v)
		case []interface{}:
			result[i] = convertArrayKeysToCamelCase(v)
		default:
			result[i] = v
		}
	}
	return result
}