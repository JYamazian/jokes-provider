# Jokes Provider API

A high-performance REST API service for serving random jokes with built-in caching, rate limiting, and health checks. Built with Go using the Fiber web framework.

## 🚀 Features

- **Random Joke API** - Get random jokes with caching support
- **Health Checks** - Readiness probes for Kubernetes deployments
- **Rate Limiting** - Built-in request rate limiting with configurable thresholds
- **Swagger Documentation** - Auto-generated API documentation
- **Structured Logging** - JSON and text format logging with request tracking
- **Docker Support** - Multi-stage Docker build for optimized production images
- **Caching** - Optional Caching integration for improved performance
- **Configurable** - Extensive environment variable configuration

## 📋 Prerequisites

- Go 1.25.5 or higher
- Docker & Docker Compose (for containerized deployment)
- Redis (optional, for caching functionality)

## 🛠️ Dependencies

- **Fiber v2** - High-performance web framework for Go
- **Fiber Healthcheck** - Built-in liveness/readiness probe middleware
- **Fiber Contrib Swagger** - Swagger UI integration for Fiber
- **Swaggo/Swag** - Swagger documentation generator from annotations
- **Cache Storage** - Support for Fiber storage (Redis)

See [go.mod](go.mod) for complete dependency list.

## 📁 Project Structure

```bash
.
├── api/                 # API initialization and startup logic
├── config/              # Configuration and environment variable handling
├── controllers/         # HTTP request handlers with Swagger annotations
│   ├── health.go        # Health check endpoints (readiness, liveness)
│   ├── jokes.go         # Joke endpoints
│   └── metadata.go      # Metadata endpoints
├── docs/                # Auto-generated Swagger documentation
├── helpers/             # Utility functions (cache, jokes loading, random selection)
├── middleware/          # Request middleware (caching, rate limiting)
├── models/              # Data structures and types
│   ├── appConfig.go     # Application configuration model
│   ├── cacheConfig.go   # Cache configuration model
│   ├── fiberConfig.go   # Fiber framework configuration model
│   ├── joke.go          # Joke data model
│   ├── metadata.go      # Metadata response models
│   └── readinessHealthStatus.go  # Health status model
├── router/              # Route definitions and grouping
│   └── routers.go       # Route registration with versioned API groups
├── services/            # Business logic layer
│   ├── health.go        # Health check business logic
│   ├── jokes.go         # Joke retrieval and caching logic
│   ├── metadata.go      # Metadata assembly logic
│   ├── rateLimiter.go   # Rate limiting configuration
│   └── swagger.go       # Swagger UI setup
├── utils/               # Utility functions and knobs
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── Dockerfile           # Multi-stage Docker build configuration
├── compose.yml          # Docker Compose configuration for local development
└── README.md            # This file
```

## 🚀 Quick Start

### Local Development

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd jokes-provider
   ```

2. **Download dependencies**

   ```bash
   go mod download
   ```

3. **Run locally**

   ```bash
   go run main.go
   ```

   The API will start on `http://localhost:3000`

### Docker Compose

1. **Start all services**

   ```bash
   docker-compose -f compose.yml up -d --build --remove-orphans
   ```

   This will:
   - Download jokes data from CDN
   - Build the Jokes Provider application
   - Start the jokes-provider service on port 3000
   - Start redis-client for caching (if configured)

2. **Stop services**

   ```bash
   docker-compose -f compose.yml down
   ```

## 🔌 API Endpoints

### Jokes (API v1)

- `GET /v1/jokes/random` - Returns a random joke with optional caching

### Health Status

