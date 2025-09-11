package entities

import "time"

type User struct {
	Id                int       `json:"id" db:"id"`
	PublicId          string    `json:"publicId" db:"public_id"`
	Email             *string   `json:"email" db:"email"`
	Phone             *string   `json:"phone" db:"phone"`
	Password          string    `json:"password" db:"password"`
	FileId            *string   `json:"fileId" db:"file_id"`
	FileUri           *string   `json:"fileUri" db:"file_uri"`
	FileThumbnailUri  *string   `json:"fileThumbnailUri" db:"file_thumbnail_uri"`
	BankAccountName   *string   `json:"bankAccountName" db:"bank_account_name"`
	BankAccountHolder *string   `json:"bankAccountHolder" db:"bank_account_holder"`
	BankAccountNumber *string   `json:"bankAccountNumber" db:"bank_account_number"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
}
