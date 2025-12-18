# Multi-stage build for optimized production image
# Stage 1: Builder
FROM golang:1.25-alpine AS builder

# Build Version Argument
ARG BUILD_VERSION
ARG BUILD_FLAVOR

# Expose some build arguments as Env variables
ENV BUILD_VERSION=$BUILD_VERSION
ENV BUILD_FLAVOR=$BUILD_FLAVOR
RUN echo "Version is: $BUILD_FLAVOR-$BUILD_VERSION"

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations and version info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${BUILD_VERSION} -X main.Flavor=${BUILD_FLAVOR}" \
    -o /build/jokes-provider \
    ./main.go

# Stage 2: Runtime
FROM alpine:3.19

# Build Version Argument
ARG BUILD_VERSION
ARG BUILD_FLAVOR

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Create non-root user
RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=appuser:appuser /build/jokes-provider /app/jokes-provider

# Copy any config or data files if needed
COPY --chown=appuser:appuser docs/ /app/docs/
COPY --chown=appuser:appuser .env.example /app/

# Set environment variables for application
ENV BUILD_VERSION=${BUILD_VERSION}
ENV BUILD_FLAVOR=${BUILD_FLAVOR}

# Switch to non-root user
USER appuser

EXPOSE 3000

# Healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:3000/health || exit 1

# Labels for metadata
LABEL org.opencontainers.image.version=${BUILD_VERSION}
LABEL org.opencontainers.image.flavor=${BUILD_FLAVOR}
LABEL maintainer="your-email@example.com"
LABEL description="Production-ready Jokes Provider API with Fiber"

CMD ["/app/jokes-provider"]
