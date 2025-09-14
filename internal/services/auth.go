package services

import (
	"context"
	"errors"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/entities"
	"tutuplapak-user/internal/utils"
)

var ErrInvalidCredentials = errors.New("invalid identity or password")
var ErrUserAlreadyExists = errors.New("user already exists")

type repoContract interface {
	GetUserByEmail(ctx context.Context, email string) (entities.User, error)
	GetUserByPhone(ctx context.Context, phone string) (entities.User, error)
	RegisterByEmail(ctx context.Context, req dto.AuthEmailRequest) (entities.User, error)
	RegisterByPhone(ctx context.Context, req dto.AuthPhoneRequest) (entities.User, error)
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

	tkn, err = utils.GenerateJWTToken(user)
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

	tkn, err = utils.GenerateJWTToken(user)
	if err != nil {
		return "", user, err
	}

	return tkn, user, nil
}

func (s AuthService) RegisterByEmail(c context.Context, req dto.AuthEmailRequest) (tkn string, user entities.User, err error) {

	user, _ = s.repo.GetUserByEmail(c, req.Email)
	if user.Email != nil {
		return "", entities.User{}, ErrUserAlreadyExists
	}

	user, err = s.repo.RegisterByEmail(c, req)
	if err != nil {
		return "", entities.User{}, err
	}

	tkn, err = utils.GenerateJWTToken(user)
	if err != nil {
		return "", user, err
	}

	return tkn, user, nil
}

func (s AuthService) RegisterByPhone(c context.Context, req dto.AuthPhoneRequest) (tkn string, user entities.User, err error) {

	user, _ = s.repo.GetUserByPhone(c, req.Phone)
	if user.Phone != nil {
		return "", entities.User{}, ErrUserAlreadyExists
	}

	user, err = s.repo.RegisterByPhone(c, req)
	if err != nil {
		return "", entities.User{}, err
	}

	tkn, err = utils.GenerateJWTToken(user)
	if err != nil {
		return "", user, err
	}

	return tkn, user, nil
}
