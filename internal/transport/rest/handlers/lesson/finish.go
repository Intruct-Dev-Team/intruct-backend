package lesson

import (
	"errors"
	"net/http"

	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	contextKeys "github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/context"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/httputils"
	"go.uber.org/zap"
)

const finishLessonRoute = "/{lesson_id}/finish"

// FinishLesson marks lesson as completed and updates user progress
// @Summary      Finish lesson
// @Description  Mark lesson as completed and advance user progress in the course
// @Description  Prerequisites: User must be authenticated. User must have access to the lesson and not skipping lessons
// @Tags         lesson
// @Produce      json
// @Security     BearerAuth
// @Param        lesson_id path int true "Lesson ID"
// @Success      200
// @Failure      400 {object} httputils.ErrorStruct "Invalid lesson ID"
// @Failure      401 {object} httputils.ErrorStruct "Authentication required or user not registered"
// @Failure      403 {object} httputils.ErrorStruct "Access forbidden - lesson locked or course access denied"
// @Failure      404 {object} httputils.ErrorStruct "Lesson not found"
// @Failure      500 {object} httputils.ErrorStruct
// @Router       /lessons/{lesson_id}/finish [put]
func (h *LessonHandlers) FinishLesson() http.HandlerFunc {
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

		err = h.lessonService.FinishLesson(r.Context(), lessonID, userID)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrs.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorLessonNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorLessonNotReached):
				httputils.RespondWith403(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorLessonFinished):
				httputils.RespondWith403(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorForbidden):
				httputils.RespondWith403(w, err.Error(), h.log)
			default:
				h.log.Error("failed to finish lesson", zap.Int("lesson_id", lessonID), zap.Error(err))
				httputils.RespondWith500(w, h.log)
			}
			return
		}

		httputils.RespondWith200(w, struct{}{}, h.log)
	}
}
