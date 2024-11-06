package main

import (
	"myapi/database"
	"myapi/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	database.ConnectDB()

	routes.Routes(app)

	app.Listen(":3000")
}
