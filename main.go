package main

import (
	"fmt"
	"log"
	"myapi/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	port := "5000"

	routes.Routes(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
