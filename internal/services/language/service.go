package language

import (
	"context"
	"fmt"
)

type Repository interface {
	GetLanguages(ctx context.Context) ([]string, error)
}

type LanguageService struct {
	repo Repository
}

func NewService(repo Repository) *LanguageService {
	return &LanguageService{
		repo: repo,
	}
}

// GetLanguages returns all languages.
func (s *LanguageService) GetLanguages(ctx context.Context) ([]string, error) {
	categories, err := s.repo.GetLanguages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get languages from db: %w", err)
	}

	return categories, nil
}
