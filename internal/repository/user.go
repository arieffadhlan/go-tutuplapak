package repository

import (
	"context"
	"database/sql"
	"errors"

	"tutuplapak-user/internal/entities"

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
