package context

type contextKey string

const (
	UserIDKey           contextKey = "user_id"
	ExternalUserUUIDKey contextKey = "external_user_uuid"
	EmailKey            contextKey = "email"
)
