package caches

import (
	"context"
	"time"
)

// KeyScanResult represents the result of a key scan operation.
type KeyScanResult struct {
	// Cursor is the cursor to use in the next scan call.
	// A cursor value of 0 indicates the iteration is complete.
	Cursor uint64
	// Keys contains the scanned keys.
	Keys []string
}

// KeyCommand defines operations for key management and lifecycle in the cache.
// This includes key creation, deletion, expiration, scanning, and metadata operations.
type KeyCommand interface {
	// DBSize returns the number of keys in the current database.
	DBSize(ctx context.Context) Result[int64]

	// Del deletes one or more keys.
	// Non-existing keys are ignored.
	// Returns the number of keys that were deleted.
	Del(ctx context.Context, keys ...string) Result[int64]

	// Exists checks if one or more keys exist.
	// Returns the number of keys that exist among the given keys.
	Exists(ctx context.Context, keys ...string) Result[int64]

	// Expire sets a timeout on a key using seconds.
	// After the timeout expires, the key will be automatically deleted.
	// Returns true if the timeout was set, false if the key does not exist.
	Expire(ctx context.Context, key string, expiration time.Duration) Result[bool]

	// ExpireNX sets a timeout on a key only if the key has no existing expiration.
	// Returns true if the timeout was set, false if the key does not exist or already has an expiration.
	ExpireNX(ctx context.Context, key string, expiration time.Duration) Result[bool]

	// ExpireXX sets a timeout on a key only if the key has an existing expiration.
	// Returns true if the timeout was set, false if the key does not exist or has no expiration.
	ExpireXX(ctx context.Context, key string, expiration time.Duration) Result[bool]

	// ExpireGT sets a timeout on a key only if the new expiration is greater than the current one.
	// Returns true if the timeout was set, false otherwise.
	ExpireGT(ctx context.Context, key string, expiration time.Duration) Result[bool]

	// ExpireLT sets a timeout on a key only if the new expiration is less than the current one.
	// Returns true if the timeout was set, false otherwise.
	ExpireLT(ctx context.Context, key string, expiration time.Duration) Result[bool]

	// ExpireAt sets an expiration timestamp on a key.
	// The key will be automatically deleted at the specified UNIX time in seconds.
	// Returns true if the expiration was set, false if the key does not exist.
	ExpireAt(ctx context.Context, key string, tm time.Time) Result[bool]

	// ExpireTime returns the absolute UNIX timestamp at which the key will expire.
	// Returns -1 if the key exists but has no associated expiration.
	// Returns -2 if the key does not exist.
	ExpireTime(ctx context.Context, key string) Result[time.Duration]

	// PExpire sets a timeout on a key using milliseconds.
	// After the timeout expires, the key will be automatically deleted.
	// Returns true if the timeout was set, false if the key does not exist.
	PExpire(ctx context.Context, key string, expiration time.Duration) Result[bool]

	// PExpireAt sets an expiration timestamp on a key using milliseconds.
	// The key will be automatically deleted at the specified UNIX time in milliseconds.
	// Returns true if the expiration was set, false if the key does not exist.
	PExpireAt(ctx context.Context, key string, tm time.Time) Result[bool]

	// PExpireTime returns the absolute UNIX timestamp in milliseconds at which the key will expire.
	// Returns -1 if the key exists but has no associated expiration.
	// Returns -2 if the key does not exist.
	PExpireTime(ctx context.Context, key string) Result[time.Duration]

	// FlushAll deletes all keys from the current database.
	// This operation is irreversible.
	FlushAll(ctx context.Context) StatusResult

	// Persist removes the expiration timeout from a key, making it persistent.
	// Returns true if the timeout was removed, false if the key does not exist or has no expiration.
	Persist(ctx context.Context, key string) Result[bool]

	// Keys returns all keys matching a pattern.
	// Pattern syntax: * matches any number of characters, ? matches a single character, [] matches character ranges.
	// Use with caution on large databases as this is a O(N) operation.
	Keys(ctx context.Context, pattern string) Result[[]string]

	// Rename renames a key to a new key.
	// If the new key already exists, it will be overwritten.
	// Returns an error if the source key does not exist.
	Rename(ctx context.Context, key string, newKey string) StatusResult

	// RenameNX renames a key to a new key only if the new key does not exist.
	// Returns true if the key was renamed, false if the new key already exists or the source key does not exist.
	RenameNX(ctx context.Context, key string, newKey string) Result[bool]

	// TTL returns the remaining time to live of a key in seconds.
	// Returns -1 if the key exists but has no associated expiration.
	// Returns -2 if the key does not exist.
	TTL(ctx context.Context, key string) Result[time.Duration]

	// PTTL returns the remaining time to live of a key in milliseconds.
	// Returns -1 if the key exists but has no associated expiration.
	// Returns -2 if the key does not exist.
	PTTL(ctx context.Context, key string) Result[time.Duration]

	// Type returns the string representation of the type of the value stored at key.
	// Returns "none" if the key does not exist.
	Type(ctx context.Context, key string) Result[string]

	// RandomKey returns a random key from the current database.
	// Returns an empty string if the database is empty.
	RandomKey(ctx context.Context) Result[string]

	// Scan iterates over keys in the database.
	// cursor is the cursor to start iteration from (0 to start).
	// match is a glob-style pattern to filter keys (empty string for no filter).
	// count is a hint for how many keys to return per iteration.
	// Returns the next cursor and a slice of keys.
	Scan(ctx context.Context, cursor uint64, match string, count int64) Result[KeyScanResult]
}
