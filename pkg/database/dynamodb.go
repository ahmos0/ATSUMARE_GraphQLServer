package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ahmos0/DyanamodbConnectMobile/pkg/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func GetAllItems() ([]models.Item, error) {
	svc, err := createDynamoDBService()
	if err != nil {
		log.Fatalf("Error creating DynamoDB service: %v", err)
	}
	input := &dynamodb.ScanInput{
		TableName: aws.String("DepatureManageTable"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	var items []models.Item
	for _, dbItem := range result.Items {
		capacity, err := strconv.Atoi(*dbItem["Capacity"].N)
		if err != nil {
			return nil, fmt.Errorf("Error converting capacity to int: %v", err)
		}
		passenger, err := strconv.Atoi(*dbItem["Passenger"].N)
		if err != nil {
			return nil, fmt.Errorf("Error converting capacity to int: %v", err)
		}

		var passengers []models.PassengerModel
		err = json.Unmarshal([]byte(*dbItem["Passengers"].S), &passengers)
		if err != nil {
			return nil, fmt.Errorf("Error decoding passengers: %v", err)
		}

		item := models.Item{
			UUID:        *dbItem["uuid"].S,
			Name:        *dbItem["name"].S,
			Departure:   *dbItem["Departure"].S,
			Destination: *dbItem["Destination"].S,
			Time:        *dbItem["Time"].S,
			Capacity:    capacity,
			Passenger:   passenger,
			Passengers:  passengers,
		}

		items = append(items, item)
	}

	return items, nil
}

func UpdatePassenger(uuid string, name string, namelist string, comment string) (*models.Item, error) {
	svc, err := createDynamoDBService()
	if err != nil {
		log.Fatalf("Error creating DynamoDB service: %v", err)
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String("DepatureManageTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {S: aws.String(uuid)},
			"name": {S: aws.String(name)},
		},
	}
	getItemOutput, err := svc.GetItem(getItemInput)
	if err != nil {
		return nil, err
	}

	passenger, _ := strconv.Atoi(*getItemOutput.Item["Passenger"].N)
	passenger++

	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("DepatureManageTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {S: aws.String(uuid)},
			"name": {S: aws.String(name)},
		},
		UpdateExpression: aws.String("SET Passenger = :p, Namelist = :nl, Comment = :c"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p":  {N: aws.String(fmt.Sprintf("%d", passenger))},
			":nl": {S: aws.String(namelist)},
			":c":  {S: aws.String(comment)},
		},
	}

	_, err = svc.UpdateItem(updateInput)
	if err != nil {
		return nil, err
	}

	updatedItem := &models.Item{
		UUID:      uuid,
		Passenger: passenger,
		Passengers: []models.PassengerModel{
			{
				Namelist: namelist,
				Comment:  comment,
			},
		},
	}
	return updatedItem, nil
}

func SaveItem(uuid, name, departure, destination, time string, capacity int, passenger int, passengers []models.PassengerModel) (*models.Item, error) {
	svc, err := createDynamoDBService()
	if err != nil {
		log.Fatalf("Error creating DynamoDB service: %v", err)
	}

	//debug
	log.Printf("UUID: %s", uuid)
	log.Printf("Name: %s", name)
	log.Printf("Departure: %s", departure)
	log.Printf("Destination: %s", destination)
	log.Printf("Time: %s", time)
	log.Printf("Capacity: %d", capacity)
	log.Printf("Passenger: %d", passenger)
	log.Printf("Input Passengers: %v", passengers)

	var passengerAttributes []*dynamodb.AttributeValue
	for _, p := range passengers {
		passengerAttributes = append(passengerAttributes, &dynamodb.AttributeValue{
			M: convertPassengerToAttributeValue(p),
		})
	}
	log.Printf("Passenger Attributes: %v", passengerAttributes)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("DepatureManageTable"),
		Item: map[string]*dynamodb.AttributeValue{
			"uuid":        {S: aws.String(uuid)},
			"name":        {S: aws.String(name)},
			"Departure":   {S: aws.String(departure)},
			"Destination": {S: aws.String(destination)},
			"Time":        {S: aws.String(time)},
			"Capacity":    {N: aws.String(fmt.Sprintf("%d", capacity))},
			"Passenger":   {N: aws.String(fmt.Sprintf("%d", passenger))},
			"Passengers":  {L: passengerAttributes},
		},
	}
	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("Error putting item: %v", err)
		return nil, err
	}

	item := &models.Item{
		UUID:        uuid,
		Name:        name,
		Departure:   departure,
		Destination: destination,
		Time:        time,
		Capacity:    capacity,
		Passenger:   passenger,
		Passengers:  passengers,
	}
	return item, nil
}

func createDynamoDBService() (*dynamodb.DynamoDB, error) {
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
	return svc, nil
}

func convertPassengerToAttributeValue(passenger models.PassengerModel) map[string]*dynamodb.AttributeValue {
	result := map[string]*dynamodb.AttributeValue{
		"Namelist": {S: aws.String(passenger.Namelist)},
		"Comment":  {S: aws.String(passenger.Comment)},
	}
	log.Printf("Converted Passenger: %v", result)
	return result
}
