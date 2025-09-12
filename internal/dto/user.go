package dto

type (
	LinkEmailRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	LinkPhoneRequest struct {
		Phone string `json:"phone" validate:"required,e164"`
	}

	UpdateUserRequest struct {
		FileId            *string `json:"fileId"`
		BankAccountName   *string `json:"bankAccountName" validate:"omitempty,min=4,max=32"`
		BankAccountHolder *string `json:"bankAccountHolder" validate:"omitempty,min=4,max=32"`
		BankAccountNumber *string `json:"bankAccountNumber" validate:"omitempty,min=4,max=32"`
	}

	UserResponse struct {
		Email             *string `json:"email"`
		Phone             *string `json:"phone"`
		FileId            *string `json:"fileId"`
		FileUri           *string `json:"fileUri"`
		FileThumbnailUri  *string `json:"fileThumbnailUri"`
		BankAccountName   *string `json:"bankAccountName"`
		BankAccountHolder *string `json:"bankAccountHolder"`
		BankAccountNumber *string `json:"bankAccountNumber"`
	}
)
