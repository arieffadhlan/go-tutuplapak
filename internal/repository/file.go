package repository

import (
	"context"
	"database/sql"
	"errors"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FileRepositoryInterface interface {
	GetFileById(ctx *fiber.Ctx, id uuid.UUID) (*dto.FileResponse, error)
	Post(c context.Context, files dto.File) (*dto.File, error)
}

type repository struct {
	db *sqlx.DB
}

func NewFileRepository(db *sqlx.DB) FileRepositoryInterface {
	return &repository{db: db}
}

func (r repository) GetFileById(ctx *fiber.Ctx, id uuid.UUID) (*dto.FileResponse, error) {
	var dbFile dto.File
	query := `
		SELECT id, url, thumbnail_url, created_at
		FROM files
		WHERE id = $1
	`

	err := r.db.GetContext(ctx.Context(), &dbFile, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFound("file not found")
		}
		return nil, utils.NewInternal("failed to get file")
	}

	// mapping ke response DTO
	return &dto.FileResponse{
		ID:           dbFile.ID,
		Url:          dbFile.Url,
		ThumbnailUrl: dbFile.ThumbnailUrl,
	}, nil
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
