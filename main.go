package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

const portNumber = ":8080"

func GetEnv(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Could not load environment")
	}
}

// type Inputs struct {
// 	Azimuth        string `json:"azimuth"`
// 	SystemCapacity string `json:"system_capacity"`
// 	Losses         string `json:"losses"`
// 	ArrayType      string `json:"array_type"`
// 	ModuleType     string `json:"module_type"`
// 	Tilt           string `json:"tilt"`
// 	Adress         string `json:"address"`
// }

// type Outputs struct {
// 	DcMonthly []float64 `json:"dc_monthly"`
// 	AcMonthly []float64 `json:"ac_monthly"`
// 	AcAnnual float64 `json:"ac_annual"`

// }

// type PowerEstimate struct {
// 	Inputs  Inputs  `json:"inputs"`
// 	Outputs Outputs `json:"outputs"`
// }

type PVWattsResponse struct {
	Inputs  map[string]interface{} `json:"inputs"`
	Outputs map[string]interface{} `json:"outputs"`
}

const (
	baseURL = "https://developer.nrel.gov/api/pvwatts/v8.json"
)

func makeAPIRequest(address string) (*PVWattsResponse, error) {
	LoadEnv()
	apiKey := GetEnv("API_KEY")
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

	var pvWattsResponse PVWattsResponse
	err = json.NewDecoder(resp.Body).Decode(&pvWattsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode API response: %v", err)
	}

	return &pvWattsResponse, nil
}

func main() {
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/render", func(c *fiber.Ctx) error {
		address := "boulder, co" // You can change this to any desired location
		pvWattsResponse, err := makeAPIRequest(address)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error fetching data from the API")
		}

		return c.JSON(pvWattsResponse)
	})

	app.Listen(":3000")
}

// LoadEnv()
// apiKey := GetEnv("API_KEY")
// log.Println(apiKey)
// baseURL := "https://developer.nrel.gov/api/pvwatts/v8.json"
// address := "boulder, co"

// // Encode the query parameters, including the properly encoded address
// queryParams := url.Values{}
// queryParams.Add("api_key", apiKey)
// queryParams.Add("azimuth", "180")
// queryParams.Add("system_capacity", "4")
// queryParams.Add("losses", "14")
// queryParams.Add("array_type", "1")
// queryParams.Add("module_type", "0")
// queryParams.Add("gcr", "0.4")
// queryParams.Add("dc_ac_ratio", "1.2")
// queryParams.Add("inv_eff", "96.0")
// queryParams.Add("radius", "0")
// queryParams.Add("dataset", "nsrdb")
// queryParams.Add("tilt", "10")
// queryParams.Add("address", address)
// queryParams.Add("soiling", "12|4|45|23|9|99|67|12.54|54|9|0|7.6")
// queryParams.Add("albedo", "0.3")
// queryParams.Add("bifaciality", "0.7")

// apiEndpoint := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

// resp, err := http.Get(apiEndpoint)
// if err != nil {
// 	fmt.Println("Failed to make the request:", err)
// 	return
// }
// defer resp.Body.Close()

// if resp.StatusCode != http.StatusOK {
// 	fmt.Println("API request failed with status code:", resp.StatusCode)
// 	return
// }

// var pvWattsResponse PVWattsResponse
// err = json.NewDecoder(resp.Body).Decode(&pvWattsResponse)
// if err != nil {
// 	fmt.Println("Failed to decode API response:", err)
// 	return
// }

// fmt.Println("Input Parameters:")
// for key, value := range pvWattsResponse.Inputs {
// 	fmt.Printf("%s: %v\n", key, value)
// }

// fmt.Println("\nOutput Parameters:")
// for key, value := range pvWattsResponse.Outputs {
// 	fmt.Printf("%s: %v\n", key, value)
// }
