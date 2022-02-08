package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/simrie/go-grpc-car-service/cars/carspb"
	"github.com/simrie/go-grpc-car-service/cars/data"

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

	//make, model := req.GetReq().Make, req.GetReq().Model

	rec, err := data.GetRecordById(3)
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf(`In stock: make: %s, model: %s`, rec.Make, rec.Model)
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

	//make, model := req.GetReq().Make, req.GetReq().Model
	recs, err := data.GetAllRecords()
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf(`cars: %v`, recs)
	res := &carspb.CarWithDeadlineResponse{
		Result: result,
	}
	return res, nil

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
