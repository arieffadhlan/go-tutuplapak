// package route

// import "github.com/gofiber/fiber/v2"

// func RegisterUserRoutes(router fiber.Router, userHandler handler.UserHandler) {
// 	router.Use(middleware.Authenticate())
// 	router.Get("/user", userHandler.GetUser)
// 	router.Put("/user", userHandler.UpdateUser)
// 	router.Post("/user/link/phone", userHandler.CreateLinkPhone)
// 	router.Post("/user/link/email", userHandler.UpdateLinkEmail)
// }
