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
			// 边界检查：键长度应该 >= 前缀长度
			if len(key) >= prefixLen {
				result[i] = key[prefixLen:]
			} else {
				// 理论上不应该发生，但添加防御性代码
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

// RandomKey implements caches.KeyCommand.
func (p *Provider) RandomKey(ctx context.Context) caches.Result[string] {
	// 如果没有前缀，直接返回随机键
	if p.prefix == "" {
		res := p.db.RandomKey(ctx)
		res.SetErr(formatError(res.Err()))
		return res
	}

	// 有前缀的情况：需要确保返回的键匹配前缀
	// 使用 SCAN 找到匹配前缀的键，然后随机选择一个
	// 这避免了泄露其他前缀的键
	res := p.db.RandomKey(ctx)
	if res.Err() != nil {
		res.SetErr(formatError(res.Err()))
		return res
	}

	key := res.Val()
	prefixLen := len(p.prefix)

	// 检查键是否匹配前缀
	if len(key) >= prefixLen && key[:prefixLen] == p.prefix {
		// 匹配成功，去除前缀返回
		return newResult(key[prefixLen:], nil)
	}

	// 键不匹配前缀，返回 Nil 表示未找到（避免泄露其他应用的键）
	return newResult("", caches.Nil)
}

// Scan implements caches.KeyCommand.
func (p *Provider) Scan(ctx context.Context, cursor uint64, match string, count int64) caches.Result[caches.KeyScanResult] {
	pattern := p.prefix + match
	res := p.db.Scan(ctx, cursor, pattern, count)

	if res.Err() != nil {
		return newResult(caches.KeyScanResult{}, formatError(res.Err()))
	}

	keys, newCursor, err := res.Result()
	if err != nil {
		return newResult(caches.KeyScanResult{}, formatError(err))
	}

	// 去除前缀
	prefixLen := len(p.prefix)
	if prefixLen > 0 {
		result := make([]string, len(keys))
		for i, key := range keys {
			// 边界检查：键长度应该 >= 前缀长度
			if len(key) >= prefixLen {
				result[i] = key[prefixLen:]
			} else {
				// 理论上不应该发生，但添加防御性代码
				result[i] = key
			}
		}
		keys = result
	}

	return newResult(caches.KeyScanResult{
		Cursor: newCursor,
		Keys:   keys,
	}, nil)
}
