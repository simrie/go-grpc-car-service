package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/simrie/go-grpc-car-service/cars/carspb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("REST service starting.")

	// Create a connection object to the microservice
	clientConnectionObject, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial error %v", err)
	}

	defer clientConnectionObject.Close()

	// This is the client object for making microservice gRPC calls
	client := carspb.NewCarServiceClient(clientConnectionObject)

	fmt.Printf("Created the gRPCclient %f", client)
	doUnary(client)

	doUnaryWithDeadline(client, 5*time.Second) // should complete
	//doUnaryWithDeadline(client, 1*time.Millisecond) // should timeout

	router := mux.NewRouter()
	router.HandleFunc("/car/microservice", MicroserviceHandlerSelector(client, "car/microservice")).Methods("GET")
	router.HandleFunc("/cars", MicroserviceHandlerSelector(client, "car/microservice")).Methods("GET")
	router.HandleFunc("/car/{id}", MicroserviceHandlerSelector(client, "car/microservice")).Methods("GET")
	http.ListenAndServe(":8080", router)

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

/*
	Call unary gRPC from microservice
*/
func GetCarMicroserviceHandler(c carspb.CarServiceClient, response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	fmt.Println("GetCarMicroserviceHandler")
	//var searchCar models.TradeIn
	var carReq carspb.CarRequest

	request.Body = http.MaxBytesReader(response, request.Body, 1048576)

	//decoder := json.NewDecoder(request.Body)

	//err := decoder.Decode(&searchCar)
	//if err != nil {
	//	log.Fatalf("\nerror decoding request body")
	//}

	fmt.Println("\n...creating gRPC req...")
	carReq = carspb.CarRequest{
		Req: &carspb.Car{
			Make:  "Trixie",
			Model: "Racer",
		},
	}
	// context.Background() initializes a new, non-nil context
	// to be passed between server APIs
	res, err := c.Car(context.Background(), &carReq)
	if err != nil {
		log.Fatalf("\nerror while calling Car RPC: %v", err)
	}
	log.Printf("\nResponse from Car: %v", res.Result)

	if err != nil {
		status := http.StatusBadRequest
		response.WriteHeader(status)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(res)
}

/*
HandlerPlaceholder is a placeholder
*/
func HandlerPlaceholder(response http.ResponseWriter, request *http.Request) {
	fmt.Printf("handler placeholder %s\n", request.RequestURI)
}

/*
MicroserviceHandlerSelector ties the api endpoint to a function
*/
func MicroserviceHandlerSelector(c carspb.CarServiceClient, endpoint string) http.HandlerFunc {
	var fn http.HandlerFunc
	fmt.Println("\nMicroserviceHandlerSelector ", endpoint)
	switch endpoint {
	case "car/microservice":
		fn = func(w http.ResponseWriter, r *http.Request) {
			GetCarMicroserviceHandler(c, w, r)
		}
	default:
		fn = func(w http.ResponseWriter, r *http.Request) {
			HandlerPlaceholder(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
