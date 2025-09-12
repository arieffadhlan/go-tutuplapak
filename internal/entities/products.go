package entities

import "time"

type Product struct {
	ID               int       `db:"id" json:"id"`
	Name             string    `db:"name" json:"name"`
	UserID           int       `db:"user_id" json:"user_id"`
	FileID           int       `db:"file_id" json:"file_id"`
	SKU              string    `db:"sku" json:"sku"`
	Qty              int       `db:"qty" json:"qty"`
	Price            int       `db:"prc" json:"price"`
	Category         string    `db:"ctg" json:"category"`
	FileURI          string    `db:"file_uri" json:"file_uri"`
	FileThumbnailURI string    `db:"file_thumbnail_uri" json:"file_thumbnail_uri"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
