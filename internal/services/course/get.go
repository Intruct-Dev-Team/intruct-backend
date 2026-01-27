package course

import (
	"context"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
)

func (s *CourseService) GetCourse(ctx context.Context, courseID int, userID int) (*entities.Course, error) {
	// check if user exists
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return nil, serviceErrs.ErrorUserNotFound
	}

	course, err := s.repo.GetCourseByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, serviceErrs.ErrorCourseNotFound
		}
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	// check access: course must be published or user must be the owner
	if !course.IsPublic && course.OwnerID != userID {
		return nil, serviceErrs.ErrorForbidden
	}

	// get modules with lessons in one query
	modules, err := s.repo.GetModulesWithLessonsByCourseID(ctx, courseID)
	if err != nil && !errors.Is(err, serviceErrs.ErrorSelectEmpty) {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}
	course.Modules = modules

	// get user's course progression if exists
	progression, err := s.repo.GetCourseProgressionByUserAndCourse(ctx, userID, courseID)
	if err != nil && !errors.Is(err, serviceErrs.ErrorSelectEmpty) {
		return nil, fmt.Errorf("failed to get course progression: %w", err)
	}
	if progression != nil {
		course.CourseProgression = progression
	}

	// get statistic
	stat, err := s.repo.GetCourseStatisticByCourseID(ctx, courseID)
	if err != nil || stat == nil {
		return nil, fmt.Errorf("failed to get course statistic: %w", err)
	}
	course.Statistic = *stat

	return course, nil
}
