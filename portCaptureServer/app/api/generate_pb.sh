#!/usr/bin/sh
protoc --go_out=./pb --go_opt=paths=source_relative  --go-grpc_opt=paths=source_relative --go-grpc_out=./pb ./*.proto
