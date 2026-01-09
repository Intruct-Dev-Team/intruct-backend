package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	internalErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	"github.com/Masterminds/squirrel"
)

func (r *Repository) GetLanguageIDByName(ctx context.Context, languageName string) (int, error) {
	query, args, err := r.sqlBuilder.Select(
		"language_id",
	).
		From("languages").
		Where(squirrel.Eq{"name": languageName}).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("failed to build sql query: %w", err)
	}

	var languageID int
	err = r.db.GetContext(ctx, &languageID, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, internalErrs.ErrorSelectEmpty
		}

		return 0, fmt.Errorf("failed to execute language_id by language_name: %w", err)
	}

	return languageID, nil
}
