package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db         *sqlx.DB
	sqlBuilder squirrel.StatementBuilderType
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db:         db,
		sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
