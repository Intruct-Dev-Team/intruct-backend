package object

import (
	"context"
	"time"
)

type StorageService interface {
	UploadFile(ctx context.Context, file *File) error
	GetFile(ctx context.Context, fileName string) (*File, error)
	IsFileExists(ctx context.Context, fileName string) (bool, error)
	GetPresignedURL(ctx context.Context, fileName string, expiry time.Duration) (string, error)
	GetPublicURL(fileName string) string
	DeleteFile(ctx context.Context, fileName string) error
}
