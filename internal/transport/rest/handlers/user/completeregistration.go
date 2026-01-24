package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/imgutils"
	contextKeys "github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/context"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/httputils"
	"go.uber.org/zap"
)

const compliteRegistrationRoute = "/complete-registration"

// CompleteRegistration completes user registration after External auth
// @Summary      Complete registration
// @Description  Complete registration by providing additional profile information
// @Description  Prerequisites: User must be authenticated via External service (Supabase Auth)
// @Description  This endpoint creates user profile in application database
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body completeRegistrationRequest true "Profile Data"
// @Success      201
// @Failure      400 {object} httputils.ErrorStruct "Invalid request or validation failed"
// @Failure      401 {object} httputils.ErrorStruct "External auth is missing or invalid"
// @Failure      409 {object} httputils.ErrorStruct "User already completed registration"
// @Failure      500 {object} httputils.ErrorStruct
// @Router       /auth/complete-registration [post]
func (h *UserHandlers) CompleteRegistration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		externalUserUUID, ok := r.Context().Value(contextKeys.ExternalUserUUIDKey).(string)
		if !ok || externalUserUUID == "" {
			h.log.Debug("external user uuid not found in context")
			httputils.RespondWith401(w, "authentication in sso required", h.log)

			return
		}

		email, ok := r.Context().Value(contextKeys.EmailKey).(string)
		if !ok || email == "" {
			h.log.Debug("email not found in context")
			httputils.RespondWith401(w, "authentication in sso required", h.log)

			return
		}

		_, exists, err := h.userService.TryGetUserIDByExternalUUID(r.Context(), externalUserUUID)
		if err != nil {
			h.log.Error("failed to check user existence by external uuid", zap.Error(err))
			httputils.RespondWith500(w, h.log)

			return
		}

		if exists {
			h.log.Debug("user already completed registration",
				zap.String("external_uuid", externalUserUUID),
			)
			httputils.RespondWith409(w, "Registration already completed", h.log)

			return
		}

		var req completeRegistrationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputils.RespondWith400(w, "failed to decode body", h.log)

			return
		}

		if req.Name == "" || req.Surname == "" {
			httputils.RespondWith400(w, "name or surname is empty", h.log)

			return
		}

		if req.Birthdate.Before(time.Date(1910, 01, 01, 0, 0, 0, 0, time.UTC)) {
			httputils.RespondWith400(w, "birthdate is missed or too old", h.log)

			return
		}

		var avatarReader io.Reader
		var avatarSize int64

		// if upload avatar
		if req.Avatar != "" {

			imageBytes, err := imgutils.DecodeImage(req.Avatar)
			if err != nil {
				httputils.RespondWith400(w, err.Error(), h.log)

				return
			}

			width, height, err := imgutils.GetImageDimensions(imageBytes)
			if err != nil {
				h.log.Error("failed to get image dimension", zap.Error(err))
				httputils.RespondWith500(w, h.log)

				return
			}

			if err = imgutils.CheckDimension(1, 1, width, height); err != nil {
				httputils.RespondWith400(w, err.Error(), h.log)

				return
			}
			avatarReader = bytes.NewReader(imageBytes)
			defer func() {
				if closer, ok := avatarReader.(io.Closer); ok {
					if err := closer.Close(); err != nil {
						h.log.Error("failed to close reader", zap.Error(err))
					}
				}
			}()
			avatarSize = int64(len(imageBytes))
		}

		user := &entities.User{
			ExternalUUID: externalUserUUID,
			Email:        email,
			Name:         req.Name,
			Surname:      req.Surname,
			Birthdate:    req.Birthdate,
		}

		_, err = h.userService.CreateUser(r.Context(), user, avatarReader, avatarSize)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrs.ErrorUserExists):
				httputils.RespondWith409(w, err.Error(), h.log)
			default:
				h.log.Error(err.Error())
				httputils.RespondWith500(w, h.log)
			}

			return

		}

		httputils.RespondWith201(w, struct{}{}, h.log)
	}
}

// @Description completeRegistrationRequest to complete user registration.
type completeRegistrationRequest struct {
	Name      string    `json:"name"      example:"John"                 binding:"required"`
	Surname   string    `json:"surname"   example:"Smith"                binding:"required"`
	Birthdate time.Time `json:"birthdate" example:"2000-01-01T00:00:00Z" binding:"required"`
	Avatar    string    `json:"avatar"    example:"base64 encoded image"`
}
