package route

import (
	"tutuplapak-user/internal/handlers"
	"tutuplapak-user/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterMerchantBeliRoutes(router fiber.Router, merchantBeliHandler handlers.MerchantBeliHandler) {
	router.Post("/admin/merchants", middleware.ProtectedBeli(true), merchantBeliHandler.Create)
	router.Get("/admin/merchants", middleware.ProtectedBeli(true), merchantBeliHandler.Get)
}
