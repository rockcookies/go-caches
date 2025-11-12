package caches

import (
	"context"
	"time"
)

// SetArgs provides arguments for the SetArgs function.
type SetArgs struct {
	// Mode can be `NX` or `XX` or empty.
	Mode string

	// Zero `TTL` or `Expiration` means that the key has no expiration time.
	TTL      time.Duration
	ExpireAt time.Time

	// When Get is true, the command returns the old value stored at key, or nil when key did not exist.
	Get bool

	// KeepTTL is a Redis KEEPTTL option to keep existing TTL, it requires your redis-server version >= 6.0,
	// otherwise you will receive an error: (error) ERR syntax error.
	KeepTTL bool
}

// StringCommand defines operations for string values in the cache.
// Strings are the most basic Redis data type and can contain text, JSON, serialized objects,
// or binary data. All string operations are atomic.
type StringCommand interface {
	// Decr decrements the integer value of a key by one.
	// If the key does not exist, it is set to 0 before performing the operation.
	// Returns an error if the key contains a value that cannot be interpreted as an integer.
	Decr(ctx context.Context, key string) Result[int64]

	// DecrBy decrements the integer value of a key by the specified amount.
	// If the key does not exist, it is set to 0 before performing the operation.
	// Returns an error if the key contains a value that cannot be interpreted as an integer.
	DecrBy(ctx context.Context, key string, value int64) Result[int64]

	// Get retrieves the value of a key.
	// Returns nil if the key does not exist.
	Get(ctx context.Context, key string) Result[[]byte]

	// Incr increments the integer value of a key by one.
	// If the key does not exist, it is set to 0 before performing the operation.
	// Returns an error if the key contains a value that cannot be interpreted as an integer.
	Incr(ctx context.Context, key string) Result[int64]

	// IncrBy increments the integer value of a key by the specified amount.
	// If the key does not exist, it is set to 0 before performing the operation.
	// Returns an error if the key contains a value that cannot be interpreted as an integer.
	IncrBy(ctx context.Context, key string, value int64) Result[int64]

	// IncrByFloat increments the float value of a key by the specified amount.
	// If the key does not exist, it is set to 0 before performing the operation.
	// Returns an error if the key contains a value that cannot be interpreted as a float.
	IncrByFloat(ctx context.Context, key string, value float64) Result[float64]

	// Set sets the value of a key with an optional expiration time.
	// expiration of 0 means the key has no expiration time.
	// Overwrites any existing value and clears any existing TTL.
	Set(ctx context.Context, key string, value any, expiration time.Duration) StatusResult

	// SetArgs sets the value of a key with advanced options specified in SetArgs.
	// Provides fine-grained control over set operations including mode (NX/XX), TTL, and Get options.
	SetArgs(ctx context.Context, key string, value any, args SetArgs) StatusResult

	// SetNX sets the value of a key only if the key does not exist.
	// Returns true if the key was set, false if the key already exists.
	// expiration of 0 means the key has no expiration time.
	SetNX(ctx context.Context, key string, value any, expiration time.Duration) Result[bool]

	// SetXX sets the value of a key only if the key already exists.
	// Returns true if the key was set, false if the key does not exist.
	// expiration of 0 means the key has no expiration time.
	SetXX(ctx context.Context, key string, value any, expiration time.Duration) Result[bool]

	// StrLen returns the length of the string value stored at a key.
	// Returns 0 if the key does not exist.
	StrLen(ctx context.Context, key string) Result[int64]

	// MGet retrieves the values of multiple keys.
	// For each key that does not exist, the corresponding value in the result map will be nil.
	// Returns a map where keys are the requested keys and values are their corresponding values.
	MGet(ctx context.Context, keys ...string) Result[map[string][]byte]

	// MSet sets the values of multiple keys to their corresponding values.
	// This operation is atomic: either all keys are set or none are.
	// Overwrites any existing values and clears any existing TTLs.
	MSet(ctx context.Context, values map[string]any) StatusResult

	// MSetNX sets the values of multiple keys to their corresponding values only if none of the keys exist.
	// This operation is atomic: either all keys are set or none are.
	// Returns true if all keys were set, false if any key already exists.
	MSetNX(ctx context.Context, values map[string]any) Result[bool]
}
