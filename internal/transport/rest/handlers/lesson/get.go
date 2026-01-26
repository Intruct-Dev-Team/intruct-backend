//transport/lesson/get.go

package lesson

import (
	"errors"
	"net/http"
	"time"

	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	contextKeys "github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/context"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/httputils"
	"go.uber.org/zap"
)

const getLessonRoute = "/{lesson_id}"

// GetLesson returns lesson details with quizzes
// @Summary      Get lesson
// @Description  Get lesson details including content and quizzes with options
// @Description  Prerequisites: User must be authenticated. User must have access to the course and have reached this lesson in progression
// @Tags         lesson
// @Produce      json
// @Security     BearerAuth
// @Param        lesson_id path int true "Lesson ID"
// @Success      200 {object} getLessonResponse
// @Failure      400 {object} httputils.ErrorStruct "Invalid lesson ID"
// @Failure      401 {object} httputils.ErrorStruct "Authentication required or user not registered"
// @Failure      403 {object} httputils.ErrorStruct "Access forbidden - lesson locked or course access denied"
// @Failure      404 {object} httputils.ErrorStruct "Lesson not found"
// @Failure      500 {object} httputils.ErrorStruct
// @Router       /lessons/{lesson_id} [get]
func (h *LessonHandlers) GetLesson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(contextKeys.UserIDKey).(int)
		if !ok || userID == 0 {
			h.log.Debug("user id not found in context")
			httputils.RespondWith401(w, "authentication required", h.log)
			return
		}

		lessonID, err := httputils.GetIntParamFromRequestPath(r, "lesson_id")
		if err != nil {
			h.log.Debug("failed to parse id from URL path", zap.Error(err))
			httputils.RespondWith400(w, "incorrect or missed {lesson_id} param in url path", h.log)
			return
		}

		lesson, err := h.lessonService.GetLesson(r.Context(), lessonID, userID)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrs.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorLessonNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorLessonNotReached):
				httputils.RespondWith403(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorForbidden):
				httputils.RespondWith403(w, err.Error(), h.log)
			default:
				h.log.Error("failed to get lesson", zap.Int("lesson_id", lessonID), zap.Error(err))
				httputils.RespondWith500(w, h.log)
			}
			return
		}

		response := getLessonResponse{
			ID:           lesson.ID,
			CourseID:     lesson.CourseID,
			ModuleID:     lesson.ModuleID,
			SerialNumber: lesson.SerialNumber,
			Title:        lesson.Title,
			Description:  lesson.Description,
			Content:      lesson.Content,
			CreatedAt:    lesson.CreatedAt,
			UpdatedAt:    lesson.UpdatedAt,
			Quizzes:      make([]quizItem, 0, len(lesson.Quizzes)),
		}

		for _, quiz := range lesson.Quizzes {
			quizItem := quizItem{
				ID:           quiz.ID,
				SerialNumber: quiz.SerialNumber,
				Question:     quiz.Question,
				CreatedAt:    quiz.CreatedAt,
				UpdatedAt:    quiz.UpdatedAt,
				Options:      make([]string, 0, len(quiz.Options)),
				CorrectIndex: -1,
			}

			for i, option := range quiz.Options {
				quizItem.Options = append(quizItem.Options, option.Content)
				if option.IsAnswer {
					quizItem.CorrectIndex = i
				}
			}

			response.Quizzes = append(response.Quizzes, quizItem)
		}

		httputils.RespondWith200(w, response, h.log)
	}
}

// @Description getLessonResponse returns lesson details with quizzes and options
type getLessonResponse struct {
	ID           int        `json:"id"            example:"1"`
	CourseID     int        `json:"course_id"     example:"1"`
	ModuleID     int        `json:"module_id"     example:"1"`
	SerialNumber int        `json:"serial_number" example:"1"`
	Title        string     `json:"title"         example:"Variables and Types"`
	Description  string     `json:"description"   example:"Learn about basic data types"`
	Content      string     `json:"content"       example:"# Variables in Go\n\nIn Go, variables are..."`
	CreatedAt    time.Time  `json:"created_at"    example:"2022-09-09T10:10:10Z"`
	UpdatedAt    time.Time  `json:"updated_at"    example:"2022-09-09T10:10:10Z"`
	Quizzes      []quizItem `json:"quizzes"`
}

// @Description quizItem represents a quiz in the lesson
type quizItem struct {
	ID           int       `json:"id"            example:"1"`
	SerialNumber int       `json:"serial_number" example:"1"`
	Question     string    `json:"question"      example:"What is the zero value of an int?"`
	CreatedAt    time.Time `json:"created_at"    example:"2022-09-09T10:10:10Z"`
	UpdatedAt    time.Time `json:"updated_at"    example:"2022-09-09T10:10:10Z"`
	Options      []string  `json:"options"       example:"0,null,undefined,nil"`
	CorrectIndex int       `json:"correct_index" example:"0"`
}
