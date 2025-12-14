package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIDKey        = "user_id"
	ExternalUserUUID = "sub"
	defaultTTL       = time.Hour * 24 * 7
)

var ErrorTokenExpired = errors.New("token is expired")

type JWTService struct {
	secretKey []byte
	issuer    string
	duration  time.Duration
}

type Option func(*JWTService)

func WithDuration(duration time.Duration) Option {
	return func(s *JWTService) {
		s.duration = duration
	}
}

func WithIssuer(issuer string) Option {
	return func(s *JWTService) {
		s.issuer = issuer
	}
}

func NewService(secretKey string, opts ...Option) *JWTService {
	s := &JWTService{
		secretKey: []byte(secretKey),
		duration:  defaultTTL,
		issuer:    "default",
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// GenerateJWTToken creates a JWT token for a user.
func (s *JWTService) GenerateJWTToken(userID int) (string, error) {
	// Set token expiration time
	expirationTime := time.Now().Add(s.duration)

	// Create claims
	claims := jwt.MapClaims{
		UserIDKey: userID,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
		"iss":     s.issuer,
	}

	// Create token with signing algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with secret key
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWTToken validates the JWT token.
func (s *JWTService) ValidateJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Additional expiration time check
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, ErrorTokenExpired
		}
	}

	return claims, nil
}

func (s *JWTService) ExtractExternalUserUUID(claims jwt.MapClaims) (string, error) {
	userUUID, ok := claims[ExternalUserUUID].(string)
	if !ok {
		return "", errors.New("invalid or missing user ID in claims")
	}

	return userUUID, nil
}

func (v *JWTService) ExtractEmail(claims jwt.MapClaims) (string, error) {
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return "", errors.New("email not found in claims")
	}
	return email, nil
}

func (s *JWTService) GetUserKey() string {
	return UserIDKey
}

func (s *JWTService) GetExternalUserUUIDKey() string {
	return ExternalUserUUID
}

func (s *JWTService) GetExpiredError() error {
	return ErrorTokenExpired
}
