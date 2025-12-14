package s3

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewClient(config *Config) (*minio.Client, error) {
	address := config.Host
	if config.Port != 443 && config.Port != 80 {
		address = fmt.Sprintf("%s:%d", config.Host, config.Port)
	}

	client, err := minio.New(address, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	return client, nil
}
