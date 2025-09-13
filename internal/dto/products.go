package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProductResponse struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	FileID           uuid.UUID `json:"fileId" db:"file_id"`
	SKU              string    `json:"sku" db:"sku"`
	Qty              int       `json:"qty" db:"qty"`
	Price            int       `json:"price" db:"price"`
	Category         string    `json:"category" db:"category"`
	FileURI          string    `json:"file_uri" db:"file_uri"`
	FileThumbnailURI string    `json:"file_thumbnail_uri" db:"file_thumbnail_uri"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type CreateProductRequest struct {
	Name     string    `json:"name" validate:"required,min=4,max=32"`
	SKU      string    `json:"sku" validate:"required,max=32"`
	Qty      int       `json:"qty" validate:"required,min=1"`
	Price    int       `json:"price" validate:"required,min=100"`
	Category string    `json:"category" validate:"required,oneof=Foods Tools Clothes Beverages Furniture"`
	ProdID   uuid.UUID `json:"-"`
	UserID   uuid.UUID `json:"-"`
	FileID   uuid.UUID `json:"fileId" validate:"required"`
}

type UpdateProductRequest struct {
	Name     string    `json:"name" validate:"required,min=4,max=32"`
	SKU      string    `json:"sku" validate:"required,max=32"`
	Qty      int       `json:"qty" validate:"required,min=1"`
	Price    int       `json:"price" validate:"required,min=100"`
	Category string    `json:"category" validate:"required,oneof=Food Tools Clothes Beverage Furniture"`
	ProdID   uuid.UUID `json:"-"`
	UserID   uuid.UUID `json:"-"`
	FileID   uuid.UUID `json:"fileId" validate:"required"`
}

type CreateProductResponse ProductResponse
type UpdateProductResponse ProductResponse

type GetAllProductsParams struct {
	Limit     int
	Offset    int
	ProductID uuid.UUID
	SKU       string
	Category  string
	SortBy    string
}
