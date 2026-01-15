package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	contextkeys "github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/context"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/httputils"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type TokenValidator interface {
	ValidateJWTToken(tokenString string) (jwt.MapClaims, error)
	ExtractExternalUserUUID(claims jwt.MapClaims) (string, error)
	ExtractEmail(claims jwt.MapClaims) (string, error)
	GetExternalUserUUIDKey() string
	GetExpiredError() error
}

type UserService interface {
	TryGetUserIDByExternalUUID(ctx context.Context, ExternalUserUUID string) (id int, exists bool, err error)
}

// JWTAuthMiddleware middleware for JWT token validation
func JWTAuthMiddleware(validator TokenValidator, log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httputils.RespondWith401(w, "missed Authorization header (required)", log)

				return
			}

			// check header format (Bearer Token)
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				httputils.RespondWith401(w, "Invalid token format", log)

				return
			}

			// extract token
			tokenString := parts[1]

			// validate token using provided validator
			claims, err := validator.ValidateJWTToken(tokenString)
			if err != nil {
				if errors.Is(err, validator.GetExpiredError()) {
					log.Error("token expired", zap.Error(err))
					httputils.RespondWith401(w, "token expired", log)

					return
				}

				log.Error("failed to validate token", zap.Error(err))
				httputils.RespondWith401(w, "Failed to validate token", log)

				return
			}

			// extract ID from token
			userUUID, err := validator.ExtractExternalUserUUID(claims)
			if err != nil {
				log.Error("failed to extract external user's UUID", zap.Error(err))
				httputils.RespondWith401(w,
					fmt.Sprintf("Invalid token: missing field: %s", validator.GetExternalUserUUIDKey()),
					log)
				return
			}

			email, err := validator.ExtractEmail(claims)
			if err != nil {
				log.Error("failed to extract email", zap.Error(err))
				httputils.RespondWith401(w,
					"Invalid token: missing field: email",
					log)
				return
			}

			ctx := context.WithValue(r.Context(), contextkeys.ExternalUserUUIDKey, userUUID)
			ctx = context.WithValue(ctx, contextkeys.EmailKey, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SystemAuthMiddleware specially middlware for checking auth in own system
func SystemAuthMiddleware(userService UserService, log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userUUID, ok := r.Context().Value(contextkeys.ExternalUserUUIDKey).(string) // get external user's UUID after JWT Auth middlware
			if !ok {
				log.Error("external user UUID not found in context")
				httputils.RespondWith401(w, "Unauthorized", log)
				return
			}

			userID, exists, err := userService.TryGetUserIDByExternalUUID(r.Context(), userUUID)
			if err != nil {
				log.Error("failed to sync user with database",
					zap.String("external_uuid", userUUID),
					zap.Error(err),
				)

				httputils.RespondWith500(w, log)
				return
			}

			if !exists {
				log.Warn("user not found", zap.String("external_uuid", userUUID))
				httputils.RespondWith401(w, "registration was not completed", log)
				return
			}

			ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, userID)

			// log.Debug("user authenticated successfully",
			// 	zap.Int("user_id", userID),
			// 	zap.String("path", r.URL.Path),
			// )

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
