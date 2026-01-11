package repository

import (
	"context"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	"github.com/jmoiron/sqlx"
)

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

	var lessonID int
	if err := tx.GetContext(ctx, &lessonID, query, args...); err != nil {
		return 0, fmt.Errorf("failed to execute insert lesson: %w", err)
	}

	return lessonID, nil
}
