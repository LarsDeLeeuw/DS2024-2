package auth

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func RunApp() {
	// Initialize a new Fiber app
	app := fiber.New()

	api := app.Group("/api") // /api

	v1 := api.Group("/v1")       // /api/v1
	v1.Post("/login", postLogin) // /api/v1/login

	// Start the server on port 3001
	log.Fatal(app.Listen(":3001"))
}

func postLogin(c *fiber.Ctx) (err error) {
	return c.Send(c.BodyRaw())
}
