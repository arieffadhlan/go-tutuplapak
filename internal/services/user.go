package services

import (
	"context"
	"errors"
	"strings"

	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/repository"

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

func (s *UserService) GetUserProfile(ctx context.Context, publicId string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByPublicId(ctx, publicId)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		Email:             user.Email,
		Phone:             user.Phone,
		FileId:            user.FileId,
		FileUri:           user.FileUri,
		FileThumbnailUri:  user.FileThumbnailUri,
		BankAccountName:   user.BankAccountName,
		BankAccountHolder: user.BankAccountHolder,
		BankAccountNumber: user.BankAccountNumber,
	}, nil
}

func (s *UserService) LinkEmail(ctx context.Context, publicId string, req dto.LinkEmailRequest) error {
	// Check if email already exists
	exists, err := s.userRepo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("email is taken")
	}

	err = s.userRepo.LinkEmail(ctx, publicId, req.Email)
	if err != nil {
		// Handle PostgreSQL unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if strings.Contains(pqErr.Message, "email") {
				return errors.New("email is taken")
			}
		}
		return err
	}

	return nil
}

func (s *UserService) LinkPhone(ctx context.Context, publicId string, req dto.LinkPhoneRequest) error {
	// Check if phone already exists
	exists, err := s.userRepo.CheckPhoneExists(ctx, req.Phone)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("phone is taken")
	}

	err = s.userRepo.LinkPhone(ctx, publicId, req.Phone)
	if err != nil {
		// Handle PostgreSQL unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if strings.Contains(pqErr.Message, "phone") {
				return errors.New("phone is taken")
			}
		}
		return err
	}

	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, publicId string, req dto.UpdateUserRequest) error {
	return s.userRepo.UpdateUser(ctx, publicId, req)
}
