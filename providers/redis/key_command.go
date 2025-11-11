package redis

import (
	"context"
	"time"

	"github.com/rockcookies/go-caches"
)

var _ caches.KeyCommand = (*Provider)(nil)

// DBSize implements caches.KeyCommand.
func (p *Provider) DBSize(ctx context.Context) caches.Result[int64] {
	res := p.db.DBSize(ctx)
	res.SetErr(formatError(res.Err()))
	return res
}

// Del implements caches.KeyCommand.
func (p *Provider) Del(ctx context.Context, keys ...string) caches.Result[int64] {
	keys = prefixKeys(p.prefix, keys)
	res := p.db.Del(ctx, keys...)
	res.SetErr(formatError(res.Err()))
	return res
}

// Exists implements caches.KeyCommand.
func (p *Provider) Exists(ctx context.Context, keys ...string) caches.Result[int64] {
	keys = prefixKeys(p.prefix, keys)
	res := p.db.Exists(ctx, keys...)
	res.SetErr(formatError(res.Err()))
	return res
}

// Expire implements caches.KeyCommand.
func (p *Provider) Expire(ctx context.Context, key string, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.Expire(ctx, key, expiration)
	res.SetErr(formatError(res.Err()))
	return res
}

// ExpireAt implements caches.KeyCommand.
func (p *Provider) ExpireAt(ctx context.Context, key string, tm time.Time) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.ExpireAt(ctx, key, tm)
	res.SetErr(formatError(res.Err()))
	return res
}

// ExpireGT implements caches.KeyCommand.
func (p *Provider) ExpireGT(ctx context.Context, key string, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.ExpireGT(ctx, key, expiration)
	res.SetErr(formatError(res.Err()))
	return res
}

// ExpireLT implements caches.KeyCommand.
func (p *Provider) ExpireLT(ctx context.Context, key string, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.ExpireLT(ctx, key, expiration)
	res.SetErr(formatError(res.Err()))
	return res
}

// ExpireNX implements caches.KeyCommand.
func (p *Provider) ExpireNX(ctx context.Context, key string, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.ExpireNX(ctx, key, expiration)
	res.SetErr(formatError(res.Err()))
	return res
}

// ExpireTime implements caches.KeyCommand.
func (p *Provider) ExpireTime(ctx context.Context, key string) caches.Result[time.Duration] {
	key = p.prefix + key
	res := p.db.ExpireTime(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// ExpireXX implements caches.KeyCommand.
func (p *Provider) ExpireXX(ctx context.Context, key string, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.ExpireXX(ctx, key, expiration)
	res.SetErr(formatError(res.Err()))
	return res
}

// FlushAll implements caches.KeyCommand.
func (p *Provider) FlushAll(ctx context.Context) caches.StatusResult {
	res := p.db.FlushAll(ctx)
	res.SetErr(formatError(res.Err()))
	return res
}

// Keys implements caches.KeyCommand.
func (p *Provider) Keys(ctx context.Context, pattern string) caches.Result[[]string] {
	pattern = p.prefix + pattern
	res := p.db.Keys(ctx, pattern)

	// 去除前缀
	prefixLen := len(p.prefix)
	if prefixLen > 0 && res.Err() == nil {
		keys := res.Val()
		result := make([]string, len(keys))
		for i, key := range keys {
			if len(key) > prefixLen {
				result[i] = key[prefixLen:]
			} else {
				result[i] = key
			}
		}
		return newResult(result, formatError(res.Err()))
	}

	res.SetErr(formatError(res.Err()))
	return res
}

// PExpire implements caches.KeyCommand.
func (p *Provider) PExpire(ctx context.Context, key string, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.PExpire(ctx, key, expiration)
	res.SetErr(formatError(res.Err()))
	return res
}

// PExpireAt implements caches.KeyCommand.
func (p *Provider) PExpireAt(ctx context.Context, key string, tm time.Time) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.PExpireAt(ctx, key, tm)
	res.SetErr(formatError(res.Err()))
	return res
}

// PExpireTime implements caches.KeyCommand.
func (p *Provider) PExpireTime(ctx context.Context, key string) caches.Result[time.Duration] {
	key = p.prefix + key
	res := p.db.PExpireTime(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// PTTL implements caches.KeyCommand.
func (p *Provider) PTTL(ctx context.Context, key string) caches.Result[time.Duration] {
	key = p.prefix + key
	res := p.db.PTTL(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// Persist implements caches.KeyCommand.
func (p *Provider) Persist(ctx context.Context, key string) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.Persist(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// Rename implements caches.KeyCommand.
func (p *Provider) Rename(ctx context.Context, key string, newKey string) caches.StatusResult {
	key = p.prefix + key
	newKey = p.prefix + newKey
	res := p.db.Rename(ctx, key, newKey)
	res.SetErr(formatError(res.Err()))
	return res
}

// RenameNX implements caches.KeyCommand.
func (p *Provider) RenameNX(ctx context.Context, key string, newKey string) caches.Result[bool] {
	key = p.prefix + key
	newKey = p.prefix + newKey
	res := p.db.RenameNX(ctx, key, newKey)
	res.SetErr(formatError(res.Err()))
	return res
}

// TTL implements caches.KeyCommand.
func (p *Provider) TTL(ctx context.Context, key string) caches.Result[time.Duration] {
	key = p.prefix + key
	res := p.db.TTL(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// Type implements caches.KeyCommand.
func (p *Provider) Type(ctx context.Context, key string) caches.Result[string] {
	key = p.prefix + key
	res := p.db.Type(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}
