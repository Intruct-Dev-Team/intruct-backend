package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Intruct-Dev-Team/intruct-backend/pkg/storage/object"
)

type Service struct {
	projectURL string
	serviceKey string
	bucket     string
	httpClient *http.Client
}

func NewService(config *Config) *Service {
	// Build project URL from host and SSL
	scheme := "http"
	if config.UseSSL {
		scheme = "https"
	}
	projectURL := fmt.Sprintf("%s://%s", scheme, config.Host)

	return &Service{
		projectURL: projectURL,
		serviceKey: config.ServiceKey,
		bucket:     config.Bucket,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Service) UploadFile(ctx context.Context, file *object.File) error {
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		s.projectURL, s.bucket, file.Name)

	// Read file content
	fileContent, err := io.ReadAll(file.FileReader)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(fileContent))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Determine content type from extension
	contentType := s.getContentType(file.Extension)
	req.Header.Set("Authorization", "Bearer "+s.serviceKey)
	req.Header.Set("Content-Type", contentType)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *Service) GetFile(ctx context.Context, fileName string) (*object.File, error) {
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		s.projectURL, s.bucket, fileName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.serviceKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("file not found: %s", fileName)
		}
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get file failed with status %d: %s", resp.StatusCode, string(body))
	}

	return &object.File{
		Name:       fileName,
		Extension:  strings.TrimPrefix(filepath.Ext(fileName), "."),
		FileReader: resp.Body,
		Size:       resp.ContentLength,
	}, nil
}

func (s *Service) IsFileExists(ctx context.Context, fileName string) (bool, error) {
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		s.projectURL, s.bucket, fileName)

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.serviceKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

func (s *Service) GetPresignedURL(ctx context.Context, fileName string, expiry time.Duration) (string, error) {
	url := fmt.Sprintf("%s/storage/v1/object/sign/%s/%s",
		s.projectURL, s.bucket, fileName)

	expiresIn := int(expiry.Seconds())
	payload := map[string]interface{}{
		"expiresIn": expiresIn,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.serviceKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to create presigned URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("presigned URL failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		SignedURL string `json:"signedURL"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Supabase returns relative path, need to prepend project URL
	if strings.HasPrefix(result.SignedURL, "/") {
		return s.projectURL + result.SignedURL, nil
	}

	return result.SignedURL, nil
}

func (s *Service) GetPublicURL(fileName string) string {
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
		s.projectURL, s.bucket, fileName)
}

func (s *Service) DeleteFile(ctx context.Context, fileName string) error {
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		s.projectURL, s.bucket, fileName)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.serviceKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *Service) getContentType(extension string) string {
	contentTypes := map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"webp": "image/webp",
		"svg":  "image/svg+xml",
		"pdf":  "application/pdf",
		"mp4":  "video/mp4",
		"webm": "video/webm",
		"mp3":  "audio/mpeg",
		"wav":  "audio/wav",
		"txt":  "text/plain",
		"json": "application/json",
		"xml":  "application/xml",
		"zip":  "application/zip",
	}

	if contentType, ok := contentTypes[strings.ToLower(extension)]; ok {
		return contentType
	}

	return "application/octet-stream"
}
