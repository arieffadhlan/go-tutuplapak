package repository

import (
	"context"
	"database/sql"
	"errors"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/entities"
	"tutuplapak-user/internal/utils"

	"github.com/jmoiron/sqlx"
)

type UserBeliRepositoryInterface interface {
	RegisterUserBeli(ctx context.Context, req dto.AuthUserBeliRequest, isAdmin bool) (user entities.UserBeli, err error)
	GetUserByUsername(ctx context.Context, username string) (user entities.UserBeli, err error)
}

type UserBeliRepository struct {
	db *sqlx.DB
}

func NewUserBeliRepository(db *sqlx.DB) UserBeliRepository {
	return UserBeliRepository{db: db}
}

var ErrUserBeliNotFound = errors.New("user not found")

func (r UserBeliRepository) GetUserByUsername(ctx context.Context, username string) (user entities.UserBeli, err error) {
	query := `
				SELECT username, password, email, isAdmin
				FROM users_beli WHERE username = $1
	`
	err = r.db.QueryRowContext(ctx, query, username).Scan(&user.Username, &user.Password, &user.Email, user.IsAdmin)

	if err == sql.ErrNoRows {
		return entities.UserBeli{}, ErrUserBeliNotFound
	}
	if err != nil {
		return entities.UserBeli{}, err
	}

	return user, nil
}

func (r UserBeliRepository) RegisterUserBeli(ctx context.Context, req dto.AuthUserBeliRequest, isAdmin bool) (user entities.UserBeli, err error) {
	hashedPassword, err := utils.HashPasswordBeli(req.Password)

	if err != nil {
		return
	}

	query := `
        INSERT INTO users_beli (username, email, password, isAdmin) 
        VALUES ($1, $2, $3, $4)
        RETURNING username, email, isAdmin
    `
	err = r.db.QueryRowContext(ctx, query, req.Username, req.Email, hashedPassword, isAdmin).Scan(&user.Username, &user.Email, &user.IsAdmin)
	if err != nil {
		return entities.UserBeli{}, err
	}

	return user, nil
}
