package redis

import (
	"context"

	"github.com/rockcookies/go-caches"
)

var _ caches.HashCommand = (*Provider)(nil)

// HDel implements caches.HashCommand.
func (p *Provider) HDel(ctx context.Context, key string, fields ...string) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.HDel(ctx, key, fields...)
	res.SetErr(formatError(res.Err()))
	return res
}

// HExists implements caches.HashCommand.
func (p *Provider) HExists(ctx context.Context, key string, field string) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.HExists(ctx, key, field)
	res.SetErr(formatError(res.Err()))
	return res
}

// HGet implements caches.HashCommand.
func (p *Provider) HGet(ctx context.Context, key string, field string) caches.Result[[]byte] {
	key = p.prefix + key
	res := p.db.HGet(ctx, key, field)
	return newResult(res.Bytes())
}

// HGetAll implements caches.HashCommand.
func (p *Provider) HGetAll(ctx context.Context, key string) caches.Result[map[string][]byte] {
	key = p.prefix + key
	res := p.db.HGetAll(ctx, key)

	if res.Err() != nil {
		return newResult(map[string][]byte(nil), formatError(res.Err()))
	}

	result := make(map[string][]byte, len(res.Val()))
	for field, value := range res.Val() {
		result[field] = []byte(value)
	}

	return newResult(result, nil)
}

// HIncrBy implements caches.HashCommand.
func (p *Provider) HIncrBy(ctx context.Context, key string, field string, increment int64) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.HIncrBy(ctx, key, field, increment)
	res.SetErr(formatError(res.Err()))
	return res
}

// HIncrByFloat implements caches.HashCommand.
func (p *Provider) HIncrByFloat(ctx context.Context, key string, field string, increment float64) caches.Result[float64] {
	key = p.prefix + key
	res := p.db.HIncrByFloat(ctx, key, field, increment)
	res.SetErr(formatError(res.Err()))
	return res
}

// HKeys implements caches.HashCommand.
func (p *Provider) HKeys(ctx context.Context, key string) caches.Result[[]string] {
	key = p.prefix + key
	res := p.db.HKeys(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// HLen implements caches.HashCommand.
func (p *Provider) HLen(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.HLen(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// HMGet implements caches.HashCommand.
func (p *Provider) HMGet(ctx context.Context, key string, fields ...string) caches.Result[map[string][]byte] {
	key = p.prefix + key
	res := p.db.HMGet(ctx, key, fields...)

	if res.Err() != nil {
		return newResult(map[string][]byte(nil), formatError(res.Err()))
	}

	result := make(map[string][]byte, len(fields))
	for i, value := range res.Val() {
		if value != nil {
			if strVal, ok := value.(string); ok {
				result[fields[i]] = []byte(strVal)
			}
		}
	}

	return newResult(result, nil)
}

// HMSet implements caches.HashCommand.
func (p *Provider) HMSet(ctx context.Context, key string, values map[string]any) caches.StatusResult {
	key = p.prefix + key
	res := p.db.HSet(ctx, key, values)

	// Convert IntCmd to StatusResult
	if res.Err() != nil {
		return newStatusResult(nil, formatError(res.Err()))
	}
	return newStatusResult([]byte("OK"), nil)
}

// HScan implements caches.HashCommand.
func (p *Provider) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) caches.Result[caches.HScanResult] {
	key = p.prefix + key
	res := p.db.HScan(ctx, key, cursor, match, count)

	if res.Err() != nil {
		return newResult(caches.HScanResult{}, formatError(res.Err()))
	}

	keys, newCursor := res.Val()
	fields := make(map[string][]byte)

	// HScan returns field-value pairs as a flat slice
	for i := 0; i < len(keys); i += 2 {
		if i+1 < len(keys) {
			fields[keys[i]] = []byte(keys[i+1])
		}
	}

	return newResult(caches.HScanResult{
		Cursor: newCursor,
		Fields: fields,
	}, nil)
}

// HSet implements caches.HashCommand.
func (p *Provider) HSet(ctx context.Context, key string, values map[string]any) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.HSet(ctx, key, values)
	res.SetErr(formatError(res.Err()))
	return res
}

// HSetNX implements caches.HashCommand.
func (p *Provider) HSetNX(ctx context.Context, key string, field string, value any) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.HSetNX(ctx, key, field, value)
	res.SetErr(formatError(res.Err()))
	return res
}

// HVals implements caches.HashCommand.
func (p *Provider) HVals(ctx context.Context, key string) caches.Result[[][]byte] {
	key = p.prefix + key
	res := p.db.HVals(ctx, key)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}
