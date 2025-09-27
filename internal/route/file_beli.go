package route

import (
	"tutuplapak-user/internal/handlers"
	"tutuplapak-user/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterFileBeliRoutes(router fiber.Router, authHandler handlers.FileBeliHandler) {
	router.Post("/image", middleware.ProtectedBeli(true), authHandler.Post)
}
