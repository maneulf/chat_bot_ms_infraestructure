# Use Go 1.23 bookworm as base image
FROM golang:1.23-bookworm AS base

# Move to working directory /build
WORKDIR /build

# Copy the go.mod and go.sum files to the /build directory
COPY go.mod go.sum ./
#COPY go.docker.mod ./go.mod

# Install dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

RUN go mod tidy

# Build the application
RUN go build -o web_adapter_ms cmd/server/main.go

# Document the port that may need to be published
EXPOSE 8000

# Start the application
CMD ["/build/web_adapter_ms"]