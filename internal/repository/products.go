package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/utils"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductsRepository struct {
	db *sqlx.DB
}

func NewProductsRepository(db *sqlx.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(ctx context.Context, params dto.GetAllProductsParams) ([]dto.ProductResponse, error) {
	query := `
		SELECT id, name, sku, qty, price, category, file_id, file_uri, file_thumbnail_uri, created_at, updated_at
		FROM products
	`

	conditions := []string{}
	args := []any{}
	i := 1

	if params.ProductID != uuid.Nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", i))
		args = append(args, params.ProductID)
		i++
	}

	if params.SKU != "" {
		conditions = append(conditions, fmt.Sprintf("sku = $%d", i))
		args = append(args, params.SKU)
		i++
	}

	if params.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", i))
		args = append(args, params.Category)
		i++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	sortMap := map[string]string{
		"newest":    "updated_at DESC",
		"oldest":    "updated_at ASC",
		"expensive": "price DESC",
		"cheapest":  "price ASC",
	}
	orderBy, ok := sortMap[params.SortBy]
	if !ok {
		orderBy = "created_at DESC"
	}

	query += " ORDER BY " + orderBy

	if params.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", i)
		args = append(args, params.Limit)
		i++
	}
	if params.Offset > 0 {
		query += fmt.Sprintf(" Offset $%d", i)
		args = append(args, params.Offset)
		i++
	}

	var products []dto.ProductResponse
	err := r.db.SelectContext(ctx, &products, query, args...)
	if err != nil {
		 return nil, utils.NewInternal("failed get product")
	}

	return products, nil
}

func (r *ProductsRepository) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (dto.CreateProductResponse, error) {
	query := `
		INSERT INTO products (
      name, 
      user_id,
      file_id,
      sku,
      qty,
      price,
      category,
			file_uri,
			file_thumbnail_uri,
      created_at,
      updated_at
    ) VALUES (
      $1,
      $2,
      $3,
      $4,
      $5,
      $6,
      $7,
      $8,
      $9,
      NOW(),
      NOW()
    )
		RETURNING id, name, file_id, sku, qty, price, category, file_uri, file_thumbnail_uri, created_at, updated_at
	`

	var product dto.CreateProductResponse
	err := r.db.GetContext(ctx, &product, query,
		req.Name,
		req.UserID,
		req.FileID,
		req.SKU,
		req.Qty,
		req.Price,
		req.Category,
		req.FileURI,
		req.FileThumbnailURI,
	)

	if err != nil {
			fmt.Println("err", err.Error())
		 return dto.CreateProductResponse{}, utils.NewInternal("failed create product")
	}

	return product, nil
}

func (r *ProductsRepository) UpdateProduct(ctx context.Context, req dto.UpdateProductRequest) (dto.UpdateProductResponse, error) {
	query := `
		UPDATE products
		SET name = $1,
		    user_id = $2,
		    file_id = $3,
		    sku = $4,
		    qty = $5,
		    price = $6,
		    category = $7,
				file_uri = $8,
				file_thumbnail_uri = $9,
		    updated_at = NOW()
		WHERE id = $10
		RETURNING id, name, file_id, sku, qty, price, category, file_uri, file_thumbnail_uri, created_at, updated_at
	`

	var product dto.UpdateProductResponse
	err := r.db.GetContext(ctx, &product, query,
		req.Name,
		req.UserID,
		req.FileID,
		req.SKU,
		req.Qty,
		req.Price,
		req.Category,
		req.FileURI,
		req.FileThumbnailURI,
		req.ProdID,
	)

	if err != nil {
		 return dto.UpdateProductResponse{}, utils.NewInternal("failed update product")
	}

	return product, nil
}

func (r *ProductsRepository) DeleteProduct(ctx context.Context, userId uuid.UUID, prodId uuid.UUID) error {
	query := `
		DELETE FROM products
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, prodId, userId)
	if err != nil {
		 return utils.NewInternal("failed delete product")
	}

	a, err := result.RowsAffected()
	if err != nil {
		 return utils.NewInternal("failed delete product")
	}
	
	if a == 0 {
		 return utils.NewInternal("failed delete product")
	}

	return nil
}

func (r *ProductsRepository) CheckSKUExist(ctx context.Context, userId uuid.UUID, prodId uuid.UUID, sku string) error {
	query := `
		SELECT id
		FROM products 
		WHERE user_id = $1 AND sku = $2
		LIMIT 1
	`

	var ext uuid.UUID
	err := r.db.GetContext(ctx, &ext, query, userId, sku)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return utils.NewInternal("fail to check sku")
	}

	if ext != prodId {
		return utils.NewConflict("sku already exist")
	}

	return nil
}

func (r *ProductsRepository) CheckPrdOwner(ctx context.Context, userId uuid.UUID, productId uuid.UUID) error {
	query := `SELECT user_id FROM products WHERE id = $1`

	var ownerID uuid.UUID
	err := r.db.GetContext(ctx, &ownerID, query, productId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewNotFound("product data not found")
		} else {
			return utils.NewInternal("failed to check owners")
		}
	}

	if ownerID != userId {
		return utils.NewConflict("product has owned by another user")
	}

	return nil
}
