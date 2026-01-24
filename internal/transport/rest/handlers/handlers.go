package handlers

import (
	"net/http"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/handlers/course"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/handlers/language"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/handlers/lesson"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/handlers/user"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Services interface {
	user.UserService
	language.LanguageService
	course.CourseService
	course.N8NService
	lesson.LessonService
}

type Handlers struct {
	services Services
	log      *zap.Logger
}

func New(services Services, log *zap.Logger) *Handlers {
	return &Handlers{
		services: services,
		log:      log,
	}
}

func (h *Handlers) SetupRoutes(router *chi.Mux,
	jwtAuthMiddleware func(http.Handler) http.Handler,
	systemAuthMiddleware func(http.Handler) http.Handler) {
	//recomendation from AI about downcast to certain interfaces (ISP)

	var userService user.UserService = h.services
	userHandlers := user.NewUserHandlers(userService, h.log)
	userHandlers.SetupUserRoutes(router, jwtAuthMiddleware, systemAuthMiddleware)

	var languageService language.LanguageService = h.services
	languageHandlers := language.NewLanguageHandlers(languageService, h.log)
	languageHandlers.SetupLanguageRoutes(router)

	var courseService course.CourseService = h.services
	var n8nService course.N8NService = h.services
	courseHandlers := course.NewCourseHandlers(courseService, n8nService, h.log)
	courseHandlers.SetupCourseRoutes(router, jwtAuthMiddleware, systemAuthMiddleware)

	var lessonService lesson.LessonService = h.services
	lessonHandlers := lesson.NewLessonHandlers(lessonService, h.log)
	lessonHandlers.SetupLessonRoutes(router, jwtAuthMiddleware, systemAuthMiddleware)

	// var teacherService teacher.TeacherService = h.services
	// teacherHandlers := teacher.NewTeacherHandlers(teacherService, h.log)
	// teacherHandlers.SetupTeacherRoutes(router, middlwares.ClassicAuth)

	// var scheduleService schedule.ScheduleService = h.services
	// scheduleHandlers := schedule.NewScheduleHandlers(scheduleService, h.log)
	// scheduleHandlers.SetupScheduleRoutes(router, middlwares.ClassicAuth)

	// var reviewService review.ReviewService = h.services
	// reviewHandlers := review.NewReviewHandlers(reviewService, h.log)
	// reviewHandlers.SetupReviewRoutes(router, middlwares.ClassicAuth)

	// var lessonService lesson.LessonService = h.services
	// lessonHandlers := lesson.NewLessonHandlers(lessonService, h.log)
	// lessonHandlers.SetupLessonRoutes(router, middlwares.ClassicAuth)

	// var imageService image.ImageService = h.services
	// imageHandlers := image.NewImageHandlers(imageService, h.log)
	// imageHandlers.SetupImageRoutes(router)

	// var categoryService category.CategoryService = h.services
	// categoryHandlers := category.NewCategoryHandlers(categoryService, h.log)
	// categoryHandlers.SetupCategoryRoutes(router)

	// var complaintService complaint.ComplaintService = h.services
	// complaintHandlers := complaint.NewComplaintHandlers(complaintService, h.log)
	// complaintHandlers.SetupComplaintRoutes(router, middlwares.ClassicAuth)

	// var adminService admin.AdminService = h.services
	// adminHandlers := admin.NewAdminHandlers(adminService, h.log)
	// adminHandlers.SetupAdminRoutes(router, middlwares.ClassicAuth)
}
