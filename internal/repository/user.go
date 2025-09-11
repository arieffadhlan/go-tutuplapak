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
        SELECT email, phone, password, public_id FROM users WHERE email = $1
    `
	// err = r.db.GetContext(ctx, &user, query, email)
	err = r.db.QueryRowContext(ctx, query, email).Scan(&user.Email, &user.Phone, &user.Password, &user.PublicId)

	if err == sql.ErrNoRows {
		return entities.User{}, ErrUserNotFound
	}

	if err != nil {
		return
	}

	return
}

func (r UserRepository) GetUserByPhone(ctx context.Context, phone string) (user entities.User, err error) {
	query := `
        SELECT email, phone, password, public_id FROM users WHERE phone = $1
    `

	err = r.db.QueryRowContext(ctx, query, phone).Scan(&user.Email, &user.Phone, &user.Password, &user.PublicId)

	if err == sql.ErrNoRows {
		return entities.User{}, ErrUserNotFound
	}

	if err != nil {
		return
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
        RETURNING email, phone, public_id
    `

	err = r.db.QueryRowContext(ctx, query, publicID, req.Email, hashedPassword).Scan(&user.Email, &user.Phone, &user.PublicId)

	if err != nil {
		return
	}

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
        RETURNING email, phone, public_id
    `

	err = r.db.QueryRowContext(ctx, query, publicID, req.Phone, hashedPassword).Scan(&user.Email, &user.Phone, &user.PublicId)

	if err != nil {
		return
	}

	return
}
