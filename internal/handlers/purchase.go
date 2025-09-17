package handlers

import (
	"errors"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PurchaseHandler struct {
	service *services.PurchaseService
}

func NewPurchaseHandler(service *services.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{service: service}
}

func (h *PurchaseHandler) Purchase(c *fiber.Ctx) error {
	var request dto.CreatePurchaseRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	r, err := h.service.Purchase(c, request)
	if err != nil {
		if errors.Is(err, services.ErrQtyExceedsStock) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if errors.Is(err, services.ErrProductIdNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(r)
}

func (h *PurchaseHandler) PurchasePaymentProof(c *fiber.Ctx) error {
	// validate purchaseId param
	purchaseId := c.Params("purchaseId")

	var request dto.CreatePurchasePaymentProofRequest
	uuidPurchaseId, err := uuid.Parse(purchaseId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid purchaseId format"})
	}
	request.PurchaseId = uuidPurchaseId

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.service.PurchasePaymentProof(c, request)
	if err != nil {
		if errors.Is(err, services.ErrFileIdsCountMismatch) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if errors.Is(err, services.ErrFileIdNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{})
}
