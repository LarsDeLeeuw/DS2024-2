package calendar

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"database/sql"

	_ "github.com/lib/pq"
)

type Event struct {
	Id        int
	Title     string
	Date      string
	Organizer string
	Public    bool
}

var DBHandle *sql.DB

func RunApp() {
	// Initialize a new Fiber app
	app := fiber.New()
	log.SetOutput(os.Stderr)

	// Setup connection to AuthDB
	dat, err := os.ReadFile(os.Getenv("POSTGRES_PASSWORD_FILE"))
	if err != nil {
		log.Fatal(err)
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), string(dat), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))
	DBHandle, err = sql.Open("postgres", connStr)
	fmt.Println(connStr)
	if err != nil {
		log.Fatal(err)
	}

	api := app.Group("/api") // /api

	v1 := api.Group("/v1") // /api/v1

	v1.Get("/calendar", getCalendar)
	v1.Post("/events", postEvent)

	// Start the server on port 3002
	log.Fatal(app.Listen(":3002"))
}

func getCalendar(c *fiber.Ctx) (err error) {
	// Check if handle on DB still valid
	pingErr := DBHandle.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	a := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(a)

	a.Request().Header.SetMethod(fiber.MethodPost)
	a.Request().SetRequestURI("http://gateway:3000/authenticate")
	args := fiber.AcquireArgs()
	args.Set("username", c.FormValue("username"))
	args.Set("password", c.FormValue("password"))

	a.Form(args)
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
	if statusCode == fiber.StatusUnauthorized {
		// Unauthorized
		return c.Status(statusCode).Send(body)
	} else if statusCode != fiber.StatusOK {
		// Unable to authorize atm
		fmt.Println(statusCode)
		return c.Status(fiber.StatusServiceUnavailable).SendString("Unable to authorize, try again later.")
	}
	// Authorized
	// Check if calendar exists, if not make one.
	var count int
	if err := DBHandle.QueryRow("SELECT COUNT(owner) FROM calendars WHERE owner = $1::TEXT;", c.FormValue("username")).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Count of owners somehow failed")
			return c.Status(500).SendString("Internal server error")
		}
		fmt.Println(err)
		return c.Status(500).SendString("Internal server error")
	}
	if count == 0 {
		fmt.Printf("Making new calendar for: %s", c.FormValue("username"))
		result, err := DBHandle.Exec("INSERT INTO calendars(owner) VALUES ($1::TEXT);", c.FormValue("username"))
		if err != nil {
			fmt.Println(err)
			return c.Status(500).SendString("Internal server error")
		}
		fmt.Println(result)
	}

	rows, err := DBHandle.Query("SELECT e.event_id::INT, e.title::TEXT, e.date::TEXT, e.organizer::TEXT, e.public::BOOLEAN FROM events AS e NATURAL JOIN calendars AS c NATURAL JOIN calendareventsrelationship WHERE owner = $1::TEXT;", c.FormValue("username"))
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Count of owners somehow failed")
			return c.Status(500).SendString("Internal server error")
		}
		fmt.Println(err)
		return c.Status(500).SendString("Internal server error")
	}
	defer rows.Close()

	var events = make([]Event, 0)

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.Id, &e.Title, &e.Date,
			&e.Organizer, &e.Public); err != nil {
			fmt.Println(err)
			return c.Status(500).SendString("Internal server error")
		}
		events = append(events, e)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Internal server error")
	}
	return c.Status(fiber.StatusOK).JSON(events)
}

func postEvent(c *fiber.Ctx) (err error) {
	// Check if handle on DB still valid
	pingErr := DBHandle.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	a := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(a)

	a.Request().Header.SetMethod(fiber.MethodPost)
	a.Request().SetRequestURI("http://gateway:3000/authenticate")
	args := fiber.AcquireArgs()
	args.Set("username", c.FormValue("username"))
	args.Set("password", c.FormValue("password"))

	a.Form(args)
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
	if statusCode == fiber.StatusUnauthorized {
		// Unauthorized
		return c.Status(statusCode).Send(body)
	} else if statusCode != fiber.StatusOK {
		// Unable to authorize atm
		fmt.Println(statusCode)
		return c.Status(fiber.StatusServiceUnavailable).SendString("Unable to authorize, try again later.")
	}
	// Authorized
	var calendar_id int
	if err := DBHandle.QueryRow("SELECT calendar_id FROM calendars WHERE owner = $1::TEXT;", c.FormValue("username")).Scan(&calendar_id); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No calendar exist for user")
			return c.Status(500).SendString("Internal server error")
		}
		fmt.Println(err)
		return c.Status(500).SendString("Internal server error")
	}
	result, err := DBHandle.Exec("INSERT INTO events(title, date, organizer, public) VALUES ($1::TEXT, $2::TEXT, $3::TEXT, $4::BOOLEAN);", c.FormValue("title"), c.FormValue("date"), c.FormValue("username"), c.FormValue("public") != "true")
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Internal server error")
	}
	fmt.Println(result)
	var event_id int
	if err := DBHandle.QueryRow("SELECT e.event_id::INT FROM events AS e WHERE organizer = $1::TEXT AND title = $2::TEXT AND date = $3::TEXT;", c.FormValue("username"), c.FormValue("title"), c.FormValue("date")).Scan(&event_id); err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Error retrieving event_id of event with title: %s", c.FormValue("title"))
			return c.Status(400).SendString("Username doesn't exist.")
		}
		fmt.Println(err)
		return c.Send(c.BodyRaw())
	}
	// Create relation organizer and new event
	result, err = DBHandle.Exec("INSERT INTO calendareventsrelationship(calendar_id, event_id) VALUES ($1::INT, $2::INT);", calendar_id, event_id)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Internal server error")
	}
	fmt.Println(result)

	return c.Status(fiber.StatusOK).SendString(string(event_id))
}
