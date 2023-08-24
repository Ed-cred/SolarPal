package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Ed-cred/SolarPal/config"
	"github.com/Ed-cred/SolarPal/internal/helpers"
	"github.com/Ed-cred/SolarPal/internal/models"
	"github.com/Ed-cred/SolarPal/repository"
	"github.com/Ed-cred/SolarPal/repository/database"
	"github.com/gofiber/fiber/v2"
)

const baseURL = "https://developer.nrel.gov/api/pvwatts/v8.json"

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
	if err != nil {
		return err
	}
	err = currSession.Regenerate()
	if err != nil {
		return err
	}
	currSession.Set("User", fiber.Map{"ID": user.ID})
	currSession.Save()

	return c.SendString("Successfully logged in")
}

type Response struct {
	value *models.PowerEstimate
	error error
}

// GetPowerEstimate makes the API request and sens the response as JSON
func (r *Repository) GetPowerEstimate(c *fiber.Ctx) error {
	c.SetUserContext(r.Cfg.Ctx)
	ctx := c.UserContext()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	sarrayId := c.Params("array_id")
	arrayId, err := strconv.Atoi(sarrayId)
	if err != nil {
		return err
	}
	currSession, err := r.Cfg.Session.Get(c)
	if err != nil {
		return err
	}
	sessionUser := currSession.Get("User").(fiber.Map)
	id := sessionUser["ID"]
	inputs, opts, err := r.DB.FetchSolarArrayData(id.(uint), arrayId)
	log.Println("These are the database params for my solar array:", inputs, opts)
	respch := make(chan Response, 1)
	if err != nil {
		log.Println("Unable to fetch solar array data: ", err)
	}

	go func() {
		pvWattsResponse, err := makeAPIRequest(inputs, opts)
		respch <- Response{
			value: pvWattsResponse,
			error: err,
		}
		log.Println("Solar array data for array:", arrayId)
	}()

	for {
		select {
		case <-ctx.Done():
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
	return c.SendString("Account created successfully!")
}

func (r *Repository) LogoutUser(c *fiber.Ctx) error {
	currSession, err := r.Cfg.Session.Get(c)
	if err != nil {
		return err
	}

	// Clear the session data.
	currSession.Destroy()
	return c.SendString("Logged out successfully.")
}

func (r *Repository) AddSolarArray(c *fiber.Ctx) error {
	currSession, err := r.Cfg.Session.Get(c)
	if err != nil {
		log.Println("Unable to access session: ", err)
	}
	sessionUser := currSession.Get("User").(fiber.Map)
	id := sessionUser["ID"]
	inputs := &models.RequiredInputs{}
	opts := &models.OptionalInputs{}
	err = c.BodyParser(inputs)
	if err != nil {
		return err
	}
	err = c.BodyParser(opts)
	if err != nil {
		return err
	}
	if inputs.Azimuth == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Azimuth is required.")
	}
	if inputs.SystemCapacity == "" {
		return c.Status(fiber.StatusBadRequest).SendString("System capacity is required.")
	}
	if inputs.Losses == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Losses are required.")
	}
	if inputs.ArrayType == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Array type is required.")
	}
	if inputs.ModuleType == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Module type is required.")
	}
	if inputs.Tilt == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Tilt is required.")
	}
	if inputs.Address == "" && opts.Latitude == "" && opts.Longitude == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Location data is required.")
	}
	arrayId, err := r.DB.AddSolarArray(id.(uint), *inputs, *opts)
	if err != nil {
		return err
	}

	return c.SendString("Solar array has been added with id:" + fmt.Sprint(arrayId))
}

