// Example of cache middleware with real Redis storage in Quick
package main

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/jeffotoni/quick"
// 	"github.com/jeffotoni/quick/middleware/cache"
// 	"github.com/redis/go-redis/v9"
// )

// // RedisAdapter adapts the go-redis client to implement the cache.RedisClient interface
// type RedisAdapter struct {
// 	client *redis.Client
// }

// // NewRedisAdapter creates a new Redis adapter
// func NewRedisAdapter(addr, password string, db int) *RedisAdapter {
// 	return &RedisAdapter{
// 		client: redis.NewClient(&redis.Options{
// 			Addr:     addr,
// 			Password: password,
// 			DB:       db,
// 		}),
// 	}
// }

// // Get implements the cache.RedisClient interface
// func (r *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
// 	return r.client.Get(ctx, key).Result()
// }

// // Set implements the cache.RedisClient interface
// func (r *RedisAdapter) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
// 	return r.client.Set(ctx, key, value, expiration).Err()
// }

// // Del implements the cache.RedisClient interface
// func (r *RedisAdapter) Del(ctx context.Context, keys ...string) (int64, error) {
// 	return r.client.Del(ctx, keys...).Result()
// }

// func main() {
// 	// Create a new Quick app
// 	q := quick.New()

// 	// Create Redis adapter
// 	redisAdapter := NewRedisAdapter("localhost:6379", "", 0)

// 	// Create Redis storage
// 	redisStorage, err := cache.NewRedisStorage(cache.RedisConfig{
// 		Client:    redisAdapter,
// 		TTL:       1 * time.Minute,
// 		KeyPrefix: "quick:cache",
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Use the cache middleware with Redis storage
// 	q.Use(cache.New(cache.Config{
// 		Expiration:           1 * time.Minute,
// 		CacheHeader:          "X-Cache-Status",
// 		CacheControl:         true,
// 		StoreResponseHeaders: true,
// 		Storage:              redisStorage,
// 	}))

// 	// Route 1: Returns the current time
// 	q.Get("/time", func(c *quick.Ctx) error {
// 		return c.String("Hora atual: " + time.Now().Format(time.RFC1123))
// 	})

// 	// Route 2: Returns a random number
// 	q.Get("/random", func(c *quick.Ctx) error {
// 		return c.String("Número aleatório: " + time.Now().Format("15:04:05.000"))
// 	})

// 	// Route 3: Returns JSON data
// 	q.Get("/profile", func(c *quick.Ctx) error {
// 		return c.JSON(quick.M{
// 			"user":  "jeffotoni",
// 			"since": time.Now().Format("2006-01-02 15:04:05"),
// 		})
// 	})

// 	// Start the server
// 	fmt.Println("Server running on http://localhost:3000")
// 	fmt.Println("Using real Redis storage (requires Redis server running on localhost:6379)")
// 	fmt.Println("Try these endpoints:")
// 	fmt.Println("  - GET /time (cached for 1 minute)")
// 	fmt.Println("  - GET /random (cached for 1 minute)")
// 	fmt.Println("  - GET /profile (cached for 1 minute)")
// 	fmt.Println("Check the X-Cache-Status header in the response to see if it's a HIT or MISS")
// 	q.Listen(":3000")
// }
