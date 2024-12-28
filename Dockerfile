# Use the official Golang image as the base image
FROM golang:1.21.0-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app for the API server
RUN go build -o api ./cmd/api

# Build the Go app for the Kafka consumer
RUN go build -o processor ./cmd/processor

# Use a minimal base image to run the application
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the built executables from the builder stage
COPY --from=builder /app/api .
COPY --from=builder /app/processor .

# Expose port 8080 for the API server
EXPOSE 8080

# Command to run both executables
CMD ["sh", "-c", "./api & ./processor"]