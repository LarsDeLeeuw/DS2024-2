package gateway

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func RunApp() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c *fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, gateway!")
	})

	app.Post("/login", postLogin)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}

func postLogin(c *fiber.Ctx) (err error) {
	a := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(a)

	a.Request().Header.SetMethod(fiber.MethodPost)
	a.Request().SetRequestURI("http://auth:3001/api/v1/login")
	args := fiber.AcquireArgs()
	args.Set("username", c.FormValue("username"))
	args.Set("password", c.FormValue("password"))

	a.Form(args)
	// ...
	if err := a.Parse(); err != nil {
		panic(err)
	}
	fiber.ReleaseArgs(args)

	statusCode, body, errs := a.Bytes() // ..
	if len(errs) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	// pass status code and body received by the proxy
	return c.Status(statusCode).Send(body)
}
