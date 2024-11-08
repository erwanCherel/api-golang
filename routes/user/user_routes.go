package userRoutes

import (
	userHandler "myapi/handlers/user"
	"myapi/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(router fiber.Router) {

	router.Get("/users", userHandler.GetUsers)
	router.Get("/user/:id", userHandler.GetUser)
	router.Post("/user", userHandler.CreateUser)

	userProtected := router.Group("/protected", middleware.AuthMiddleware)

	userProtected.Put("/user/:id", userHandler.UpdateUser)
	userProtected.Delete("/user/:id", userHandler.DeleteUser)
}
