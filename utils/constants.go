package utils

// HTTP Headers
const (
	HeaderCacheControl = "Cache-Control"
	HeaderContentType  = "Content-Type"
)

// Cache Control Values
const (
	CacheControlNoCache = "no-cache"
)

// Cache Key Prefixes
const (
	CacheKeyPrefixJoke = "joke:"
)

// CSV Column Names
const (
	CSVColumnID = "ID"
)

// API Route Prefixes
const (
	APIVersionV1       = "/v1"
	RouteJokes         = "/jokes"
	RandomJokeEndpoint = "/random"
	RouteHealth        = "/health"
	JokeByIDEndpoint   = "/:id"
	MetadataEndpoint   = "/metadata"
	LivenessEndpoint   = "/liveness"
	ReadinessEndpoint  = "/readiness"
)

// Route Parameters
const (
	ParamID = "id"
)

// Error Messages
const (
	ErrMsgJokeNotFound     = "Joke not found"
	ErrMsgJokeIDRequired   = "Joke ID is required"
	ErrMsgFailedToRetrieve = "Failed to retrieve joke"
	ErrMsgIDColumnNotFound = "id column not found"
	ErrMsgNoJokesAvailable = "No jokes available in CSV file"
)

// JSON Response Keys
const (
	JSONKeyError = "error"
	JSONKeyID    = "id"
)
