package handlers

import (
	"strconv"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductsHandler struct {
	service *services.ProductsService
}

func NewProductsHandler(service *services.ProductsService) *ProductsHandler {
	return &ProductsHandler{service: service}
}

func (h *ProductsHandler) GetAllProducts(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	productID, _ := strconv.Atoi(c.Query("product_id", "0"))

	params := dto.GetAllProductsParams{
		Limit:     limit,
		Offset:    offset,
		ProductID: productID,
		SKU:       c.Query("sku", ""),
		SortBy:    c.Query("sort_by", ""),
		Category:  c.Query("category", ""),
	}

	products, err := h.service.GetAllProducts(c, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	} else {
		return c.Status(fiber.StatusOK).JSON(products)
	}
}

func (h *ProductsHandler) CreateProduct(c *fiber.Ctx) error {
	var request dto.CreateProductRequest
	
	if err := c.BodyParser(&request); err != nil {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	
	userId := 1
	r, err := h.service.CreateProduct(c, request, userId)
	if err != nil {
		 return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	} else {
		 return c.Status(fiber.StatusCreated).JSON(r)
	}
}

func (h *ProductsHandler) UpdateProduct(c *fiber.Ctx) error {
	id,err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid product data"})
	}

	var request dto.UpdateProductRequest

	if err := c.BodyParser(&request); err != nil {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	
	requestValidator := validator.New()
	if err := requestValidator.Struct(request); err != nil {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userId := 1
	r, err := h.service.UpdateProduct(c, request, userId, id)
	if err != nil {
		 return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	} else {
		 return c.Status(fiber.StatusOK).JSON(r)
	}
}

func (h *ProductsHandler) DeleteProduct(c *fiber.Ctx) error {
	id,err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid product id"})
	}

	if err := h.service.DeleteProduct(c, id); err != nil {
		 return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
