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
			"uuid":        &graphql.Field{Type: graphql.String},
			"name":        &graphql.Field{Type: graphql.String},
			"departure":   &graphql.Field{Type: graphql.String},
			"destination": &graphql.Field{Type: graphql.String},
			"time":        &graphql.Field{Type: graphql.String},
			"capacity":    &graphql.Field{Type: graphql.Int},
			"passenger":   &graphql.Field{Type: graphql.Int},
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
			"Namelist": &graphql.Field{
				Type: graphql.String,
			},
			"Comment": &graphql.Field{
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

var mutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"putItem": &graphql.Field{
				Type: itemType,
				Args: graphql.FieldConfigArgument{
					"uuid":              &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"name":              &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"departure":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"destination":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"time":              &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"capacity":          &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"passenger":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"passengerNames":    &graphql.ArgumentConfig{Type: graphql.NewList(graphql.String)},
					"passengerComments": &graphql.ArgumentConfig{Type: graphql.NewList(graphql.String)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					uuid, _ := params.Args["uuid"].(string)
					name, _ := params.Args["name"].(string)
					departure, _ := params.Args["departure"].(string)
					destination, _ := params.Args["destination"].(string)
					time, _ := params.Args["time"].(string)
					capacity, _ := params.Args["capacity"].(int)
					passenger, _ := params.Args["passenger"].(int)

					passengerNames, ok := params.Args["passengerNames"].([]interface{})
					if !ok {
						return nil, errors.New("passengerNames must be an array of strings")
					}

					passengerComments, ok := params.Args["passengerComments"].([]interface{})
					if !ok {
						return nil, errors.New("passengerComments must be an array of strings")
					}

					var passengers []models.PassengerModel
					for i := range passengerNames {
						name, ok1 := passengerNames[i].(string)
						comment, ok2 := passengerComments[i].(string)
						if !ok1 || !ok2 {
							return nil, errors.New("Error converting to string")
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
					"uuid":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"name":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"namelist": &graphql.ArgumentConfig{Type: graphql.String},
					"comment":  &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					uuid, _ := params.Args["uuid"].(string)
					name, _ := params.Args["name"].(string)
					namelist, _ := params.Args["namelist"].(string)
					comment, _ := params.Args["comment"].(string)
					return updatePassenger(uuid, name, namelist, comment)
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
