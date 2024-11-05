package routes

import (
	"myapi/handlers"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Get("/health", handlers.HealthCheck)
	// api.Get("/users", handlers.GetUsers)
}
