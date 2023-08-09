package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Ed-cred/SolarPal/config"
	"github.com/Ed-cred/SolarPal/internal/handlers"
	"github.com/Ed-cred/SolarPal/repository/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var (
	sessionStore *session.Store
	cfg          config.AppConfig
)

func main() {
	app := fiber.New()
	db, err := run()
	if err != nil {
		panic(err)
	}
	log.Println("Connected to database")
	defer db.SQL.Close()
	app.Use(logger.New())
	setupRoutes(app)
	fmt.Printf("Server started and listening at localhost:5000 - csrfActive: %v\n", len(os.Args) > 4 && os.Args[4] == "withoutCsrf")
	// Start server
	log.Fatal(app.Listen(":5000"))
}

func run() (*database.DB, error) {
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
		return nil, err
	}
	repo := handlers.NewRepository(&cfg, db)
	handlers.NewHandlers(repo)
	return db, nil
}
