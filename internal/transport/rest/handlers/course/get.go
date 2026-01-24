//transport/course/get.go

package course

import (
	"errors"
	"net/http"
	"time"

	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	contextKeys "github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/context"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/httputils"
	"go.uber.org/zap"
)

const getCourseRoute = "/{course_id}"

// GetCourse returns course details with modules and lessons
// @Summary      Get course
// @Description  Get course details including modules and lessons structure
// @Description  Prerequisites: User must be authenticated. Course must be published or user must be the owner
// @Tags         course
// @Produce      json
// @Security     BearerAuth
// @Param        course_id path int true "Course ID"
// @Success      200 {object} getCourseResponse
// @Failure      400 {object} httputils.ErrorStruct "Invalid course ID"
// @Failure      401 {object} httputils.ErrorStruct "Authentication required or user not registered"
// @Failure      403 {object} httputils.ErrorStruct "Access forbidden - course not published"
// @Failure      404 {object} httputils.ErrorStruct "Course not found"
// @Failure      500 {object} httputils.ErrorStruct
// @Router       /courses/{course_id} [get]
func (h *CourseHandlers) GetCourse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(contextKeys.UserIDKey).(int)
		if !ok || userID == 0 {
			h.log.Debug("user id not found in context")
			httputils.RespondWith401(w, "authentication required", h.log)
			return
		}

		courseID, err := httputils.GetIntParamFromRequestPath(r, "course_id")
		if err != nil {
			h.log.Debug("failed to parse id from URL path", zap.Error(err))
			httputils.RespondWith400(w, "incorrect or missed {course_id} param in url path", h.log)
			return
		}

		course, err := h.courseService.GetCourse(r.Context(), courseID, userID)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrs.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorCourseNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorForbidden):
				httputils.RespondWith403(w, err.Error(), h.log)
			default:
				h.log.Error("failed to get course", zap.Int("course_id", courseID), zap.Error(err))
				httputils.RespondWith500(w, h.log)
			}
			return
		}

		var finishedLessonsCount int
		var isInMine bool
		if course.CourseProgression != nil {
			finishedLessonsCount = course.CourseProgression.FinishedLessonsCount
			isInMine = true
		}

		response := getCourseResponse{
			ID:              course.ID,
			Title:           course.Title,
			Description:     course.Description,
			LessonsNumber:   course.LessonsCount,
			FinishedLessons: finishedLessonsCount,
			CreatedAt:       course.CreatedAt,
			UpdatedAt:       course.UpdatedAt,
			AuthorID:        course.OwnerID,
			IsPublic:        course.IsPublic,
			IsMine:          course.OwnerID == userID,
			IsInMine:        isInMine,
			Modules:         make([]moduleItem, 0, len(course.Modules)),
		}

		for _, module := range course.Modules {
			moduleItem := moduleItem{
				ID:           module.ID,
				SerialNumber: module.SerialNumber,
				Title:        module.Title,
				Description:  module.Description,
				CreatedAt:    module.CreatedAt,
				UpdatedAt:    module.UpdatedAt,
				Lessons:      make([]lessonItem, 0, len(module.Lessons)),
			}

			for _, lesson := range module.Lessons {
				moduleItem.Lessons = append(moduleItem.Lessons, lessonItem{
					ID:           lesson.ID,
					SerialNumber: lesson.SerialNumber,
					Title:        lesson.Title,
					Description:  lesson.Description,
					CreatedAt:    lesson.CreatedAt,
					UpdatedAt:    lesson.UpdatedAt,
				})
			}

			response.Modules = append(response.Modules, moduleItem)
		}

		httputils.RespondWith200(w, response, h.log)
	}
}

// @Description getCourseResponse returns course details with modules and lessons
type getCourseResponse struct {
	ID              int          `json:"id"               example:"1"`
	Title           string       `json:"title"            example:"Introduction to Go"`
	Description     string       `json:"description"      example:"Learn Go programming language"`
	LessonsNumber   int          `json:"lessons_number"   example:"10"`
	FinishedLessons int          `json:"finished_lessons" example:"5"`
	CreatedAt       time.Time    `json:"created_at"       example:"2022-09-09T10:10:10Z"`
	UpdatedAt       time.Time    `json:"updated_at"       example:"2022-09-09T10:10:10Z"`
	AuthorID        int          `json:"author_id"        example:"1"`
	IsPublic        bool         `json:"is_public"        example:"true"`
	IsMine          bool         `json:"is_mine"          example:"false"`
	IsInMine        bool         `json:"is_in_mine"       example:"false"`
	Modules         []moduleItem `json:"modules"`
}

// @Description moduleItem represents a module in the course
type moduleItem struct {
	ID           int          `json:"id"            example:"1"`
	SerialNumber int          `json:"serial_number" example:"1"`
	Title        string       `json:"title"         example:"Getting Started"`
	Description  string       `json:"description"   example:"Introduction to the basics"`
	CreatedAt    time.Time    `json:"created_at"    example:"2022-09-09T10:10:10Z"`
	UpdatedAt    time.Time    `json:"updated_at"    example:"2022-09-09T10:10:10Z"`
	Lessons      []lessonItem `json:"lessons"`
}

// @Description lessonItem represents a lesson in the module
type lessonItem struct {
	ID           int       `json:"id"            example:"1"`
	SerialNumber int       `json:"serial_number" example:"1"`
	Title        string    `json:"title"         example:"Variables and Types"`
	Description  string    `json:"description"   example:"Learn about basic data types"`
	CreatedAt    time.Time `json:"created_at"    example:"2022-09-09T10:10:10Z"`
	UpdatedAt    time.Time `json:"updated_at"    example:"2022-09-09T10:10:10Z"`
}
