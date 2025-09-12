package route

import (
	"tutuplapak-user/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterFileRoutes(router fiber.Router, fileHandler handlers.FileHandler) {
	router.Post("/file", fileHandler.Post)
}
