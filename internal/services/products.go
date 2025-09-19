package services

import (
	"context"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/repository"

	"github.com/google/uuid"
)

type ProductsService struct {
	repository *repository.ProductsRepository
	fs         *FileService
}

func NewProductsService(repository *repository.ProductsRepository, fs *FileService) *ProductsService {
	return &ProductsService{
		repository: repository,
		fs:         fs,
	}
}

func (s *ProductsService) GetAllProducts(ctx context.Context, params dto.GetAllProductsParams) ([]dto.ProductResponse, error) {
	rs, err := s.repository.GetAllProducts(ctx, params) 
	if err != nil {
		 return nil, err
	}

	if len(rs) == 0 {
		 return []dto.ProductResponse{}, nil
	}

	return rs, nil
}

func (s *ProductsService)CreateProduct(ctx context.Context, req dto.CreateProductRequest) (dto.CreateProductResponse, error) {
	if err := s.repository.CheckSKUExist(ctx, req.UserID, req.ProdID, req.SKU); err != nil {
		 return dto.CreateProductResponse{}, err
	}

	file, err := s.fs.GetFileById(ctx, req.FileID)
	if err != nil {
		 return dto.CreateProductResponse{}, err
	}

	req.FileURI = file.Url
	req.FileThumbnailURI = file.ThumbnailUrl
	
	rs,err := s.repository.CreateProduct(ctx, req)
	if err != nil {
		 return rs, err
	}
	return rs, nil
}

func (s *ProductsService)UpdateProduct(ctx context.Context, req dto.UpdateProductRequest) (dto.UpdateProductResponse, error) {
	if err := s.repository.CheckPrdOwner(ctx, req.UserID, req.ProdID); err != nil {
		 return dto.UpdateProductResponse{}, err
	}

	if err := s.repository.CheckSKUExist(ctx, req.UserID, req.ProdID, req.SKU); err != nil {
		 return dto.UpdateProductResponse{}, err
	}

	file, err := s.fs.GetFileById(ctx, req.FileID)
	if err != nil {
		 return dto.UpdateProductResponse{}, err
	}

	req.FileURI = file.Url
	req.FileThumbnailURI = file.ThumbnailUrl
	
	rs,err := s.repository.UpdateProduct(ctx, req)
	if err != nil {
		 return rs, err
	}
	return rs, nil
}

func (s *ProductsService)DeleteProduct(ctx context.Context, prodId uuid.UUID, userId uuid.UUID) error {
	if err := s.repository.CheckPrdOwner(ctx, userId, prodId); err != nil {
		 return err
	}
	
	err := s.repository.DeleteProduct(ctx, userId, prodId)
	if err != nil {
		 return err
	}
	return nil
}