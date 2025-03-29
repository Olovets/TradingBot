# Step 1: Use a Go base image to build the application
FROM golang:1.24 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the entire project into the container
COPY . .

# Install necessary dependencies and update glibc
RUN apt-get update && apt-get install -y \
    libc6-dev \
    wget \
    && wget http://ftp.gnu.org/gnu/libc/glibc-2.34.tar.gz \
    && tar -xvzf glibc-2.34.tar.gz \
    && cd glibc-2.34 \
    && mkdir build \
    && cd build \
    && ../configure --prefix=/usr \
    && make -j$(nproc) \
    && make install

# Step 2: Build the application
WORKDIR /app/cmd/rest-server
RUN go build -o /bin/rest-server .

# Step 3: Create a minimal image to run the app
FROM debian:bullseye-slim

WORKDIR /root/

# Copy the compiled binary from the builder image
COPY --from=builder /bin/rest-server .

# Expose the port the app will run on
EXPOSE 8080

# Run the Go app
CMD ["./rest-server"]
