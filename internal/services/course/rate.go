package course

import (
	"context"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
)

// RateCourse add rating to the course
func (s *CourseService) RateCourse(ctx context.Context, courseID int, userID int, rating int) error {
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existence by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// check course existss
	exists, err = s.repo.IsCourseExistsByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("failed to check course existence by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorCourseNotFound
	}

	// check user progress
	_, err = s.repo.GetCourseProgressionByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorUserHasNoCourseProgression
		}
	}

	ratingEntity := &entities.Rating{
		CourseID: courseID,
		UserID:   userID,
		Rating:   rating,
	}

	err = s.repo.CreateRating(ctx, ratingEntity)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorNonUniqueData) {
			return serviceErrs.ErrorRatingExists
		}
		return fmt.Errorf("failed to ctreate rating: %w", err)
	}

	return nil
}
