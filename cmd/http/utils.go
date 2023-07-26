package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
)

var csrfProtection = csrf.New(csrf.Config{
	// only to control the switch whether csrf is activated or not
	Next: func(c *fiber.Ctx) bool {
		return csrfActivated
	},
	KeyLookup:      "form:_csrf",
	CookieName:     "csrf_",
	CookieSameSite: "Strict",
	Expiration:     1 * time.Hour,
	KeyGenerator:   utils.UUID,
	ContextKey:     "token",
})

func requireLogin(c *fiber.Ctx) error {
	currSession, err := sessionStore.Get(c)
	if err != nil {
		return err
	}
	user := currSession.Get("User")
	defer currSession.Save()

	if user == nil {
		// This request is from a user that is not logged in.
		// Send them to the login page.
		return c.Redirect("/login")
	}

	// If we got this far, the request is from a logged-in user.
	// Continue on to other middleware or routes.
	return c.Next()
}



