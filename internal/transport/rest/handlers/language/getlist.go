package language

import (
	"net/http"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/httputils"
)

const getLanguageListRoute = "/languages"

// GetLanguageList returns http.HandlerFunc
// @Summary Get languages
// @Description Get list of all languages
// @Tags languages
// @Produce json
// @Success 200 {object} getLanguagesResponse
// @Failure 500 {object} httputils.ErrorStruct
// @Router /languages [get]
func (h *LanguageHandlers) GetLanguageList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		languages, err := h.languageService.GetLanguages(r.Context())
		if err != nil {
			h.log.Error(err.Error())
			httputils.RespondWith500(w, h.log)
			return
		}

		resp := &getLanguagesResponse{
			Languages: languages,
		}

		httputils.RespondWith200(w, resp, h.log)
	}
}

// @Description get languages getLanguagesResponse.
type getLanguagesResponse struct {
	Languages []string `json:"languages"`
}
