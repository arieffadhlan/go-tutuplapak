package services

import (
	"errors"
	"fmt"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PurchaseService struct {
	repository *repository.PurchaseRepository
	db         *sqlx.DB
}

func NewPurchaseService(repository *repository.PurchaseRepository, db *sqlx.DB) *PurchaseService {
	return &PurchaseService{
		repository: repository,
		db:         db,
	}
}

type ServiceError struct {
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

var ErrQtyExceedsStock = errors.New("requested quantity exceeds stock")
var ErrFileIdsCountMismatch = errors.New("file ids count does not match order payment count")
var ErrFileIdNotFound = errors.New("one or more file ids not found")
var ErrProductIdNotFound = errors.New("one or more product ids not found")

func (s *PurchaseService) Purchase(ctx *fiber.Ctx, req dto.CreatePurchaseRequest) (dto.CreatePurchaseResponse, error) {
	productIds := make([]uuid.UUID, 0, len(req.PurchasedItems))

	for _, product := range req.PurchasedItems {
		productIds = append(productIds, product.ProductId)
	}

	// get product by ids
	products, err := s.repository.GetProductByIds(ctx, productIds)
	if err != nil {
		return dto.CreatePurchaseResponse{}, err
	}

	productMap := make(map[uuid.UUID]dto.PurchaseProductDetail)
	for _, p := range products {
		productMap[p.ID] = p
	}

	// validate if product id in request exists and qty does not exceed stock
	for _, item := range req.PurchasedItems {
		p, ok := productMap[item.ProductId]
		if !ok {
			return dto.CreatePurchaseResponse{}, fmt.Errorf("%w: productId %s not found", ErrProductIdNotFound, item.ProductId)
		}

		if item.Qty > p.Qty {
			return dto.CreatePurchaseResponse{}, fmt.Errorf("%w productId %s not enough stock", ErrQtyExceedsStock, p.ID)
		}
	}

	// get seller by user ids and calculate total price per seller
	userIds := make([]uuid.UUID, 0, len(products))
	userIdMap := make(map[uuid.UUID]struct{})
	for _, p := range products {
		if _, exists := userIdMap[p.UserId]; !exists {
			userIds = append(userIds, p.UserId)
			userIdMap[p.UserId] = struct{}{}
		}
	}

	sellers, err := s.repository.GetSellerByIds(ctx, userIds)
	if err != nil {
		return dto.CreatePurchaseResponse{}, err
	}

	sellerMap := make(map[uuid.UUID]dto.SellerData)
	for _, s := range sellers {
		sellerMap[s.UserId] = s
	}

	// calculate total price
	totalPrice := 0
	for _, item := range req.PurchasedItems {
		p, ok := productMap[item.ProductId]
		if !ok {
			continue
		}

		totalPrice += p.Price * item.Qty
	}

	// calculate total price per seller
	totalPricePerSeller := make(map[uuid.UUID]int)
	for _, item := range req.PurchasedItems {
		p, ok := productMap[item.ProductId]
		if !ok {
			continue
		}

		seller, ok := sellerMap[p.UserId]
		if !ok {
			continue
		}

		totalPricePerSeller[seller.UserId] += p.Price * item.Qty
	}

	// start transaction
	trx, err := s.db.BeginTxx(ctx.Context(), nil)
	if err != nil {
		return dto.CreatePurchaseResponse{}, fmt.Errorf("failed start transaction: %w", err)
	}

	purchasePayload := dto.PurchaseRequest{
		SenderName:          req.SenderName,
		SenderContactType:   req.SenderContactType,
		SenderContactDetail: req.SenderContactDetail,
		TotalPrice:          totalPrice,
	}

	purchaseId, err := s.repository.CreatePurchase(ctx, purchasePayload, trx)
	if err != nil {
		return dto.CreatePurchaseResponse{}, err
	}

	purchaseItem := make([]dto.PurchaseItemResponse, 0, len(req.PurchasedItems))
	purchaseItemPayload := make([]dto.PurchaseItemRequest, 0, len(req.PurchasedItems))
	for _, item := range req.PurchasedItems {
		p, ok := productMap[item.ProductId]
		if !ok {
			continue
		}

		purchaseItemPayload = append(purchaseItemPayload, dto.PurchaseItemRequest{
			OrderId:          purchaseId,
			ProductId:        p.ID,
			Name:             p.Name,
			SKU:              p.SKU,
			Qty:              p.Qty,
			Price:            p.Price,
			Category:         p.Category,
			FileURI:          p.FileURI,
			FileThumbnailURI: p.FileThumbnailURI,
			PurchaseQty:      item.Qty,
		})

		purchaseItem = append(purchaseItem, dto.PurchaseItemResponse{
			ProductId:        p.ID,
			Name:             p.Name,
			SKU:              p.SKU,
			Qty:              p.Qty,
			Price:            p.Price,
			Category:         p.Category,
			FileURI:          p.FileURI,
			FileThumbnailURI: p.FileThumbnailURI,
			FileID:           p.FileID,
			CreatedAt:        p.CreatedAt,
			UpdatedAt:        p.UpdatedAt,
		})
	}

	purchasePayment := make([]dto.PurchasePaymentResponse, 0, len(totalPricePerSeller))
	purchasePaymentPayload := make([]dto.PurchasePaymentRequest, 0, len(totalPricePerSeller))
	for userId, totalPrice := range totalPricePerSeller {
		seller, ok := sellerMap[userId]
		if !ok {
			continue
		}

		purchasePaymentPayload = append(purchasePaymentPayload, dto.PurchasePaymentRequest{
			BankAccountName:   seller.BankAccountName,
			BankAccountHolder: seller.BankAccountHolder,
			BankAccountNumber: seller.BankAccountNumber,
			TotalPrice:        totalPrice,
			SellerId:          seller.UserId,
			OrderId:           purchaseId,
		})

		purchasePayment = append(purchasePayment, dto.PurchasePaymentResponse{
			BankAccountName:   seller.BankAccountName,
			BankAccountHolder: seller.BankAccountHolder,
			BankAccountNumber: seller.BankAccountNumber,
			TotalPrice:        totalPrice,
		})
	}

	// insert purchase items
	err = s.repository.CreatePurchaseItems(ctx, purchaseItemPayload, trx)
	if err != nil {
		_ = trx.Rollback()
		return dto.CreatePurchaseResponse{}, err
	}

	result := dto.CreatePurchaseResponse{
		PurchaseId:     purchaseId,
		TotalPrice:     totalPrice,
		PurchasedItems: purchaseItem,
		PaymentDetails: purchasePayment,
	}

	// insert purchase payments
	err = s.repository.CreatePurchasePayments(ctx, purchasePaymentPayload, trx)
	if err != nil {
		_ = trx.Rollback()
		return dto.CreatePurchaseResponse{}, err
	}

	// commit transaction
	if err := trx.Commit(); err != nil {
		return dto.CreatePurchaseResponse{}, fmt.Errorf("failed commit transaction: %w", err)
	}

	return result, err
}

func (s *PurchaseService) PurchasePaymentProof(ctx *fiber.Ctx, req dto.CreatePurchasePaymentProofRequest) error {
	// check if file id is exists
	for _, fileId := range req.FileIds {
		exists, err := s.repository.CheckFileIdExists(ctx, fileId)
		if err != nil {
			return fmt.Errorf("failed check file id %d exists: %w", fileId, err)
		}

		if !exists {
			return fmt.Errorf("%w: fileId %d not found", ErrFileIdNotFound, fileId)
		}
	}

	// get order payment by purchase id and compare order payment count with file ids count
	orderPayment, err := s.repository.GetOrderPaymentByPurchaseId(ctx, req.PurchaseId)
	if err != nil {
		return fmt.Errorf("failed get order payment: %w", err)
	}

	if len(orderPayment) != len(req.FileIds) {
		return fmt.Errorf("%w: expected %d file ids, got %d", ErrFileIdsCountMismatch, len(orderPayment), len(req.FileIds))
	}

	// start transaction
	trx, err := s.db.BeginTxx(ctx.Context(), nil)
	if err != nil {
		return fmt.Errorf("failed start transaction: %w", err)
	}

	err = s.repository.CreatePurchasePaymentProof(ctx, req.PurchaseId, req.FileIds, trx)
	if err != nil {
		_ = trx.Rollback()
		return err
	}

	// get purchased item by purchase id and reduce stock based on purchase qty
	purchasedItems, err := s.repository.GetOrderItemByPurchaseId(ctx, req.PurchaseId)
	if err != nil {
		_ = trx.Rollback()
		return err
	}

	for _, item := range purchasedItems {
		err = s.repository.UpdateProductQty(ctx, item.ProductId, item.Qty-item.PurchaseQty, trx)
		if err != nil {
			_ = trx.Rollback()
			return fmt.Errorf("failed update product id %d qty: %w", item.ProductId, err)
		}
	}

	// commit transaction
	if err := trx.Commit(); err != nil {
		return fmt.Errorf("failed commit transaction: %w", err)
	}

	return nil
}
