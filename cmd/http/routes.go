package main

import (
	"os"

	"github.com/Ed-cred/SolarPal/internal/handlers"
	"github.com/Ed-cred/SolarPal/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	
	csrfActivated = true
)

func init() {
	// this mean, csrf is activated
	csrfActivated = len(os.Args) > 1 && os.Args[1] == "withoutCsrf"
}

func setupRoutes(app *fiber.App) {
	validLogins := []models.User{
		{Username: "bob",Email:"bob@test.com", Password: "test"},
		{Username: "alice",Email:"alice@test.com", Password: "test"},
	}
	app.Use(recover.New())

	app.Get("/", requireLogin, csrfProtection, func(c *fiber.Ctx) error {
		currSession, err := sessionStore.Get(c)
		if err != nil {
			return err
		}
		sessionUser := currSession.Get("User").(fiber.Map)
		// release the currSession
		err = currSession.Save()
		if err != nil {
			return err
		}

		if sessionUser["Name"] == "" {
			return c.Status(fiber.StatusBadRequest).SendString("User is empty")
		}
		username := sessionUser["Name"].(string)

		return c.JSON(fiber.Map{
			"username":  username,
			"csrfToken": c.Locals("token"),
		})
	})
	app.Get("/login", func(c *fiber.Ctx) error {
		c.SendString("Please enter your credentials")
		return c.Redirect("/login")
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		user := &models.User{}
		err := c.BodyParser(user)
		if err != nil {
			return err
		}

		if user.Username == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Username is required.")
		}

		if user.Password == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Password is required.")
		}

		if !findUser(validLogins, user) {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid username or password.")
		}

		// Valid login.
		// Create a new currSession and save their user data in the currSession.
		currSession, err := sessionStore.Get(c)
		defer currSession.Save()
		if err != nil {
			return err
		}
		err = currSession.Regenerate()
		if err != nil {
			return err
		}
		currSession.Set("User", fiber.Map{"Name": user.Username})

		return c.Redirect("/")
	})
	app.Get("/render",requireLogin, handlers.Repo.GetPowerEstimate)

	app.Post("/add",requireLogin,csrfProtection, handlers.Repo.AddSolarArray)
}
