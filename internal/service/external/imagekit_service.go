package external

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"beautyessentials.com/internal/config"
)

// ImageKitService handles interactions with the ImageKit API
type ImageKitService struct {
	publicKey   string
	privateKey  string
	urlEndpoint string
	client      *http.Client
}

// ImageKitUploadResponse represents the response from ImageKit upload API
type ImageKitUploadResponse struct {
	FileID       string `json:"fileId"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Size         int64  `json:"size"`
	FileType     string `json:"fileType"`
}

// ImageKitErrorResponse represents an error response from ImageKit API
type ImageKitErrorResponse struct {
	Message string `json:"message"`
}

// NewImageKitService creates a new instance of ImageKitService
func NewImageKitService(cfg *config.Config) *ImageKitService {
	imageKitConfig := cfg.ImageKit()
	return &ImageKitService{
		publicKey:   imageKitConfig.PublicKey,
		privateKey:  imageKitConfig.PrivateKey,
		urlEndpoint: imageKitConfig.URLEndpoint,
		client:      &http.Client{},
	}
}

// UploadFile uploads a file to ImageKit
func (s *ImageKitService) UploadFile(file *multipart.FileHeader) (map[string]interface{}, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read the file content
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Encode file content to base64
	encodedFile := base64.StdEncoding.EncodeToString(fileBytes)

	// Prepare form data
	formData := url.Values{}
	formData.Set("file", encodedFile)
	formData.Set("fileName", file.Filename)
	formData.Set("useUniqueFileName", "true")

	// Create request
	req, err := http.NewRequest("POST", "https://upload.imagekit.io/api/v1/files/upload", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.privateKey, "")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResp ImageKitErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Message != "" {
			return nil, errors.New(errorResp.Message)
		}
		return nil, fmt.Errorf("failed to upload file: status code %d", resp.StatusCode)
	}

	// Parse response
	var uploadResp ImageKitUploadResponse
	if err := json.Unmarshal(body, &uploadResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Return result
	return map[string]interface{}{
		"fileId":    uploadResp.FileID,
		"url":       uploadResp.URL,
		"thumbUrl":  uploadResp.ThumbnailURL,
		"size":      uploadResp.Size,
		"file_type": uploadResp.FileType,
	}, nil
}

// UploadFromURL uploads a file from a URL to ImageKit
func (s *ImageKitService) UploadFromURL(sourceURL string, fileName string) (map[string]interface{}, error) {
	// If fileName is empty, extract it from URL
	if fileName == "" {
		parsedURL, err := url.Parse(sourceURL)
		if err == nil {
			fileName = filepath.Base(parsedURL.Path)
		} else {
			fileName = "file"
		}
	}

	// Prepare form data
	formData := url.Values{}
	formData.Set("file", sourceURL)
	formData.Set("fileName", fileName)
	formData.Set("useUniqueFileName", "true")

	// Create request
	req, err := http.NewRequest("POST", "https://upload.imagekit.io/api/v1/files/upload", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.privateKey, "")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResp ImageKitErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Message != "" {
			return nil, errors.New(errorResp.Message)
		}
		return nil, fmt.Errorf("failed to upload file from URL: status code %d", resp.StatusCode)
	}

	// Parse response
	var uploadResp ImageKitUploadResponse
	if err := json.Unmarshal(body, &uploadResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Return result
	return map[string]interface{}{
		"fileId":    uploadResp.FileID,
		"url":       uploadResp.URL,
		"thumbUrl":  uploadResp.ThumbnailURL,
		"size":      uploadResp.Size,
		"file_type": uploadResp.FileType,
	}, nil
}

// DeleteFile deletes a file from ImageKit
func (s *ImageKitService) DeleteFile(fileID string) error {
	// Create request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.imagekit.io/v1/files/%s", fileID), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.SetBasicAuth(s.privateKey, "")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for success (204 No Content)
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		var errorResp ImageKitErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Message != "" {
			return errors.New(errorResp.Message)
		}
		return fmt.Errorf("failed to delete file: status code %d", resp.StatusCode)
	}

	return nil
}

// DeleteBulkFiles deletes multiple files from ImageKit
func (s *ImageKitService) DeleteBulkFiles(fileIDs []string) error {
	// If more than 100 files, chunk them into batches of 100
	if len(fileIDs) > 100 {
		chunks := make([][]string, 0)
		for i := 0; i < len(fileIDs); i += 100 {
			end := i + 100
			if end > len(fileIDs) {
				end = len(fileIDs)
			}
			chunks = append(chunks, fileIDs[i:end])
		}
		
		for _, chunk := range chunks {
			// Process each chunk with a single API call
			if err := s.bulkDeleteFilesRequest(chunk); err != nil {
				return err
			}
		}
		return nil
	}
	
	// For 100 or fewer files, process them in a single request
	return s.bulkDeleteFilesRequest(fileIDs)
}

// bulkDeleteFilesRequest makes the actual API request to delete files in bulk
func (s *ImageKitService) bulkDeleteFilesRequest(fileIDs []string) error {
	// Prepare request body
	type requestBody struct {
		FileIDs []string `json:"fileIds"`
	}
	reqBody := requestBody{FileIDs: fileIDs}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", "https://api.imagekit.io/v1/files/batch/deleteByFileIds", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(s.privateKey, "")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for success (204 No Content)
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		var errorResp ImageKitErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Message != "" {
			return errors.New(errorResp.Message)
		}
		return fmt.Errorf("failed to delete files in bulk: status code %d", resp.StatusCode)
	}

	return nil
}

// ListAllFiles lists all files from ImageKit
func (s *ImageKitService) ListAllFiles() ([]map[string]interface{}, error) {
	// Create request
	req, err := http.NewRequest("GET", "https://api.imagekit.io/v1/files", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.SetBasicAuth(s.privateKey, "")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		var errorResp ImageKitErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Message != "" {
			return nil, errors.New(errorResp.Message)
		}
		return nil, fmt.Errorf("failed to list files: status code %d", resp.StatusCode)
	}

	// Parse response
	var files []map[string]interface{}
	if err := json.Unmarshal(body, &files); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return files, nil
}