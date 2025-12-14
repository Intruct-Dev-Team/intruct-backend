package minio

type Config struct {
	Host      string `env:"MINIO_HOST"       env-required:"true"`
	Port      int    `env:"MINIO_PORT"       env-required:"true"`
	AccessKey string `env:"MINIO_ACCESS_KEY" env-required:"true"`
	SecretKey string `env:"MINIO_SECRET_KEY" env-required:"true"`
	Bucket    string `env:"MINIO_BUCKET"     env-required:"true"`
	UseSSL    bool   `env:"MINIO_USE_SSL"     env-required:"true"`
}
