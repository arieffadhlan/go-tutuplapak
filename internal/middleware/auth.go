package middleware

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ErrorHandler: jwtError,
		SuccessHandler: func(c *fiber.Ctx) error {
			// Mendapatkan token yang sudah diverifikasi dari konteks
			user := c.Locals("user").(*jwt.Token)

			// Mengekstrak klaim dari token
			claims := user.Claims.(jwt.MapClaims)

			// Menyimpan klaim di c.locals untuk digunakan di handler berikutnya
			c.Locals("userId", claims["user_id"])
			// Tambahkan klaim lain yang Anda butuhkan

			return c.Next()
		},
	})
}

// Protected returns a middleware that checks if user is admin or not
func ProtectedBeli(requireAdmin bool) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ErrorHandler: jwtError,
		SuccessHandler: func(c *fiber.Ctx) error {
			// Get verified token
			user := c.Locals("user").(*jwt.Token)

			// Extract claims
			claims := user.Claims.(jwt.MapClaims)

			// Store username in locals
			c.Locals("username", claims["username"])

			// Extract isAdmin claim
			isAdmin, ok := claims["isAdmin"].(bool)
			if !ok {
				return fiber.NewError(fiber.StatusForbidden, "Invalid role claim")
			}

			// Check against required role
			if requireAdmin && !isAdmin {
				return fiber.NewError(fiber.StatusForbidden, "Admin access required")
			}
			if !requireAdmin && isAdmin {
				return fiber.NewError(fiber.StatusForbidden, "User access only")
			}

			return c.Next()
		},
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
