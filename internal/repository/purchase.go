package repository

import (
	"fmt"
	"tutuplapak-user/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PurchaseRepository struct {
	db *sqlx.DB
}

func NewPurchaseRepository(db *sqlx.DB) *PurchaseRepository {
	return &PurchaseRepository{
		db: db,
	}
}

func (r *PurchaseRepository) CreatePurchase(ctx *fiber.Ctx, req dto.PurchaseRequest, trx *sqlx.Tx) (uuid.UUID, error) {
	query := `
		INSERT INTO orders (
      total_price, 
      sender_name,
      sender_contact_type,
      sender_contact_detail,
      created_at
    ) VALUES (
      $1,
      $2,
      $3,
      $4,
      NOW()
    )
		RETURNING id
	`

	var purchaseId uuid.UUID
	err := trx.QueryRowContext(ctx.Context(), query,
		req.TotalPrice,
		req.SenderName,
		req.SenderContactType,
		req.SenderContactDetail,
	).Scan(&purchaseId)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed create purchase: %w", err)
	}

	return purchaseId, nil
}

func (r *PurchaseRepository) GetProductByIds(ctx *fiber.Ctx, productIds []uuid.UUID) ([]dto.PurchaseProductDetail, error) {
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
			user_id,
			created_at,
			updated_at
		FROM products
		WHERE id = ANY($1)
	`

	var products []dto.PurchaseProductDetail
	err := r.db.SelectContext(ctx.Context(), &products, query, pq.Array(productIds))
	if err != nil {
		return nil, fmt.Errorf("failed get product by ids: %w", err)
	}

	return products, nil
}

func (r *PurchaseRepository) GetSellerByIds(ctx *fiber.Ctx, userIds []uuid.UUID) ([]dto.SellerData, error) {
	query := `
		SELECT 
      id, 
      bank_account_name,
			bank_account_holder,
			bank_account_number
		FROM users
		WHERE id = ANY($1)
	`

	var sellers []dto.SellerData
	err := r.db.SelectContext(ctx.Context(), &sellers, query, pq.Array(userIds))
	if err != nil {
		return nil, fmt.Errorf("failed get seller by ids: %w", err)
	}

	return sellers, nil
}

func (r *PurchaseRepository) CreatePurchaseItems(ctx *fiber.Ctx, items []dto.PurchaseItemRequest, trx *sqlx.Tx) error {
	query := `
		INSERT INTO order_items (
			order_id, 
			product_id,
			name,
			sku,
			qty,
			price,
			category,
			file_uri,
			file_thumbnail_uri,
			purchase_qty
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
			$10
		)
	`

	for _, item := range items {
		_, err := trx.ExecContext(ctx.Context(), query,
			item.OrderId,
			item.ProductId,
			item.Name,
			item.SKU,
			item.Qty,
			item.Price,
			item.Category,
			item.FileURI,
			item.FileThumbnailURI,
			item.PurchaseQty,
		)
		if err != nil {
			return fmt.Errorf("failed create purchase items: %w", err)
		}
	}

	return nil
}

func (r *PurchaseRepository) CreatePurchasePayments(ctx *fiber.Ctx, payments []dto.PurchasePaymentRequest, trx *sqlx.Tx) error {
	query := `
		INSERT INTO order_payments (
			order_id, 
			seller_id,
			bank_account_name,
			bank_account_holder,
			bank_account_number,
			total_price
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
	`

	for _, payment := range payments {
		_, err := trx.ExecContext(ctx.Context(), query,
			payment.OrderId,
			payment.SellerId,
			payment.BankAccountName,
			payment.BankAccountHolder,
			payment.BankAccountNumber,
			payment.TotalPrice,
		)

		if err != nil {
			return fmt.Errorf("failed create purchase payments: %w", err)
		}
	}

	return nil
}

func (r *PurchaseRepository) CreatePurchasePaymentProof(ctx *fiber.Ctx, purchaseId uuid.UUID, fileIds []uuid.UUID, trx *sqlx.Tx) error {
	query := `
		INSERT INTO order_payment_proofs (
			file_id, 
			order_id
		) VALUES (
			$1,
			$2
		)
	`

	for _, fileId := range fileIds {
		fmt.Println(fileId, purchaseId)
		_, err := trx.ExecContext(ctx.Context(), query, fileId, purchaseId)
		if err != nil {
			return fmt.Errorf("failed add purchase payment proof: %w", err)
		}
	}

	return nil
}

func (r *PurchaseRepository) GetOrderPaymentByPurchaseId(ctx *fiber.Ctx, purchaseId uuid.UUID) ([]dto.PurchasePaymentRequest, error) {
	query := `
		SELECT
			seller_id,
			bank_account_name,
			bank_account_holder,
			bank_account_number,
			total_price
		FROM order_payments
		WHERE order_id = $1
	`

	var payments []dto.PurchasePaymentRequest
	err := r.db.SelectContext(ctx.Context(), &payments, query, purchaseId)
	if err != nil {
		return nil, fmt.Errorf("failed get order payment by purchase id: %w", err)
	}

	return payments, nil
}

func (r *PurchaseRepository) GetOrderItemByPurchaseId(ctx *fiber.Ctx, purchaseId uuid.UUID) ([]dto.PurchaseItemRequest, error) {
	query := `
		SELECT
			product_id,
			name,
			sku,
			qty,
			price,
			category,
			file_uri,
			file_thumbnail_uri,
			purchase_qty
		FROM order_items
		WHERE order_id = $1
	`

	var items []dto.PurchaseItemRequest
	err := r.db.SelectContext(ctx.Context(), &items, query, purchaseId)
	if err != nil {
		return nil, fmt.Errorf("failed get order item by purchase id: %w", err)
	}

	return items, nil
}

func (r *PurchaseRepository) UpdateProductQty(ctx *fiber.Ctx, productId uuid.UUID, qty int, trx *sqlx.Tx) error {
	query := `
		UPDATE products
		SET qty = $1
		WHERE id = $2
	`

	_, err := trx.ExecContext(ctx.Context(), query, qty, productId)
	if err != nil {
		return fmt.Errorf("failed update product qty: %w", err)
	}

	return nil
}

func (r *PurchaseRepository) CheckFileIdExists(ctx *fiber.Ctx, fileId uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM files
		WHERE id = $1
	`

	var count int
	err := r.db.GetContext(ctx.Context(), &count, query, fileId)
	if err != nil {
		return false, fmt.Errorf("failed check file id exists: %w", err)
	}

	return count > 0, nil
}
