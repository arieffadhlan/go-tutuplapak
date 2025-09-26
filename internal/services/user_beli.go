package services

import (
	"context"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/repository"
	"tutuplapak-user/internal/utils"
)

type AuthBeliService struct {
	repo repository.UserBeliRepositoryInterface
}

func NewAuthBeliService(repo repository.UserBeliRepositoryInterface) AuthBeliService {
	return AuthBeliService{repo: repo}
}

func (s AuthBeliService) Login(ctx context.Context, req dto.AuthLoginUserBeliRequest) (tkn string, err error) {
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return "", err
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", ErrInvalidCredentials
	}

	tkn, err = utils.GenerateJWTTokenBeli(user)
	if err != nil {
		return "", err
	}

	return tkn, nil
}

func (s AuthBeliService) Register(ctx context.Context, req dto.AuthUserBeliRequest, isAdmin bool) (tkn string, err error) {
	user, err := s.repo.RegisterUserBeli(ctx, req, isAdmin)
	if err != nil {
		return "", err
	}

	tkn, err = utils.GenerateJWTTokenBeli(user)
	if err != nil {
		return "", err
	}

	return tkn, nil
}
