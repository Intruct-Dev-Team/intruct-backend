package course

import (
	"context"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
)

type Repository interface {
	IsUserExistsByID(ctx context.Context, id int) (bool, error)
	GetLanguageIDByName(ctx context.Context, languageName string) (int, error)
	IsCourseExistsByOwnerAndTitle(ctx context.Context, ownerID int, title string) (bool, error)
	GetCourseByID(ctx context.Context, id int) (*entities.Course, error)
	CreateCourse(ctx context.Context, course *entities.Course) (int, error)
	ImplementCourse(ctx context.Context, course *entities.Course, nextStateID int) error

	GetStateIDByName(ctx context.Context, name entities.StateName) (int, error)
	GetStateMachineItemByID(ctx context.Context, id int) (*entities.StateMachineItem, error)
	CheckIsTransitionAvailable(ctx context.Context, stateMachineID, currentStateID, nextStateID int) (bool, error)
	UpdateStateMachineItemState(ctx context.Context, stateMachineItemID, newStateID int) error
	SetIsPublicField(ctx context.Context, courseID int, isPublic bool) error
}

type CourseService struct {
	repo Repository
}

func NewService(repo Repository) *CourseService {
	return &CourseService{
		repo: repo,
	}
}
