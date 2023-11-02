package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ahmos0/DyanamodbConnectMobile/pkg/graphql"
	"github.com/graphql-go/handler"
)

func main() {

	http.HandleFunc("/graphql", handlerFunc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, nil)
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {

	h := handler.New(&handler.Config{
		Schema:   &graphql.Schema,
		Pretty:   true,
		GraphiQL: true,
	})

	log.Printf("Received request: %+v", r)
	h.ServeHTTP(w, r)
	log.Printf("Sent response: %+v", w)
}
