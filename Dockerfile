# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git bash ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy dependency files first (better caching)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -a -installsuffix cgo \
    -o main ./cmd/main.go

# Stage 2: Create minimal production image
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    wget \
    && rm -rf /var/cache/apk/*

# Set timezone (optional, change as needed)
ENV TZ=Asia/Phnom_Penh

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy migration files if needed
COPY --from=builder /app/migrate ./migrate

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose API port
EXPOSE 8090

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8090/ping || exit 1

# Run the API
CMD ["./main"]