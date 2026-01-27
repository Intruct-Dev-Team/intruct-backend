package course

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	contextKeys "github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/context"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/httputils"
	"go.uber.org/zap"
)

const listCoursesRoute = "/"

// ListCourses returns list of courses
// @Summary      List courses
// @Description  Get list of courses with optional filter for user's courses
// @Description  If 'in_mine' filter is TRUE method will return your created courses + courses which you learn
// @Description  If 'in_mine' filter is FALSE method will return only public courses
// @Description  Prerequisites: User must be authenticated
// @Tags         course
// @Produce      json
// @Security     BearerAuth
// @Param        in_mine query boolean false "Filter user's courses (true/false)"
// @Success      200 {object} listCoursesResponse
// @Failure      400 {object} httputils.ErrorStruct "Invalid query parameter"
// @Failure      401 {object} httputils.ErrorStruct "Authentication required"
// @Failure      500 {object} httputils.ErrorStruct
// @Router       /courses [get]
func (h *CourseHandlers) ListCourses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(contextKeys.UserIDKey).(int)
		if !ok || userID == 0 {
			h.log.Debug("user id not found in context")
			httputils.RespondWith401(w, "authentication required", h.log)
			return
		}

		// get query as filters
		inMineStr := r.URL.Query().Get("in_mine")
		inMine, err := strconv.ParseBool(inMineStr)
		if err != nil {
			inMine = false
		}

		courses, err := h.courseService.GetCourseList(r.Context(), userID, inMine)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrs.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			default:
				h.log.Error("failed to list courses", zap.Error(err))
				httputils.RespondWith500(w, h.log)
			}
			return
		}

		response := listCoursesResponse{
			Courses: make([]courseItem, 0, len(courses)),
		}

		var finishedLessonsCount int
		var isInMine bool
		for _, course := range courses {
			if course.CourseProgression != nil {
				finishedLessonsCount = course.CourseProgression.FinishedLessonsCount
				isInMine = true
			} else {
				finishedLessonsCount = 0
				isInMine = false
			}
			response.Courses = append(response.Courses, courseItem{
				ID:              course.ID,
				Title:           course.Title,
				Description:     course.Description,
				LessonsNumber:   course.LessonsCount,
				FinishedLessons: finishedLessonsCount,
				CreatedAt:       course.CreatedAt,
				UpdatedAt:       course.UpdatedAt,
				AuthorID:        course.OwnerID,
				State:           course.State,
				IsPublic:        course.IsPublic,
				IsMine:          course.OwnerID == userID,
				IsInMine:        isInMine,
				StudentsCount:   course.Statistic.StudentsCount,
				AverageRating:   course.Statistic.AverageRating,
				RatingsCount:    course.Statistic.RatingsCount,
			})
		}

		httputils.RespondWith200(w, response, h.log)
	}
}

// @Description listCoursesResponse returns list of courses
type listCoursesResponse struct {
	Courses []courseItem `json:"courses"`
}

// @Description courseItem represents a course in the list
type courseItem struct {
	ID              int       `json:"id"               example:"1"`
	Title           string    `json:"title"            example:"Introduction to Go"`
	Description     string    `json:"description"      example:"Learn Go programming language"`
	LessonsNumber   int       `json:"lessons_number"   example:"10"`
	FinishedLessons int       `json:"finished_lessons" example:"5"`
	CreatedAt       time.Time `json:"created_at"       example:"2022-09-09T10:10:10Z"`
	UpdatedAt       time.Time `json:"updated_at"       example:"2022-09-09T10:10:10Z"`
	AuthorID        int       `json:"author_id"        example:"1"`
	State           string    `json:"state"            example:"published"`
	IsPublic        bool      `json:"is_public"        example:"true"`
	IsMine          bool      `json:"is_mine"          example:"false"`
	IsInMine        bool      `json:"is_in_mine"       example:"false"`
	StudentsCount   int       `json:"students_count"   example:"0"`
	AverageRating   float64   `json:"average_rating"   example:"5"`
	RatingsCount    int       `json:"ratings_count"    example:"1"`
}
