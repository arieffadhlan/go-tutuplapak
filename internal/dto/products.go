package dto

import "time"

type ProductResponse struct {
	ID               int       `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	SKU              string    `json:"sku" db:"sku"`
	Qty              int       `json:"qty" db:"qty"`
	Price            int       `json:"price" db:"price"`
	FileID           int       `json:"fileId" db:"file_id"`
	Category         string    `json:"category" db:"category"`
	FileURI          string    `json:"file_uri" db:"file_uri"`
	FileThumbnailURI string    `json:"file_thumbnail_uri" db:"file_thumbnail_uri"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type CreateProductRequest struct {
	Name     string `json:"name" validate:"required,min=4,max=32"`
	SKU      string `json:"sku" validate:"required,max=32"`
	Qty      int    `json:"qty" validate:"required,min=1"`
	Price    int    `json:"price" validate:"required,min=100"`
	FileID   int    `json:"fileId" validate:"required"`
	Category string `json:"category" validate:"required,oneof=Food Tools Clothes Beverage Furniture"`
}

type UpdateProductRequest struct {
	Name     string `json:"name" validate:"required,min=4,max=32"`
	SKU      string `json:"sku" validate:"required,max=32"`
	Qty      int    `json:"qty" validate:"required,min=1"`
	Price    int    `json:"price" validate:"required,min=100"`
	FileID   int    `json:"fileId" validate:"required"`
	Category string `json:"category" validate:"required,oneof=Food Tools Clothes Beverage Furniture"`
}

type CreateProductResponse ProductResponse
type UpdateProductResponse ProductResponse

type GetAllProductsParams struct {
	Limit     int
	Offset    int
	ProductID int
	SKU       string
	SortBy    string
	Category  string
}
