package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	url := "https://graphqlserver-k5fklwkbdq-an.a.run.app/graphql"
	//url := "http://localhost:8080/graphql"
	/*mutation := `
		mutation {
			putItem(
				uuid: "hogehoge",
				name: "New Item",
				departure: "Tokyo",
				destination: "Osaka",
				time: "10:00",
				capacity: 100
			) {
				uuid
				name
			}
		}
	`

	requestBody, err := json.Marshal(map[string]string{"query": mutation})
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

	fmt.Printf("Response: %+v\n", result)*/

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
