#!/bin/bash

pkg=cars

if [ -z $pkg ];
then
    echo "Need package name as parameter"
    exit 0
fi

pkg_path=${pkg}/${pkg}pb/${pkg}.proto

# --go_out flag generates in ${pkg}.pb.go
# --go_grpc_out flag generates ${pkg}_grpc.pb.go

protoc --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative ${pkg_path}
echo "make sure two files in ${pkg}/${pkg}pb/ were generated: ${pkg}_grpc.pb.go and ${pkg}.pb.go"


# Command from Udemy tutorial not working with Windows protoc.exe
# protoc greet.proto --go_out=plugins=grpc:.

# new: plugin no loger does grpc code generation
# use --go-grpc_out and --go_out flags
# https://github.com/golang/protobuf/issues/1070
#
#"gRPC generation is no longer handled as a plugin 
#of this protoc-gen-go package. We only generate 
#the protobuf definitions, the protoc-gen-go-grpc 
#is now responsible for generating grpc code."
#

# --go_out flag generates in greet.pb.go
# --go_grpc_out flag generates greet_grpc.pb.go

# protoc --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative greet/greetpb/greet.proto
# echo "make sure two files in greet/greetpb/ were generated: greet_grpc.pb.go and greet.pb.go"


# Example that worked on command line executed in folder /greet/greetpb
# FROM:  https://grpc.io/docs/languages/go/quickstart/
# protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=paths=source_relative greet.proto
