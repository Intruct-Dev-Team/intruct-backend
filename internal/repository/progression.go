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
		Select("p.course_id").
		From("course_progressions p").
		InnerJoin("courses c ON p.course_id = c.course_id").
		Where(squirrel.Eq{"user_id": userID}).
		Where(squirrel.Eq{"c.deleted_at": nil}).
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

func (r *Repository) CreateCourseProgression(ctx context.Context, userID, courseID, currentLessonID int) error {
	query, args, err := r.sqlBuilder.
		Insert("course_progressions").
		Columns(
			"user_id",
			"course_id",
			"current_lesson_id",
		).
		Values(
			userID,
			courseID,
			currentLessonID,
		).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create course progression for user [%d] and course [%d]: %w", userID, courseID, err)
	}

	return nil
}

func (r *Repository) UpdateCourseProgression(ctx context.Context, userID int, courseID int, newCurrentLessonID int, finishedLessonsCount int, isFinished bool) error {
	query, args, err := r.sqlBuilder.
		Update("course_progressions").
		Set("current_lesson_id", newCurrentLessonID).
		Set("finished_lessons_count", finishedLessonsCount).
		Set("is_finished", isFinished).
		Where(squirrel.Eq{
			"user_id":   userID,
			"course_id": courseID,
		}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update course progression for user [%d] and course [%d]: %w", userID, courseID, err)
	}

	return nil
}

func (r *Repository) nullifyCurrentLessonInProgressions(ctx context.Context, tx *sqlx.Tx, courseID int) error {
	query, args, err := r.sqlBuilder.
		Update("course_progressions").
		Set("current_lesson_id", nil).
		Where(squirrel.Eq{"course_id": courseID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to nullify current lesson: %w", err)
	}

	return nil
}
