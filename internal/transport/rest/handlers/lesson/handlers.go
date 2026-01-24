package lesson

import (
	"context"
	"net/http"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	lessonsRoute = "/lessons"
)

type LessonService interface {
	GetLesson(ctx context.Context, lessonID int, userID int) (*entities.Lesson, error)
}

type LessonHandlers struct {
	lessonService LessonService
	log           *zap.Logger
}

func NewLessonHandlers(lessonService LessonService, log *zap.Logger) *LessonHandlers {
	return &LessonHandlers{
		lessonService: lessonService,
		log:           log,
	}
}

func (h *LessonHandlers) SetupLessonRoutes(router *chi.Mux, jwtAuthMiddleware func(http.Handler) http.Handler, systemAuthMiddleware func(http.Handler) http.Handler) {
	lessonsRouter := chi.NewRouter()
	lessonsRouter.Group(func(r chi.Router) {
		r.Use(jwtAuthMiddleware)
		r.Use(systemAuthMiddleware)
		r.Get(getLessonRoute, h.GetLesson())
	})

	router.Mount(lessonsRoute, lessonsRouter)
}
