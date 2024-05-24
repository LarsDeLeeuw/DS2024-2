package auth

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"database/sql"

	_ "github.com/lib/pq"
)

type User struct {
	id       int
	username string
	password string
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

	v1.Post("/login", postLogin)
	v1.Post("/register", postRegister)

	// Start the server on port 3001
	log.Fatal(app.Listen(":3001"))
}

func postRegister(c *fiber.Ctx) (err error) {
	// Check if handle on DB still valid
	pingErr := DBHandle.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	var count int
	if err := DBHandle.QueryRow("SELECT COUNT(username) FROM users WHERE username = $1::TEXT;", c.FormValue("username")).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			log.Println("Count of users somehow failed")
			return c.Status(500).SendString("Internal server error")
		}
		log.Println(err)
		return c.Status(500).SendString("Internal server error")
	}

	if count != 0 {
		log.Printf("Tried registering with existing username: %s", c.FormValue("username"))
		return c.Status(400).SendString("Username not availabe.")
	}

	result, err := DBHandle.Exec("INSERT INTO users(username, password) VALUES ($1::TEXT, $2::TEXT);", c.FormValue("username"), c.FormValue("password"))
	if err != nil {
		log.Println(err)
		return c.Status(500).SendString("Internal server error")
	}
	log.Println(result)
	return c.Status(fiber.StatusOK).SendString("Registration succeeded")
}

func postLogin(c *fiber.Ctx) (err error) {
	// Check if handle on DB still valid
	pingErr := DBHandle.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	var password string
	if err := DBHandle.QueryRow("SELECT password FROM users WHERE username = $1::TEXT;", c.FormValue("username")).Scan(&password); err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Tried non existing username: %s", c.FormValue("username"))
			return c.Status(400).SendString("Username doesn't exist.")
		}
		fmt.Println(err)
		return c.Send(c.BodyRaw())
	}
	if password != c.FormValue("password") {
		fmt.Printf("Tried wrong password %s was expecting %s", c.FormValue("password"), password)
		return c.Status(400).SendString("Incorrect password.")
	}

	return c.Status(fiber.StatusOK).SendString("Login succeeded")
}
