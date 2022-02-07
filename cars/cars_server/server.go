package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/simrie/go-grpc-car-service/cars/carspb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// This code was created following along the Udemy grpc course

// To begin with create a server type to which we will add services
// however this may be replaced later in the course

// Error if dummy struct does not implement unimplementedGreetServiceServer
type server struct {
	carspb.UnimplementedCarServiceServer
}

func (*server) Car(ctx context.Context, req *carspb.CarRequest) (*carspb.CarResponse, error) {
	fmt.Printf("Car function was invoked with %v\n", req)
	//firstName := req.GetGreeting().GetFirstName()
	make, model := req.GetReq().Make, req.GetReq().Model
	result := fmt.Sprintf(`In stock: make: %s, model: %s`, make, model)
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
	//firstName := req.GetGreeting().GetFirstName()
	make, model := req.GetReq().Make, req.GetReq().Model
	result := fmt.Sprintf(`Out of stock: make: %s, model: %s`, make, model)
	res := &carspb.CarWithDeadlineResponse{
		Result: result,
	}
	return res, nil

}

func main() {
	fmt.Println("Ohayoo-san! Robo-greeta de gozaimasu.")

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
