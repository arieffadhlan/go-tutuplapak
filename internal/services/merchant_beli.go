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

func (s MerchantBeliService) Get(ctx context.Context, filter dto.MerchantFilter) (list dto.ListMerchantResponse, err error) {
	res, err := s.repo.GetMerchant(ctx, filter)
	if err != nil {
		return dto.ListMerchantResponse{}, err
	}
	return res, nil
}

func (s MerchantBeliService) CreateItem(ctx context.Context, merchantId string, req dto.ItemCreateRequest) (item dto.ItemCreateResponse, err error) {
	res, err := s.repo.CreateItem(ctx, merchantId, req)
	if err != nil {
		return dto.ItemCreateResponse{}, err
	}
	return res, nil
}

func (s MerchantBeliService) GetItem(ctx context.Context, merchnatId string, filter dto.ItemFilter) (list dto.ListItemResponse, err error) {
	res, err := s.repo.GetItem(ctx, merchnatId, filter)
	if err != nil {
		return dto.ListItemResponse{}, err
	}
	return res, nil
}
