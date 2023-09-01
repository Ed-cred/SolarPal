package handlers

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Ed-cred/SolarPal/config"
	"github.com/Ed-cred/SolarPal/repository/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
)

const HeaderName = "X-Csrf-Token"

var (
	sessionStore  *session.Store
	cfg           config.AppConfig
	csrfActivated = true
)

func TestMain(m *testing.M) {
	sessionStore = session.New()
	sessionStore.RegisterType(fiber.Map{})
	cfg.Session = sessionStore
	ctx := context.Background()
	cfg.Ctx = ctx
	config.LoadEnv()
	dbPath := config.GetEnv("SQLITE_PATH")
	log.Println("Connecting to database...")
	db, err := database.ConnectDb(dbPath)
	if err != nil {
		log.Fatal("Couldn't connect to database:", err)
	}
	repo := NewRepository(&cfg, db)
	NewHandlers(repo)
	os.Exit(m.Run())
}

func setupRoutes(app *fiber.App) {
	app.Use(recover.New())
	app.Post("/signup", Repo.RegisterUser)
	app.Post("/login", Repo.LoginUser)
	app.Get("/logout", Repo.LogoutUser)

	app.Get("/", requireLogin, csrfProtection, Repo.DisplayAvailableData)

	app.Get("/render/:array_id", requireLogin, Repo.GetPowerEstimate)
	app.Post("/add", requireLogin, csrfProtection, Repo.AddSolarArray)
	app.Put("/update/:array_id", requireLogin, csrfProtection, Repo.UpdateSolarArrayParams)
	app.Delete("/remove/:array_id", requireLogin, csrfProtection, Repo.RemoveSolarArray)
}
func init() {
	csrfActivated = len(os.Args) > 4 && os.Args[4] == "withoutCsrf"
}

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
	Extractor:      csrf.CsrfFromHeader(HeaderName),
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
