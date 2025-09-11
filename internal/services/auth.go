package services

import (
	"context"
	"errors"
	"time"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/entities"
	"tutuplapak-user/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidCredentials = errors.New("invalid identity or password")

type repoContract interface {
	GetUserByEmail(ctx context.Context, email string) (entities.User, error)
	GetUserByPhone(ctx context.Context, phone string) (entities.User, error)
}

type AuthService struct {
	repo repoContract
}

func NewAuthService(repo repoContract) AuthService {
	return AuthService{repo: repo}
}

func (s AuthService) LoginByEmail(c context.Context, req dto.AuthEmailRequest) (tkn string, user entities.User, err error) {

	user, err = s.repo.GetUserByEmail(c, req.Email)
	if err != nil {
		return "", entities.User{}, err
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", entities.User{}, ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"user_id": user.PublicId,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tkn, err = token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", user, err
	}

	return tkn, user, nil
}

func (s AuthService) LoginByPhone(c context.Context, req dto.AuthPhoneRequest) (tkn string, user entities.User, err error) {

	user, err = s.repo.GetUserByPhone(c, req.Phone)
	if err != nil {
		return "", entities.User{}, err
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", entities.User{}, ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"user_id": user.PublicId,
		"phone":   user.Phone,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tkn, err = token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", user, err
	}

	return tkn, user, nil
}
