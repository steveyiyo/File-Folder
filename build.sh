#!/bin/bash

echo "Starting for build!\n"

echo "Linux amd64\n"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/Linux_amd64
echo "macOS amd64\m"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/macOS_amd64
echo "Windows amd64\n"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/Windows_amd64.exe