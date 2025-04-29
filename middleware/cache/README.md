# Cache Middleware for Quick

This middleware provides a high-performance caching system for HTTP responses in the Quick web framework. It intercepts responses and stores them in memory, serving subsequent identical requests directly from the cache without executing handlers.

## Features

- In-memory caching with high-performance algorithm (based on gocache)
- Optional Redis storage support
- Configurable TTL (Time-To-Live) for cached items
- Custom key generation for fine-grained cache control
- Conditional cache invalidation
- HTTP method filtering
- Response headers caching
- Cache-Control header support
- Cache status headers (X-Cache-Status, X-Cache-Source)
- Maximum response size limit

## Installation

The cache middleware is included with the Quick framework. No additional installation is required.

## Usage

```go
package main

import (
    "time"
    
    "github.com/jeffotoni/quick"
    "github.com/jeffotoni/quick/middleware/cache"
)

func main() {
    app := quick.New()
    
    // Use cache middleware with default settings
    app.Use(cache.New())
    
    // Define routes
    app.Get("/time", func(c *quick.Ctx) error {
        return c.SendString("Current time: " + time.Now().Format(time.RFC3339))
    })
    
    app.Listen(":3000")
}
```

## Configuration

The cache middleware can be customized with various options:

```go
app.Use(cache.New(cache.Config{
    // Default duration after which cached items will expire
    Expiration: 5 * time.Minute,
    
    // Function that returns a custom TTL for each request
    ExpirationGenerator: func(c *quick.Ctx, cfg *cache.Config) time.Duration {
        if c.Path() == "/api/data" {
            return 2 * time.Minute
        }
        return cfg.Expiration
    },
    
    // Function that generates a unique cache key for each request
    KeyGenerator: func(c *quick.Ctx) string {
        return c.Path() + "?user=" + c.Query("user")
    },
    
    // Name of the header that indicates cache status
    CacheHeader: "X-Cache-Status",
    
    // Enable respecting Cache-Control headers from clients
    CacheControl: true,
    
    // Cache and restore response headers
    StoreResponseHeaders: true,
    
    // Maximum size in bytes for a response to be cached
    MaxBytes: 1024 * 512, // 512KB
    
    // List of HTTP methods to cache
    Methods: []string{quick.MethodGet, quick.MethodHead},
    
    // Function that determines whether to skip the cache
    CacheInvalidator: func(c *quick.Ctx) bool {
        return c.Query("clear") == "1"
    },
    
    // Function that determines whether to skip the middleware
    Next: func(c *quick.Ctx) bool {
        return c.Path() == "/no-cache"
    },
}))
```

## Redis Storage

The middleware supports Redis as an alternative storage backend.

### Design Pattern

The cache middleware uses a **dependency inversion** approach for Redis integration. Instead of directly depending on a specific Redis client library, it defines a `RedisClient` interface that any Redis implementation can adapt to. This design provides several benefits:

1. **Flexibility**: You can use any Redis client library (go-redis, redigo, etc.) by creating an adapter that implements the interface.

2. **Decoupling**: The middleware doesn't need to know the specific details of how Redis is accessed.

3. **Testability**: Makes it easy to create mocks for unit testing, as demonstrated in the `redis.mock` example.

4. **Adaptability**: If the Redis API changes or if you want to switch to another storage system, you only need to modify the adapter.

5. **Dependency Control**: Avoids direct dependencies on external libraries, preventing version conflicts.

### Implementation Example

```go
import (
    "context"
    "time"
    
    "github.com/go-redis/redis/v8"
    "github.com/jeffotoni/quick"
    "github.com/jeffotoni/quick/middleware/cache"
)

// Create a Redis client adapter that implements the RedisClient interface
type redisAdapter struct {
    client *redis.Client
}

func (r *redisAdapter) Get(ctx context.Context, key string) (string, error) {
    return r.client.Get(ctx, key).Result()
}

func (r *redisAdapter) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
    return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisAdapter) Del(ctx context.Context, keys ...string) (int64, error) {
    return r.client.Del(ctx, keys...).Result()
}

func main() {
    app := quick.New()
    
    // Create Redis client
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    // Create Redis adapter
    adapter := &redisAdapter{client: redisClient}
    
    // Create Redis storage
    redisStorage, err := cache.NewRedisStorage(cache.RedisConfig{
        Client:    adapter,
        TTL:       5 * time.Minute,
        KeyPrefix: "quick:cache",
    })
    if err != nil {
        panic(err)
    }
    
    // Use cache middleware with Redis storage
    app.Use(cache.New(cache.Config{
        Expiration:          5 * time.Minute,
        CacheHeader:         "X-Cache-Status",
        StoreResponseHeaders: true,
        Storage:             redisStorage,
    }))
    
    // Define routes
    app.Get("/time", func(c *quick.Ctx) error {
        return c.SendString("Current time: " + time.Now().Format(time.RFC3339))
    })
    
    app.Listen(":3000")
}
```

## Cache Status Headers

The middleware adds the following headers to responses:

- `X-Cache-Status`: Indicates the cache status of the response
  - `HIT`: Response was served from cache
  - `MISS`: Response was not found in cache
  - `EXPIRED`: Response was found in cache but had expired
  - `INVALIDATED`: Cache was invalidated by the CacheInvalidator
  - `BYPASS`: Cache was bypassed due to Cache-Control header
  
- `X-Cache-Source`: Indicates the source of the cached response
  - `memory`: Response was served from in-memory cache
  - `redis`: Response was served from Redis cache

## Examples

### Basic Usage

```go
app := quick.New()
app.Use(cache.New())

app.Get("/time", func(c *quick.Ctx) error {
    return c.SendString("Current time: " + time.Now().Format(time.RFC3339))
})
```

### Custom Key Generator

```go
app := quick.New()
app.Use(cache.New(cache.Config{
    KeyGenerator: func(c *quick.Ctx) string {
        return c.Path() + "?lang=" + c.Query("lang")
    },
}))

app.Get("/greeting", func(c *quick.Ctx) error {
    lang := c.Query("lang")
    greeting := "Hello"
    if lang == "es" {
        greeting = "Hola"
    } else if lang == "fr" {
        greeting = "Bonjour"
    }
    return c.SendString(greeting + ", World!")
})
```

### Cache Invalidation

```go
app := quick.New()
app.Use(cache.New(cache.Config{
    CacheInvalidator: func(c *quick.Ctx) bool {
        return c.Query("clear") == "1"
    },
}))

app.Get("/data", func(c *quick.Ctx) error {
    return c.JSON(map[string]interface{}{
        "timestamp": time.Now().Unix(),
        "data":      "Some data",
    })
})
```

### Dynamic Expiration

```go
app := quick.New()
app.Use(cache.New(cache.Config{
    ExpirationGenerator: func(c *quick.Ctx, cfg *cache.Config) time.Duration {
        if c.Path() == "/api/data" {
            return 30 * time.Second
        }
        return 5 * time.Minute
    },
}))
```

## License

This middleware is part of the Quick framework and is licensed under the same terms.
