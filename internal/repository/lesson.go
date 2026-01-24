package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	internalErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) GetLessonByID(ctx context.Context, lessonID int) (*entities.Lesson, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"lesson_id",
			"course_id",
			"module_id",
			"serial_number",
			"title",
			"description",
			"content",
			"created_at",
			"updated_at",
		).
		From("lessons").
		Where(squirrel.Eq{"lesson_id": lessonID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var lesson entities.Lesson
	err = r.db.GetContext(ctx, &lesson, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}
		return nil, fmt.Errorf("failed to get lesson [%d]: %w", lessonID, err)
	}

	return &lesson, nil
}

func (r *Repository) GetQuizzesWithOptionsByLessonID(ctx context.Context, lessonID int) ([]*entities.Quiz, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"q.quiz_id",
			"q.lesson_id",
			"q.serial_number",
			"q.question",
			"q.created_at",
			"q.updated_at",
			"o.option_id",
			"o.quiz_id as option_quiz_id",
			"o.content",
			"o.is_answer",
			"o.created_at as option_created_at",
			"o.updated_at as option_updated_at",
		).
		From("quizes q").
		InnerJoin("quizes_options o ON q.quiz_id = o.quiz_id").
		Where(squirrel.Eq{"q.lesson_id": lessonID}).
		OrderBy("q.serial_number", "o.option_id").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	type row struct {
		entities.Quiz
		OptionID        int       `db:"option_id"`
		OptionQuizID    int       `db:"option_quiz_id"`
		Content         string    `db:"content"`
		IsAnswer        bool      `db:"is_answer"`
		OptionCreatedAt time.Time `db:"option_created_at"`
		OptionUpdatedAt time.Time `db:"option_updated_at"`
	}

	var rows []row
	err = r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}
		return nil, fmt.Errorf("failed to get quizzes with options for lesson [%d]: %w", lessonID, err)
	}

	if len(rows) == 0 {
		return nil, internalErrs.ErrorSelectEmpty
	}

	quizzesMap := make(map[int]*entities.Quiz)
	var quizOrder []int

	for _, r := range rows {
		// get or create quiz
		quiz, exists := quizzesMap[r.ID]
		if !exists {
			quiz = &entities.Quiz{
				ID:           r.ID,
				LessonID:     r.LessonID,
				SerialNumber: r.SerialNumber,
				Question:     r.Question,
				CreatedAt:    r.CreatedAt,
				UpdatedAt:    r.UpdatedAt,
				Options:      []*entities.QuizOption{},
			}
			quizzesMap[r.ID] = quiz
			quizOrder = append(quizOrder, r.ID)
		}

		// add option
		option := &entities.QuizOption{
			ID:        r.OptionID,
			QuizID:    r.OptionQuizID,
			Content:   r.Content,
			IsAnswer:  r.IsAnswer,
			CreatedAt: r.OptionCreatedAt,
			UpdatedAt: r.OptionUpdatedAt,
		}
		quiz.Options = append(quiz.Options, option)
	}

	// preserve order
	quizzes := make([]*entities.Quiz, 0, len(quizOrder))
	for _, quizID := range quizOrder {
		quizzes = append(quizzes, quizzesMap[quizID])
	}

	return quizzes, nil
}

func (r *Repository) GetAllLessonIDsByCourseIDOrderedBySerial(ctx context.Context, courseID int) ([]int, error) {
	query, args, err := r.sqlBuilder.
		Select("lesson_id").
		From("lessons").
		Where(squirrel.Eq{"course_id": courseID}).
		OrderBy("serial_number").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var lessonIDs []int
	err = r.db.SelectContext(ctx, &lessonIDs, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}
		return nil, fmt.Errorf("failed to get lesson IDs for course [%d]: %w", courseID, err)
	}

	if len(lessonIDs) == 0 {
		return nil, internalErrs.ErrorSelectEmpty
	}

	return lessonIDs, nil
}

func (r *Repository) insertLesson(ctx context.Context, tx *sqlx.Tx, courseID, moduleID int, lesson *entities.Lesson) (int, error) {
	query, args, err := r.sqlBuilder.
		Insert("lessons").
		Columns(
			"course_id",
			"module_id",
			"serial_number",
			"title",
			"description",
			"content",
		).
		Values(
			courseID,
			moduleID,
			lesson.SerialNumber,
			lesson.Title,
			lesson.Description,
			lesson.Content,
		).
		Suffix("RETURNING lesson_id").
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("failed to build insert lesson query: %w", err)
	}

	// logger.NewDefault().Info("lesson insert query", zap.String("query", query), zap.Any("args", args))

	var lessonID int
	if err := tx.GetContext(ctx, &lessonID, query, args...); err != nil {
		return 0, fmt.Errorf("failed to execute insert lesson: %w", err)
	}

	return lessonID, nil
}
