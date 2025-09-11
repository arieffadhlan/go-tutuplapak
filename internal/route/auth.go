package route

import (
	"tutuplapak-user/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(router fiber.Router, authHandler handlers.AuthHandler) {
	router.Post("/login/email", authHandler.LoginByEmail)
	router.Post("/login/phone", authHandler.LoginByPhone)

	router.Post("/register/email", authHandler.RegisterByEmail)
	// router.Post("/register/phone", userHandler.RegisterByPhone)
	// router.Post("/login/email", userHandler.LoginByEmail)
}
