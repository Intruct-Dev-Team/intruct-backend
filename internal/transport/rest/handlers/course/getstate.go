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

const getCourseStateRoute = "/{course_id}/state"

// DeleteCourse  Return a course state
// @Summary      Get course state
// @Description  Get a course state return course state
// @Description  Prerequisites: User must be authenticated and must be the owner of the course
// @Description  Course must exist
// @Tags         course
// @Produce      json
// @Security     BearerAuth
// @Param        course_id path int true "Course ID"
// @Success      200 {object} getCourseStateResponse
// @Failure      400 {object} httputils.ErrorStruct "Invalid course ID"
// @Failure      401 {object} httputils.ErrorStruct "Authentication required"
// @Failure      403 {object} httputils.ErrorStruct "User is not the owner of the course"
// @Failure      404 {object} httputils.ErrorStruct "Course not found"
// @Failure      500 {object} httputils.ErrorStruct
// @Router       /courses/{course_id}/state [get]
func (h *CourseHandlers) GetCourseState() http.HandlerFunc {
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

		state, updated, err := h.courseService.GetCourseStateAndUpdatedTime(r.Context(), courseID, userID)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrs.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorCourseNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorNotCourseOwner):
				httputils.RespondWith403(w, err.Error(), h.log)
			default:
				h.log.Error("failed to delete course", zap.Error(err))
				httputils.RespondWith500(w, h.log)
			}
			return
		}

		httputils.RespondWith200(w, getCourseStateResponse{
			State:     state,
			UpdatedAt: updated,
		}, h.log)
	}
}

// @Description getCourseStateResponse returns course state
type getCourseStateResponse struct {
	State     string    `json:"state"         example:"in creation"`
	UpdatedAt time.Time `json:"updated_at"    example:"2022-09-09T10:10:10Z"`
}
