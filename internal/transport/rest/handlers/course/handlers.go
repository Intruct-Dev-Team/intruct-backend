package course

import (
	"context"
	"io"
	"net/http"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	coursesRoute = "/courses"
)

type CourseService interface {
	CreateCourse(ctx context.Context, course *entities.Course) (int, error)
	UploadCourseContent(ctx context.Context, course *entities.Course) error
	PublishCourse(ctx context.Context, courseID int, userID int) error
	GetCourseList(ctx context.Context, userID int, inMine bool) ([]*entities.Course, error)
	GetCourse(ctx context.Context, courseID int, userID int) (*entities.Course, error)
	DeleteCourse(ctx context.Context, courseID int, userID int) error
}

type N8NService interface {
	SendCourse(ctx context.Context, course *entities.Course, fileReader io.Reader, fileSize int64, fileName string) error
}

type CourseHandlers struct {
	courseService CourseService
	n8nService    N8NService
	log           *zap.Logger
}

func NewCourseHandlers(courseService CourseService, n8nService N8NService, log *zap.Logger) *CourseHandlers {
	return &CourseHandlers{
		courseService: courseService,
		n8nService:    n8nService,
		log:           log,
	}
}

func (h *CourseHandlers) SetupCourseRoutes(router *chi.Mux, jwtAuthMiddleware func(http.Handler) http.Handler, systemAuthMiddleware func(http.Handler) http.Handler) {
	router.Group(func(r chi.Router) {
		r.Use(jwtAuthMiddleware)
		r.Use(systemAuthMiddleware)
		r.Post(CreateCourseRoute, h.CreateCourse())
	})

	coursesRouter := chi.NewRouter()
	coursesRouter.Group(func(r chi.Router) {
		r.Use(jwtAuthMiddleware)
		r.Use(systemAuthMiddleware)
		r.Get(listCoursesRoute, h.ListCourses())
		r.Get(getCourseRoute, h.GetCourse())
		r.Put(publishCourseRoute, h.PublishCourse())
		r.Delete(deleteCourseRoute, h.DeleteCourse())
	})
	coursesRouter.Patch(UploadCourseRoute, h.UploadCourseContent())

	router.Mount(coursesRoute, coursesRouter)
}
