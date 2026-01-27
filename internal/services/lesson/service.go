package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
)

type Repository interface {
	IsUserExistsByID(ctx context.Context, id int) (bool, error)
	GetLessonByID(ctx context.Context, lessonID int) (*entities.Lesson, error)
	GetQuizzesWithOptionsByLessonID(ctx context.Context, lessonID int) ([]*entities.Quiz, error)
	GetCourseByID(ctx context.Context, courseID int) (*entities.Course, error)
	GetCourseProgressionByUserAndCourse(ctx context.Context, userID int, courseID int) (*entities.CourseProgression, error)
	GetAllLessonIDsByCourseIDOrderedBySerial(ctx context.Context, courseID int) ([]int, error)
	CreateCourseProgression(ctx context.Context, userID, courseID, currentLessonID int) error
	UpdateCourseProgression(ctx context.Context, userID int, courseID int, newCurrentLessonID int, finishedLessonsCount int, isFinished bool) error
}

type LessonService struct {
	repo Repository
}

func NewService(repo Repository) *LessonService {
	return &LessonService{
		repo: repo,
	}
}

func (s *LessonService) checkUserAccess(ctx context.Context, lessonID, userID int) (*entities.Lesson, error) {
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

	return lesson, nil
}
