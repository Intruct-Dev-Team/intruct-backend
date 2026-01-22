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

func (r *Repository) GetModulesWithLessonsByCourseID(ctx context.Context, courseID int) ([]*entities.Module, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"m.module_id",
			"m.course_id",
			"m.serial_number",
			"m.title",
			"m.description",
			"m.created_at",
			"m.updated_at",
			"l.lesson_id",
			"l.course_id as lesson_course_id",
			"l.module_id as lesson_module_id",
			"l.serial_number as lesson_serial_number",
			"l.title as lesson_title",
			"l.description as lesson_description",
			"l.created_at as lesson_created_at",
			"l.updated_at as lesson_updated_at",
		).
		From("modules m").
		InnerJoin("lessons l ON m.module_id = l.module_id").
		Where(squirrel.Eq{"m.course_id": courseID}).
		OrderBy("m.serial_number", "l.serial_number").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	type row struct {
		entities.Module
		LessonID        int       `db:"lesson_id"`
		LessonCourseID  int       `db:"lesson_course_id"`
		LessonModuleID  int       `db:"lesson_module_id"`
		LessonSerial    int       `db:"lesson_serial_number"`
		LessonTitle     string    `db:"lesson_title"`
		LessonDesc      string    `db:"lesson_description"`
		LessonCreatedAt time.Time `db:"lesson_created_at"`
		LessonUpdatedAt time.Time `db:"lesson_updated_at"`
	}

	var rows []row
	err = r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}
		return nil, fmt.Errorf("failed to get modules with lessons for course [%d]: %w", courseID, err)
	}

	if len(rows) == 0 {
		return nil, internalErrs.ErrorSelectEmpty
	}

	modulesMap := make(map[int]*entities.Module)
	var moduleOrder []int

	for _, r := range rows {
		// get or create module
		module, exists := modulesMap[r.ID]
		if !exists {
			module = &entities.Module{
				ID:           r.ID,
				CourseID:     r.CourseID,
				SerialNumber: r.SerialNumber,
				Title:        r.Title,
				Description:  r.Description,
				CreatedAt:    r.CreatedAt,
				UpdatedAt:    r.UpdatedAt,
				Lessons:      []*entities.Lesson{},
			}
			modulesMap[r.ID] = module
			moduleOrder = append(moduleOrder, r.ID)
		}

		// add lesson
		lesson := &entities.Lesson{
			ID:           r.LessonID,
			CourseID:     r.LessonCourseID,
			ModuleID:     r.LessonModuleID,
			SerialNumber: r.LessonSerial,
			Title:        r.LessonTitle,
			Description:  r.LessonDesc,
			CreatedAt:    r.LessonCreatedAt,
			UpdatedAt:    r.LessonUpdatedAt,
		}
		module.Lessons = append(module.Lessons, lesson)
	}

	// preserve order
	modules := make([]*entities.Module, 0, len(moduleOrder))
	for _, moduleID := range moduleOrder {
		modules = append(modules, modulesMap[moduleID])
	}

	return modules, nil
}

func (r *Repository) insertModule(ctx context.Context, tx *sqlx.Tx, courseID int, module *entities.Module) (int, error) {
	query, args, err := r.sqlBuilder.
		Insert("modules").
		Columns(
			"course_id",
			"serial_number",
			"title",
			"description",
		).
		Values(
			courseID,
			module.SerialNumber,
			module.Title,
			module.Description,
		).
		Suffix("RETURNING module_id").
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("failed to build insert module query: %w", err)
	}

	var moduleID int
	if err := tx.GetContext(ctx, &moduleID, query, args...); err != nil {
		return 0, fmt.Errorf("failed to execute insert module: %w", err)
	}

	return moduleID, nil
}
