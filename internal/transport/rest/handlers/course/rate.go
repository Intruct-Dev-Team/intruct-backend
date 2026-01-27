package course

import (
	"encoding/json"
	"errors"
	"net/http"

	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	contextKeys "github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/context"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/httputils"
	"go.uber.org/zap"
)

const rateCourseRoute = "/{course_id}/rate"

// RateCourse  Add rating to a course
// @Summary      Rate course
// @Description  Add rating to a course form 1 to 5 grate
// @Description  Prerequisites: User must be authenticated and must learn the course
// @Description  Course must exist
// @Tags         course
// @Produce      json
// @Security     BearerAuth
// @Param        course_id path int true "Course ID"
// @Param        addRateRequest body addRateRequest true "rating"
// @Success      201
// @Failure      400 {object} httputils.ErrorStruct "Invalid course ID or body"
// @Failure      401 {object} httputils.ErrorStruct "Authentication required"
// @Failure      403 {object} httputils.ErrorStruct "User is not the 'student' of the course"
// @Failure      404 {object} httputils.ErrorStruct "Course not found"
// @Failure      500 {object} httputils.ErrorStruct
// @Router       /courses/{course_id}/rate [post]
func (h *CourseHandlers) RateCourse() http.HandlerFunc {
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

		var req addRateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Debug("failed to decode body", zap.Error(err))
			httputils.RespondWith400(w, "failed to decode body", h.log)
			return
		}

		if req.Rating < 1 || req.Rating > 5 {
			h.log.Debug("rating must be in range form 1 to 5", zap.Int("rate", req.Rating))
			httputils.RespondWith400(w, "rating must be in range form 1 to 5", h.log)
			return
		}

		err = h.courseService.RateCourse(r.Context(), courseID, userID, req.Rating)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrs.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorCourseNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorUserHasNoCourseProgression):
				httputils.RespondWith403(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorRatingExists):
				httputils.RespondWith409(w, err.Error(), h.log)
			default:
				h.log.Error("failed to delete course", zap.Error(err))
				httputils.RespondWith500(w, h.log)
			}
			return
		}

		httputils.RespondWith201(w, struct{}{}, h.log)
	}
}

// @Description addRateRequest body of request
type addRateRequest struct {
	Rating int `json:"rating"        example:"5"            binding:"required"`
}
