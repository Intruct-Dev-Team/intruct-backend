package supabase

type Config struct {
	Host       string `env:"SUPABASE_HOST"        env-required:"true"` // sdhlcmhzujflfrfdtuyc.supabase.co
	UseSSL     bool   `env:"SUPABASE_USE_SSL"     env-default:"true"`  // всегда true для Supabase
	ServiceKey string `env:"SUPABASE_SERVICE_KEY" env-required:"true"` // service_role key
	Bucket     string `env:"SUPABASE_BUCKET"      env-required:"true"` // avatars
}
