# Start from the official Golang image
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY homemade-redis/go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY homemade-redis/*.go ./

# Build the Go app
RUN go build -o redis-server

# Use a minimal image for the final container
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder
COPY --from=builder /app/redis-server .

# Expose Redis default port
EXPOSE 6379

# Run the server
CMD ["./redis-server"]