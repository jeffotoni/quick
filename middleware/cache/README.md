# 🚀 Cache Middleware for Quick

[![GoDoc](https://godoc.org/github.com/jeffotoni/quick/middleware/cache?status.svg)](https://godoc.org/github.com/jeffotoni/quick/middleware/cache) 
[![Go Report](https://goreportcard.com/badge/github.com/jeffotoni/quick)](https://goreportcard.com/badge/github.com/jeffotoni/quick)
[![License](https://img.shields.io/github/license/jeffotoni/quick)](https://img.shields.io/github/license/jeffotoni/quick)

This middleware provides a high-performance caching system for HTTP responses in the Quick web framework. It intercepts responses and stores them in memory, serving subsequent identical requests directly from the cache without executing handlers, significantly improving response times and reducing server load.

```
   ██████╗ █████╗  ██████╗██╗  ██╗███████╗
  ██╔════╝██╔══██╗██╔════╝██║  ██║██╔════╝
  ██║     ███████║██║     ███████║█████╗  
  ██║     ██╔══██║██║     ██╔══██║██╔══╝  
  ╚██████╗██║  ██║╚██████╗██║  ██║███████╗
   ╚═════╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚══════╝

```

## 📑 Table of Contents

- [Features](#-features)
- [Installation](#-installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Redis Storage](#redis-storage)
- [Cache Status Headers](#cache-status-headers)
- [Examples](#examples)
  - [Basic Usage](#basic-usage)
  - [Custom Key Generator](#custom-key-generator)
  - [Cache Invalidation](#cache-invalidation)
  - [Dynamic Expiration](#dynamic-expiration)
- [Performance](#-performance)
- [Troubleshooting](#-troubleshooting)
- [License](#-license)
- [Contributing](#-contributing)

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🧠 **In-memory caching** | High-performance algorithm (based on gocache) |
| 🔄 **Redis support** | Optional Redis storage backend |
| ⏱️ **Configurable TTL** | Time-To-Live for cached items |
| 🔑 **Custom key generation** | Fine-grained cache control |
| 🚫 **Conditional invalidation** | Invalidate cache based on custom conditions |
| 🔍 **HTTP method filtering** | Cache only specific HTTP methods |
| 📋 **Headers caching** | Store and restore response headers |
| 🔒 **Cache-Control support** | Respect standard HTTP cache headers |
| 📊 **Status headers** | X-Cache-Status and X-Cache-Source headers |
| 📏 **Size limiting** | Maximum response size limit |

## 📦 Installation

The cache middleware is included with the Quick framework. No additional installation is required.

```bash
# If you're using Quick, you already have access to the cache middleware
go get -u github.com/jeffotoni/quick
```

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
cache.New(cache.Config{
    // ...
    OnHit: func(key string) {
        fmt.Println("Cache hit:", key)
    },
    OnMiss: func(key string) {
        fmt.Println("Cache miss:", key)
    },
})
```

```go
cache.New(cache.Config{
    OnCacheHit: func(c *quick.Ctx, key string) {
        log.Printf("[HIT] %s from %s", key, c.IP())
    },
    OnCacheSet: func(c *quick.Ctx, key string) {
        log.Printf("[SET] %s status %d", key, c.Response.StatusCode())
    },
})
```

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

    // OnHit is triggered when a cache hit occurs 
    OnHit: func(key string) {
        fmt.Println("Cache hit:", key)
    },

    // OnMiss is triggered when a cache miss occurs
    OnMiss: func(key string) {
        fmt.Println("Cache miss:", key)
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
package main

import (
    "time"

    "github.com/jeffotoni/quick"
    "github.com/jeffotoni/quick/middleware/cache"
)

func main() {

    app := quick.New()
    app.Use(cache.New())

    // Rota 1
    app.Get("/hello", func(c *quick.Ctx) error {
        return c.SendString("Hello World!")
    })

    // Rota 2
    app.Get("/now", func(c *quick.Ctx) error {
        return c.SendString("Now: " + time.Now().Format(time.RFC3339Nano))
    })

    // Rota 3
    app.Get("/random", func(c *quick.Ctx) error {
        return c.SendString("Random: " + fmt.Sprint(rand.Intn(10000)))
    })

    app.Listen(":3000")
}
```

#### Testando o Cache na Prática

Acesse as rotas múltiplas vezes e observe o cabeçalho `X-Cache-Status`:

```bash
# Primeira chamada (MISS)
curl -i http://localhost:3000/hello
# ... X-Cache-Status: MISS

# Segunda chamada (HIT)
curl -i http://localhost:3000/hello
# ... X-Cache-Status: HIT

# Outro path, sempre começa com MISS
curl -i http://localhost:3000/now
# ... X-Cache-Status: MISS
curl -i http://localhost:3000/now
# ... X-Cache-Status: HIT

# Each path has its own cache!
curl -i http://localhost:3000/random
# ... X-Cache-Status: MISS
curl -i http://localhost:3000/random
# ... X-Cache-Status: HIT
```

> Tip: To see the cache expire, simply wait for the configured TTL or restart the server.

These examples clearly demonstrate the cache behavior for different routes, illustrating HIT/MISS in a simple and functional way.

### Custom Key Generator

```go
package main

import (
    "time"

    "github.com/jeffotoni/quick"
    "github.com/jeffotoni/quick/middleware/cache"
)

func main() {

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

    app.Listen(":3000")
}
```

### Cache Invalidation

```go
package main

import (
    "time"

    "github.com/jeffotoni/quick"
    "github.com/jeffotoni/quick/middleware/cache"
)

func main() {

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

    app.Listen(":3000")
}
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

## 📊 Performance

The cache middleware is designed for high performance with minimal overhead:

- **Memory Efficiency**: Optimized storage to minimize memory footprint
- **Low Latency**: Sub-millisecond overhead for cache lookups
- **Concurrent Access**: Thread-safe implementation for high concurrency
- **Minimal GC Impact**: Careful memory management to reduce garbage collection pressure

### Benchmarks

```
BenchmarkCacheMiddleware/Cache_Hit-8         	 5000000	       247 ns/op	      32 B/op	       1 allocs/op
BenchmarkCacheMiddleware/Cache_Miss-8        	 1000000	      1021 ns/op	     112 B/op	       3 allocs/op
BenchmarkCacheMiddleware/Redis_Hit-8         	 2000000	       634 ns/op	      48 B/op	       2 allocs/op
```

## 🔍 Troubleshooting

### Common Issues

| Issue | Solution |
|-------|----------|
| Cache not working | Ensure the middleware is registered before route handlers |
| Responses not being cached | Check if the response size exceeds `MaxBytes` or if the HTTP method is not in the `Methods` list |
| Cache not invalidating | Verify your `CacheInvalidator` function logic |
| Redis connection issues | Check Redis connection parameters and ensure the Redis server is running |

### Debugging

You can enable debug mode to get more information about cache operations:

```go
app.Use(cache.New(cache.Config{
    Debug: true,
}))
```

## 📄 License

This middleware is part of the Quick framework and is licensed under the MIT License.

## 👥 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
