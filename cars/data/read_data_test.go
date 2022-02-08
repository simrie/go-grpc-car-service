package data

import (
	"testing"

	"github.com/simrie/go-grpc-car-service/cars/models"
)

func TestGetAllRecords(t *testing.T) {
	var cars []models.Car
	var err error
	if cars, err = GetAllRecords(); &cars == nil || err != nil || len(cars) == 0 {
		t.Errorf("Failed! %v :", err)
	}
}

func TestGetRecordById(t *testing.T) {
	var car models.Car
	var err error
	if car, err = GetRecordById(2); &car == nil || err != nil || car.Id != 2 {
		t.Errorf("Failed! %v :", err)
	}
}
