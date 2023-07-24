package main

import (
	"github.com/Ed-cred/SolarPal/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/render", handlers.GetPowerEstimate)

	app.Post("/add", func(c *fiber.Ctx) error {
		return c.SendString("I'm a POST request for creating a new PV object to estimate")
	})
}