package graphql

import (
	"errors"

	"github.com/ahmos0/DyanamodbConnectMobile/pkg/models"
	"github.com/graphql-go/graphql"
)

var itemType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Item",
		Fields: graphql.Fields{
			"uuid": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"departure": &graphql.Field{
				Type: graphql.String,
			},
			"destination": &graphql.Field{
				Type: graphql.String,
			},
			"time": &graphql.Field{
				Type: graphql.String,
			},
			"capacity": &graphql.Field{
				Type: graphql.Int,
			},
			"passenger": &graphql.Field{
				Type: graphql.Int,
			},
			"passengers": &graphql.Field{
				Type: graphql.NewList(passengerType),
			},
		},
	},
)

var passengerType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Passenger",
		Fields: graphql.Fields{
			"namelist": &graphql.Field{
				Type: graphql.String,
			},
			"comment": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"allItems": &graphql.Field{
				Type: graphql.NewList(itemType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return getAllItems()
				},
			},
		},
	},
)

var passengerInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "PassengerInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"namelist": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"comment": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
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
					"passenger":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"passengers":  &graphql.ArgumentConfig{Type: graphql.NewList(passengerInputType)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					uuid, _ := params.Args["uuid"].(string)
					name, _ := params.Args["name"].(string)
					departure, _ := params.Args["departure"].(string)
					destination, _ := params.Args["destination"].(string)
					time, _ := params.Args["time"].(string)
					capacity, _ := params.Args["capacity"].(int)
					passenger, _ := params.Args["passenger"].(int)

					passengerInputs, ok := params.Args["passengers"].([]interface{})
					if !ok {
						return nil, errors.New("passengers must be an array of PassengerInput")
					}

					var passengers []models.PassengerModel
					for _, p := range passengerInputs {
						input, ok := p.(map[string]interface{})
						if !ok {
							return nil, errors.New("Error converting passenger input")
						}
						name, ok1 := input["namelist"].(string)
						comment, ok2 := input["comment"].(string)
						if !ok1 || !ok2 {
							return nil, errors.New("Error extracting name and comment from passenger input")
						}
						passenger := models.PassengerModel{
							Namelist: name,
							Comment:  comment,
						}
						passengers = append(passengers, passenger)
					}

					item, err := saveItem(uuid, name, departure, destination, time, capacity, passenger, passengers)
					if err != nil {
						return nil, err
					}

					return item, nil
				},
			},

			"incrementPassenger": &graphql.Field{
				Type: itemType,
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"passengers": &graphql.ArgumentConfig{
						Type: graphql.NewList(passengerInputType), // passengerInputTypeはPassengerInputのGraphQL型です
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					uuid, _ := params.Args["uuid"].(string)
					name, _ := params.Args["name"].(string)
					passengers := make([]models.PassengerModel, 0)
					if ps, ok := params.Args["passengers"].([]interface{}); ok {
						for _, p := range ps {
							if pmap, ok := p.(map[string]interface{}); ok {
								passengers = append(passengers, models.PassengerModel{
									Namelist: pmap["namelist"].(string),
									Comment:  pmap["comment"].(string),
								})
							}
						}
					}
					return updatePassenger(uuid, name, passengers)
				},
			},
		},
	},
)

var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)
