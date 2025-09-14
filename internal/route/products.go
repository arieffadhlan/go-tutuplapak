package route

import (
	"tutuplapak-user/internal/handlers"
	"tutuplapak-user/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterProductsRoutes(router fiber.Router, productHandler *handlers.ProductsHandler) {
	router.Get("/product", productHandler.GetAllProducts)

	router.Post("/product", middleware.Protected(), productHandler.CreateProduct)
	router.Put("/product/:productId", middleware.Protected(), productHandler.UpdateProduct)
	router.Delete("/product/:productId", middleware.Protected(), productHandler.DeleteProduct)
}
