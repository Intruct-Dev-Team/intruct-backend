package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	internalErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) GetUserStreakByUserID(ctx context.Context, id int) (*entities.Streak, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"user_id",
			"days_streak_count",
			"created_at",
			"updated_at",
		).
		From("public.streaks").
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var streak entities.Streak
	err = r.db.GetContext(ctx, &streak, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}
		return nil, fmt.Errorf("failed to find user's streak by user_id: %w", err)
	}

	return &streak, nil
}

func (r *Repository) upsertStreak(ctx context.Context, tx *sqlx.Tx, userID int, newDaysStreak int) error {
	// Check if streak exists
	streakExists, err := r.checkStreakExists(ctx, tx, userID)
	if err != nil {
		return fmt.Errorf("failed to check streak existence: %w", err)
	}

	if streakExists {
		return r.updateStreak(ctx, tx, userID, newDaysStreak)
	}
	return r.createStreak(ctx, tx, userID, newDaysStreak)
}

func (r *Repository) checkStreakExists(ctx context.Context, tx *sqlx.Tx, userID int) (bool, error) {
	query, args, err := r.sqlBuilder.
		Select("COUNT(*)").
		From("streaks").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var count int
	if err := tx.GetContext(ctx, &count, query, args...); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) createStreak(ctx context.Context, tx *sqlx.Tx, userID int, daysStreak int) error {
	query, args, err := r.sqlBuilder.
		Insert("streaks").
		Columns(
			"user_id",
			"days_streak_count",
		).
		Values(
			userID,
			daysStreak,
		).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute insert: %w", err)
	}

	return nil
}

func (r *Repository) updateStreak(ctx context.Context, tx *sqlx.Tx, userID int, newDaysStreak int) error {
	query, args, err := r.sqlBuilder.
		Update("streaks").
		Set("days_streak_count", newDaysStreak).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}

	return nil
}
