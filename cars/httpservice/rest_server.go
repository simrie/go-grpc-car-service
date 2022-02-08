package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	router.HandleFunc("/cars", MicroserviceHandlerSelector(client, "cars")).Methods("GET")
	router.HandleFunc("/car/{id}", MicroserviceHandlerSelector(client, "car/{id}")).Methods("GET")
	http.ListenAndServe(":8080", router)

}

func doUnary(c carspb.CarServiceClient) {
	fmt.Println("\n...Starting to do a Unary RPC...")
	req := &carspb.CarRequest{
		Id: int64(2),
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
		Id: int64(0),
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
	GetCarMicroserviceHandler sends a reqeust id to the gRPC service
	and returns a Car item with that id if found
*/
func GetCarMicroserviceHandler(c carspb.CarServiceClient, response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	vars := mux.Vars(request)
	idParam, ok := vars["id"]
	if !ok {
		response.WriteHeader(400)
		response.Write([]byte(`{ "message": "id parameter not defined" }`))
		return
	}
	//convert idParam from string to int64
	id, _ := strconv.Atoi(idParam)
	if !ok {
		response.WriteHeader(400)
		response.Write([]byte(`{ "message": "id parameter is not an integer" }`))
		return
	}

	carReq := carspb.CarRequest{
		Id: int64(id),
	}
	// context.Background() initializes a new, non-nil context
	// to be passed between server APIs
	res, err := c.Car(context.Background(), &carReq)
	if err != nil {
		log.Printf("\nerror while calling Car RPC: %v", err)
		status := http.StatusBadRequest
		response.WriteHeader(status)
		response.Write([]byte(`{ "message": "error retrieving Cars information"}`))
		return
	}

	json.NewEncoder(response).Encode(res)
}

func GetCarsMicroserviceHandler(c carspb.CarServiceClient, response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	carReq := &carspb.CarWithDeadlineRequest{
		Id: int64(0),
	}
	timeout := 4 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	res, err := c.CarWithDeadline(ctx, carReq)
	if err != nil {

		statusErr, ok := status.FromError(err)
		humanMsg := ""
		if ok {
			// this is a gRPC error
			if statusErr.Code() == codes.DeadlineExceeded {
				humanMsg = "Timeout was hit.  Deadline exceeded."
				log.Printf("%s: %v", humanMsg, statusErr)
			} else {
				humanMsg = "Unexpected status error"
				log.Printf("%s: %v", "Unexpected gRPC status error from Cars service", statusErr)
			}

		} else {
			// regular error
			humanMsg = "error retrieving Cars information"
			log.Printf("\nerror while calling Cars RPC: %v", err)
		}
		// return on any err so we do not try to print a non-existant res.Result
		response.WriteHeader(400)
		msg := fmt.Sprintf(`{ "message": %s }`, humanMsg)
		response.Write([]byte(msg))
		return
	}

	json.NewEncoder(response).Encode(res)
}

/*
HandlerPlaceholder is a placeholder
*/
func HandlerPlaceholder(response http.ResponseWriter, request *http.Request) {
	log.Printf("handler placeholder %s\n", request.RequestURI)
}

/*
MicroserviceHandlerSelector ties the api endpoint to a function
*/
func MicroserviceHandlerSelector(c carspb.CarServiceClient, endpoint string) http.HandlerFunc {
	var fn http.HandlerFunc
	switch endpoint {
	case "cars":
		fn = func(w http.ResponseWriter, r *http.Request) {
			GetCarsMicroserviceHandler(c, w, r)
		}
	case "car/{id}":
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
