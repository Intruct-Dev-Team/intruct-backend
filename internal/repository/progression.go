package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	internalErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	"github.com/Masterminds/squirrel"
)

func (r *Repository) GetUsersCourseProgressions(ctx context.Context, userID int) ([]*entities.CourseProgression, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"user_id",
			"course_id",
			"current_lesson_id",
			"finished_lessons_count",
			"is_finished",
		).
		From("course_progressions").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var progressions []*entities.CourseProgression
	err = r.db.SelectContext(ctx, &progressions, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}
		return nil, fmt.Errorf("failed to get course progressions for user [%d]: %w", userID, err)
	}

	if len(progressions) == 0 {
		return nil, internalErrs.ErrorSelectEmpty
	}

	return progressions, nil
}

func (r *Repository) GetCourseProgressionByUserAndCourse(ctx context.Context, userID int, courseID int) (*entities.CourseProgression, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"user_id",
			"course_id",
			"current_lesson_id",
			"finished_lessons_count",
			"is_finished",
		).
		From("course_progressions").
		Where(squirrel.Eq{
			"user_id":   userID,
			"course_id": courseID,
		}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var progression entities.CourseProgression
	err = r.db.GetContext(ctx, &progression, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}
		return nil, fmt.Errorf("failed to get course progression for user [%d] and course [%d]: %w", userID, courseID, err)
	}

	return &progression, nil
}

func (r *Repository) GetCourseIDsFromUserProgress(ctx context.Context, userID int) ([]int, error) {
	query, args, err := r.sqlBuilder.
		Select("course_id").
		From("course_progressions").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var courseIDs []int
	err = r.db.SelectContext(ctx, &courseIDs, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}
		return nil, fmt.Errorf("failed to get course IDs from user progress for user [%d]: %w", userID, err)
	}

	if len(courseIDs) == 0 {
		return nil, internalErrs.ErrorSelectEmpty
	}

	return courseIDs, nil
}
