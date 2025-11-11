# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go caching library that provides a unified interface for multiple cache backends (Redis and Redka). The project uses an interface-driven design with provider implementations.

## Build and Development Commands

### Testing
```bash
# Run all tests with coverage and race detection
go test -race -coverprofile=coverage.out ./tests/

# Run specific provider tests
go test -v -run TestRedis ./tests/
go test -v -run TestRedka ./tests/

# Run tests with memory optimizations (for CI)
go test -short ./tests/
```

### Module Management
This is a Go workspace with multiple modules:
```bash
# Sync workspace dependencies
go work sync

# Build all modules
go build ./...
go build ./providers/redis
go build ./providers/redka
```

### Linting and Formatting
The project uses GitHub Actions for CI/CD with golangci-lint. No local Makefile is present.

## Architecture Overview

### Core Design Pattern
- **Interface-first**: All functionality is defined through interfaces in the root `caches` package
- **Provider pattern**: Multiple backends implement the same interfaces
- **Generic results**: Type-safe result wrappers with error handling using `Result[T]`

### Main Interfaces

#### StringCommand (15 methods)
String operations like Get, Set, Incr, MGet, MSet, etc.
- Location: `string_command.go:25-40`
- Fully implemented for both Redis and Redka providers

#### KeyCommand (29 methods)
Key management like Del, Exists, TTL, Expire, etc.
- Location: `key_command.go`
- **Status**: Redka provider fully implements, Redis provider has stubs only (panics)

### Provider Structure
```
caches package (interfaces)
├── providers/redis/     # Redis wrapper using go-redis/v9
└── providers/redka/     # Redka wrapper (SQLite-based)
```

### Key Abstractions

#### Result Pattern
```go
type Result[T any] interface {
    Result() (T, error)
    Val() T
    Err() error
}
```

#### SetArgs Configuration
```go
type SetArgs struct {
    Mode     string        // "NX", "XX", or empty
    TTL      time.Duration
    ExpireAt time.Time
    Get      bool          // Return old value
    KeepTTL  bool          // Preserve existing TTL (Redis >= 6.0)
}
```

## Provider Implementation Details

### Redis Provider (`providers/redis/`)
- Uses `github.com/redis/go-redis/v9`
- Supports key prefixing for namespace isolation
- **StringCommand**: Fully implemented
- **KeyCommand**: Not implemented (all methods panic)
- Error normalization: converts Redis nil responses to `caches.Nil`

### Redka Provider (`providers/redka/`)
- Uses `github.com/nalgeon/redka` (SQLite-based)
- Supports in-memory databases for testing
- **StringCommand + KeyCommand**: Fully implemented
- Advanced expiration modes (NX, XX, GT, LT)
- Transaction-based operations with `viewAndReturn`/`updateAndReturn` pattern

## Testing Strategy

### Test Architecture
- Unified test suite in `tests/` package using testify
- Provider interfaces defined in tests, providers implement them
- Integration tests with real backends (Redis requires localhost:6379)
- Memory-based testing for Redka (no external dependencies)

### Test Patterns
```go
// Tests expect providers to implement these interfaces
type StringCommandProvider interface {
    GetStringCommand() caches.StringCommand
    GetContext() context.Context
}

// Key prefixing to avoid test collisions
testKeyPrefix := "test:go-caches:" + uuid.New().String() + ":"
```

## Configuration

### Provider Options
Both providers support the same options structure:
```go
type Options struct {
    Prefix string  // Key prefix for namespace isolation
}
```

### Error Handling
- Custom error: `caches.CachesError`
- Standardized nil: `caches.Nil`
- Providers convert backend-specific errors to standard format

## Development Notes

### Current State
- ✅ StringCommand fully implemented and tested
- ✅ Redka KeyCommand fully implemented
- ❌ Redis KeyCommand needs implementation (currently panics)

### Go Workspace
- Root module: `go 1.21`
- Provider modules: `go 1.23`
- Uses `go.work` for multi-module development

### Dependencies
- Redis: `github.com/redis/go-redis/v9 v9.16.0`
- Redka: `github.com/nalgeon/redka v0.6.0`
- Testing: `github.com/stretchr/testify v1.11.1`