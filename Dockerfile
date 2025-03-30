# Step 1: Use a Go base image to build the application
FROM golang:1.23-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies (this will cache dependencies if the go.mod and go.sum files are not changed)
RUN go mod tidy

# Copy the entire project into the container
COPY . .

# Build the Go app (Assuming main.go is located in cmd/rest-server/)
WORKDIR /app/cmd/rest-server

# Build the binary executable inside /bin directory
RUN go build -o /bin/rest-server .

# Step 2: Create a minimal image to run the app
FROM debian:bullseye-slim

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder image
COPY --from=builder /bin/rest-server .

# Expose the port the app will run on
EXPOSE 8071

# Run the Go app
CMD ["./rest-server"]
