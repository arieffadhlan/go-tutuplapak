package repository

import (
	"context"
	"tutuplapak-user/internal/dto"

	"github.com/jmoiron/sqlx"
)

type FileRepositoryInterface interface {
	Post(c context.Context, files dto.File) (*dto.File, error)
}

type repository struct {
	db *sqlx.DB
}

func NewFileRepository(db *sqlx.DB) FileRepositoryInterface {
	return &repository{db: db}
}

func (r repository) Post(ctx context.Context, files dto.File) (*dto.File, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO files
		(url, thumbnail_url)
		VALUES($1, $2)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query, files.Url, files.ThumbnailUrl).Scan(&files.ID, &files.CreateAt)

	if err != nil {
		return nil, err
	}

	return &files, nil
}
