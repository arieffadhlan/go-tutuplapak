package route

import (
	"tutuplapak-user/internal/handlers"
	"tutuplapak-user/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterProductsRoutes(router fiber.Router, productHandler *handlers.ProductsHandler) {
	router.Get("/product", productHandler.GetAllProducts)

	privateRouter := router.Use(middleware.Protected())
	privateRouter.Post("/product", productHandler.CreateProduct)
	privateRouter.Put("/product/:productId", productHandler.UpdateProduct)
	privateRouter.Delete("/product/:productId", productHandler.DeleteProduct)
}