func (r *Repository) UpdateSolarArrayParams(c *fiber.Ctx) error {
	currSession, err := r.Cfg.Session.Get(c)
	if err != nil {
		log.Println("Unable to access session storage: ", err)
	}
	sessionUser := currSession.Get("User").(fiber.Map)
	id := sessionUser["ID"]
	sarrayId := c.Params("array_id")
	arrayId, err := strconv.Atoi(sarrayId)
	if err != nil {
		return err
	}
	existingInputs, existingOpts, err := r.DB.FetchSolarArrayData(id.(uint), arrayId)
	if err != nil {
		return err
	}
	inputs := &models.RequiredInputs{}
	opts := &models.OptionalInputs{}
	err = c.BodyParser(inputs)
	if err != nil {
		return err
	}
	if inputs.Azimuth == "" {
		inputs.Azimuth = existingInputs.Azimuth
	}
	if inputs.SystemCapacity == "" {
		inputs.SystemCapacity = existingInputs.SystemCapacity
	}
	if inputs.Losses == "" {
		inputs.Losses = existingInputs.Losses
	}
	if inputs.ArrayType == "" {
		inputs.ArrayType = existingInputs.ArrayType
	}
	if inputs.ModuleType == "" {
		inputs.ModuleType = existingInputs.ModuleType
	}
	if inputs.Tilt == "" {
		inputs.Tilt = existingInputs.Tilt
	}
	if inputs.Address == "" {
		inputs.Address = existingInputs.Address
	}
	err = c.BodyParser(opts)
	if err != nil {
		return err
	}
	if opts.Gcr == "" {
		opts.Gcr = existingOpts.Gcr
	}
	if opts.DcAcRatio == "" {
		opts.DcAcRatio = existingOpts.DcAcRatio
	}
	if opts.InvEff == "" {
		opts.InvEff = existingOpts.InvEff
	}
	if opts.Radius == "" {
		opts.Radius = existingOpts.Radius
	}
	if opts.Dataset == "" {
		opts.Dataset = existingOpts.Dataset
	}
	if opts.Soiling == "" {
		opts.Soiling = existingOpts.Soiling
	}
	if opts.Albedo == "" {
		opts.Albedo = existingOpts.Albedo
	}
	if opts.Bifaciality == "" {
		opts.Bifaciality = existingOpts.Bifaciality
	}
	err = r.DB.UpdateSolarArrayData(arrayId, id.(uint), inputs, opts)
	if err != nil {
		log.Println("Error updating solar array parameters: ", err)
		return err
	}

	return c.SendString("Array has been updated")
}

func (r *Repository) DisplayAvailableData(c *fiber.Ctx) error {
	currSession, err := r.Cfg.Session.Get(c)
	if err != nil {
		return err
	}
	sessionUser := currSession.Get("User").(fiber.Map)
	// release the currSession
	err = currSession.Save()
	if err != nil {
		return err
	}

	if sessionUser["ID"] == nil {
		return c.Status(fiber.StatusBadRequest).SendString("User is empty")
	}
	id := sessionUser["ID"]
	arrayIds, err := r.DB.FetchUserArrays(id.(uint))
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"ID":               id,
		"Available arrays": arrayIds,
		"csrfToken":        c.Locals("token"),
	})
}

func (r *Repository) RemoveSolarArray(c *fiber.Ctx) error {
	currSession, err := r.Cfg.Session.Get(c)
	if err != nil {
		log.Println("Unable to access session storage: ", err)
	}
	sessionUser := currSession.Get("User").(fiber.Map)
	id := sessionUser["ID"]
	sarrayId := c.Params("array_id")
	arrayId, err := strconv.Atoi(sarrayId)
	if err != nil {
		return err
	}
	err = r.DB.RemoveSolarArrayData(id.(uint), arrayId)
	if err != nil {
		return err
	}
	return c.SendString("Successfully removed solar array")
}

func makeAPIRequest(inputs models.RequiredInputs, opts models.OptionalInputs) (*models.PowerEstimate, error) {
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
	if inputs.Address != "" {
		queryParams.Add("address", inputs.Address)
	} else {
		queryParams.Add("lat", opts.Latitude)
		queryParams.Add("lon", opts.Longitude)
	}
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
