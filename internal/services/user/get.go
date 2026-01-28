package user

import (
	"context"
	"errors"
	"fmt"
	"time"

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

	streak, err := s.repo.GetUserStreakByUserID(ctx, userID)
	if err != nil || streak == nil {
		// if streak is empty just return user with 0 streak
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return user, nil
		}
		return nil, fmt.Errorf("failed to get user streak: %w", err)
	}

	if streak.DaysStreakCount > 0 {
		now := time.Now().UTC()
		today := time.Date(
			now.Year(), now.Month(), now.Day(),
			0, 0, 0, 0,
			time.UTC,
		)

		updatedDay := time.Date(
			streak.UpdatedAt.UTC().Year(),
			streak.UpdatedAt.UTC().Month(),
			streak.UpdatedAt.UTC().Day(),
			0, 0, 0, 0,
			time.UTC,
		)

		if updatedDay.Equal(today) {
			streak.IsStreakActiveToday = true
		} else if updatedDay.Before(today.AddDate(0, 0, -1)) {
			streak.DaysStreakCount = 0
		}
	}

	user.Statistic = *stat
	user.Streak = *streak

	return user, nil
}
