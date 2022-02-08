package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/simrie/go-grpc-car-service/cars/carspb"
	"github.com/simrie/go-grpc-car-service/cars/data"
	"github.com/simrie/go-grpc-car-service/cars/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error if dummy struct does not implement unimplementedGreetServiceServer
type server struct {
	carspb.UnimplementedCarServiceServer
}

func (*server) Car(ctx context.Context, req *carspb.CarRequest) (*carspb.CarResponse, error) {
	fmt.Printf("Car function was invoked with %v\n", req)

	id := req.Id

	rec, err := data.GetRecordById(id)
	if err != nil {
		return nil, err
	}

	// convert result to *carspb.Car
	var result *carspb.Car
	result, err = ConvertCarToCarpb(rec)
	if err != nil {
		return nil, err
	}

	res := &carspb.CarResponse{
		Result: result,
	}
	return res, nil
}

func (*server) CarWithDeadline(ctx context.Context, req *carspb.CarWithDeadlineRequest) (*carspb.CarWithDeadlineResponse, error) {
	fmt.Printf("CarsWithDeadline function was invoked with %v\n", req)

	// Check to see if the timeout occurred
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			msg := "Client cancelled the request (timeout)"
			fmt.Println(msg)
			return nil, status.Error(codes.DeadlineExceeded, msg)
		}
		time.Sleep(1 * time.Second)
	}

	recs, err := data.GetAllRecords()
	if err != nil {
		return nil, err
	}

	// convert result to *carspb.Car
	var results []*carspb.Car
	for _, v := range recs {
		var result *carspb.Car
		result, err = ConvertCarToCarpb(v)
		if err == nil {
			results = append(results, result)
		}
	}

	res := &carspb.CarWithDeadlineResponse{
		Result: results,
	}
	return res, nil
}

func ConvertCarToCarpb(car models.Car) (*carspb.Car, error) {
	var carpb carspb.Car
	carpb.Id = car.Id
	carpb.Make = car.Make
	carpb.Model = car.Model
	return &carpb, nil
}

func main() {
	fmt.Println("Microservice starting.")

	// Here we test the grpc code generated from cars.proto

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("cannot listen to grpc port for tcp: %v", err)
	}

	s := grpc.NewServer()

	carspb.RegisterCarServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}

}
