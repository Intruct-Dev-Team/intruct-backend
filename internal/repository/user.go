package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
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

func (r *Repository) CreateUser(ctx context.Context, user *entities.User) (int, error) {
	query, args, err := r.sqlBuilder.
		Insert("users").
		Columns(
			"external_uuid",
			"email",
			"name",
			"surname",
			"birthdate",
			"avatar").
		Values(
			user.ExternalUUID,
			user.Email,
			user.Name,
			user.Surname,
			user.Birthdate,
			user.Avatar).
		Suffix("RETURNING user_id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build insert query: %w", err)
	}

	var userID int
	err = r.db.GetContext(ctx, &userID, query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return 0, internalErrs.ErrorNonUniqueData
			}
		}

		return 0, fmt.Errorf("failed to execute insert: %w", err)
	}

	return userID, nil
}

// func (r *Repository) IsUserExistsByID(ctx context.Context, id int) (bool, error) {
// 	const query = `SELECT EXISTS(SELECT 1 FROM public.users WHERE user_id = $1)`

// 	var exists bool

// 	err := r.db.GetContext(ctx, &exists, query, id)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to check user existence: %w", err)
// 	}

// 	return exists, nil
// }

// func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
// 	const query = `SELECT user_id, email, password, name, surname, birthdate FROM public.users WHERE email = $1`

// 	var user entities.User
// 	err := r.db.GetContext(ctx, &user, query, email)

// 	if errors.Is(err, sql.ErrNoRows) {
// 		return nil, internalErrs.ErrorSelectEmpty
// 	}

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to find user by email: %w", err)
// 	}

// 	return &user, nil
// }

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

// func (r *Repository) GetUserStatByUserID(ctx context.Context, id int) (*entities.StudentStatistic, error) {
// 	const query = `
//     SELECT
//         COUNT(DISTINCT CASE WHEN s.name = $1 THEN l.lesson_id END) as count_of_finished_lesson,
//         COUNT(DISTINCT CASE WHEN s.name = $2 THEN l.lesson_id END) as count_of_verification_lesson,
//         COUNT(DISTINCT CASE WHEN s.name = $3 THEN l.lesson_id END) as count_of_waiting_lesson,
//         COUNT(DISTINCT CASE WHEN s.name = $1 THEN l.teacher_id END) as count_of_teachers
//     FROM lessons l
//     INNER JOIN statuses s ON s.status_id = l.status_id
//     WHERE l.student_id = $4
//     `

// 	var stat entities.StudentStatistic

// 	err := r.db.GetContext(ctx, &stat, query,
// 		entities.FinishedStatusName,     // $1
// 		entities.VerificationStatusName, // $2
// 		entities.WaitingStatusName,      // $3
// 		id,                              // $4
// 	)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, internalErrs.ErrorSelectEmpty
// 		}

// 		return nil, fmt.Errorf("failed to find user's statistic by user id: %w", err)
// 	}

// 	return &stat, nil
// }

// func (r *Repository) UpdateUser(ctx context.Context, userID int, user *entities.User) error {
// 	query, args, err := r.sqlBuilder.
// 		Update("users").
// 		Set("name", user.Name).
// 		Set("surname", user.Surname).
// 		Set("birthdate", user.Birthdate).
// 		Set("avatar", user.Avatar).
// 		Where(squirrel.Eq{"user_id": userID}).
// 		ToSql()
// 	if err != nil {
// 		return fmt.Errorf("failed to build update query: %w", err)
// 	}

// 	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
// 		return fmt.Errorf("failed to execute update: %w", err)
// 	}

// 	return nil
// }
