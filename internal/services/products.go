package services

import (
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type ProductsService struct {
	repository *repository.ProductsRepository
}

func NewProductsService(repository *repository.ProductsRepository) *ProductsService {
	return &ProductsService{
		repository: repository,
	}
}

func (s *ProductsService)GetAllProducts(ctx *fiber.Ctx, params dto.GetAllProductsParams) ([]dto.ProductResponse, error) {
	r, err := s.repository.GetAllProducts(ctx, params)
	if err != nil {
		return nil, err
	} else {
    return r, nil
  }
}

func (s *ProductsService)CreateProduct(ctx *fiber.Ctx, req dto.CreateProductRequest, userId int) (dto.CreateProductResponse, error) {
	r, err := s.repository.CreateProduct(ctx, req, userId)
	if err != nil {
		return dto.CreateProductResponse{}, err
	} else {
    return r, nil
  }
}

func (s *ProductsService)UpdateProduct(ctx *fiber.Ctx, req dto.UpdateProductRequest, userId int, id int) (dto.UpdateProductResponse, error) {
	r, err := s.repository.UpdateProduct(ctx, req, userId, id)
	if err != nil {
		return dto.UpdateProductResponse{}, err
	} else {
    return r, nil
  }
}


func (s *ProductsService) DeleteProduct(ctx *fiber.Ctx, id int) error {
	err := s.repository.DeleteProduct(ctx, id)
	if err != nil {
		return err
	} else {
		return nil
	}
}