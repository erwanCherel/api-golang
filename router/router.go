package router

import (
	authRoutes "myapi/routes/auth"
	healthRoutes "myapi/routes/health"
	userRoutes "myapi/routes/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api/", logger.New())

	healthRoutes.SetupHealthRoutes(api)
	userRoutes.SetupUserRoutes(api)
	authRoutes.SetupAuthRoutes(api)
}
