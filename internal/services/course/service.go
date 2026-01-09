package course

import (
	"context"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
)

type Repository interface {
	IsUserExistsByID(ctx context.Context, id int) (bool, error)
	GetLanguageIDByName(ctx context.Context, languageName string) (int, error)
	IsCourseExistsByOwnerAndTitle(ctx context.Context, ownerID int, title string) (bool, error)
	CreateCourse(ctx context.Context, course *entities.Course) (int, error)
	// GetCourseByID(ctx context.Context, id int) (*entities.Course, error)
}

type CourseService struct {
	repo Repository
}

func NewService(repo Repository) *CourseService {
	return &CourseService{
		repo: repo,
	}
}
