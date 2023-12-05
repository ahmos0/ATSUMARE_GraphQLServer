package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error reading env file: %v", err)
	}
	url := os.Getenv("ENDPOINT")

	query := `
	query GetAllItems {
		allItems {
			uuid
			name
			departure
			destination
			time
			capacity
			passenger
			passengers {
				namelist
				comment
			}
		}
	}`

	requestBody, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		log.Fatalf("Error creating request body: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			AllItems []struct {
				UUID        string `json:"uuid"`
				Name        string `json:"name"`
				Departure   string `json:"departure"`
				Destination string `json:"destination"`
				Time        string `json:"time"`
				Capacity    int    `json:"capacity"`
				Passenger   int    `json:"passenger"`
				Passengers  []struct {
					Namelist string `json:"namelist"`
					Comment  string `json:"comment"`
				} `json:"passengers"`
			} `json:"allItems"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	fmt.Printf("Response: %+v\n", result.Data.AllItems)
}
