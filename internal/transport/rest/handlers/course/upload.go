package course

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/httputils"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const UploadCourseRoute = "/{course_id}/upload"

// UploadCourseContent uploads processed course content from n8n
// @Summary      Upload course content
// @Description  Upload processed course content with modules, lessons and quizzes from n8n workflow
// @Description  This endpoint fills the course with educational content after file processing
// @Tags         course
// @Accept       json
// @Produce      json
// @Param        course_id path int true "Course ID"
// @Param        request body uploadCourseRequest true "Course Content Data"
// @Success      200
// @Failure      400 {object} httputils.ErrorStruct "Invalid request or validation failed"
// @Failure      404 {object} httputils.ErrorStruct "Course not found"
// @Failure      500 {object} httputils.ErrorStruct
// @Router       /courses/{course_id}/upload [patch]
func (h *CourseHandlers) UploadCourseContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		courseID, err := httputils.GetIntParamFromRequestPath(r, "course_id")
		if err != nil {
			h.log.Debug("failed to parse id from URL path", zap.Error(err))
			httputils.RespondWith400(w, "incorrect or missed {course_id} param in url path", h.log)

			return
		}

		var req uploadCourseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputils.RespondWith400(w, "failed to decode body", h.log)

			return
		}

		if err := validateUploadRequest(&req); err != nil {
			httputils.RespondWith400(w, err.Error(), h.log)

			return
		}

		modules := make([]*entities.Module, 0, len(req.Modules))
		for moduleIdx, modReq := range req.Modules {
			lessons := make([]*entities.Lesson, 0, len(modReq.Lessons))

			for lessonIdx, lessonReq := range modReq.Lessons {
				quizzes := make([]*entities.Quiz, 0, len(lessonReq.Quiz))

				for quizIdx, quizReq := range lessonReq.Quiz {
					options := make([]*entities.QuizOption, 0, len(quizReq.Options))

					for optionIdx, optionText := range quizReq.Options {
						option := &entities.QuizOption{
							Content:  optionText,
							IsAnswer: optionIdx == quizReq.CorrectIndex,
						}
						options = append(options, option)
					}

					quiz := &entities.Quiz{
						SerialNumber: quizIdx + 1,
						Question:     quizReq.Question,
						Options:      options,
					}
					quizzes = append(quizzes, quiz)
				}

				lesson := &entities.Lesson{
					SerialNumber: lessonIdx + 1,
					Title:        lessonReq.LessonTitle,
					Content:      lessonReq.Content,
					Quizzes:      quizzes,
				}
				lessons = append(lessons, lesson)
			}

			module := &entities.Module{
				SerialNumber: moduleIdx + 1,
				Title:        modReq.ModuleTitle,
				Lessons:      lessons,
			}
			modules = append(modules, module)
		}

		course := &entities.Course{
			ID:       courseID,
			Title:    req.CourseTitle,
			Language: req.Language,
			Modules:  modules,
		}

		err = h.courseService.UploadCourseContent(r.Context(), course)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrs.ErrorCourseNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorUnavailableStateTransition):
				httputils.RespondWith403(w, "course is not 'in process' state", h.log)
			default:
				h.log.Error("failed to upload course content", zap.Error(err))
				httputils.RespondWith500(w, h.log)
			}

			return
		}

		httputils.RespondWith200(w, struct{}{}, h.log)
	}
}

func validateUploadRequest(req *uploadCourseRequest) error {
	if req.CourseTitle == "" {
		return errors.New("course_title is required")
	}

	if req.Language == "" {
		return errors.New("language is required")
	}

	if req.TotalModules <= 0 {
		return errors.New("total_modules must be positive")
	}

	if req.TotalLessons <= 0 {
		return errors.New("total_lessons must be positive")
	}

	if len(req.Modules) == 0 {
		return errors.New("modules are required")
	}

	if len(req.Modules) != req.TotalModules {
		return errors.New("modules count doesn't match total_modules")
	}

	totalLessonsCount := 0
	var eg errgroup.Group

	for _, mod := range req.Modules {
		if mod.ModuleID <= 0 {
			return errors.New("invalid module_id")
		}

		if mod.ModuleTitle == "" {
			return errors.New("module_title is required")
		}

		if len(mod.Lessons) == 0 {
			return errors.New("module must contain at least one lesson")
		}

		totalLessonsCount += len(mod.Lessons)

		for _, lesson := range mod.Lessons {
			lesson := lesson // capture for goroutine
			eg.Go(func() error {
				if lesson.LessonID <= 0 {
					return errors.New("invalid lesson_id")
				}

				if lesson.LessonTitle == "" {
					return errors.New("lesson_title is required")
				}

				if lesson.Content == "" {
					return errors.New("lesson content is required")
				}

				for _, quiz := range lesson.Quiz {
					if quiz.Question == "" {
						return errors.New("quiz question is required")
					}

					if len(quiz.Options) < 2 {
						return errors.New("quiz must have at least 2 options")
					}

					if quiz.CorrectIndex < 0 || quiz.CorrectIndex >= len(quiz.Options) {
						return errors.New("quiz correct_index is out of bounds")
					}
				}

				return nil
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	if totalLessonsCount != req.TotalLessons {
		return errors.New("total lessons count doesn't match total_lessons")
	}

	return nil
}

// @Description uploadCourseRequest contains processed course content from n8n
type uploadCourseRequest struct {
	CourseTitle  string                `json:"course_title"  example:"Introduction to Programming"   binding:"required"`
	Language     string                `json:"language"      example:"Russian"                       binding:"required"`
	TotalModules int                   `json:"total_modules" example:"10"                            binding:"required"`
	TotalLessons int                   `json:"total_lessons" example:"60"                            binding:"required"`
	LastUpdated  time.Time             `json:"last_updated"  example:"2026-01-06T16:57:59.272Z"      binding:"required"`
	Modules      []uploadModuleRequest `json:"modules"       binding:"required"`
}

type uploadModuleRequest struct {
	ModuleID    int                   `json:"module_id"    example:"1"         binding:"required"`
	ModuleTitle string                `json:"module_title" example:"Module 1"  binding:"required"`
	Lessons     []uploadLessonRequest `json:"lessons"      binding:"required"`
}

type uploadLessonRequest struct {
	LessonID    int                 `json:"lesson_id"    example:"1"          binding:"required"`
	LessonTitle string              `json:"lesson_title" example:"Lesson 1.1" binding:"required"`
	Content     string              `json:"content"      example:"..."        binding:"required"`
	Quiz        []uploadQuizRequest `json:"quiz"`
}

type uploadQuizRequest struct {
	Question     string   `json:"question"      example:"What is 2+2?" binding:"required"`
	Options      []string `json:"options"       example:"1,2,3,4"      binding:"required"`
	CorrectIndex int      `json:"correct_index" example:"1"            binding:"required"`
}
