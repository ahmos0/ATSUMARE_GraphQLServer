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
		fmt.Printf("Error reading env file: %v", err)
	}
	url := os.Getenv("ENDPOINT")

	query := `
	{
		allItems {
		  uuid
		  name
		  departure
		  destination
		  time
		  capacity
		}
	  }
	`

	requestBody, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		log.Fatalf("Error creating request body: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	fmt.Printf("Response: %+v\n", result)
}
