package repository

import (
	"context"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	"github.com/jmoiron/sqlx"
)

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
