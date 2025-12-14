package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest"
	"github.com/Intruct-Dev-Team/intruct-backend/pkg/migrator"
	"github.com/Intruct-Dev-Team/intruct-backend/pkg/storage/db/postgres"
	"github.com/Intruct-Dev-Team/intruct-backend/pkg/storage/object/s3"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	maxPort      = 1<<16 - 1
	maskedString = "********"
)

type Config struct {
	DB           postgres.Config
	Migrator     migrator.Config
	Server       rest.Config
	S3           s3.Config
	IsInitDb     bool   `env:"IS_INIT_DB" env-required:"true"`
	JwtSecretKey string `env:"JWT_SECRET_KEY" env-required:"true"`
}

func LoadConfig(paths []string) (*Config, error) {
	var config Config

	// Looking for .env file in different directories
	var envPath string
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			envPath = path
			break
		}
	}

	if envPath != "" {
		err := cleanenv.ReadConfig(envPath, &config)
		if err != nil {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	} else {
		// if .env unexists - read from environment (for Docker)
		err := cleanenv.ReadEnv(&config)
		if err != nil {
			return nil, fmt.Errorf("failed to read env: %w", err)
		}
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// // Validate config validation.
func (c *Config) Validate() error {
	if !checkPortValidation(c.Server.Port) {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if !checkPortValidation(c.DB.Port) {
		return fmt.Errorf("invalid database port: %d", c.DB.Port)
	}

	if !checkPortValidation(c.Migrator.Port) {
		return fmt.Errorf("invalid database (migrator) port: %d", c.DB.Port)
	}

	if !checkPortValidation(c.S3.Port) {
		return fmt.Errorf("invalid minio port: %d", c.S3.Port)
	}

	return nil
}

// LogConfig logs configuration with sensitive data masking.
func (c *Config) LogConfig() (string, error) {
	// Create a copy of config for logging
	logConfig := *c

	// Mask passwords

	if logConfig.DB.Password != "" {
		logConfig.DB.Password = maskedString
	}

	if logConfig.Migrator.Password != "" {
		logConfig.Migrator.Password = maskedString
	}

	if logConfig.JwtSecretKey != "" {
		logConfig.JwtSecretKey = maskedString
	}

	if logConfig.S3.AccessKey != "" {
		logConfig.S3.AccessKey = maskedString
	}

	if logConfig.S3.SecretKey != "" {
		logConfig.S3.SecretKey = maskedString
	}

	// Convert to JSON with indents for readability
	jsonBytes, err := json.MarshalIndent(logConfig, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling config: %w", err)
	}

	return "Application Configuration:\n" + string(jsonBytes), nil
}

func checkPortValidation(port int) bool {
	return port >= 1 && port <= maxPort
}
