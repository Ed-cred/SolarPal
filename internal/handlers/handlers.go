package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Ed-cred/SolarPal/config"
	"github.com/Ed-cred/SolarPal/internal/models"
	"github.com/gofiber/fiber/v2"
)

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
func GetPowerEstimate (c *fiber.Ctx) error {
	address := "boulder, co" // You can change this to any desired location
		pvWattsResponse, err := MakeAPIRequest(address)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error fetching data from the API")
		}

		c.JSON(pvWattsResponse)
		return nil
}
