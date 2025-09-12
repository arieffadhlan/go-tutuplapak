package repository

import (
	"fmt"
	"strings"
	"tutuplapak-user/internal/dto"

	"github.com/gofiber/fiber/v2"
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

func (r *ProductsRepository) GetAllProducts(ctx *fiber.Ctx, params dto.GetAllProductsParams) ([]dto.ProductResponse, error) {
	query := `
		SELECT 
      id, 
      name, 
      sku, 
      qty, 
      price, 
      category, 
      file_id, 
      file_uri, 
      file_thumbnail_uri, 
      created_at, 
      updated_at
		FROM products
	`

	conditions := []string{}
	args := []any{}
	i := 1

	if params.ProductID > 0 {
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
		"newest":   "updated_at DESC",
		"oldest":   "updated_at ASC",
		"expensive": "prc DESC",
		"cheapest": "prc ASC",
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
	err := r.db.SelectContext(ctx.Context(), &products, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed get product: %w", err)
	}

	return products, nil
}

func (r *ProductsRepository) CreateProduct(ctx *fiber.Ctx, req dto.CreateProductRequest, userId int) (dto.CreateProductResponse, error) {
	query := `
		INSERT INTO products (
      name, 
      user_id,
      file_id,
      sku,
      qty,
      price,
      category,
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
      NOW(),
      NOW()
    )
		RETURNING id, name, file_id, sku, qty, price, category, file_uri, file_thumbnail_uri, created_at, updated_at
	`

	var product dto.CreateProductResponse
	err := r.db.GetContext(ctx.Context(), &product, query,
		req.Name,
		userId,
		req.FileID,
		req.SKU,
		req.Qty,
		req.Price,
		req.Category,
	)

	if err != nil {
		return dto.CreateProductResponse{}, fmt.Errorf("failed create product: %w", err)
	}

	return product, nil
}

func (r *ProductsRepository) UpdateProduct(ctx *fiber.Ctx, req dto.UpdateProductRequest, userId int, id int) (dto.UpdateProductResponse, error) {
	query := `
		UPDATE products
		SET name = $1,
		    user_id = $2,
		    file_id = $3,
		    sku = $4,
		    qty = $5,
		    price = $6,
		    category = $7,
		    updated_at = NOW()
		WHERE id = $8
		RETURNING id, name, file_id, sku, qty, price, category, file_uri, file_thumbnail_uri, created_at, updated_at
	`

	var product dto.UpdateProductResponse
	err := r.db.GetContext(ctx.Context(), &product, query,
		req.Name,
		userId,
		req.FileID,
		req.SKU,
		req.Qty,
		req.Price,
		req.Category,
		id,
	)

	if err != nil {
		return dto.UpdateProductResponse{}, fmt.Errorf("failed update product: %w", err)
	}

	return product, nil
}

func (r *ProductsRepository) DeleteProduct(ctx *fiber.Ctx, id int) error {
	query := `
		DELETE FROM products
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx.Context(), query, id)
	if err != nil {
		return fmt.Errorf("failed delete product: %w", err)
	}

	a, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed delete product: %w", err)
	}
	
	if a == 0 {
		return fmt.Errorf("failed delete product: %w", err)
	}

	return nil
}
