package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

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

func MakeAPIRequest(inputs models.RequiredInputs, opts models.OptionalInputs) (*models.PowerEstimate, error) {
	config.LoadEnv()
	apiKey := config.GetEnv("API_KEY")
	queryParams := url.Values{}
	queryParams.Add("api_key", apiKey)
	queryParams.Add("azimuth", inputs.Azimuth)
	queryParams.Add("system_capacity", inputs.SystemCapacity)
	queryParams.Add("losses", inputs.Losses)
	queryParams.Add("array_type", inputs.ArrayType)
	queryParams.Add("module_type", inputs.ModuleType)
	queryParams.Add("tilt", inputs.Tilt)
	queryParams.Add("address", inputs.Adress)
	if (models.OptionalInputs{}) != opts {
		queryParams.Add("gcr", opts.Gcr)
		queryParams.Add("dc_ac_ratio", opts.DcAcRatio)
		queryParams.Add("inv_eff", opts.InvEff)
		queryParams.Add("radius", opts.Radius)
		queryParams.Add("dataset", opts.Dataset)
		queryParams.Add("soiling", opts.Soiling)
		queryParams.Add("albedo", opts.Albedo)
		queryParams.Add("bifaciality", opts.Bifaciality)
	}

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

type Response struct {
	value *models.PowerEstimate
	error   error
}

// GetPowerEstimate makes the API request and sens the response as JSON
func (r *Repository) GetPowerEstimate(c *fiber.Ctx) error {
	c.SetUserContext(r.Cfg.Ctx)
	ctx := c.UserContext()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*2500)
	defer cancel()
	respch := make(chan Response)
	inputs := models.RequiredInputs{
		Azimuth:        "180",
		SystemCapacity: "4",
		Losses:         "14",
		ArrayType:      "1",
		ModuleType:     "0",
		Tilt:           "10",
		Adress:         "boulder, co",
	}
	opts := models.OptionalInputs{
		Gcr:         "0.4",
		DcAcRatio:   "1.2",
		InvEff:      "96.0",
		Radius:      "0",
		Dataset:     "nsrdb",
		Soiling:     "12|4|45|23|9|99|67|12.54|54|9|0|7.6",
		Albedo:      "0.3",
		Bifaciality: "0.7",
	}
	go func () {
		pvWattsResponse, err := MakeAPIRequest(inputs, opts)
		respch <- Response{
			value: pvWattsResponse,
			error: err,
		}
	}()
	
	for {
		select  {
		case <- ctx.Done():
			return errors.New("fetching api data took too long")
		case resp := <-respch:
			c.JSON(resp.value)
			return resp.error
		}
	}
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

	if helpers.FindUser(validLogins, user) != 0 {
		return c.Status(fiber.StatusBadRequest).SendString("This user is already registered.")
	}
	err = r.DB.CreateUser(user)
	if err != nil {
		log.Println("Error creating user")
		return err
	}
	return c.Redirect("/login")
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

	if user.ID = helpers.FindUser(validLogins, user); user.ID == 0 {
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
	currSession.Set("User", fiber.Map{"ID": user.ID})

	return c.Redirect("/", fiber.StatusSeeOther)
}

func (r *Repository) AddSolarArray(c *fiber.Ctx) error {
	// ctx, cancel := c.WithTimeout(c.Background(), 2*time.Second)
	// defer cancel()
	// stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
	// 		values ($1, $2, $3, $4, $5, $6, $7)`
	return c.SendString("Solar array has been added ")
}
