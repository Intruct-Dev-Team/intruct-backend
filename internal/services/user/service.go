package user

import (
	"context"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	"github.com/Intruct-Dev-Team/intruct-backend/pkg/storage/object"
)

type ObjectStorage interface {
	UploadFile(ctx context.Context, file *object.File) error
}

type Repository interface {
	GetUserIDByExternalUUID(ctx context.Context, ExternalUUID string) (int, error)
	IsUserExistsByEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user *entities.User) (int, error)
	GetUserByID(ctx context.Context, id int) (*entities.User, error)
	GetUserStatisticByUserID(ctx context.Context, id int) (*entities.UserStatistic, error)
	// IsTeacherExistsByUserID(ctx context.Context, id int) (bool, error)
	// GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	// UpdateUser(ctx context.Context, userID int, user *entities.User) error
}

type UserService struct {
	repo          Repository
	objectStorage ObjectStorage
}

func NewService(repo Repository, objectStorage ObjectStorage) *UserService {
	return &UserService{
		repo:          repo,
		objectStorage: objectStorage,
	}
}
