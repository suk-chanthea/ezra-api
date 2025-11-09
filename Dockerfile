# ============================================
# Stage 1: Build Stage
# ============================================
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    bash \
    ca-certificates \
    tzdata \
    make

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build arguments for version info
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-w -s \
    -X 'main.Version=${VERSION}' \
    -X 'main.BuildTime=${BUILD_TIME}' \
    -X 'main.GitCommit=${GIT_COMMIT}'" \
    -a -installsuffix cgo \
    -o ezra-api \
    ./cmd/main.go

# Verify the binary was created
RUN ls -lh ezra-api

# ============================================
# Stage 2: Runtime Stage
# ============================================
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/cache/apk/*

# Set timezone
ENV TZ=Asia/Phnom_Penh
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=appuser:appuser /build/ezra-api .

# Copy migration files (if needed)
COPY --from=builder --chown=appuser:appuser /build/migrate ./migrate

# Copy config template (optional)
COPY --from=builder --chown=appuser:appuser /build/config ./config

# Create directories for runtime data
RUN mkdir -p /app/logs /app/tmp && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose API port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Set environment variables
ENV GIN_MODE=release \
    APP_ENV=production

# Run the binary
ENTRYPOINT ["./ezra-api"]

# Labels for metadata
LABEL maintainer="your-email@example.com" \
      version="${VERSION}" \
      description="Ezra API - Clean Architecture RESTful API" \
      org.opencontainers.image.source="https://github.com/your-org/ezra"