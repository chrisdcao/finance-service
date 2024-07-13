#!/bin/bash

# Navigate to the root of the Go project (one level up from the bashscripts directory)
cd ../..

# Define the output directory for the Swagger docs
OUTPUT_DIR="./controllers/docs"

# Check if the output directory exists
if [ ! -d "$OUTPUT_DIR" ]; then
  echo "Creating output directory $OUTPUT_DIR"
  mkdir -p "$OUTPUT_DIR"
fi

# Generate the Swagger documentation
swag init --output "$OUTPUT_DIR"

# Check if the Swagger generation was successful
if [ $? -eq 0 ]; then
  echo "Swagger documentation generated successfully in $OUTPUT_DIR"
else
  echo "Failed to generate Swagger documentation"
  exit 1
fi
