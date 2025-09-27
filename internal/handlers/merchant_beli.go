package handlers

import (
	"strconv"
	"strings"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/services"

	"github.com/gofiber/fiber/v2"
)

type MerchantBeliHandler struct {
	service services.MerchantBeliService
}

func NewMerchantBeliHandler(service services.MerchantBeliService) MerchantBeliHandler {
	return MerchantBeliHandler{service: service}
}

func (h MerchantBeliHandler) Create(ctx *fiber.Ctx) error {
	input := dto.MerchantCreateRequest{}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	res, err := h.service.Create(ctx.Context(), input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h MerchantBeliHandler) Get(ctx *fiber.Ctx) error {
	// Username comes from JWT middleware (c.Locals)
	usernameVal := ctx.Locals("username")
	if usernameVal == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	username, ok := usernameVal.(string)
	if !ok || username == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid username in token"})
	}

	// Parse filters from query params
	filter := dto.MerchantFilter{
		MerchantID:       ctx.Query("merchantId"),
		Name:             ctx.Query("name"),
		MerchantCategory: ctx.Query("merchantCategory"),
		SortCreatedAt:    strings.ToLower(ctx.Query("createdAt")),
	}

	// Limit & offset
	if l := ctx.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			filter.Limit = val
		}
	}
	if o := ctx.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			filter.Offset = val
		}
	}

	// Call service
	res, err := h.service.Get(ctx.Context(), username, filter)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Always return 200 with list (can be empty)
	return ctx.Status(fiber.StatusOK).JSON(res)
}
