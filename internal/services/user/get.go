package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
)

func (s *UserService) GetUser(ctx context.Context, userID int) (*entities.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, serviceErrs.ErrorUserNotFound
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	stat, err := s.repo.GetUserStatisticByUserID(ctx, userID)
	if err != nil || stat == nil {
		return nil, fmt.Errorf("failed to get user statistic: %w", err)
	}
	user.Statistic = *stat

	return user, nil
}
