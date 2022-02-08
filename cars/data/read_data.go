package data

import (
	"encoding/json"

	"github.com/simrie/go-grpc-car-service/cars/models"
)

/*
	GetAllRecords returns all the Car records
	from the hard-coded data standing in for a database
*/
func GetAllRecords() ([]models.Car, error) {
	json_string := `[
    {
        "id": 1,
        "make": "Ford",
        "model": "F10"
    },
    {
        "id": 2,
        "make": "Toyota",
        "model": "Camry"
    },
    {
        "id": 3,
        "make": "Toyota",
        "model": "Rav4"
    },
    {
        "id": 4,
        "make": "Ford",
        "model": "Bronco"
    },
    {
        "id": 5,
        "make": "Toyota",
        "model": "Tundra"
    },
    {
    	"id": 6,
        "make": "Honda",
		"model": "Fit"
    }
	]`
	var cars []models.Car
	err := json.Unmarshal([]byte(json_string), &cars)
	if err != nil {
		return nil, err
	}
	return cars, nil
}

/*
	GetRecordById returns a matching record from among
	the hard-coded data substituting for a database
*/
func GetRecordById(searchId int64) (models.Car, error) {
	var car models.Car
	data, err := GetAllRecords()
	if err != nil {
		return car, err
	}
	for i := range data {
		if data[i].Id == searchId {
			car = data[i]
			return car, nil
		}
	}
	return car, nil
}
