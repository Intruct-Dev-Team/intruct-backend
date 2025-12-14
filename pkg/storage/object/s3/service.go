package s3

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Intruct-Dev-Team/intruct-backend/pkg/storage/object"
	"github.com/minio/minio-go/v7"
)

type Service struct {
	client *minio.Client
	bucket string
	host   string
}

func NewService(client *minio.Client, bucket, host string) *Service {
	return &Service{
		client: client,
		bucket: bucket,
		host:   host,
	}
}

func (s *Service) UploadFile(ctx context.Context, file *object.File) error {
	_, err := s.client.PutObject(ctx, s.bucket, file.Name, file.FileReader, file.Size,
		minio.PutObjectOptions{
			ContentType: "application/octet-stream", // optional
		})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (s *Service) GetFile(ctx context.Context, fileName string) (*object.File, error) {
	objectReader, err := s.client.GetObject(ctx, s.bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	stat, err := objectReader.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stat: %w", err)
	}

	return &object.File{
		Name:       fileName,
		Extension:  strings.TrimPrefix(filepath.Ext(fileName), "."),
		FileReader: objectReader,
		Size:       stat.Size,
	}, nil
}

func (s *Service) IsFileExists(ctx context.Context, fileName string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucket, fileName, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" || errResponse.Code == "NotFound" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

func (s *Service) GetPresignedURL(ctx context.Context, fileName string, expiry time.Duration) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, s.bucket, fileName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}

// for public buckets
func (s *Service) GetPublicURL(fileName string) string {
	// Supabase format
	return fmt.Sprintf("https://%s/storage/v1/object/public/%s/%s",
		s.host, s.bucket, fileName)
}

func (s *Service) DeleteFile(ctx context.Context, fileName string) error {
	err := s.client.RemoveObject(ctx, s.bucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
