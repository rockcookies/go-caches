# go-caches

A unified caching interface for Go applications with pluggable backends (Redis and Redka). This library provides a consistent API for cache operations while allowing you to switch between different storage backends without changing your application code.

## Features

- **🔄 Unified Interface**: Single API for multiple cache backends
- **📦 Multiple Providers**: Support for Redis and Redka (SQLite-based) backends
- **🎯 Type-Safe**: Generic results with proper error handling using `Result[T]`
- **⚡ Full Redis API Support**: Complete implementation of Redis commands
- **🔧 Namespace Isolation**: Key prefixing for multi-tenant applications
- **🧪 Well Tested**: Comprehensive test suite with integration tests
- **📈 High Performance**: Minimal overhead and efficient backend implementations

## Quick Start

### Installation

```bash
go get github.com/rockcookies/go-caches
go get github.com/rockcookies/go-caches/providers/redis
go get github.com/rockcookies/go-caches/providers/redka
```

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/rockcookies/go-caches"
    "github.com/rockcookies/go-caches/providers/redis"
)

func main() {
    // Redis Provider
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    cache := redis.New(rdb)
    ctx := context.Background()

    // Set a value
    result := cache.Set(ctx, "key", "value", time.Hour)
    if result.Err() != nil {
        panic(result.Err())
    }

    // Get a value
    value := cache.Get(ctx, "key")
    if value.Err() != nil {
        panic(value.Err())
    }

    fmt.Printf("Value: %s\n", string(value.Val()))
}
```

### Switching Backends

```go
// Redis Provider
redisCache := redis.New(redisClient)

// Redka Provider (SQLite-based)
redkaCache, _ := redka.New("file:cache.db?mode=memory")

// Both implement the same interface!
func useCache(cache caches.StringCommand) {
    // Your cache logic works with any provider
    result := cache.Get(ctx, "key")
    // ...
}
```

## API Overview

The library provides six main command interfaces:

### StringCommand
String value operations with atomic increments/decrements:

```go
// Basic operations
cache.Set(ctx, "counter", "0", 0)           // Set value
result := cache.Get(ctx, "counter")         // Get value

// Atomic operations
cache.Incr(ctx, "counter")                  // Increment by 1
cache.IncrBy(ctx, "counter", 10)            // Increment by 10
cache.DecrBy(ctx, "counter", 5)             // Decrement by 5

// Multi-key operations
values := map[string]any{
    "key1": "value1",
    "key2": "value2",
}
cache.MSet(ctx, values)                     // Set multiple keys
```

### KeyCommand
Key lifecycle management:

```go
// Expiration
cache.Expire(ctx, "key", time.Hour)         // Set expiration
cache.TTL(ctx, "key")                       // Get remaining TTL
cache.Persist(ctx, "key")                   // Remove expiration

// Key operations
cache.Exists(ctx, "key1", "key2")           // Check if keys exist
cache.Del(ctx, "key1", "key2")              // Delete keys
cache.Type(ctx, "key")                      // Get key type
```

### SetCommand, HashCommand, ListCommand, SortedSetCommand
Full Redis data structure support:

```go
// Sets
cache.SAdd(ctx, "myset", "member1", "member2")
cache.SMembers(ctx, "myset")

// Hashes
cache.HSet(ctx, "myhash", "field", "value")
cache.HGet(ctx, "myhash", "field")

// Lists
cache.LPush(ctx, "mylist", "item1", "item2")
cache.LRange(ctx, "mylist", 0, -1)

// Sorted Sets
cache.ZAdd(ctx, "myzset", &caches.ZMember{
    Member: []byte("member"),
    Score:  1.0,
})
cache.ZRange(ctx, "myzset", 0, -1)
```

## Configuration

### Provider Options

```go
// With key prefixing for namespace isolation
options := &redis.Options{
    Prefix: "myapp:v1:",
}

cache := redis.NewWithOptions(redisClient, options)
```

### Advanced Set Operations

```go
args := caches.SetArgs{
    Mode:    "NX",        // Set only if key doesn't exist
    TTL:     time.Hour,   // Auto-expire after 1 hour
    Get:     true,        // Return old value
    KeepTTL: false,       // Don't preserve existing TTL
}

result := cache.SetArgs(ctx, "key", "value", args)
```

## Error Handling

The library uses a consistent `Result[T]` pattern for error handling:

```go
result := cache.Get(ctx, "key")

// Check for errors
if result.Err() != nil {
    if result.Err() == caches.Nil {
        fmt.Println("Key does not exist")
    } else {
        fmt.Printf("Error: %v\n", result.Err())
    }
    return
}

// Get the value (zero value if error occurred)
value := result.Val()

// Or get both at once
value, err := result.Result()
```

## Testing

The library includes comprehensive tests:

```bash
# Run all tests
go test -race -coverprofile=coverage.out ./tests/

# Run specific provider tests
go test -v -run TestRedis ./tests/
go test -v -run TestRedka ./tests/

# Skip integration tests (CI environments)
go test -short ./tests/
```

## Providers

### Redis Provider
- **Backend**: Redis server using `github.com/redis/go-redis/v9`
- **Features**: Full Redis API support, clustering, sentinel support
- **Requirements**: Redis server running on `localhost:6379` (for development)

### Redka Provider
- **Backend**: SQLite using `github.com/nalgeon/redka`
- **Features**: Zero dependencies, in-memory databases, transactions
- **Requirements**: None (uses embedded SQLite)

## Architecture

The library follows an interface-driven design:

```
caches package (interfaces)
├── StringCommand    # String operations
├── KeyCommand       # Key management
├── SetCommand       # Set data structure
├── HashCommand      # Hash data structure
├── ListCommand      # List data structure
└── SortedSetCommand # Sorted set data structure

providers/
├── redis/           # Redis provider implementation
└── redka/           # Redka provider implementation
```

## Dependencies

- **Go**: 1.21+ (root module), 1.21+ (provider modules)
- **Redis Provider**: `github.com/redis/go-redis/v9 v9.16.0`
- **Redka Provider**: `github.com/nalgeon/redka v0.6.0`
- **Testing**: `github.com/stretchr/testify v1.11.1`

## License

[Add your license information here]

## Contributing

[Add contributing guidelines here]
