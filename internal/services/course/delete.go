package course

import (
	"context"
	"errors"
	"fmt"

	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
)

// DeleteCourse remove modules, lessons and quizzes (with options) and set soft delete to course
func (s *CourseService) DeleteCourse(ctx context.Context, courseID int, userID int) error {
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existence by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// Check course exists and get its current state
	course, err := s.repo.GetCourseByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorCourseNotFound
		}

		return fmt.Errorf("failed to get course: %w", err)
	}

	if course.OwnerID != userID {
		return serviceErrs.ErrorNotCourseOwner
	}

	if err := s.repo.DeleteCourseDataAndSoftDelete(ctx, courseID); err != nil {
		return fmt.Errorf("failed to delete course data: %w", err)
	}

	return nil
}
