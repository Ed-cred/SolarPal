package main

import (
	"os"

	"github.com/Ed-cred/SolarPal/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var csrfActivated = true

func init() {
	// this mean, csrf is activated
	csrfActivated = len(os.Args) > 1 && os.Args[1] == "withoutCsrf"
}

func setupRoutes(app *fiber.App) {
	
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
		return c.Redirect("/login", fiber.StatusNetworkAuthenticationRequired)
	})

	app.Post("/signup", handlers.Repo.RegisterUser)
	app.Post("/login", handlers.Repo.LoginUser)

	app.Get("/render", requireLogin, handlers.Repo.GetPowerEstimate)
	app.Post("/add", requireLogin, csrfProtection, handlers.Repo.AddSolarArray)
}
