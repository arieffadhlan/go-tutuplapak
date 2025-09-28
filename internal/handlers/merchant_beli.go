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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on create merchant request", "data": err})
	}

	res, err := h.service.Create(ctx.Context(), input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h MerchantBeliHandler) Get(ctx *fiber.Ctx) error {
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
	res, err := h.service.Get(ctx.Context(), filter)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Always return 200 with list (can be empty)
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h MerchantBeliHandler) CreateItem(ctx *fiber.Ctx) error {
	input := dto.ItemCreateRequest{}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on create merchant request", "data": err})
	}
	var merchantId string

	if pidStr := ctx.Params("merchantId"); pidStr != "" {
		merchantId = pidStr
	}
	res, err := h.service.CreateItem(ctx.Context(), merchantId, input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h MerchantBeliHandler) GetItem(ctx *fiber.Ctx) error {
	// Parse filters from query params
	filter := dto.ItemFilter{
		Name:            ctx.Query("name"),
		SortCreatedAt:   strings.ToLower(ctx.Query("createdAt")),
		ItemID:          ctx.Query("itemId"),
		ProductCategory: ctx.Query("productCategory"),
	}

	var merchantId string
	if pidStr := ctx.Params("merchantId"); pidStr != "" {
		merchantId = pidStr
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
	res, err := h.service.GetItem(ctx.Context(), merchantId, filter)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Always return 200 with list (can be empty)
	return ctx.Status(fiber.StatusOK).JSON(res)
}
