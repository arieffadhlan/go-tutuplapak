package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProductResponse struct {
	ID               uuid.UUID `json:"productId" db:"id"`
	FileID           uuid.UUID `json:"fileId" db:"file_id"`
	SKU              string    `json:"sku" db:"sku"`
	Qty              int       `json:"qty" db:"qty"`
	Price            int       `json:"price" db:"price"`
	Category         string    `json:"category" db:"category"`
	FileURI          string    `json:"fileUri" db:"file_uri"`
	FileThumbnailURI string    `json:"fileThumbnailUri" db:"file_thumbnail_uri"`
	Name             string    `json:"name" db:"name"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type CreateProductRequest struct {
	ProdID           uuid.UUID `json:"-"`
	UserID           uuid.UUID `json:"-"`
	FileID           uuid.UUID `json:"fileId" validate:"required"`
	Name             string    `json:"name" validate:"required,min=4,max=32"`
	SKU              string    `json:"sku" validate:"required,max=32"`
	Qty              int       `json:"qty" validate:"required,min=1"`
	Price            int       `json:"price" validate:"required,min=100"`
	Category         string    `json:"category" validate:"required,oneof=Food Tools Clothes Beverage Furniture"`
	FileURI          string    `json:"-"`
	FileThumbnailURI string    `json:"-"`
}

type UpdateProductRequest struct {
	ProdID           uuid.UUID `json:"-"`
	UserID           uuid.UUID `json:"-"`
	FileID           uuid.UUID `json:"fileId" validate:"required"`
	Name             string    `json:"name" validate:"required,min=4,max=32"`
	SKU              string    `json:"sku" validate:"required,max=32"`
	Qty              int       `json:"qty" validate:"required,min=1"`
	Price            int       `json:"price" validate:"required,min=100"`
	Category         string    `json:"category" validate:"required,oneof=Food Tools Clothes Beverag Furniture"`
	FileURI          string    `json:"-"`
	FileThumbnailURI string    `json:"-"`
}

type CreateProductResponse ProductResponse
type UpdateProductResponse ProductResponse

type GetAllProductsParams struct {
	Limit     int
	Offset    int
	ProductID uuid.UUID
	SKU       string
	SortBy    string
	Category  string
}
