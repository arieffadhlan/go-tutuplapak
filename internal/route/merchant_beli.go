package route

import (
	"tutuplapak-user/internal/handlers"
	"tutuplapak-user/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterMerchantBeliRoutes(router fiber.Router, merchantBeliHandler handlers.MerchantBeliHandler) {
	router.Post("/admin/merchants", middleware.ProtectedBeli(true), merchantBeliHandler.Create)
	router.Get("/admin/merchants", middleware.ProtectedBeli(true), merchantBeliHandler.Get)

	router.Post("/admin/merchants/:merchantId/items", middleware.ProtectedBeli(true), merchantBeliHandler.CreateItem)
	router.Get("/admin/merchants/:merchantId/items", middleware.ProtectedBeli(true), merchantBeliHandler.GetItem)
}
