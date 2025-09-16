package repository

import (
	"context"
	"database/sql"
	"errors"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/utils"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FileRepositoryInterface interface {
	GetFileById(ctx context.Context, id uuid.UUID) (dto.FileResponse, error)
	Post(c context.Context, files dto.File) (dto.File, error)
}

type repository struct {
	db *sqlx.DB
}

func NewFileRepository(db *sqlx.DB) FileRepositoryInterface {
	return &repository{db: db}
}

func (r repository) GetFileById(ctx context.Context, id uuid.UUID) (dto.FileResponse, error) {
	var dbFile dto.File
	query := `
		SELECT id, url, thumbnail_url, created_at
		FROM files
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &dbFile, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.FileResponse{}, utils.NewNotFound("file not found")
		}
		return dto.FileResponse{}, utils.NewInternal("failed to get file")
	}

	// mapping ke response DTO
	return dto.FileResponse{
		ID:           dbFile.ID,
		Url:          dbFile.Url,
		ThumbnailUrl: dbFile.ThumbnailUrl,
	}, nil
}

func (r repository) Post(ctx context.Context, files dto.File) (dto.File, error) {
	if err := ctx.Err(); err != nil {
		return dto.File{}, err
	}

	query := `
		INSERT INTO files
		(url, thumbnail_url)
		VALUES($1, $2)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query, files.Url, files.ThumbnailUrl).Scan(&files.ID, &files.CreateAt)

	if err != nil {
		return dto.File{}, err
	}

	return files, nil
}
