package gateway

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func RunApp() {
	// Initialize a new Fiber app
	app := fiber.New()
	log.SetOutput(os.Stderr)

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c *fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, gateway!")
	})

	app.Post("/login", postLogin)
	app.Post("/register", postRegister)
	app.Post("/authenticate", postAuthenticate)
	app.Get("/calendar", getCalendar)
	app.Post("/event", postEvent)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}

func postRegister(c *fiber.Ctx) (err error) {
	a := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(a)

	a.Request().Header.SetMethod(fiber.MethodPost)
	a.Request().SetRequestURI("http://auth:3001/api/v1/register")
	args := fiber.AcquireArgs()
	args.Set("username", c.FormValue("username"))
	args.Set("password", c.FormValue("password"))

	a.Form(args)
	// ...
	if err := a.Parse(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error in API gateway")
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
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error in API gateway")
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

func postAuthenticate(c *fiber.Ctx) (err error) {
	// For simplicity of project repurpose the login endpoint.
	a := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(a)

	a.Request().Header.SetMethod(fiber.MethodPost)
	a.Request().SetRequestURI("http://auth:3001/api/v1/login")
	args := fiber.AcquireArgs()
	args.Set("username", c.FormValue("username"))
	args.Set("password", c.FormValue("password"))

	a.Form(args)
	if err := a.Parse(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error in API gateway")
	}
	fiber.ReleaseArgs(args)

	statusCode, _, errs := a.Bytes() // ..
	if len(errs) > 0 {
		return c.Status(fiber.StatusServiceUnavailable).SendString("Authentication service is currently unavailable.")
	}
	if statusCode == fiber.StatusOK {
		return c.Status(fiber.StatusOK).SendString("Credentials are valid.")
	} else {
		return c.Status(fiber.StatusUnauthorized).SendString("Credentials are invalid.")
	}
}

func getCalendar(c *fiber.Ctx) (err error) {
	a := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(a)

	a.Request().Header.SetMethod(fiber.MethodGet)
	a.Request().SetRequestURI("http://calendar:3002/api/v1/calendar")
	args := fiber.AcquireArgs()
	args.Set("username", c.FormValue("username"))
	args.Set("password", c.FormValue("password"))

	a.Form(args)
	if err := a.Parse(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error in API gateway")
	}
	fiber.ReleaseArgs(args)

	statusCode, body, errs := a.Bytes()
	fmt.Println(errs)
	if len(errs) > 0 {
		return c.Status(fiber.StatusServiceUnavailable).SendString("Authentication service is currently unavailable.")
	}
	// pass status code and body received by the proxy
	return c.Status(statusCode).Send(body)
}

func postEvent(c *fiber.Ctx) (err error) {
	a := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(a)

	a.Request().Header.SetMethod(fiber.MethodPost)
	a.Request().SetRequestURI("http://calendar:3002/api/v1/events")
	args := fiber.AcquireArgs()
	args.Set("username", c.FormValue("username"))
	args.Set("password", c.FormValue("password"))

	a.Form(args)
	if err := a.Parse(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error in API gateway")
	}
	fiber.ReleaseArgs(args)

	statusCode, body, errs := a.Bytes()
	fmt.Println(errs)
	if len(errs) > 0 {
		return c.Status(fiber.StatusServiceUnavailable).SendString("Authentication service is currently unavailable.")
	}
	// pass status code and body received by the proxy
	return c.Status(statusCode).Send(body)
}
