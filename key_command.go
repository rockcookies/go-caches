package caches

import (
	"context"
	"time"
)

type KeyCommand interface {
	DBSize(ctx context.Context) Result[int64]
	Del(ctx context.Context, keys ...string) Result[int64]
	Exists(ctx context.Context, keys ...string) Result[int64]
	Expire(ctx context.Context, key string, expiration time.Duration) Result[bool]
	ExpireNX(ctx context.Context, key string, expiration time.Duration) Result[bool]
	ExpireXX(ctx context.Context, key string, expiration time.Duration) Result[bool]
	ExpireGT(ctx context.Context, key string, expiration time.Duration) Result[bool]
	ExpireLT(ctx context.Context, key string, expiration time.Duration) Result[bool]
	ExpireAt(ctx context.Context, key string, tm time.Time) Result[bool]
	ExpireTime(ctx context.Context, key string) Result[time.Duration]
	PExpire(ctx context.Context, key string, expiration time.Duration) Result[bool]
	PExpireAt(ctx context.Context, key string, tm time.Time) Result[bool]
	PExpireTime(ctx context.Context, key string) Result[time.Duration]
	FlushAll(ctx context.Context) StatusResult
	Persist(ctx context.Context, key string) Result[bool]
	Keys(ctx context.Context, pattern string) Result[[]string]
	Rename(ctx context.Context, key string, newKey string) StatusResult
	RenameNX(ctx context.Context, key string, newKey string) Result[bool]
	TTL(ctx context.Context, key string) Result[time.Duration]
	PTTL(ctx context.Context, key string) Result[time.Duration]
	Type(ctx context.Context, key string) Result[string]
}
