#!/bin/bash

# build for lambda arch arm64 to compressed zip file 
go clean
GOOS=linux GOARCH=arm64 go build -o main
zip deployment.zip bootstrap main
