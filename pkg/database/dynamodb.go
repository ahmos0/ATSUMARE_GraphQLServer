package database

import (
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
		return nil, fmt.Errorf("Error creating DynamoDB service: %v", err)
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String("DepatureManageTable"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("Failed to scan DynamoDB: %v", err)
	}

	fmt.Printf("Found %d items in the table.\n", len(result.Items))
	var items []models.Item
	for _, dbItem := range result.Items {
		fmt.Println("Processing an item...")
		item := models.Item{}

		// Handle string and numeric attributes
		if dbItem["uuid"] != nil {
			item.UUID = *dbItem["uuid"].S
		}
		if dbItem["name"] != nil {
			item.Name = *dbItem["name"].S
		}
		if dbItem["Departure"] != nil {
			item.Departure = *dbItem["Departure"].S
		}
		if dbItem["Destination"] != nil {
			item.Destination = *dbItem["Destination"].S
		}
		if dbItem["Time"] != nil {
			item.Time = *dbItem["Time"].S
		}
		if dbItem["Capacity"] != nil {
			item.Capacity, err = strconv.Atoi(*dbItem["Capacity"].N)
			if err != nil {
				return nil, fmt.Errorf("Error converting Capacity to int: %v", err)
			}
		}
		if dbItem["Passenger"] != nil {
			item.Passenger, err = strconv.Atoi(*dbItem["Passenger"].N)
			if err != nil {
				return nil, fmt.Errorf("Error converting Passenger to int: %v", err)
			}
		}

		// Handle the Passengers attribute
		if dbItem["Passengers"] != nil && dbItem["Passengers"].L != nil {
			for _, p := range dbItem["Passengers"].L {
				passenger := models.PassengerModel{}
				if cmt, ok := p.M["Comment"]; ok && cmt.S != nil {
					passenger.Comment = *cmt.S
				}
				if nmlst, ok := p.M["Namelist"]; ok && nmlst.S != nil {
					passenger.Namelist = *nmlst.S
				}
				item.Passengers = append(item.Passengers, passenger)
			}
		}

		items = append(items, item)
	}
	for i, item := range result.Items {
		fmt.Printf("Item %d: %v\n", i, item)
	}

	return items, nil
}

func UpdatePassenger(uuid string, name string, passengers []models.PassengerModel) (*models.Item, error) {
	svc, err := createDynamoDBService()
	if err != nil {
		return nil, fmt.Errorf("error creating DynamoDB service: %v", err)
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
		return nil, fmt.Errorf("error retrieving item: %v", err)
	}

	passengerCount, _ := strconv.Atoi(*getItemOutput.Item["Passenger"].N)
	passengerCount += len(passengers)

	dbPassengers := make([]*dynamodb.AttributeValue, len(passengers))
	for i, passenger := range passengers {
		dbPassengers[i] = &dynamodb.AttributeValue{M: map[string]*dynamodb.AttributeValue{
			"Namelist": {S: aws.String(passenger.Namelist)},
			"Comment":  {S: aws.String(passenger.Comment)},
		}}
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("DepatureManageTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {S: aws.String(uuid)},
			"name": {S: aws.String(name)},
		},
		UpdateExpression: aws.String("SET Passenger = :p, Passengers = list_append(Passengers, :new_passengers)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p":              {N: aws.String(fmt.Sprintf("%d", passengerCount))},
			":new_passengers": {L: dbPassengers},
		},
	}

	_, err = svc.UpdateItem(updateInput)
	if err != nil {
		return nil, fmt.Errorf("error updating item: %v", err)
	}

	updatedItem := &models.Item{
		UUID:       uuid,
		Name:       name,
		Passenger:  passengerCount,
		Passengers: passengers,
	}

	return updatedItem, nil
}

func SaveItem(uuid, name, departure, destination, time string, capacity int, passenger int, passengers []models.PassengerModel) (*models.Item, error) {
	svc, err := createDynamoDBService()
	if err != nil {
		log.Fatalf("Error creating DynamoDB service: %v", err)
	}

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
