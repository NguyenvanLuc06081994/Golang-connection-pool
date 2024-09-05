# Stage 1: Build the Go application
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o myapp .

# Stage 2: Create a smaller image to run the app
FROM alpine:latest

# Set environment variables
ENV GIN_MODE=release

# Create a directory for the application
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/myapp .

# Expose the port the app will run on
EXPOSE 8080

# Command to run the binary
CMD ["./myapp"]