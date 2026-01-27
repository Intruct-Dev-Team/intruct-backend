package repository

import (
	"context"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) insertQuiz(ctx context.Context, tx *sqlx.Tx, lessonID int, quiz *entities.Quiz) (int, error) {
	query, args, err := r.sqlBuilder.
		Insert("quizes").
		Columns(
			"lesson_id",
			"serial_number",
			"question",
		).
		Values(
			lessonID,
			quiz.SerialNumber,
			quiz.Question,
		).
		Suffix("RETURNING quiz_id").
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("failed to build insert quiz query: %w", err)
	}

	var quizID int
	if err := tx.GetContext(ctx, &quizID, query, args...); err != nil {
		return 0, fmt.Errorf("failed to execute insert quiz: %w", err)
	}

	return quizID, nil
}

func (r *Repository) insertQuizOptions(ctx context.Context, tx *sqlx.Tx, quizID int, options []*entities.QuizOption) error {
	if len(options) == 0 {
		return nil
	}

	insertBuilder := r.sqlBuilder.
		Insert("quizes_options").
		Columns(
			"quiz_id",
			"content",
			"is_answer",
		)

	for _, option := range options {
		insertBuilder = insertBuilder.Values(
			quizID,
			option.Content,
			option.IsAnswer,
		)
	}

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert quiz options query: %w", err)
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute insert quiz options: %w", err)
	}

	return nil
}

func (r *Repository) deleteQuizOptionsByLessonIDs(ctx context.Context, tx *sqlx.Tx, lessonIDs []int) error {
	// First get all quiz IDs for these lessons
	quizQuery, quizArgs, err := r.sqlBuilder.
		Select("quiz_id").
		From("quizes").
		Where(squirrel.Eq{"lesson_id": lessonIDs}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build get quiz IDs query: %w", err)
	}

	var quizIDs []int
	if err := tx.SelectContext(ctx, &quizIDs, quizQuery, quizArgs...); err != nil {
		return fmt.Errorf("failed to get quiz IDs: %w", err)
	}

	if len(quizIDs) == 0 {
		return nil
	}

	// Delete quiz options
	query, args, err := r.sqlBuilder.
		Delete("quizes_options").
		Where(squirrel.Eq{"quiz_id": quizIDs}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to delete quiz options: %w", err)
	}

	return nil
}

func (r *Repository) deleteQuizzesByLessonIDs(ctx context.Context, tx *sqlx.Tx, lessonIDs []int) error {
	query, args, err := r.sqlBuilder.
		Delete("quizes").
		Where(squirrel.Eq{"lesson_id": lessonIDs}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to delete quizzes: %w", err)
	}

	return nil
}
