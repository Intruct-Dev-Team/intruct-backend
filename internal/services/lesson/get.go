package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
)

func (s *LessonService) GetLesson(ctx context.Context, lessonID int, userID int) (*entities.Lesson, error) {
	// check if user exists
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return nil, serviceErrs.ErrorUserNotFound
	}

	// get lesson
	lesson, err := s.repo.GetLessonByID(ctx, lessonID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, serviceErrs.ErrorLessonNotFound
		}
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	// get course to check access
	course, err := s.repo.GetCourseByID(ctx, lesson.CourseID)
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

	// if course.OwnerID != userID { // if user is not the owner, check lesson progression
	// get user's course progression
	progression, err := s.repo.GetCourseProgressionByUserAndCourse(ctx, userID, lesson.CourseID)
	if err != nil && !errors.Is(err, serviceErrs.ErrorSelectEmpty) {
		return nil, fmt.Errorf("failed to get course progression: %w", err)
	}

	// get all lesson IDs ordered by serial number
	lessonIDs, err := s.repo.GetAllLessonIDsByCourseIDOrderedBySerial(ctx, lesson.CourseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lesson IDs: %w", err)
	}

	// find requested lesson index
	requestedLessonIndex := -1
	for i, id := range lessonIDs {
		if id == lessonID {
			requestedLessonIndex = i
			break
		}
	}

	if requestedLessonIndex == -1 {
		return nil, serviceErrs.ErrorLessonNotFound
	}

	// if no progression, user can only access first lesson
	if progression == nil {
		if requestedLessonIndex != 0 {
			return nil, serviceErrs.ErrorForbidden
		}
	} else {
		// find current lesson index
		currentLessonIndex := -1
		for i, id := range lessonIDs {
			if id == progression.CurrentLessonID {
				currentLessonIndex = i
				break
			}
		}

		if currentLessonIndex == -1 {
			return nil, fmt.Errorf("current lesson not found in course")
		}

		// user can access lessons up to current + 1
		// if requestedLessonIndex > currentLessonIndex + 1, deny access
		if requestedLessonIndex > currentLessonIndex+1 {
			return nil, serviceErrs.ErrorForbidden
		}
	}
	// }

	// get quizzes with options
	quizzes, err := s.repo.GetQuizzesWithOptionsByLessonID(ctx, lessonID)
	if err != nil && !errors.Is(err, serviceErrs.ErrorSelectEmpty) {
		return nil, fmt.Errorf("failed to get quizzes: %w", err)
	}
	lesson.Quizzes = quizzes

	return lesson, nil
}
