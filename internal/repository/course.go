package repository

import (
	"context"
	"database/sql"
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
		From("courses").
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

func (r *Repository) GetOwnCourseIDs(ctx context.Context, userID int) ([]int, error) {
	query, args, err := r.sqlBuilder.
		Select("course_id").
		From("courses").
		Where(squirrel.Eq{"owner_id": userID}).
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
		return nil, fmt.Errorf("failed to get own course IDs for user [%d]: %w", userID, err)
	}

	if len(courseIDs) == 0 {
		return nil, internalErrs.ErrorSelectEmpty
	}

	return courseIDs, nil
}

func (r *Repository) GetPublicCourseIDs(ctx context.Context) ([]int, error) {
	query, args, err := r.sqlBuilder.
		Select("course_id").
		From("courses").
		Where(squirrel.Eq{"is_public": true}).
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
		return nil, fmt.Errorf("failed to get public course IDs: %w", err)
	}

	if len(courseIDs) == 0 {
		return nil, internalErrs.ErrorSelectEmpty
	}

	return courseIDs, nil
}

func (r *Repository) GetCourseByID(ctx context.Context, id int) (*entities.Course, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"c.course_id",
			"c.owner_id",
			"c.title",
			"c.description",
			"c.lessons_count",
			"c.language_id",
			"c.state_machine_item_id",
			"c.created_at",
			"c.updated_at",
			"l.name AS language",
		).
		From("courses c").
		LeftJoin("languages l ON c.language_id = l.language_id").
		Where(squirrel.Eq{"c.course_id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var course entities.Course
	err = r.db.GetContext(ctx, &course, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}
		return nil, fmt.Errorf("failed to get course by id [%d]: %w", id, err)
	}

	return &course, nil
}

func (r *Repository) CreateCourse(ctx context.Context, course *entities.Course) (int, error) {
	stateMachine, err := r.getStateMachineByName(ctx, entities.CourseStateMachineName)
	if err != nil {
		return 0, fmt.Errorf("failed to get course's statemachine: %w", err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
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

func (r *Repository) ImplementCourse(ctx context.Context, course *entities.Course, nextStateID int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. update state machine item state
	if err := r.updateStateMachineItemState(ctx, tx, course.StateMachineItemID, nextStateID); err != nil {
		return fmt.Errorf("failed to update state machine item: %w", err)
	}

	// 2. update course record
	query, args, err := r.sqlBuilder.
		Update("courses").
		Set("title", course.Title).
		Set("description", course.Description).
		Set("lessons_count", course.LessonsCount).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"course_id": course.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update course query: %w", err)
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}

	// 3. Insert modules
	for _, module := range course.Modules {
		moduleID, err := r.insertModule(ctx, tx, course.ID, module)
		if err != nil {
			return fmt.Errorf("failed to insert module: %w", err)
		}

		// 4. Insert lessons for each module
		for _, lesson := range module.Lessons {
			lessonID, err := r.insertLesson(ctx, tx, course.ID, moduleID, lesson)
			if err != nil {
				return fmt.Errorf("failed to insert lesson: %w", err)
			}

			// 5. Insert quizzes for each lesson
			for _, quiz := range lesson.Quizzes {
				quizID, err := r.insertQuiz(ctx, tx, lessonID, quiz)
				if err != nil {
					return fmt.Errorf("failed to insert quiz: %w", err)
				}

				// 6. Insert quiz options
				if err := r.insertQuizOptions(ctx, tx, quizID, quiz.Options); err != nil {
					return fmt.Errorf("failed to insert quiz options: %w", err)
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) SetIsPublicField(ctx context.Context, courseID int, isPublic bool) error {

	query, args, err := r.sqlBuilder.
		Update("courses").
		Set("is_public", isPublic).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"course_id": isPublic}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update course query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}

	return nil
}
