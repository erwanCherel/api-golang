package userRoutes

import (
	userHandler "myapi/handlers/user"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(router fiber.Router) {
	user := router.Group("/user")

	user.Post("/", userHandler.CreateUser)
	user.Get("/", userHandler.GetUsers)
	user.Get("/:id", userHandler.GetUser)
	// User.Put("/:userId", userHandler.UpdateUser)
	// User.Delete("/:userId", userHandler.DeleteUser)
}
