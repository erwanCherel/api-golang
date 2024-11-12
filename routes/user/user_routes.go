package userRoutes

import (
	userHandler "myapi/handlers/user"
	"myapi/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(router fiber.Router) {
	public := router.Group("/")

	public.Get("/users", userHandler.GetUsers)
	public.Get("/user/:id", userHandler.GetUser)
	public.Post("/user", userHandler.CreateUser)
	public.Get("/user/:id/videos", userHandler.GetVideosByUser)

	private := router.Group("/private", middleware.AuthMiddleware)

	private.Put("/user/:id", userHandler.UpdateUser)
	private.Delete("/user/:id", userHandler.DeleteUser)
	private.Post("/user/:id/video", userHandler.CreateVideo)
}
