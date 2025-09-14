package services

import (
	"context"
	"errors"
	"strings"

	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/repository"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// helper: convert *string → string
func strOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func (s *UserService) GetUserProfile(ctx context.Context, userId uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		Email:             strOrEmpty(user.Email),
		Phone:             strOrEmpty(user.Phone),
		FileId:            user.FileId,
		FileUri:           strOrEmpty(user.FileUri),
		FileThumbnailUri:  strOrEmpty(user.FileThumbnailUri),
		BankAccountName:   strOrEmpty(user.BankAccountName),
		BankAccountHolder: strOrEmpty(user.BankAccountHolder),
		BankAccountNumber: strOrEmpty(user.BankAccountNumber),
	}, nil
}

func (s *UserService) LinkEmail(ctx context.Context, userId uuid.UUID, req dto.LinkEmailRequest) error {
	exists, err := s.userRepo.CheckEmailExists(ctx, req.Email, userId)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("email is taken")
	}

	err = s.userRepo.LinkEmail(ctx, userId, req.Email)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if strings.Contains(pqErr.Message, "email") {
				return errors.New("email is taken")
			}
		}
		return err
	}

	return nil
}

func (s *UserService) LinkPhone(ctx context.Context, userId uuid.UUID, req dto.LinkPhoneRequest) error {
	exists, err := s.userRepo.CheckPhoneExists(ctx, req.Phone, userId)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("phone is taken")
	}

	err = s.userRepo.LinkPhone(ctx, userId, req.Phone)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if strings.Contains(pqErr.Message, "phone") {
				return errors.New("phone is taken")
			}
		}
		return err
	}

	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, userId uuid.UUID, req dto.UpdateUserRequest) error {
	return s.userRepo.UpdateUser(ctx, userId, req)
}
