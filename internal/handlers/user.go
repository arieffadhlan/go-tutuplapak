package handlers

import (
	"errors"
	"net/http"

	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/repository"
	"tutuplapak-user/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	validator   *validator.Validate
	userService *services.UserService
	fileService *services.FileService
}

func NewUserHandler(userService *services.UserService, fileService *services.FileService) *UserHandler {
	return &UserHandler{
		userService: userService,
		fileService: fileService,
		validator:   validator.New(),
	}
}

func parseUserId(c *fiber.Ctx) (uuid.UUID, error) {
	userIdStr, ok := c.Locals("userId").(string)
	if !ok || userIdStr == "" {
		return uuid.Nil, errors.New("unauthorized")
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, errors.New("unauthorized")
	}

	return userId, nil
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userId, err := parseUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	user, err := h.userService.GetUserProfile(c.Context(), userId)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	return c.JSON(user)
}

func (h *UserHandler) LinkEmail(c *fiber.Ctx) error {
	userId, err := parseUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req dto.LinkEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Bad Request"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Validation error"})
	}

	err = h.userService.LinkEmail(c.Context(), userId, req)
	if err != nil {
		if err.Error() == "email is taken" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "email is taken"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	user, err := h.userService.GetUserProfile(c.Context(), userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	return c.JSON(user)
}

func (h *UserHandler) LinkPhone(c *fiber.Ctx) error {
	userId, err := parseUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req dto.LinkPhoneRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Bad Request"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Validation error"})
	}

	err = h.userService.LinkPhone(c.Context(), userId, req)
	if err != nil {
		if err.Error() == "phone is taken" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "phone is taken"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	user, err := h.userService.GetUserProfile(c.Context(), userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Server Error"})
	}

	return c.JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userId, err := parseUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Bad Request"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Validation error"})
	}

	file, err := h.fileService.GetFileById(c.Context(), req.FileId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req.FileURI = file.Url
	req.FileThumbnailURI = file.ThumbnailUrl

	err = h.userService.UpdateUser(c.Context(), userId, req)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "fileId is not valid / exists"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.userService.GetUserProfile(c.Context(), userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}
