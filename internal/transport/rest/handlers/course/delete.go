package course

import (
	"errors"
	"net/http"

	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	contextKeys "github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/context"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/httputils"
	"go.uber.org/zap"
)

const deleteCourseRoute = "/{course_id}/delete"

// DeleteCourse  Delete a course
// @Summary      Delete course
// @Description  Delete a course by removing all content (modules, lessons, quizzes) and soft delete course from db
// @Description  Prerequisites: User must be authenticated and must be the owner of the course
// @Description  Course must exist
// @Tags         course
// @Produce      json
// @Security     BearerAuth
// @Param        course_id path int true "Course ID"
// @Success      200
// @Failure      400 {object} httputils.ErrorStruct "Invalid course ID"
// @Failure      401 {object} httputils.ErrorStruct "Authentication required"
// @Failure      403 {object} httputils.ErrorStruct "User is not the owner of the course"
// @Failure      404 {object} httputils.ErrorStruct "Course not found"
// @Failure      500 {object} httputils.ErrorStruct
// @Router       /courses/{course_id}/delete [delete]
func (h *CourseHandlers) DeleteCourse() http.HandlerFunc {
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

		err = h.courseService.DeleteCourse(r.Context(), courseID, userID)
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

		httputils.RespondWith200(w, struct{}{}, h.log)
	}
}
