package dto

import "github.com/google/uuid"

type (
	LinkEmailRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	LinkPhoneRequest struct {
		Phone string `json:"phone" validate:"required,e164"`
	}

	UpdateUserRequest struct {
		FileId            uuid.UUID `json:"fileId"`
		BankAccountName   string    `json:"bankAccountName" validate:"omitempty,min=4,max=32"`
		BankAccountHolder string    `json:"bankAccountHolder" validate:"omitempty,min=4,max=32"`
		BankAccountNumber string    `json:"bankAccountNumber" validate:"omitempty,min=4,max=32"`
		FileURI           string    `json:"-"`
		FileThumbnailURI  string    `json:"-"`
	}

	UserResponse struct {
		Email             string    `json:"email"`
		Phone             string    `json:"phone"`
		FileId            uuid.UUID `json:"fileId"`
		FileUri           string    `json:"fileUri"`
		FileThumbnailUri  string    `json:"fileThumbnailUri"`
		BankAccountName   string    `json:"bankAccountName"`
		BankAccountHolder string    `json:"bankAccountHolder"`
		BankAccountNumber string    `json:"bankAccountNumber"`
	}
)
