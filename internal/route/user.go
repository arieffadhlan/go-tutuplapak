package route

// func RegisterUserRoutes(router fiber.Router, userHandler handler.UserHandler) {
// 	router.Use(middleware.Authenticate())
// 	router.Get("/user", userHandler.GetUser)
// 	router.Put("/user", userHandler.UpdateUser)
// 	router.Post("/user/link/phone", userHandler.CreateLinkPhone)
// 	router.Post("/user/link/email", userHandler.UpdateLinkEmail)
// }

import (
	"tutuplapak-user/internal/handlers"
	"tutuplapak-user/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(router fiber.Router, userHandler *handlers.UserHandler) {
	router.Get("/user", middleware.Protected(), userHandler.GetUser)
	router.Post("/user/link/email", middleware.Protected(), userHandler.LinkEmail)
	router.Post("/user/link/phone", middleware.Protected(), userHandler.LinkPhone)
	router.Put("/user", middleware.Protected(), userHandler.UpdateUser)
}
