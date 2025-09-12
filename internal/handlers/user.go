package handlers

import (
	"errors"
	"net/http"

	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/repository"
	"tutuplapak-user/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *services.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.New(),
	}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	publicId := c.Locals("publicId").(string)

	user, err := h.userService.GetUserProfile(c.Context(), publicId)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server Error",
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) LinkEmail(c *fiber.Ctx) error {
	publicId := c.Locals("publicId").(string)

	var req dto.LinkEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation error",
		})
	}

	err := h.userService.LinkEmail(c.Context(), publicId, req)
	if err != nil {
		if err.Error() == "email is taken" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"error": "email is taken",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server Error",
		})
	}

	// Return updated user profile
	user, err := h.userService.GetUserProfile(c.Context(), publicId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server Error",
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) LinkPhone(c *fiber.Ctx) error {
	publicId := c.Locals("publicId").(string)

	var req dto.LinkPhoneRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation error",
		})
	}

	err := h.userService.LinkPhone(c.Context(), publicId, req)
	if err != nil {
		if err.Error() == "phone is taken" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"error": "phone is taken",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server Error",
		})
	}

	// Return updated user profile
	user, err := h.userService.GetUserProfile(c.Context(), publicId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server Error",
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	publicId := c.Locals("publicId").(string)

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation error",
		})
	}

	err := h.userService.UpdateUser(c.Context(), publicId, req)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "fileId is not valid / exists",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server Error",
		})
	}

	// Return updated user profile
	user, err := h.userService.GetUserProfile(c.Context(), publicId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server Error",
		})
	}

	return c.JSON(user)
}
