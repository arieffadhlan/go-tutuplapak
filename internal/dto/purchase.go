package dto

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseItem struct {
	ProductId uuid.UUID `json:"productId" validate:"required,min=1"`
	Qty       int       `json:"qty" validate:"required,min=1"`
}

type PurchaseItemResponse struct {
	ProductId        uuid.UUID `json:"productId"`
	Name             string    `json:"name"`
	SKU              string    `json:"sku"`
	Qty              int       `json:"qty"`
	Price            int       `json:"price"`
	FileID           uuid.UUID `json:"fileId"`
	Category         string    `json:"category"`
	FileURI          string    `json:"fileUri"`
	FileThumbnailURI string    `json:"file_thumbnailUri"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
type CreatePurchaseRequest struct {
	SenderName          string `json:"senderName" validate:"required,min=4,max=55"`
	SenderContactType   string `json:"senderContactType" validate:"required,oneof=email phone"`
	SenderContactDetail string `json:"senderContactDetail" validate:"required,min=4,max=55"`
	PurchasedItems      []PurchaseItem
}

type PurchasePaymentResponse struct {
	BankAccountName   string `json:"bankAccountName"`
	BankAccountHolder string `json:"bankAccountHolder"`
	BankAccountNumber string `json:"bankAccountNumber"`
	TotalPrice        int    `json:"totalPrice"`
}
type CreatePurchaseResponse struct {
	PurchaseId     uuid.UUID                 `json:"purchaseId"`
	TotalPrice     int                       `json:"totalPrice"`
	PurchasedItems []PurchaseItemResponse    `json:"purchasedItems"`
	PaymentDetails []PurchasePaymentResponse `json:"paymentDetails"`
}

type PurchaseRequest struct {
	SenderName          string
	SenderContactType   string
	SenderContactDetail string
	TotalPrice          int
}

type PurchaseProductDetail struct {
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	ProductResponse
}

type PurchaseItemRequest struct {
	OrderId          uuid.UUID `db:"order_id"`
	ProductId        uuid.UUID `db:"product_id"`
	Name             string    `db:"name"`
	SKU              string    `db:"sku"`
	Qty              int       `db:"qty"`
	Price            int       `db:"price"`
	Category         string    `db:"category"`
	FileURI          string    `db:"file_uri"`
	FileThumbnailURI string    `db:"file_thumbnail_uri"`
	PurchaseQty      int       `db:"purchase_qty"`
}

type PurchasePaymentRequest struct {
	OrderId           uuid.UUID `db:"order_id"`
	SellerId          uuid.UUID `db:"seller_id"`
	BankAccountName   string    `db:"bank_account_name"`
	BankAccountHolder string    `db:"bank_account_holder"`
	BankAccountNumber string    `db:"bank_account_number"`
	TotalPrice        int       `db:"total_price"`
}

type SellerData struct {
	UserId            uuid.UUID `db:"id"`
	BankAccountName   string    `db:"bank_account_name"`
	BankAccountHolder string    `db:"bank_account_holder"`
	BankAccountNumber string    `db:"bank_account_number"`
}

type CreatePurchasePaymentProofRequest struct {
	PurchaseId uuid.UUID   `json:"purchaseId" validate:"required,min=1"`
	FileIds    []uuid.UUID `json:"fileIds" validate:"required,min=1,dive,required,min=1"`
}
