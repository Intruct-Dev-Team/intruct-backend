package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	contextKeys "github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/context"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/httputils"
	"go.uber.org/zap"
)

const (
	GetPublicRoute    = "/{id}/profile"
	getProtectedRoute = "/profile"
)

// GetUserPublic returns http.HandlerFunc
// @Summary Get user profile
// @Description Get info about user by user id in route (/api/users/{id}/profile)
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} getUserResponse
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /users/{id}/profile [get]
func (h *UserHandlers) GetUserPublic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httputils.GetIntParamFromRequestPath(r, "id")

		if err != nil {
			h.log.Debug("failed to parse id from URL path", zap.Error(err))
			httputils.RespondWith400(w, "missed {id} param in url path", h.log)

			return
		}

		user, err := h.userService.GetUser(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrs.ErrorUserNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)

			default:
				h.log.Error(err.Error())
				httputils.RespondWith500(w, h.log)
			}

			return
		}

		resp := mappingToUserResp(user)

		httputils.RespondWith200(w, resp, h.log)
	}
}

// GetUserProtected returns http.HandlerFunc
// @Summary Get user profile
// @Description Get info about user by jwt token (in Authorization enter: Bearer <your_jwt_token>)
// @Tags users
// @Produce json
// @Security     BearerAuth
// @Success 200 {object} getUserResponse
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /user/profile [get]
func (h *UserHandlers) GetUserProtected() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(contextKeys.UserIDKey).(int)
		if !ok || userID == 0 {
			h.log.Debug("user id not found in context")
			httputils.RespondWith401(w, "authentication required", h.log)

			return
		}

		user, err := h.userService.GetUser(r.Context(), userID)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrs.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)

			default:
				h.log.Error(err.Error())
				httputils.RespondWith500(w, h.log)
			}

			return
		}

		resp := mappingToUserResp(user)

		httputils.RespondWith200(w, resp, h.log)
	}
}

/* helpers */

func mappingToUserResp(user *entities.User) *getUserResponse {
	resp := getUserResponse{
		ID:               user.ID,
		ExternalUUID:     user.ExternalUUID,
		Email:            user.Email,
		Name:             user.Name,
		Surname:          user.Surname,
		RegistrationDate: user.RegistrationDate,
		Birthdate:        user.Birthdate,
		Avatar:           user.Avatar,
	}

	return &resp
}

/* Mapping struct */

type getUserResponse struct {
	ID               int       `json:"id"                   example:"1"`
	ExternalUUID     string    `json:"external_uuid"        example:"uuid"`
	Email            string    `json:"email"                example:"qwerty@example.com"`
	Name             string    `json:"name"                 example:"John"`
	Surname          string    `json:"surname"              example:"Smith"`
	RegistrationDate time.Time `json:"registration_date"    example:"2022-09-09T10:10:10+09:00"`
	Birthdate        time.Time `json:"birthdate"            example:"2002-09-09T10:10:10+09:00"`
	Avatar           string    `json:"avatar"               example:"uuid.png"`
}
