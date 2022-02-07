package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/simrie/go-grpc-car-service/cars/carspb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("We got cars")

	clientConnectionObject, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial error %v", err)
	}

	defer clientConnectionObject.Close()

	// The client generating line below worked to generate a client from the service
	// with greetpb.UnimplementedGreetServiceServer
	// before we added the (*server) Greet function definition to server.go

	// We create the client but we cannot do anything with it yet
	//client := greetpb.GreetServiceClient(clientConnectionObject)

	// Now that the service has a Greet service implemented we do:
	client := carspb.NewCarServiceClient(clientConnectionObject)

	fmt.Printf("Created client %f", client)
	doUnary(client)

	doUnaryWithDeadline(client, 5*time.Second) // should complete
	//doUnaryWithDeadline(client, 1*time.Millisecond) // should timeout
}

func doUnary(c carspb.CarServiceClient) {
	fmt.Println("\n...Starting to do a Unary RPC...")
	req := &carspb.CarRequest{
		Req: &carspb.Car{
			Make:  "Boopsie",
			Model: "McFeathers",
		},
	}
	// context.Background() initializes a new, non-nil context
	// to be passed between server APIs
	res, err := c.Car(context.Background(), req)
	if err != nil {
		log.Fatalf("\nerror while calling Car RPC: %v", err)
	}
	log.Printf("\nResponse from Car: %v", res.Result)
}

func doUnaryWithDeadline(c carspb.CarServiceClient, timeout time.Duration) {
	fmt.Println("\n...Starting to do a Unary With Deadline RPC...")
	req := &carspb.CarWithDeadlineRequest{
		Req: &carspb.Car{
			Make:  "Boopsie",
			Model: "McFeathers",
		},
	}
	// We initialize the context with the a timeout
	// to be passed between server APIs
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	res, err := c.CarWithDeadline(ctx, req)
	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			// this is a gRPC error
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit.  Deadline exceeded")
			} else {
				log.Fatalf("Unexpected gRPC status error: %v", statusErr)
			}

		} else {
			// regular error
			log.Fatalf("\nerror while calling Cars RPC: %v", err)
		}
		// return on any err so we do not try to print a non-existant res.Result
		return
	}
	log.Printf("\nResponse from Cars: %v", res.Result)
}
