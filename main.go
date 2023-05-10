package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("welcome to api")
}

func setupRoutes(app *fiber.App) {
	app.Get("/api", welcome)
}

func main() {

	app := fiber.New()
	setupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
