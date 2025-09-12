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
	// All user routes require authentication
	userRouter := router.Use(middleware.Protected())

	userRouter.Get("/user", userHandler.GetUser)
	userRouter.Post("/user/link/email", userHandler.LinkEmail)
	userRouter.Post("/user/link/phone", userHandler.LinkPhone)
	userRouter.Put("/user", userHandler.UpdateUser)
}
