# Tests Package

This package contains unified tests for all go-caches providers using the Provider Interface Pattern.

## Quick Start

```bash
# Run all tests
go test -v ./tests/

# Run specific provider tests
go test -v -run TestRedis ./tests/
go test -v -run TestRedka ./tests/

# Run specific command tests
go test -v -run StringCommand ./tests/
go test -v -run KeyCommand ./tests/

# Skip integration tests (for CI)
go test -short ./tests/
```

## Architecture

Tests use the Provider Interface Pattern:

1. **Command Test Files** (`*_test.go`): Define provider interfaces and test logic
2. **Provider Suites** (`redis_test.go`, `redka_test.go`): Implement provider interfaces
3. **Test Runners** (`RunCommandTests`): Execute tests for all providers

### Example: StringCommand Tests

```go
// Define provider interface
type StringCommandProvider interface {
    GetStringCommand() caches.StringCommand
    GetContext() context.Context
}

// Test runner
func RunStringCommandTests(t *testing.T, provider StringCommandProvider) {
    t.Run("SetAndGet", func(t *testing.T) {
        testSetAndGet(t, provider)
    })
    // ... more tests
}

// Test function
func testSetAndGet(t *testing.T, provider StringCommandProvider) {
    cmd := provider.GetStringCommand()
    ctx := provider.GetContext()
    // Test implementation
}
```

## Adding New Tests

1. **Create Command Test File** (`command_test.go`):
   - Define `CommandProvider` interface
   - Implement `RunCommandTests()` function
   - Write individual test functions

2. **Update Provider Suites** (`redis_test.go`, `redka_test.go`):
   - Implement the provider interface methods
   - Add `TestCommand()` method that calls the test runner

3. **Run Tests**: Verify both providers pass all tests

### Naming Conventions

- **Provider Interface**: `<Command>Provider` (e.g., `StringCommandProvider`)
- **Test Runner**: `Run<Command>Tests` (e.g., `RunStringCommandTests`)
- **Test Functions**: `test<Operation>` (e.g., `testSetAndGet`)
- **Test Methods**: `Test<Command>` (e.g., `TestStringCommand`)

## Environment Requirements

- **Redis Provider**: Requires Redis server on `localhost:6379` (auto-skipped with `-short` flag)
- **Redka Provider**: No external dependencies (uses in-memory database)

## Test Structure

```
tests/
├── go.mod                    # Module definition with all dependencies
├── string_test.go           # StringCommand interface tests
├── key_test.go              # KeyCommand interface tests
├── hash_test.go             # HashCommand interface tests
├── list_test.go             # ListCommand interface tests
├── set_test.go              # SetCommand interface tests
├── sorted_set_test.go       # SortedSetCommand interface tests
├── redis_test.go            # Redis provider test suite
├── redka_test.go            # Redka provider test suite
└── README.md                # This file
```

## Testing Best Practices

- Use unique key prefixes to avoid test collisions
- Each test should be independent and not rely on other tests
- Test both success and error conditions
- Verify return values and edge cases