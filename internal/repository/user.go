package repository

import (
	"context"
	"database/sql"
	"errors"

	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/entities"
	"tutuplapak-user/internal/utils"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return UserRepository{db: db}
}

var ErrUserNotFound = errors.New("user not found")

func (r UserRepository) GetUserByEmail(ctx context.Context, email string) (user entities.User, err error) {
	query := `
        SELECT id, public_id, email, phone, password
        FROM users WHERE email = $1
    `
	err = r.db.QueryRowContext(ctx, query, email).
		Scan(&user.Id, &user.PublicId, &user.Email, &user.Phone, &user.Password)

	if err == sql.ErrNoRows {
		return entities.User{}, ErrUserNotFound
	}

	return
}

func (r UserRepository) GetUserByPhone(ctx context.Context, phone string) (user entities.User, err error) {
	query := `
        SELECT id, public_id, email, phone, password
        FROM users WHERE phone = $1
    `
	err = r.db.QueryRowContext(ctx, query, phone).
		Scan(&user.Id, &user.PublicId, &user.Email, &user.Phone, &user.Password)

	if err == sql.ErrNoRows {
		return entities.User{}, ErrUserNotFound
	}

	return
}

func (r UserRepository) RegisterByEmail(ctx context.Context, req dto.AuthEmailRequest) (user entities.User, err error) {
	publicID := uuid.New().String()

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return
	}

	query := `
        INSERT INTO users (public_id, email, password) 
        VALUES ($1, $2, $3)
        RETURNING id, public_id, email, phone
    `
	err = r.db.QueryRowContext(ctx, query, publicID, req.Email, hashedPassword).
		Scan(&user.Id, &user.PublicId, &user.Email, &user.Phone)

	return
}

func (r UserRepository) RegisterByPhone(ctx context.Context, req dto.AuthPhoneRequest) (user entities.User, err error) {
	publicID := uuid.New().String()

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return
	}

	query := `
        INSERT INTO users (public_id, phone, password) 
        VALUES ($1, $2, $3)
        RETURNING id, public_id, email, phone
    `
	err = r.db.QueryRowContext(ctx, query, publicID, req.Phone, hashedPassword).
		Scan(&user.Id, &user.PublicId, &user.Email, &user.Phone)

	return
}

func (r UserRepository) GetUserById(ctx context.Context, userId uuid.UUID) (user entities.User, err error) {
	query := `
        SELECT id, public_id, email, phone, file_id, file_uri, file_thumbnail_uri,
               bank_account_name, bank_account_holder, bank_account_number,
               created_at, updated_at 
        FROM users WHERE id = $1
    `
	err = r.db.GetContext(ctx, &user, query, userId)
	if err == sql.ErrNoRows {
		return entities.User{}, ErrUserNotFound
	}
	return
}

func (r UserRepository) LinkEmail(ctx context.Context, userId uuid.UUID, email string) error {
	query := `UPDATE users SET email = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, email, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r UserRepository) LinkPhone(ctx context.Context, userId uuid.UUID, phone string) error {
	query := `UPDATE users SET phone = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, phone, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r UserRepository) UpdateUser(ctx context.Context, userId uuid.UUID, req dto.UpdateUserRequest) error {
	query := `
        UPDATE users SET 
            file_id = COALESCE($1, file_id),
            file_uri = COALESCE($2, file_uri),
            file_thumbnail_uri = COALESCE($3, file_thumbnail_uri),
            bank_account_name = COALESCE($4, bank_account_name),
            bank_account_holder = COALESCE($5, bank_account_holder),
            bank_account_number = COALESCE($6, bank_account_number),
            updated_at = NOW()
        WHERE id = $7
    `

	result, err := r.db.ExecContext(ctx, query,
		req.FileId,
		req.FileURI,
		req.FileThumbnailURI,
		req.BankAccountName,
		req.BankAccountHolder,
		req.BankAccountNumber,
		userId)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r UserRepository) CheckEmailExists(ctx context.Context, email string, userId uuid.UUID) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = $1 AND id <> $2`

	err := r.db.GetContext(ctx, &count, query, email, userId)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r UserRepository) CheckPhoneExists(ctx context.Context, phone string, userId uuid.UUID) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE phone = $1 AND id <> $2`

	err := r.db.GetContext(ctx, &count, query, phone, userId)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
