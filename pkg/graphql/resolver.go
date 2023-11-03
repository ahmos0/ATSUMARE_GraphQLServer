package graphql

import (
	"github.com/ahmos0/DyanamodbConnectMobile/pkg/database"
	"github.com/ahmos0/DyanamodbConnectMobile/pkg/models"
)

func getAllItems() ([]models.Item, error) {
	return database.GetAllItems()
}

func updatePassenger(uuid string, name string, passengers []models.PassengerModel) (*models.Item, error) {
	return database.UpdatePassenger(uuid, name, passengers)
}

func saveItem(uuid, name, departure, destination, time string, capacity int, passenger int, passengers []models.PassengerModel) (*models.Item, error) {
	return database.SaveItem(uuid, name, departure, destination, time, capacity, passenger, passengers)
}
