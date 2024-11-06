package authRoutes

import (
	authHandler "myapi/handlers/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Post("/", authHandler.Auth)
}
