package handlers

import (
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthBeliHandler struct {
	services services.AuthBeliService
}

func NewAuthBeliHandler(service services.AuthBeliService) AuthBeliHandler {
	return AuthBeliHandler{services: service}
}

func (h AuthBeliHandler) Login(ctx *fiber.Ctx) error {
	input := dto.AuthLoginUserBeliRequest{}

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	tkn, err := h.services.Login(ctx.Context(), input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	authRes := dto.AuthLoginBeliResponse{
		Token: tkn,
	}
	return ctx.Status(fiber.StatusOK).JSON(authRes)
}

func (h AuthBeliHandler) RegisterAdmin(ctx *fiber.Ctx) error {
	input := dto.AuthUserBeliRequest{}

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on register request", "data": err})
	}

	tkn, err := h.services.Register(ctx.Context(), input, true)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	authRes := dto.AuthLoginBeliResponse{
		Token: tkn,
	}
	return ctx.Status(fiber.StatusOK).JSON(authRes)
}

func (h AuthBeliHandler) RegisterUser(ctx *fiber.Ctx) error {
	input := dto.AuthUserBeliRequest{}

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on register request", "data": err})
	}

	tkn, err := h.services.Register(ctx.Context(), input, false)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	authRes := dto.AuthLoginBeliResponse{
		Token: tkn,
	}
	return ctx.Status(fiber.StatusOK).JSON(authRes)
}
