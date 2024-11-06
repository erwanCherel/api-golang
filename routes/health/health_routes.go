package healthRoutes

import (
	handlers "myapi/handlers/health"

	"github.com/gofiber/fiber/v2"
)

func SetupHealthRoutes(router fiber.Router) {
	health := router.Group("/health")

	health.Get("/", handlers.HealthCheck)
}
