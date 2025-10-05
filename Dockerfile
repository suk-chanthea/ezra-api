# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder

# Install dependencies
RUN apk add --no-cache git bash

# Set workdir
WORKDIR /app

# Copy Go modules files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o main ./cmd/main.go

# Stage 2: Run the binary
FROM alpine:3.18

# Install certificates for HTTPS if needed
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy .env file (optional)
COPY .env .

# Expose API port
EXPOSE 8090

# Run the API
CMD ["./main"]
