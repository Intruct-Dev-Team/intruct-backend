package n8n

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	"go.uber.org/zap"
)

var allowedExtensions = map[string]bool{
	".pdf": true,
	".txt": true,
}

type N8NService struct {
	baseURL    string
	httpClient *http.Client
	log        *zap.Logger
}

func NewService(baseURL string, log *zap.Logger) *N8NService {
	return &N8NService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		log: log,
	}
}

func (s *N8NService) SendCourse(ctx context.Context, course *entities.Course, fileReader io.Reader, fileSize int64, fileName string) error {
	// Validate file extension
	ext := strings.ToLower(filepath.Ext(fileName))
	if !allowedExtensions[ext] {
		return fmt.Errorf("unsupported file format: %s, only .pdf and .txt are allowed", ext)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add course_id field
	if err := writer.WriteField("course_id", strconv.Itoa(course.ID)); err != nil {
		return fmt.Errorf("failed to write course_id field: %w", err)
	}

	// Add course_title field
	if err := writer.WriteField("course_title", course.Title); err != nil {
		return fmt.Errorf("failed to write course_title field: %w", err)
	}

	// Add language field
	if err := writer.WriteField("language", course.Language); err != nil {
		return fmt.Errorf("failed to write language field: %w", err)
	}

	// Add file field with original filename
	fileWriter, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(fileWriter, fileReader); err != nil {
		return fmt.Errorf("failed to copy file data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create request
	url := s.baseURL + "/webhook/upload-material"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to n8n: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.log.Error("failed to close response body", zap.Error(err))
		}
	}()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("n8n returned non-2xx status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
