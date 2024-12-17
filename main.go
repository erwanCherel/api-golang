package main

import (
	"myapi/database"
	"myapi/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		StreamRequestBody: true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PATCH, DELETE",
	}))

	database.ConnectDB()

	router.SetupRoutes(app)

	app.Listen(":5001")
}
