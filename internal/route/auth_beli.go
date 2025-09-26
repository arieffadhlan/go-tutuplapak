package route

import (
	"tutuplapak-user/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthBeliRoutes(router fiber.Router, authBeliHendler handlers.AuthBeliHandler) {
	router.Post("/admin/register", authBeliHendler.RegisterAdmin)
	router.Post("/users/register", authBeliHendler.RegisterUser)

	router.Post("/admin/login", authBeliHendler.Login)
	router.Post("/users/login", authBeliHendler.Login)
}
