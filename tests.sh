#!/bin/bash

services=("admin" "auth" "consumer" "data")

# Loop through each service and run tests
for service in "${services[@]}"; do
  echo "Running unit tests for: $service service"
  cd "$service" && go test ./... -v
  cd ..
done
