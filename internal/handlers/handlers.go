package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/Ed-cred/SolarPal/config"
	"github.com/Ed-cred/SolarPal/internal/helpers"
	"github.com/Ed-cred/SolarPal/internal/models"
	"github.com/Ed-cred/SolarPal/repository"
	"github.com/Ed-cred/SolarPal/repository/database"
	"github.com/gofiber/fiber/v2"
)

type Repository struct {
	Cfg *config.AppConfig
	DB  repository.DBRepo
}

var Repo *Repository

func NewRepository(cfg *config.AppConfig, db *database.DB) *Repository {
	return &Repository{
		Cfg: cfg,
		DB:  database.NewSQLiteRepo(db.SQL, cfg),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

const baseURL = "https://developer.nrel.gov/api/pvwatts/v8.json"

func MakeAPIRequest(address string) (*models.PowerEstimate, error) {
	config.LoadEnv()
	apiKey := config.GetEnv("API_KEY")
	queryParams := url.Values{}
	queryParams.Add("api_key", apiKey)
	queryParams.Add("azimuth", "180")
	queryParams.Add("system_capacity", "4")
	queryParams.Add("losses", "14")
	queryParams.Add("array_type", "1")
	queryParams.Add("module_type", "0")
	queryParams.Add("gcr", "0.4")
	queryParams.Add("dc_ac_ratio", "1.2")
	queryParams.Add("inv_eff", "96.0")
	queryParams.Add("radius", "0")
	queryParams.Add("dataset", "nsrdb")
	queryParams.Add("tilt", "10")
	queryParams.Add("address", address)
	queryParams.Add("soiling", "12|4|45|23|9|99|67|12.54|54|9|0|7.6")
	queryParams.Add("albedo", "0.3")
	queryParams.Add("bifaciality", "0.7")

	apiEndpoint := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	resp, err := http.Get(apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to make the request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var pvWattsResponse models.PowerEstimate
	err = json.NewDecoder(resp.Body).Decode(&pvWattsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode API response: %v", err)
	}

	return &pvWattsResponse, nil
}

// GetPowerEstimate makes the API request and sens the response as JSON
func (r *Repository) GetPowerEstimate(c *fiber.Ctx) error {
	address := "boulder, co" // You can change this to any desired location
	pvWattsResponse, err := MakeAPIRequest(address)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Error fetching data from the API")
	}

	c.JSON(pvWattsResponse)
	return nil
}

func (r *Repository) RegisterUser(c *fiber.Ctx) error {
	user := &models.User{}
	validLogins, err := r.DB.GetUsers()
	if err != nil {
		return err
	}
	err = c.BodyParser(user)
	if err != nil {
		return err
	}

	if user.Username == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Username is required.")
	}

	if user.Password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Password is required.")
	}
	if user.Email == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Email is required.")
	}

	if helpers.FindUser(validLogins, user) {
		return c.Status(fiber.StatusBadRequest).SendString("This user is already registered.")
	}
	err = r.DB.CreateUser(user)
	if err != nil {
		log.Println("Error creating user")
		return err
	}
	return nil
}

func (r *Repository) LoginUser(c *fiber.Ctx) error {
	user := &models.User{}
	validLogins, err := r.DB.GetUsers()
	if err != nil {
		log.Println("Error getting users")
		return err
	}
	err = c.BodyParser(user)
	if err != nil {
		return err
	}

	if user.Username == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Username is required.")
	}

	if user.Password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Password is required.")
	}

	if !helpers.FindUser(validLogins, user) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid username or password.")
	}

	// Valid login.
	// Create a new currSession and save their user data in the currSession.
	currSession, err := r.Cfg.Session.Get(c)
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
}

func (r *Repository) AddSolarArray(c *fiber.Ctx) error {
	// ctx, cancel := c.WithTimeout(c.Background(), 2*time.Second)
	// defer cancel()
	// stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
	// 		values ($1, $2, $3, $4, $5, $6, $7)`
	return c.SendString("Solar array has been added ")
}
