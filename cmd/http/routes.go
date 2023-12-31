package main

import (
	"os"

	"github.com/Ed-cred/SolarPal/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var csrfActivated = true

func init() {
	csrfActivated = len(os.Args) > 4 && os.Args[4] == "withoutCsrf"
}

func setupRoutes(app *fiber.App) {
	app.Use(recover.New())
	app.Post("/signup", handlers.Repo.RegisterUser)
	app.Post("/login", handlers.Repo.LoginUser)
	app.Get("/logout", handlers.Repo.LogoutUser)

	app.Get("/", requireLogin, csrfProtection, handlers.Repo.DisplayAvailableData)

	app.Get("/render/:array_id", requireLogin, handlers.Repo.GetPowerEstimate)
	app.Post("/add", requireLogin, csrfProtection, handlers.Repo.AddSolarArray)
	app.Put("/update/:array_id", requireLogin, csrfProtection, handlers.Repo.UpdateSolarArrayParams)
	app.Delete("/remove/:array_id", requireLogin, csrfProtection, handlers.Repo.RemoveSolarArray)
}


