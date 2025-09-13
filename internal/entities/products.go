package entities

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID               uuid.UUID `db:"id" json:"id"`
	UserID           uuid.UUID `db:"user_id" json:"user_id"`
	FileID           uuid.UUID `db:"file_id" json:"file_id"`
	Name             string    `db:"name" json:"name"`
	SKU              string    `db:"sku" json:"sku"`
	Qty              int       `db:"qty" json:"qty"`
	Price            int       `db:"price" json:"price"`
	Category         string    `db:"category" json:"category"`
	FileURI          string    `db:"file_uri" json:"file_uri"`
	FileThumbnailURI string    `db:"file_thumbnail_uri" json:"file_thumbnail_uri"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
