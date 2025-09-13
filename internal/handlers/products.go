package handlers

import (
	"strconv"
	"strings"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/utils"
	"tutuplapak-user/internal/services"

	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
)

type ProductsHandler struct {
	service *services.ProductsService
}

func NewProductsHandler(service *services.ProductsService) *ProductsHandler {
	return &ProductsHandler{
		service: service,
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
	var ctg string
	var sortBy string
	
	if val := c.Query("sku"); val != "" {
		 sku = val
	}

	if val := c.Query("ctg"); val != "" {
		allowedCategories := map[string]bool{
			"Foods":     true,
			"Tools":     true,
			"Beverages": true,
			"Furniture": true,
			"Clothes":   true,
		}
		if allowedCategories[val] {
			 ctg = val
		}
	}

	if sby := c.Query("sortBy"); sby != "" {
		validVal := map[string]bool{
			"newest":   true,
			"oldest":   true,
			"cheapest": true,
			"expensiv": true,
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
		Category:  ctg,
		SortBy:    sortBy,
	}

	products, err := h.service.GetAllProducts(c, params)
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

	r, err := h.service.CreateProduct(c, request)
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

	r, err := h.service.UpdateProduct(c, request)
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

	if err := h.service.DeleteProduct(c, productId, userIdUid); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return c.Status(appErr.Code).JSON(appErr)
			} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}
