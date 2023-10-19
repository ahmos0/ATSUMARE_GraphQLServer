package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Item struct {
	UUID        string
	Name        string
	Departure   string
	Destination string
	Time        string
	Capacity    int
}

var itemType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Item",
		Fields: graphql.Fields{
			"uuid":        &graphql.Field{Type: graphql.String},
			"name":        &graphql.Field{Type: graphql.String},
			"departure":   &graphql.Field{Type: graphql.String},
			"destination": &graphql.Field{Type: graphql.String},
			"time":        &graphql.Field{Type: graphql.String},
			"capacity":    &graphql.Field{Type: graphql.Int},
		},
	},
)

// 仮実装
var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"item": &graphql.Field{
				Type: itemType,
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					uuid, isOK := params.Args["uuid"].(string)
					if !isOK {
						return nil, fmt.Errorf("uuid is required")
					}
					return map[string]interface{}{
						"uuid":        uuid,
						"name":        "Sample Name",
						"departure":   "Sample Departure",
						"destination": "Sample Destination",
						"time":        "Sample Time",
						"capacity":    5,
					}, nil
				},
			},
		},
	},
)

var mutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"putItem": &graphql.Field{
				Type: itemType,
				Args: graphql.FieldConfigArgument{
					"uuid":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"name":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"departure":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"destination": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"time":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"capacity":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					uuid, _ := params.Args["uuid"].(string)
					name, _ := params.Args["name"].(string)
					departure, _ := params.Args["departure"].(string)
					destination, _ := params.Args["destination"].(string)
					time, _ := params.Args["time"].(string)
					capacity, _ := params.Args["capacity"].(int)

					// クライアントから提供されたデータをDynamoDBに保存
					item, err := saveItem(uuid, name, departure, destination, time, capacity)
					if err != nil {
						return nil, err
					}

					return item, nil
				},
			},
		},
	},
)

func saveItem(uuid, name, departure, destination, time string, capacity int) (*Item, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		log.Printf("Error creating AWS session: %v", err)
		return nil, err
	}
	svc := dynamodb.New(sess)
	input := &dynamodb.PutItemInput{
		TableName: aws.String("DepatureManageTable"),
		Item: map[string]*dynamodb.AttributeValue{
			"uuid":        {S: aws.String(uuid)},
			"name":        {S: aws.String(name)},
			"Departure":   {S: aws.String(departure)},
			"Destination": {S: aws.String(destination)},
			"Time":        {S: aws.String(time)},
			"Capacity":    {N: aws.String(fmt.Sprintf("%d", capacity))},
		},
	}
	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("Error putting item: %v", err)
		return nil, err
	}

	item := &Item{
		UUID:        uuid,
		Name:        name,
		Departure:   departure,
		Destination: destination,
		Time:        time,
		Capacity:    capacity,
	}
	return item, nil
}

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
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
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	log.Printf("Received request: %+v", r)
	h.ServeHTTP(w, r)
	log.Printf("Sent response: %+v", w)
}
