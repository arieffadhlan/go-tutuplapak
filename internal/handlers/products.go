package handlers

import (
	"strconv"
	"strings"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/services"
	"tutuplapak-user/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductsHandler struct {
	productService *services.ProductsService
}

func NewProductsHandler(productService *services.ProductsService) *ProductsHandler {
	return &ProductsHandler{
		productService: productService,
	}
}

func (h *ProductsHandler) GetAllProducts(c *fiber.Ctx) error {
	lim := 5
	if limStr := c.Query("limit"); limStr != "" {
		if limVal, err := strconv.Atoi(limStr); err == nil && limVal > 0 {
			 lim = limVal
		}
	}

	offset := 0
	if offStr := c.Query("offset"); offStr != "" {
		if offVal, err := strconv.Atoi(offStr); err == nil && offVal > 0 {
			 offset = offVal
		}
	}

	var productId uuid.UUID
	if pidStr := c.Query("productId"); pidStr != "" {
		if pid, err := uuid.Parse(pidStr); err == nil {
			 productId = pid
		}
	}

	var sku string
	var sortBy string
	var category string

	if val := c.Query("sku"); val != "" {
		 sku = val
	}

	if val := c.Query("category"); val != "" {
		 allowedCategories := map[string]bool{
			 "Food":      true,
			 "Tools":     true,
			 "Beverage":  true,
			 "Furniture": true,
			 "Clothes":   true,
		 }
		 if allowedCategories[val] {
			  category = val
		 }
	}

	if sby := c.Query("sortBy"); sby != "" {
		 validVal := map[string]bool{
			 "newest":    true,
			 "oldest":    true,
			 "cheapest":  true,
			 "expensive": true,
		 }
		 if validVal[strings.ToLower(sby)] {
			  sortBy = strings.ToLower(sby)
		 }
	}

	params := dto.GetAllProductsParams{
		Limit:     lim,
		Offset:    offset,
		ProductID: productId,
		SKU:       sku,
		SortBy:    sortBy,
		Category:  category,
	}

	products, err := h.productService.GetAllProducts(c.Context(), params)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return c.Status(appErr.Code).JSON(appErr)
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func (h *ProductsHandler) CreateProduct(c *fiber.Ctx) error {
	var request dto.CreateProductRequest

	if err := c.BodyParser(&request); err != nil {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validator.New().Struct(request); err != nil {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userIdStr, ok := c.Locals("userId").(string)
	if !ok || userIdStr == "" {
		 return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userIdUid, err := uuid.Parse(userIdStr)
	if err != nil {
		 return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	request.UserID = userIdUid

	r, err := h.productService.CreateProduct(c.Context(), request)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return c.Status(appErr.Code).JSON(appErr)
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(r)
}

func (h *ProductsHandler) UpdateProduct(c *fiber.Ctx) error {
	var productId uuid.UUID
	if pidStr := c.Params("productId"); pidStr != "" {
		if pid, err := uuid.Parse(pidStr); err == nil {
			 productId = pid
		}
	}

	var request dto.UpdateProductRequest

	if err := c.BodyParser(&request); err != nil {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validator.New().Struct(request); err != nil {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userIdStr, ok := c.Locals("userId").(string)
	if !ok || userIdStr == "" {
		 return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userIdUid, err := uuid.Parse(userIdStr)
	if err != nil {
		 return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	request.ProdID = productId
	request.UserID = userIdUid

	r, err := h.productService.UpdateProduct(c.Context(), request)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return c.Status(appErr.Code).JSON(appErr)
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.Status(fiber.StatusOK).JSON(r)
}

func (h *ProductsHandler) DeleteProduct(c *fiber.Ctx) error {
	var productId uuid.UUID
	if pidStr := c.Params("productId"); pidStr != "" {
		if pid, err := uuid.Parse(pidStr); err == nil {
			 productId = pid
		}
	}

	userIdStr, ok := c.Locals("userId").(string)
	if !ok || userIdStr == "" {
		 return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userIdUid, err := uuid.Parse(userIdStr)
	if err != nil {
		 return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := h.productService.DeleteProduct(c.Context(), productId, userIdUid); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return c.Status(appErr.Code).JSON(appErr)
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
