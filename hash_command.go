package caches

import "context"

// HashCommand defines operations for Redis hash data structure.
// Hashes are field-value maps where both field and value are strings.
type HashCommand interface {
	// HDel deletes one or more fields from a hash.
	// Returns the number of fields that were removed from the hash,
	// not including non-existing fields.
	HDel(ctx context.Context, key string, fields ...string) Result[int64]

	// HExists determines whether a field exists in a hash.
	// Returns true if the field exists, false otherwise.
	HExists(ctx context.Context, key string, field string) Result[bool]

	// HGet returns the value of a field in a hash.
	// Returns nil if the field does not exist.
	HGet(ctx context.Context, key string, field string) Result[[]byte]

	// HGetAll returns all fields and values in a hash.
	// Returns an empty map if the key does not exist.
	HGetAll(ctx context.Context, key string) Result[map[string][]byte]

	// HIncrBy increments the integer value of a hash field by the given number.
	// If the field does not exist, it is set to 0 before performing the operation.
	// Returns the value of the field after the increment.
	HIncrBy(ctx context.Context, key string, field string, increment int64) Result[int64]

	// HIncrByFloat increments the float value of a hash field by the given amount.
	// If the field does not exist, it is set to 0 before performing the operation.
	// Returns the value of the field after the increment.
	HIncrByFloat(ctx context.Context, key string, field string, increment float64) Result[float64]

	// HKeys returns all field names in a hash.
	// Returns an empty slice if the key does not exist.
	HKeys(ctx context.Context, key string) Result[[]string]

	// HLen returns the number of fields in a hash.
	// Returns 0 if the key does not exist.
	HLen(ctx context.Context, key string) Result[int64]

	// HMGet returns the values of multiple fields in a hash.
	// For each field that does not exist, nil is returned in the corresponding position.
	// Returns a map with field names as keys and values as byte slices.
	HMGet(ctx context.Context, key string, fields ...string) Result[map[string][]byte]

	// HMSet sets the values of multiple fields in a hash.
	// If the hash does not exist, it is created.
	// If a field already exists, its value is overwritten.
	HMSet(ctx context.Context, key string, values map[string]any) StatusResult

	// HScan iterates over fields and values of a hash.
	// cursor is the cursor to start iteration from (0 to start).
	// match is a glob-style pattern to filter fields (empty string for no filter).
	// count is a hint for how many fields to return per iteration.
	// Returns the next cursor and a map of fields and values.
	HScan(ctx context.Context, key string, cursor uint64, match string, count int64) Result[HScanResult]

	// HSet sets the values of one or more fields in a hash.
	// If the hash does not exist, it is created.
	// If a field already exists, its value is overwritten.
	// Returns the number of fields that were added (not including updated fields).
	HSet(ctx context.Context, key string, values map[string]any) Result[int64]

	// HSetNX sets the value of a field in a hash only if the field does not exist.
	// If the field already exists, this operation has no effect.
	// Returns true if the field was set, false if the field already existed.
	HSetNX(ctx context.Context, key string, field string, value any) Result[bool]

	// HVals returns all values in a hash.
	// Returns an empty slice if the key does not exist.
	HVals(ctx context.Context, key string) Result[[][]byte]
}

// HScanResult represents the result of a hash scan operation.
type HScanResult struct {
	// Cursor is the cursor to use in the next scan call.
	// A cursor value of 0 indicates the iteration is complete.
	Cursor uint64
	// Fields contains the scanned fields and their values.
	Fields map[string][]byte
}
