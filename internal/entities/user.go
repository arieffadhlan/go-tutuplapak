package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id                uuid.UUID `json:"id" db:"id"`
	PublicId          uuid.UUID `json:"publicId" db:"public_id"`
	Email             *string   `json:"email" db:"email"`
	Phone             *string   `json:"phone" db:"phone"`
	Password          string    `json:"password" db:"password"`
	FileId            uuid.UUID `json:"fileId" db:"file_id"`
	FileUri           *string   `json:"fileUri" db:"file_uri"`
	FileThumbnailUri  *string   `json:"fileThumbnailUri" db:"file_thumbnail_uri"`
	BankAccountName   *string   `json:"bankAccountName" db:"bank_account_name"`
	BankAccountHolder *string   `json:"bankAccountHolder" db:"bank_account_holder"`
	BankAccountNumber *string   `json:"bankAccountNumber" db:"bank_account_number"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
}
