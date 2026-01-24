package lesson

import (
	"context"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
)

type Repository interface {
	IsUserExistsByID(ctx context.Context, id int) (bool, error)
	GetLessonByID(ctx context.Context, lessonID int) (*entities.Lesson, error)
	GetQuizzesWithOptionsByLessonID(ctx context.Context, lessonID int) ([]*entities.Quiz, error)
	GetCourseByID(ctx context.Context, courseID int) (*entities.Course, error)
	GetCourseProgressionByUserAndCourse(ctx context.Context, userID int, courseID int) (*entities.CourseProgression, error)
	GetAllLessonIDsByCourseIDOrderedBySerial(ctx context.Context, courseID int) ([]int, error)
}

type LessonService struct {
	repo Repository
}

func NewService(repo Repository) *LessonService {
	return &LessonService{
		repo: repo,
	}
}
