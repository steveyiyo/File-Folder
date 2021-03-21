#!/bin/bash

echo "Starting for build!\n"

echo "Linux amd64\n"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/Linux_arm64
echo "macOS amd64\m"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/macOS_arm64
echo "Windows amd64\n"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/Windows_arm64.exe