package userRoutes

import (
	userHandler "myapi/handlers/user"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(router fiber.Router) {
	user := router.Group("/users")

	// Create a User
	user.Post("/", userHandler.CreateUser)
	// Read all Users
	user.Get("/", userHandler.GetUsers)
	// // Read one User
	user.Get("/:id", userHandler.GetUser)
	// // Update one User
	// User.Put("/:userId", userHandler.UpdateUser)
	// // Delete one User
	// User.Delete("/:userId", userHandler.DeleteUser)
}
