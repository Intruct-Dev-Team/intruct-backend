package language

import (
	"context"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type LanguageService interface {
	GetLanguages(ctx context.Context) ([]string, error)
}

type LanguageHandlers struct {
	languageService LanguageService
	log             *zap.Logger
}

func NewLanguageHandlers(languageService LanguageService, log *zap.Logger) *LanguageHandlers {
	return &LanguageHandlers{
		languageService: languageService,
		log:             log,
	}
}

func (h *LanguageHandlers) SetupLanguageRoutes(router *chi.Mux) {
	router.Get(getLanguageListRoute, h.GetLanguageList())
}
