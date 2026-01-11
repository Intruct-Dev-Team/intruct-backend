package course

import (
	"errors"
	"net/http"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	contextKeys "github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/context"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/httputils"
	"go.uber.org/zap"
)

const CreateCourseRoute = ""

// CreateCourse creates a new course
// @Summary      Create course
// @Description  Create a new course with title, description, file and language
// @Description  Prerequisites: User must be authenticated
// @Tags         course
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        title       formData string true  "Course title"
// @Param        description formData string false  "Course description"
// @Param        file        formData file   true  "Course file"
// @Param        language    formData string true  "Course language"
// @Success      201 {object} createCourseResponse
// @Failure      400 {object} httputils.ErrorStruct "Invalid request or validation failed"
// @Failure      401 {object} httputils.ErrorStruct "Authentication required"
// @Failure      500 {object} httputils.ErrorStruct
// @Router       /course [post]
func (h *CourseHandlers) CreateCourse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(contextKeys.UserIDKey).(int)
		if !ok || userID == 0 {
			h.log.Debug("user id not found in context")
			httputils.RespondWith401(w, "authentication required", h.log)

			return
		}

		// Parse multipart form with 32MB max memory
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			h.log.Debug("failed to parse multipart form", zap.Error(err))
			httputils.RespondWith400(w, "failed to parse form data", h.log)

			return
		}

		title := r.FormValue("title")
		if title == "" {
			httputils.RespondWith400(w, "title is required", h.log)

			return
		}

		description := r.FormValue("description")

		language := r.FormValue("language")
		if language == "" {
			httputils.RespondWith400(w, "language is required", h.log)

			return
		}

		file, _, err := r.FormFile("file") //fileHeader
		if err != nil {
			h.log.Debug("failed to get file from form", zap.Error(err))
			httputils.RespondWith400(w, "file is required", h.log)

			return
		}
		defer func() {
			if err := file.Close(); err != nil {
				h.log.Error("failed to close file", zap.Error(err))
			}
		}()

		course := &entities.Course{
			OwnerID:     userID,
			Title:       title,
			Description: description,
			Language:    language,
		}

		courseID, err := h.courseService.CreateCourse(r.Context(), course)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrs.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrs.ErrorLanguageNotFound):
				httputils.RespondWith400(w, "invalid language", h.log)
			case errors.Is(err, serviceErrs.ErrorCourseExists):
				httputils.RespondWith409(w, err.Error(), h.log)

			default:
				h.log.Error("failed to create course", zap.Error(err))
				httputils.RespondWith500(w, h.log)
			}

			return
		}

		// n8n call
		//file, fileHeader.Size
		//TODO: n8n connection
		h.log.Warn("N8N REALIZATION MISSING")

		response := createCourseResponse{
			CourseID: courseID,
		}

		httputils.RespondWith201(w, response, h.log)
	}
}

// @Description createCourseResponse returns created course ID
type createCourseResponse struct {
	CourseID int `json:"course_id" example:"1"`
}