- `GET /health/liveness` - Checks if the service is UP (uses Fiber's built-in healthcheck)
- `GET /health/readiness` - Checks if the service is ready (validates Redis and data availability)

### Metadata (API v1)

- `GET /v1/metadata` - Returns comprehensive application metadata including version, configuration, and environment information

### Documentation

- `GET /swagger` - Interactive API documentation powered by Swagger UI

## ⚙️ Configuration

Configure the application using environment variables:

### Server Configuration

- `PORT` - Server port (default: `3000`)
- `ENVIRONMENT` - Environment type (default: `development`)

### Logging

- `LOG_LEVEL` - Log level: `debug`, `info`, `warn`, `error` (default: `info`)
- `LOG_FORMAT` - Custom log format pattern (default: `[${ip}]:${port} ${status} - ${method} ${path}`)
- `LOG_FORMAT_TYPE` - Format type: `json` or `text` (default: `text`)
- `LOG_DISABLE_COLORS` - Disable colored output (default: `false`)

### Build Information

- `BUILD_VERSION` - Application version (default: `dev`)
- `BUILD_FLAVOR` - Build flavor (default: `development`)

### Fiber Configuration

- `FIBER_PREFORK` - Enable prefork mode (default: `false`)
- `FIBER_CASE_SENSITIVE` - Case-sensitive routing (default: `false`)
- `FIBER_STRICT_ROUTING` - Strict routing (default: `false`)

### Files

- `JOKES_FILE_PATH` - Path to jokes CSV file (default: `/data/jokes.csv`)

### Request Headers

- `IP_HEADER_NAME` - Header name for client IP (default: `X-Forwarded-For`)
- `COUNTRY_HEADER_NAME` - Header name for country (default: `X-Country-Name`)

### Rate Limiting


- `RATE_LIMIT_ENABLED` - Control Rate Limit feature via a flag (default: `false`)
- `RATE_LIMIT_MAX_REQUESTS` - Maximum requests per time window (default: `100`)
- `RATE_LIMITER_EXPIRATION` - Rate limiter time window (default: `1m`)

### Caching

- `CACHE_URL` - Cache connection URL (default: `localhost`)
  - Supports standard Cache URL format: `redis://[:password@]host[:port]/[db]`
  - Supports TLS URLs: `rediss://[:password@]host[:port]/[db]`
- `CACHE_ENABLED` - Enable/disable caching (default: `true`)
- `CACHE_TTL` - Cache time-to-live duration (default: `5m`)
- `CACHE_CA_CERT` - Path to Cache CA certificate file (optional, for TLS)
- `CACHE_CLIENT_CERT` - Path to Cache client certificate file (optional, for mutual TLS)
- `CACHE_CLIENT_KEY` - Path to Cache client key file (optional, for mutual TLS)

## 📊 Jokes Data Format

The jokes are loaded from a CSV file with the following format:

```csv
id,joke
1,"Why did the programmer quit?"
2,"How many programmers does it take..."
```

Data is fetched from: <https://cdn.jsdelivr.net/gh/JYamazian/cdn-assets@main/assets/data/jokes.csv>

## 🔒 Features in Detail

### Caching Strategy

- Random jokes are cached in Cache client, such as Redis, when caching is enabled
- Singleton Cache connection initialized at application startup
- Connection reused across all cache operations for optimal performance
- Cache hits return previously fetched jokes
- Automatic cache invalidation based on TTL
- **TLS Support**: Secure Cache connections with optional client certificates
  - Set `CACHE_URL` to `rediss://` URL
  - Provide certificate paths via `CACHE_CA_CERT`, `CACHE_CLIENT_CERT`, `CACHE_CLIENT_KEY`

### Rate Limiting Strategy

- Global rate limiter middleware to prevent abuse
- Configurable per IP address
- Returns `429 Too Many Requests` when limit exceeded

### Structured Logging

- JSON formatted logs for better parsing
- Request tracking with unique request IDs
- Contextual information in all log entries

### Health Checks

- **Liveness probe** - Uses Fiber's built-in healthcheck middleware for optimal performance
- **Readiness probe** - Custom implementation that validates dependencies (Cache, data files)
- Follows Controller-Service pattern for consistent architecture
- Suitable for Kubernetes deployment health checks

## 🐳 Docker Deployment

### Build Docker Image

```bash
docker build -t jokes-provider:1.0.0 \
  --build-arg BUILD_VERSION=1.0.0 \
  --build-arg BUILD_FLAVOR=production .
```

### Run Docker Container

```bash
docker run -p 3000:3000 \
  -e PORT=3000 \
  -e ENVIRONMENT=production \
  -e LOG_LEVEL=info \
  jokes-provider:1.0.0
```

## 📈 Performance Considerations

- Multi-stage Docker build minimizes final image size
- Caching reduces joke service response time
- Rate limiting protects against abuse
- Structured logging enables efficient monitoring

## 🔧 Development

### Project Layout

- `main.go` - Application entry point
- `api/init.go` - Fiber app initialization and middleware setup
- `controllers/` - HTTP request handlers with Swagger annotations
- `services/` - Business logic layer (separated from HTTP concerns)
- `models/` - Data structures and type definitions
- `helpers/` - Utility functions for jokes, caching, random selection
- `middleware/` - Custom middleware implementations
- `config/` - Configuration loading and management
- `router/` - Route registration with API versioning
- `docs/` - Auto-generated Swagger documentation

### Architecture Pattern

The application follows a **Controller-Service-Model** pattern:

```
Router → Controller → Service → Helper/Middleware
                ↓
              Model
```

- **Controllers** - Handle HTTP requests, validate input, return responses, contain Swagger annotations
- **Services** - Contain business logic, separated from HTTP concerns
- **Models** - Define data structures used across the application
- **Helpers** - Provide utility functions (caching, data loading)
- **Middleware** - Handle cross-cutting concerns (rate limiting, caching)

### API Versioning

Routes are organized with API versioning:

```
/v1/jokes/random     → JokeController.GetRandomJoke
/v1/metadata         → MetadataController.GetMetadata
/health/readiness    → HealthController.Readiness
/health/liveness     → Fiber built-in healthcheck middleware
/swagger             → Swagger UI
```

### Building Locally

```bash
go build -o jokes-provider main.go
./jokes-provider
```

### Code Organization

- **Controllers Layer** - HTTP request handlers with Swagger annotations
- **Services Layer** - Business logic separated from HTTP concerns
- **Models Layer** - Data structures and type definitions
- **Helper Layer** - Reusable utility functions
- **Middleware Layer** - Request/response processing
- **Config Layer** - Environment and configuration management
- **Router Layer** - Route definitions with API versioning and grouping

## 👤 Author

Created by [JYamazian](https://github.com/JYamazian)
