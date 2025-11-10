package redis

import (
	"context"
	"time"

	rds "github.com/redis/go-redis/v9"
	"github.com/rockcookies/go-caches"
)

var _ caches.StringCommand = (*RedisCache)(nil)

// Decr implements caches.StringCommand.
func (r *RedisCache) Decr(ctx context.Context, key string) caches.Result[int64] {
	key = r.formatKey(key)
	res := r.db.Decr(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// DecrBy implements caches.StringCommand.
func (r *RedisCache) DecrBy(ctx context.Context, key string, value int64) caches.Result[int64] {
	key = r.formatKey(key)
	res := r.db.DecrBy(ctx, key, value)
	res.SetErr(formatError(res.Err()))
	return res
}

// Get implements caches.StringCommand.
func (r *RedisCache) Get(ctx context.Context, key string) caches.Result[[]byte] {
	key = r.formatKey(key)
	res := r.db.Get(ctx, key)
	return newResult(res.Bytes())
}

// GetSet implements caches.StringCommand.
func (r *RedisCache) GetSet(ctx context.Context, key string, value any) caches.Result[[]byte] {
	key = r.formatKey(key)
	res := r.db.GetSet(ctx, key, value)
	res.SetErr(formatError(res.Err()))
	return newResult(res.Bytes())
}

// Incr implements caches.StringCommand.
func (r *RedisCache) Incr(ctx context.Context, key string) caches.Result[int64] {
	key = r.formatKey(key)
	res := r.db.Incr(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// IncrBy implements caches.StringCommand.
func (r *RedisCache) IncrBy(ctx context.Context, key string, value int64) caches.Result[int64] {
	key = r.formatKey(key)
	res := r.db.IncrBy(ctx, key, value)
	res.SetErr(formatError(res.Err()))
	return res
}

// IncrByFloat implements caches.StringCommand.
func (r *RedisCache) IncrByFloat(ctx context.Context, key string, value float64) caches.Result[float64] {
	key = r.formatKey(key)
	res := r.db.IncrByFloat(ctx, key, value)
	res.SetErr(formatError(res.Err()))
	return res
}

// Set implements caches.StringCommand.
func (r *RedisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) caches.StatusResult {
	key = r.formatKey(key)
	res := r.db.Set(ctx, key, value, expiration)
	res.SetErr(formatError(res.Err()))
	return res
}

// SetArgs implements caches.StringCommand.
func (r *RedisCache) SetArgs(ctx context.Context, key string, value any, args caches.SetArgs) caches.StatusResult {
	key = r.formatKey(key)
	res := r.db.SetArgs(ctx, key, value, rds.SetArgs{
		Mode:     args.Mode,
		TTL:      args.TTL,
		ExpireAt: args.ExpireAt,
		Get:      args.Get,
		KeepTTL:  args.KeepTTL,
	})
	res.SetErr(formatError(res.Err()))
	return res
}

// SetNX implements caches.StringCommand.
func (r *RedisCache) SetNX(ctx context.Context, key string, value any, expiration time.Duration) caches.Result[bool] {
	key = r.formatKey(key)
	res := r.db.SetNX(ctx, key, value, expiration)
	res.SetErr(formatError(res.Err()))
	return res
}

// SetXX implements caches.StringCommand.
func (r *RedisCache) SetXX(ctx context.Context, key string, value any, expiration time.Duration) caches.Result[bool] {
	key = r.formatKey(key)
	res := r.db.SetXX(ctx, key, value, expiration)
	res.SetErr(formatError(res.Err()))
	return res
}

// StrLen implements caches.StringCommand.
func (r *RedisCache) StrLen(ctx context.Context, key string) caches.Result[int64] {
	key = r.formatKey(key)
	res := r.db.StrLen(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}
