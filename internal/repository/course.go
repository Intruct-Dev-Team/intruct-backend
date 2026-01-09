package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	internalErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

func (r *Repository) IsCourseExistsByOwnerAndTitle(ctx context.Context, ownerID int, title string) (bool, error) {
	query, args, err := r.sqlBuilder.
		Select("1").
		From("course").
		Where(squirrel.Eq{
			"owner_id": ownerID,
			"title":    title,
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
		return false, fmt.Errorf("failed to check course existence by owner & title: %w", err)
	}

	return exists, nil
}

func (r *Repository) CreateCourse(ctx context.Context, course *entities.Course) (int, error) {
	stateMachine, err := r.getStateMachineByName(ctx, entities.CourseStateMachineName)
	if err != nil {
		return 0, fmt.Errorf("failed to get course's statemachine: %w", err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	// create stateMachineItem
	itemID, err := r.insertStateMachineItem(ctx, tx, *stateMachine)
	if err != nil {
		return 0, fmt.Errorf("failed to create state machine item: %w", err)
	}

	// insert course
	query, args, err := r.sqlBuilder.
		Insert("courses").
		Columns(
			"owner_id",
			"title",
			"description",
			"language_id",
			"state_machine_item_id").
		Values(
			course.OwnerID,
			course.Title,
			course.Description,
			course.LanguageID,
			itemID).
		Suffix("RETURNING course_id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build insert query: %w", err)
	}

	// insert course
	var courseID int

	if err := tx.GetContext(ctx, &courseID, query, args...); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return 0, internalErrs.ErrorNonUniqueData
			}
		}

		return 0, fmt.Errorf("failed to execute insert: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("error committing transaction: %w", err)
	}

	return courseID, nil
}
