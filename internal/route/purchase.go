package route

import (
	"tutuplapak-user/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterPurchaseRoutes(router fiber.Router, handler *handlers.PurchaseHandler) {
	router.Post("/purchase", handler.Purchase)
	router.Post("/purchase/:purchaseId", handler.PurchasePaymentProof)
}
