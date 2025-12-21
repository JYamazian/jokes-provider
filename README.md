# Jokes Provider API

[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Fiber](https://img.shields.io/badge/Fiber-v2-00ACD7?logo=go&logoColor=white)](https://gofiber.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue?logo=docker&logoColor=white)](https://www.docker.com/)
[![Redis](https://img.shields.io/badge/Redis-Caching-red?logo=redis&logoColor=white)](https://redis.io/)
[![Swagger](https://img.shields.io/badge/Swagger-OpenAPI-85EA2D?logo=swagger&logoColor=black)](https://swagger.io/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

A high-performance REST API service for serving random jokes, built with Go and the Fiber framework. The service features Redis-based caching, rate limiting, comprehensive health checks, and Swagger documentation.

## Features

- **High Performance**: Built on Fiber, one of the fastest Go web frameworks
- **Redis Caching**: Configurable cache-aside pattern with TTL support
- **Rate Limiting**: Per-client IP throttling with customizable limits
- **Health Checks**: Kubernetes-ready liveness and readiness probes
- **API Documentation**: Interactive Swagger UI with OpenAPI 3.0 spec
- **Structured Logging**: JSON or text format with request tracing
- **TLS Support**: Secure Redis connections with mTLS
- **Container Ready**: Multi-stage Docker build with non-root user
- **Graceful Shutdown**: Clean resource cleanup on termination

## Directory Structure

```
jokes-provider/
├── main.go                 # Application entry point
├── Dockerfile              # Multi-stage Docker build
├── compose.yml             # Docker Compose orchestration
├── go.mod                  # Go module dependencies
├── requests.http           # HTTP request examples
├── api/
│   └── init.go             # Application initialization and lifecycle
├── config/
│   ├── envVars.go          # Environment variable loading
│   ├── fileReader.go       # CSV file operations
│   └── logger.go           # Structured logging configuration
├── controllers/
│   ├── health.go           # Health check endpoints
│   ├── jokes.go            # Joke endpoints
│   └── metadata.go         # Application metadata endpoint
├── docs/
│   ├── docs.go             # Swagger documentation generator
│   ├── swagger.json        # OpenAPI specification (JSON)
│   └── swagger.yaml        # OpenAPI specification (YAML)
├── helpers/
│   ├── cacheStatus.go      # Redis health check utilities
│   ├── loadJokes.go        # CSV validation
│   └── randomJoke.go       # Joke retrieval logic
├── middleware/
│   └── cache.go            # Redis connection and operations
├── models/
│   ├── appConfig.go        # Application configuration model
│   ├── cacheConfig.go      # Cache configuration model
│   ├── fiberConfig.go      # Fiber configuration model
│   ├── joke.go             # Joke data model
│   ├── metadata.go         # Metadata response models
│   └── readinessHealthStatus.go  # Health status model
├── router/
│   └── routers.go          # Route definitions
├── services/
│   ├── health.go           # Health check business logic
│   ├── jokes.go            # Joke service with caching
│   ├── metadata.go         # Metadata service
│   ├── rateLimiter.go      # Rate limiting configuration
│   └── swagger.go          # Swagger UI setup
├── utils/
│   ├── constants.go        # Application constants
│   └── knobs.go            # Environment utilities
└── wrapper/
    ├── cacheHandler.go     # Cache read/write operations
    └── cachePolicy.go      # Cache-Control header handling
```

## Configuration

All configuration is managed through environment variables with sensible defaults.

### Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | HTTP server port |
| `ENVIRONMENT` | `development` | Environment name (development, staging, production) |
| `BUILD_VERSION` | `dev` | Application version (set at build time) |
| `BUILD_FLAVOR` | `development` | Build flavor identifier |

### Logging Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Log level (debug, info, error) |
| `LOG_FORMAT` | `[${ip}]:${port} ${status} - ${method} ${path}` | Log message format |
| `LOG_FORMAT_TYPE` | `text` | Output format (`text` or `json`) |
| `LOG_DISABLE_COLORS` | `false` | Disable colored output |

### Cache Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `CACHE_URL` | `localhost` | Redis connection URL (e.g., `redis://host:port/db`) |
| `CACHE_ENABLED` | `true` | Enable/disable caching |
| `CACHE_TTL` | `5m` | Cache time-to-live (supports Go duration format) |
| `CACHE_CA_CERT` | - | Path to CA certificate for Redis TLS |
| `CACHE_CLIENT_CERT` | - | Path to client certificate for Redis mTLS |
| `CACHE_CLIENT_KEY` | - | Path to client key for Redis mTLS |

### Rate Limiter Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `RATE_LIMIT_ENABLED` | `false` | Enable/disable rate limiting |
| `RATE_LIMIT_MAX_REQUESTS` | `100` | Maximum requests per window |
| `RATE_LIMITER_EXPIRATION` | `1m` | Rate limit window duration |

### Fiber Framework Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `FIBER_PREFORK` | `false` | Enable prefork mode for multi-process handling |
| `FIBER_CASE_SENSITIVE` | `false` | Case-sensitive routing |
| `FIBER_STRICT_ROUTING` | `false` | Strict routing (trailing slash matters) |

### Data Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `JOKES_FILE_PATH` | `/data/jokes.csv` | Path to jokes CSV file |
| `IP_HEADER_NAME` | `X-Forwarded-For` | Header for client IP (proxy support) |
| `COUNTRY_HEADER_NAME` | `X-Country-Name` | Header for country information |

## API Endpoints

### Jokes

#### Get Random Joke

```http
GET /v1/jokes/random
```

Returns a random joke from the database.

**Response:**

```json
{
  "ID": "42",
  "Joke": "Why do programmers prefer dark mode? Because light attracts bugs!"
}
```

#### Get Joke by ID

```http
GET /v1/jokes/{id}
```

Returns a specific joke by its ID.

**Parameters:**

- `id` (path, required): The joke identifier

**Response (200):**

```json
{
  "ID": "10",
  "Joke": "What's a computer's favorite snack? Microchips!"
}
```

**Response (404):**

```json
{
  "error": "Joke not found",
  "id": "999"
}
```

### Health Checks

#### Liveness Probe

```http
GET /health/liveness
```

Simple health check indicating the service is running.

**Response (200):**

```text
OK
```

#### Readiness Probe

```http
GET /health/readiness
```

Comprehensive health check verifying all dependencies.

**Response (200):**

```json
{
  "ready": true,
  "redis": "connected",
  "csv": "accessible"
}
```

**Response (503):**

```json
{
  "ready": false,
  "reason": "Redis unavailable"
}
```

### Metadata

#### Get Application Metadata

```http
GET /v1/metadata
```

Returns comprehensive application configuration and status information.

**Response:**

```json
{
  "app": {
    "name": "Jokes Provider API",
    "version": "1.0.0",
    "flavor": "dev"
  },
  "server": {
    "port": "3000",
    "environment": "development",
    "timestamp": "2025-12-21T10:30:00Z"
  },
  "logging": {
    "level": "info",
    "format": "[${ip}]:${port} ${status} - ${method} ${path}",
    "format_type": "json",
    "disable_colors": "false"
  },
  "cache": {
    "enabled": true,
    "url": "redis://redis:6379/1",
    "ttl": "5m"
  },
  "files": {
    "jokes_path": "/data/jokes.csv"
  },
  "headers": {
    "ip_header_name": "X-Forwarded-For",
    "country_header_name": "X-Country-Name"
  },
  "rate_limiter": {
    "enabled": false,
    "max_requests": 5,
    "duration": "1m"
  },
  "fiber": {
    "prefork": false,
    "case_sensitive": false,
    "strict_routing": false
  }
}
```

### Documentation

#### Swagger UI

```http
GET /swagger
```

Interactive API documentation interface.

#### OpenAPI Specification

```http
GET /docs/swagger.json
```

OpenAPI 3.0 specification in JSON format.

## Caching

The service implements a Redis-based caching layer with the following behavior:

- **Cache-aside pattern**: Attempts to read from cache first; on miss, reads from the data source and populates cache
- **Configurable TTL**: Cache entries expire after the configured `CACHE_TTL` duration
- **Cache bypass**: Clients can skip caching by sending the `Cache-Control: no-cache` header
- **Per-request control**: Cache reads and writes respect the Cache-Control header independently

### Cache Keys

| Endpoint | Cache Key Format |
|----------|-----------------|
| Random Joke | `random` |
| Joke by ID | `joke:{id}` |

### TLS Support

Redis connections support TLS/mTLS for secure communication:

- Set `CACHE_CA_CERT` for server certificate verification
- Set `CACHE_CLIENT_CERT` and `CACHE_CLIENT_KEY` for mutual TLS authentication

## Middleware

### Rate Limiting

When enabled, the rate limiter restricts requests per client IP:

- Uses the `X-Forwarded-For` header when present (configurable via `IP_HEADER_NAME`)
- Falls back to direct client IP
- Returns HTTP 429 (Too Many Requests) when limit exceeded

### Request ID

Every request is assigned a unique identifier via the `X-Request-ID` header, enabling request tracing across logs.

### Structured Logging

All requests are logged with contextual information:

- Timestamp
- Request ID
- Client IP
- Country (if provided via header)
- HTTP method and path
- Response status and latency

## Error Handling

The API uses standard HTTP status codes with JSON error responses:

| Status Code | Condition |
|-------------|-----------|
| 200 | Successful request |
| 400 | Missing required parameters |
| 404 | Resource not found |
| 429 | Rate limit exceeded |
| 500 | Internal server error |
| 503 | Service unavailable (dependency failure) |

**Error Response Format:**

```json
{
  "error": "Error message description",
  "id": "optional-identifier"
}
```

## Logging

The application supports two log formats:

### Text Format

```
[2025-12-21T10:30:00Z] [INFO] [req-123] [192.168.1.1] Message | version=dev-1.0.0 key=value
```

### JSON Format

```json
{
  "timestamp": "2025-12-21T10:30:00Z",
  "level": "INFO",
  "message": "Cache hit",
  "request_id": "req-123",
  "ip_address": "192.168.1.1",
  "version": "dev-1.0.0",
  "cache_key": "random"
}
```

## Build and Run

### Prerequisites

- Go 1.25.5 or later
- Redis 7.x
- Docker and Docker Compose (for containerized deployment)

### Local Development

```bash
# Install dependencies
go mod download

# Run the application
go run main.go

# Generate Swagger documentation
swag init
```

### Docker Build

```bash
# Build with version information
docker build \
  --build-arg BUILD_VERSION=1.0.0 \
  --build-arg BUILD_FLAVOR=production \
  -t jokes-provider:1.0.0 .
```

### Docker Compose

The included `compose.yml` orchestrates the complete stack:

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f jokes-provider

# Stop all services
docker compose down
```

The compose setup includes:

- **jokes-init**: Downloads the jokes CSV file on first run
- **jokes-provider**: The main API service
- **redis**: Redis cache server

### Health Check Verification

```bash
# Liveness probe
curl http://localhost:3000/health/liveness

# Readiness probe
curl http://localhost:3000/health/readiness
```

## Deployment

### Container Configuration

The Docker image runs as a non-root user (`appuser`) and exposes port 3000. Configure the following for production:

```yaml
environment:
  - PORT=3000
  - ENVIRONMENT=production
  - LOG_FORMAT_TYPE=json
  - CACHE_URL=redis://your-redis-host:6379/0
  - CACHE_ENABLED=true
  - CACHE_TTL=10m
  - RATE_LIMIT_ENABLED=true
  - RATE_LIMIT_MAX_REQUESTS=100
  - RATE_LIMITER_EXPIRATION=1m
```

### Health Check Configuration

Configure container orchestrators to use the health endpoints:

- **Liveness**: `GET /health/liveness` - Basic process health
- **Readiness**: `GET /health/readiness` - Full dependency verification

Recommended probe settings:

```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:3000/health/liveness"]
  interval: 30s
  timeout: 5s
  retries: 3
  start_period: 10s
```

## Operational Considerations

### Scaling

- **Prefork mode**: Enable `FIBER_PREFORK=true` for multi-process handling on multi-core systems
- **Horizontal scaling**: The service is stateless; scale by adding instances behind a load balancer
- **Redis clustering**: Use Redis Cluster or Redis Sentinel for cache high availability

### Resource Limits

- **Rate limiting**: Protect against abuse with `RATE_LIMIT_ENABLED=true`
- **Cache TTL**: Balance freshness vs. performance with appropriate `CACHE_TTL` values
- **File descriptor limits**: Ensure sufficient limits for high-concurrency deployments

### Monitoring

Key metrics to monitor:

- Request latency (via access logs)
- Cache hit/miss ratio (logged per request)
- Redis connection health (readiness probe)
- Rate limit rejections (HTTP 429 responses)

## Security Considerations

### Network Security

- Run behind a reverse proxy (nginx, Traefik) for TLS termination
- Use Redis TLS for encrypted cache communication
- Configure `IP_HEADER_NAME` correctly when behind proxies

### Container Security

- Runs as non-root user (UID 1000)
- Minimal Alpine-based image
- No shell access in production image
- Read-only file system compatible

### Input Validation

- Joke IDs are validated before database lookup
- Cache keys are sanitized
- Request headers are parsed safely

### Rate Limiting Protection

- Enable rate limiting in production to prevent abuse
- Configure appropriate limits based on expected traffic
- Uses client IP for rate limit tracking (respects proxy headers)

## Author

### Jean Yamazian

- Email: [jeanyamazian@outlook.com](mailto:jeanyamazian@outlook.com)
