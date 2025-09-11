package handlers

import (
	"errors"
	"net/mail"
	"regexp"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/repository"
	"tutuplapak-user/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) AuthHandler {
	return AuthHandler{service: service}
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

var phoneRegex = regexp.MustCompile(`^\+[0-9]+$`)

func isPhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

func (h AuthHandler) LoginByEmail(c *fiber.Ctx) error {
	input := dto.AuthEmailRequest{}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	if !isEmail(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid email"})
	}

	tkn, user, err := h.service.LoginByEmail(c.Context(), input)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		if errors.Is(err, services.ErrInvalidCredentials) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email or password"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// In your handler
	var phoneValue string
	if user.Phone != nil {
		phoneValue = *user.Phone
	}

	authResp := dto.AuthEmailResponse{
		Email: *user.Email,
		Phone: phoneValue,
		Token: tkn,
	}

	return c.Status(fiber.StatusOK).JSON(authResp)
}

func (h AuthHandler) LoginByPhone(c *fiber.Ctx) error {
	input := dto.AuthPhoneRequest{}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	if !isPhone(input.Phone) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid phone format"})
	}

	tkn, user, err := h.service.LoginByPhone(c.Context(), input)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		if errors.Is(err, services.ErrInvalidCredentials) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid phone or password"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// In your handler
	var emailValue string
	if user.Email != nil {
		emailValue = *user.Email
	}

	authResp := dto.AuthPhoneResponse{
		Email: emailValue,
		Phone: *user.Phone,
		Token: tkn,
	}

	return c.Status(fiber.StatusOK).JSON(authResp)
}

func (h AuthHandler) RegisterByEmail(c *fiber.Ctx) error {
	input := dto.AuthEmailRequest{}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on register request", "data": err})
	}

	if !isEmail(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid email"})
	}

	tkn, user, err := h.service.RegisterByEmail(c.Context(), input)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// In your handler
	var phoneValue string
	if user.Phone != nil {
		phoneValue = *user.Phone
	}

	authResp := dto.AuthEmailResponse{
		Email: *user.Email,
		Phone: phoneValue,
		Token: tkn,
	}

	return c.Status(fiber.StatusCreated).JSON(authResp)
}
