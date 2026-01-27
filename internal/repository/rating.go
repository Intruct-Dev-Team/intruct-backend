package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	internalErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	"github.com/lib/pq"
)

func (r *Repository) CreateRating(ctx context.Context, rating *entities.Rating) error {

	query, args, err := r.sqlBuilder.
		Insert("ratings").
		Columns(
			"course_id",
			"user_id",
			"rating",
		).
		Values(
			rating.CourseID,
			rating.UserID,
			rating.Rating,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return internalErrs.ErrorNonUniqueData
			}
		}

		return fmt.Errorf("failed to execute rating insert: %w", err)
	}

	return nil
}
