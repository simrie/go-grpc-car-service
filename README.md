[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/simrie/go-grpc-car-service)

# go-grpc-car-service
Demo showing how a gRPC server can be microservice that provides responses accessible through a different service that exposes RESTful endpoints.

In this scenario it would be the gRPC microservice that would connect to a database using search criteria received from the REST service.  The retrieved data would be returned to the REST service, where it is transformed into a response to the initial GET request.

Hard-coded data is used for this demo, substituting for database access.

The carspb folder contains the "cars.proto" protobuf request and response definitions. The "generate.sh" batch file was used to automatically generate the gRPC function files in the carspb directory. It is not necessary to run generate.sh again, unless a definition in the cars.proto file is changed.

## Test the Project's Functions

```
go test -v ./cars/data
```

## Build the Project

### Build the microservice gRPC server.

```
go build cars/microservice/grpc_server.go
```

### Build the http REST service.

```
go build cars/httpservice/rest_server.go
```

## Start the gRPC and API services

These are two separate processes.  In a real-life scenario these could run in two separate containers.

### Start the gRPC Microservice

```
./grpc_server
```
The service should continue running in the terminal and log output can be seen.


### Start the REST Service
```
./rest_server
```
The service should continue running in the terminal and log output can be seen.

## gitPod browser

If you have started this and followed the instructions by clicking the Gitpod browser link, you will see a dialog box asking if you want make port 8080 public or open a browser.  If you open a browser you can append the /cars or /cars/{id} endpoints to the browser address to see the retrieved content.

## cUrl commands

The following assume that the services are running on 127.0.0.1 (localhost) and the gRPC microservice is being contacted and providing responses when the REST API service receives a GET request.


### List all the items

```
curl "http://127.0.0.1:8080/cars"
```

The response is a list of car items or a user-friendly error message.

### Return one item by id

```
curl "http://127.0.0.1:8080/cars/{id}"
```

The response is a single car item or a user-friendly error message.

