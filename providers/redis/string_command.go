package redis

import (
	"context"
	"time"

	rds "github.com/redis/go-redis/v9"
	"github.com/rockcookies/go-caches"
)

var _ caches.StringCommand = (*Provider)(nil)

// Decr implements caches.StringCommand.
func (p *Provider) Decr(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.Decr(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// DecrBy implements caches.StringCommand.
func (p *Provider) DecrBy(ctx context.Context, key string, value int64) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.DecrBy(ctx, key, value)
	res.SetErr(formatError(res.Err()))
	return res
}

// Get implements caches.StringCommand.
func (p *Provider) Get(ctx context.Context, key string) caches.Result[[]byte] {
	key = p.prefix + key
	res := p.db.Get(ctx, key)
	return newResult(res.Bytes())
}

// GetSet implements caches.StringCommand.
func (p *Provider) GetSet(ctx context.Context, key string, value any) caches.Result[[]byte] {
	key = p.prefix + key
	res := p.db.GetSet(ctx, key, value)
	res.SetErr(formatError(res.Err()))
	return newResult(res.Bytes())
}

// Incr implements caches.StringCommand.
func (p *Provider) Incr(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.Incr(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// IncrBy implements caches.StringCommand.
func (p *Provider) IncrBy(ctx context.Context, key string, value int64) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.IncrBy(ctx, key, value)
	res.SetErr(formatError(res.Err()))
	return res
}

// IncrByFloat implements caches.StringCommand.
func (p *Provider) IncrByFloat(ctx context.Context, key string, value float64) caches.Result[float64] {
	key = p.prefix + key
	res := p.db.IncrByFloat(ctx, key, value)
	res.SetErr(formatError(res.Err()))
	return res
}

// Set implements caches.StringCommand.
func (p *Provider) Set(ctx context.Context, key string, value any, expiration time.Duration) caches.StatusResult {
	key = p.prefix + key
	res := p.db.Set(ctx, key, value, expiration)
	res.SetErr(formatError(res.Err()))
	return res
}

// SetArgs implements caches.StringCommand.
func (p *Provider) SetArgs(ctx context.Context, key string, value any, args caches.SetArgs) caches.StatusResult {
	key = p.prefix + key
	res := p.db.SetArgs(ctx, key, value, rds.SetArgs{
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
func (p *Provider) SetNX(ctx context.Context, key string, value any, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.SetNX(ctx, key, value, expiration)
	res.SetErr(formatError(res.Err()))
	return res
}

// SetXX implements caches.StringCommand.
func (p *Provider) SetXX(ctx context.Context, key string, value any, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.SetXX(ctx, key, value, expiration)
	res.SetErr(formatError(res.Err()))
	return res
}

// StrLen implements caches.StringCommand.
func (p *Provider) StrLen(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.StrLen(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// MGet implements caches.StringCommand.
func (p *Provider) MGet(ctx context.Context, keys ...string) caches.Result[map[string][]byte] {
	keys = prefixKeys(p.prefix, keys)
	res := p.db.MGet(ctx, keys...)

	result := make(map[string][]byte)
	prefixLen := len(p.prefix)

	for i, val := range res.Val() {
		if val != nil {
			// 去除前缀，返回原始键名
			originalKey := keys[i]
			if prefixLen > 0 && len(keys[i]) > prefixLen {
				originalKey = keys[i][prefixLen:]
			}

			// 将值转换为字节数组
			if strVal, ok := val.(string); ok {
				result[originalKey] = []byte(strVal)
			}
		}
	}

	return newResult(result, formatError(res.Err()))
}

// MSet implements caches.StringCommand.
func (p *Provider) MSet(ctx context.Context, values map[string]any) caches.StatusResult {
	pairs := make([]any, 0, len(values)*2)
	for key, value := range values {
		pairs = append(pairs, p.prefix+key, value)
	}

	res := p.db.MSet(ctx, pairs...)
	res.SetErr(formatError(res.Err()))
	return res
}

// MSetNX implements caches.StringCommand.
func (p *Provider) MSetNX(ctx context.Context, values map[string]any) caches.Result[bool] {
	pairs := make([]any, 0, len(values)*2)
	for key, value := range values {
		pairs = append(pairs, p.prefix+key, value)
	}

	res := p.db.MSetNX(ctx, pairs...)
	res.SetErr(formatError(res.Err()))
	return res
}
