package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
)

const HeaderName = "X-Csrf-Token"
var csrfProtection = csrf.New(csrf.Config{
	// only to control the switch whether csrf is activated or not
	Next: func(c *fiber.Ctx) bool {
		return csrfActivated
	},
	KeyLookup:      "header:" + HeaderName,
	CookieName:     "csrf_",
	CookieSameSite: "Lax",
	Expiration:     6 * time.Hour,
	KeyGenerator:   utils.UUID,
	ContextKey:     "token",
	Extractor: csrf.CsrfFromHeader(HeaderName),
})

func requireLogin(c *fiber.Ctx) error {
	currSession, err := cfg.Session.Get(c)
	if err != nil {
		return err
	}
	user := currSession.Get("User")
	defer currSession.Save()

	if user == nil {
		// This request is from a user that is not logged in.
		// Send them to the login page.
		return c.Status(fiber.StatusForbidden).SendString("Please log in first.")
	}
	return c.Next()
}



