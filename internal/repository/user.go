package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	internalErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
)

func (r *Repository) GetUserIDByExternalUUID(ctx context.Context, ExternalUUID string) (int, error) {
	query, args, err := r.sqlBuilder.Select(
		"user_id",
	).
		From("users").
		Where(squirrel.Eq{"external_uuid": ExternalUUID}).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("failed to build sql query: %w", err)
	}

	var userID int
	err = r.db.GetContext(ctx, &userID, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, internalErrs.ErrorSelectEmpty
		}

		return 0, fmt.Errorf("failed to execute user_id by external_uuid: %w", err)
	}

	return userID, nil
}

func (r *Repository) IsUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	query, args, err := r.sqlBuilder.
		Select("1").
		From("users").
		Where(squirrel.Eq{
			"email": email,
		}).
		Prefix("SELECT EXISTS(").
		Suffix(")").
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var exists bool
	err = r.db.GetContext(ctx, &exists, query, args...)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence by email: %w", err)
	}

	return exists, nil
}

func (r *Repository) IsUserExistsByID(ctx context.Context, id int) (bool, error) {
	query, args, err := r.sqlBuilder.
		Select("1").
		From("users").
		Where(squirrel.Eq{
			"user_id": id,
		}).
		Prefix("SELECT EXISTS(").
		Suffix(")").
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var exists bool
	err = r.db.GetContext(ctx, &exists, query, args...)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence by email: %w", err)
	}

	return exists, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *entities.User) (int, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	userID, err := r.insertUser(ctx, tx, user)
	if err != nil {
		return 0, err
	}

	if err := r.createStreak(ctx, tx, userID, 0); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("error committing transaction: %w", err)
	}

	return userID, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"user_id",
			"external_uuid",
			"email",
			"name",
			"surname",
			"registration_date",
			"birthdate",
			"avatar",
		).
		From("public.users").
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user entities.User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetUserStatisticByUserID(ctx context.Context, id int) (*entities.UserStatistic, error) {
	query, args, err := r.sqlBuilder.
		Select(
			`COUNT(*) FILTER (WHERE is_finished = TRUE)  AS count_of_complited_courses`,
			`COUNT(*) FILTER (WHERE is_finished = FALSE) AS count_of_in_progress_courses`,
		).
		From("course_progressions").
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var stat entities.UserStatistic
	err = r.db.GetContext(ctx, &stat, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user statistic: %w", err)
	}

	return &stat, nil
}

func (r *Repository) insertUser(ctx context.Context, tx *sqlx.Tx, user *entities.User) (int, error) {
	query, args, err := r.sqlBuilder.
		Insert("users").
		Columns(
			"external_uuid",
			"email",
			"name",
			"surname",
			"birthdate",
			"avatar",
		).
		Values(
			user.ExternalUUID,
			user.Email,
			user.Name,
			user.Surname,
			user.Birthdate,
			user.Avatar,
		).
		Suffix("RETURNING user_id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build insert user query: %w", err)
	}

	var userID int
	if err := tx.GetContext(ctx, &userID, query, args...); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, internalErrs.ErrorNonUniqueData
		}
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	return userID, nil
}
