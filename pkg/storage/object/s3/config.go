package s3

type Config struct {
	Host      string `env:"S3_HOST"        env-required:"true"`
	Port      int    `env:"S3_PORT"        env-default:"443"` // Default 443
	Region    string `env:"S3_REGION"      env-default:"us-east-1"`
	AccessKey string `env:"S3_ACCESS_KEY"  env-required:"true"`
	SecretKey string `env:"S3_SECRET_KEY"  env-required:"true"`
	Bucket    string `env:"S3_BUCKET"      env-required:"true"`
	UseSSL    bool   `env:"S3_USE_SSL"     env-default:"true"`
}
