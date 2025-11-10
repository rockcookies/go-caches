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

type StringCommand interface {
	Decr(ctx context.Context, key string) Result[int64]
	DecrBy(ctx context.Context, key string, value int64) Result[int64]
	Get(ctx context.Context, key string) Result[[]byte]
	GetSet(ctx context.Context, key string, value any) Result[[]byte]
	Incr(ctx context.Context, key string) Result[int64]
	IncrBy(ctx context.Context, key string, value int64) Result[int64]
	IncrByFloat(ctx context.Context, key string, value float64) Result[float64]
	Set(ctx context.Context, key string, value any, expiration time.Duration) StatusResult
	SetArgs(ctx context.Context, key string, value any, args SetArgs) StatusResult
	SetNX(ctx context.Context, key string, value any, expiration time.Duration) Result[bool]
	SetXX(ctx context.Context, key string, value any, expiration time.Duration) Result[bool]
	StrLen(ctx context.Context, key string) Result[int64]
}
