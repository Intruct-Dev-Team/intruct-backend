package course

import (
	"context"
	"errors"
	"fmt"
	"time"

	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
)

// GetCourseState return course state (if user is owner)
func (s *CourseService) GetCourseStateAndUpdatedTime(ctx context.Context, courseID int, userID int) (string, time.Time, error) {
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to check user existence by id: %w", err)
	}
	if !exists {
		return "", time.Time{}, serviceErrs.ErrorUserNotFound
	}

	// Check course exists and get its current state
	course, err := s.repo.GetCourseByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return "", time.Time{}, serviceErrs.ErrorCourseNotFound
		}

		return "", course.UpdatedAt, fmt.Errorf("failed to get course: %w", err)
	}

	if course.OwnerID != userID {
		return "", course.UpdatedAt, serviceErrs.ErrorNotCourseOwner
	}

	return course.State, course.UpdatedAt, nil
}
