package services

import (
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductsService struct {
	repository *repository.ProductsRepository
}

func NewProductsService(repository *repository.ProductsRepository) *ProductsService {
	return &ProductsService{
		repository: repository,
	}
}

func (s *ProductsService) GetAllProducts(ctx *fiber.Ctx, params dto.GetAllProductsParams) ([]dto.ProductResponse, error) {
	rs, err := s.repository.GetAllProducts(ctx, params) 
	if err != nil {
		 return nil, err
	}

	if len(rs) == 0 {
		 return []dto.ProductResponse{}, nil
	}

	return rs, nil
}

func (s *ProductsService)CreateProduct(ctx *fiber.Ctx, req dto.CreateProductRequest) (dto.CreateProductResponse, error) {
	if err := s.repository.CheckSKUExist(ctx, req.UserID, req.ProdID, req.SKU); err != nil {
		 return dto.CreateProductResponse{}, err
	}
	
	rs,err := s.repository.CreateProduct(ctx, req)
	if err != nil {
		 return rs, err
	} else {
		 return rs, nil
	}
}

func (s *ProductsService)UpdateProduct(ctx *fiber.Ctx, req dto.UpdateProductRequest) (dto.UpdateProductResponse, error) {
	if err := s.repository.CheckPrdOwner(ctx, req.UserID, req.ProdID); err != nil {
		 return dto.UpdateProductResponse{}, err
	}

	if err := s.repository.CheckSKUExist(ctx, req.UserID, req.ProdID, req.SKU); err != nil {
		 return dto.UpdateProductResponse{}, err
	}
	
	rs,err := s.repository.UpdateProduct(ctx, req)
	if err != nil {
		 return rs, err
	} else {
		 return rs, nil
	}
}

func (s *ProductsService)DeleteProduct(ctx *fiber.Ctx, prodId uuid.UUID, userId uuid.UUID) error {
	if err := s.repository.CheckPrdOwner(ctx, userId, prodId); err != nil {
		 return err
	}
	
	err := s.repository.DeleteProduct(ctx, userId, prodId)
	if err != nil {
		 return err
	} else {
		 return nil
	}
}