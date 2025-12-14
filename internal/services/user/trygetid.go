package user

import (
	"context"
	"errors"
	"fmt"

	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
)

func (s *UserService) TryGetUserIDByExternalUUID(ctx context.Context, ExternalUserUUID string) (id int, exists bool, err error) {
	if ExternalUserUUID == "" {
		return 0, false, fmt.Errorf("ExternalUUID cannot be empty")
	}

	userID, err := s.repo.GetUserIDByExternalUUID(ctx, ExternalUserUUID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return 0, false, nil // it's not a error, because 'try' method
		}

		return 0, false, fmt.Errorf("failed to get user by external UUID %s: %w", ExternalUserUUID, err)
	}

	return userID, true, nil
}
