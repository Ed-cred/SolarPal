package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Ed-cred/SolarPal/config"
	"github.com/Ed-cred/SolarPal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//	type PVWattsResponse struct {
//		Inputs  map[string]interface{} `json:"inputs"`
//		Outputs map[string]interface{} `json:"outputs"`
//	}

func main() {
	app := fiber.New()
	db, err := run()
	if err != nil {
		panic(err)
	}
	log.Println("Connected to database")
	defer db.Close()
	app.Use(logger.New())
	// app.Use(basicauth.New(basicauth.Config{
	// 	Users: map[string]string{
	// 		"john":  "doe",
	// 		"admin": "123456",
	// 	},
	// 	Realm: "Forbidden",
	// 	Authorizer: func(user, pass string) bool {
	// 		if user == "john" && pass == "doe" {
	// 			return true
	// 		}
	// 		if user == "admin" && pass == "123456" {
	// 			return true
	// 		}
	// 		return false
	// 	},
	// 	Unauthorized: func(c *fiber.Ctx) error {
	// 		return c.SendString("Please Log In to continue")
	// 	},
	// 	ContextUsername: "_user",
	// 	ContextPassword: "_pass",
	// }))
	setupRoutes(app)
	fmt.Printf("Server started and listening at localhost:3000 - csrfActive: %v\n", len(os.Args) > 1 && os.Args[1] == "withoutCsrf")
	// Start server
	log.Fatal(app.Listen(":3000"))
}

func run() (*sql.DB, error) {
	config.LoadEnv()
	dbPath := config.GetEnv("SQLITE_PATH")
	log.Println("Connecting to database...")
	db, err := database.ConnectDb(dbPath)
	if err != nil {
		log.Fatal("Couldn't connect to database:", err)
		return nil, err
	}
	return db, nil
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
