package course

import (
	"context"
	"net/http"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	courseRoute  = "/course"
	coursesRoute = "/courses"
)

type CourseService interface {
	CreateCourse(ctx context.Context, course *entities.Course) (int, error)
	UploadCourseContent(ctx context.Context, course *entities.Course) error
}

type CourseHandlers struct {
	courseService CourseService
	log           *zap.Logger
}

func NewCourseHandlers(courseService CourseService, log *zap.Logger) *CourseHandlers {
	return &CourseHandlers{
		courseService: courseService,
		log:           log,
	}
}

func (h *CourseHandlers) SetupCourseRoutes(router *chi.Mux, jwtAuthMiddleware func(http.Handler) http.Handler, systemAuthMiddleware func(http.Handler) http.Handler) {
	router.Group(func(r chi.Router) {
		r.Use(jwtAuthMiddleware)
		r.Use(systemAuthMiddleware)

		r.Post(courseRoute, h.CreateCourse())
	})

	coursesRouter := chi.NewRouter()
	coursesRouter.Patch(UploadCourseRoute, h.UploadCourseContent())
	router.Mount(coursesRoute, coursesRouter)
}
