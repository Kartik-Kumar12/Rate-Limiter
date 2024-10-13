# Server 

This **server** package implements a rate-limiting system using the **Token Bucket Algorithm**. It ensures that requests from clients are controlled based on their IP addresses, preventing excessive requests and ensuring fair usage of server resources. **Redis** is used as a backend to manage token storage for rate limiting across distributed environments.

## Table of Contents
- [Server Configuration](#server-configuration)
- [Rate Limiting Logic](#rate-limiting-logic)
- [Design Patterns Used](#design-patterns-used)
- [Usage](#usage)

---

## Server Configuration

The server reads rate-limiting configurations from `config/ip_rate_config.json` and interacts with Redis to manage tokens. The configuration file specifies the rate limits for each IP address:

```json
{
  "ip_rate_limits": {
    "10.0.0.1": [5, 1],  // [Bucket Capacity, Refill Rate (tokens/sec)]
    "10.0.0.2": [3, 1],
    "10.0.0.3": [10, 2]
  }
}
```

Each IP address has:
- **Bucket Capacity**: Maximum number of tokens available at any given time.
- **Refill Rate**: Number of tokens added per second.

The Redis Lua script used to manage the token bucket is stored in `store/redis/script.lua`.

---

## Rate Limiting Logic

The **Token Bucket Algorithm** is used to control request rates. The Redis-backed implementation allows each IP address to have a token bucket that refills at a set rate. Requests consume tokens, and if no tokens are available, the request is rejected.

### Key Components:
- **Bucket Capacity**: Defines the maximum number of tokens.
- **Refill Rate**: Tokens are refilled at a constant rate.
- **Request Processing**: Each request consumes one token. If the bucket is empty, the request is denied with a `429 Too Many Requests` error.

The token bucket's state (tokens and refill time) is stored in Redis, ensuring that multiple server instances can share the rate-limiting state.

---

## Design Patterns Used

### 1. **Method Chaining Pattern**

The **Token Bucket** implementation uses method chaining to configure properties like capacity, refill rate, and Redis store in a readable, modular way:

```go
bucket := NewTokenBucket().
    WithCapacity(bucketCapacity).
    WithRefillRate(refillRate).
    WithStore(redis.GetStore())
```

**Benefit**: Increases code readability and enables more concise object configuration.

### 2. **Functional Options Pattern**

The Redis store initialization uses the functional options pattern to allow flexible configurations at runtime, such as passing the Redis client and Lua script:

```go
err := redisStore.WithConfigs(
    redisStore.WithScript(string(script)),
    redisStore.WithClient(client),
)
```

**Benefit**: Provides flexible and modular initialization of store configurations without changing function signatures.

### 3. **Middleware Pattern**

Rate limiting is implemented as middleware, ensuring that each request passes through the rate limiter before reaching the actual API handler:

```go
http.Handle("/ping", ratelimiter.MiddleWare(httpHandler.HandlerPing))
```

**Benefit**: Cleanly decouples rate-limiting logic from core request handling, improving maintainability.

### 4. **Store Interface for Extensibility**

The server interacts with **Redis** through an abstraction (store interface) which can easily be extended to other storage systems, such as databases or in-memory stores, without changing the core logic.

```go
type Store interface {
    Eval(ctx context.Context, ipAddress string, capacity float64, refillRate int64) (*float64, error)
}
```

**Benefits of Using the Store Interface**:
- **Extensibility**: Any data store (like Redis, SQL, or in-memory) can be swapped in as long as it implements the store interface.
- **Testability**: Makes the system easier to mock during unit tests by injecting different store implementations.
- **Flexibility**: Allows you to scale the system with different types of backends, like Redis clusters or cloud-based storage solutions.


---
## Usage

### Running the Server

After setting up Redis and configuring IP rate limits, start the server:

```bash
cd server/
go run main.go
```

The server will listen for requests on port `8080`.

### Test the API

Test the `/ping` endpoint with curl to check the rate-limiting functionality:

```bash
curl "http://localhost:8080/ping?ip=10.0.0.1"
```

The response will either be a success message:

```json
{
  "Status": "Successful",
  "Body": "Hi! You've reached the API."
}
```

Or, if the rate limit is exceeded, a `429 Too Many Requests` error:

```json
{
  "Status": "Request Failed",
  "Body": "Too Many Requests, try again later."
}
```
