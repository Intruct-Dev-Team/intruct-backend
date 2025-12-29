package user

import (
	"context"
	"io"
	"net/http"
	"path"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	authRoute  = "/auth"
	userRoute  = "/user"
	usersRoute = "/users"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entities.User, avatarReader io.Reader, avatarSize int64) (int, error)
	TryGetUserIDByExternalUUID(ctx context.Context, ExternalUserUUID string) (id int, exists bool, err error)
	GetUser(ctx context.Context, userID int) (*entities.User, error)
}

type UserHandlers struct {
	userService UserService
	log         *zap.Logger
}

func NewUserHandlers(userService UserService, log *zap.Logger) *UserHandlers {
	return &UserHandlers{
		userService: userService,
		log:         log,
	}
}

func (h *UserHandlers) SetupUserRoutes(router *chi.Mux, jwtAuthMiddleware func(http.Handler) http.Handler, systemAuthMiddleware func(http.Handler) http.Handler) {
	authRouter := chi.NewRouter()

	authRouter.Group(func(r chi.Router) {
		r.Use(jwtAuthMiddleware)

		r.Post(CompliteRegistrationRoute, h.CompleteRegistration())
	})

	router.Mount(authRoute, authRouter)

	router.Get(path.Join(usersRoute, GetPublicRoute), h.GetUserPublic())

	router.Group(func(r chi.Router) {
		r.Use(jwtAuthMiddleware)
		r.Use(systemAuthMiddleware)

		r.Get(path.Join(userRoute, getProtectedRoute), h.GetUserProtected())
		// r.Get(path.Join(userRoute, checkOnAdminRoute), h.CheckOnAdmin())
		// r.Patch(path.Join(userRoute, editRoute), h.EditUser())
	})
}
