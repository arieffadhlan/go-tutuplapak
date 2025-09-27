package services

import (
	"context"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/repository"
)

type MerchantBeliService struct {
	repo repository.MerchantsBeliRepository
}

func NewMerchantBeliService(repo repository.MerchantsBeliRepository) MerchantBeliService {
	return MerchantBeliService{repo: repo}
}

func (s MerchantBeliService) Create(ctx context.Context, req dto.MerchantCreateRequest) (merchant dto.MerchantCreateResponse, err error) {
	res, err := s.repo.CreateMerchant(ctx, req)
	if err != nil {
		return dto.MerchantCreateResponse{}, err
	}

	return res, nil
}

func (s MerchantBeliService) Get(ctx context.Context, username string, filter dto.MerchantFilter) (list dto.ListMerchantResponse, err error) {
	res, err := s.repo.GetMerchant(ctx, username, filter)
	if err != nil {
		return dto.ListMerchantResponse{}, err
	}
	return res, nil
}
