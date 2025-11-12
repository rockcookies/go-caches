package redka

import (
	"context"

	rdk "github.com/nalgeon/redka"
	"github.com/rockcookies/go-caches"
)

var _ caches.HashCommand = (*Provider)(nil)

// HSet implements caches.HashCommand.
func (p *Provider) HSet(ctx context.Context, key string, values map[string]any) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		count, e := tx.Hash().SetMany(key, values)
		return int64(count), e
	})
	return newResult(n, err)
}

// HGet implements caches.HashCommand.
func (p *Provider) HGet(ctx context.Context, key, field string) caches.Result[[]byte] {
	key = p.prefix + key
	val, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]byte, error) {
		v, e := tx.Hash().Get(key, field)
		if e != nil {
			return nil, e
		}
		return v.Bytes(), nil
	})
	return newResult(val, err)
}

// HGetAll implements caches.HashCommand.
func (p *Provider) HGetAll(ctx context.Context, key string) caches.Result[map[string][]byte] {
	key = p.prefix + key
	items, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (map[string][]byte, error) {
		vals, e := tx.Hash().Items(key)
		if e != nil {
			return nil, e
		}

		result := make(map[string][]byte, len(vals))
		for k, v := range vals {
			result[k] = v.Bytes()
		}
		return result, nil
	})
	return newResult(items, err)
}

// HDel implements caches.HashCommand.
func (p *Provider) HDel(ctx context.Context, key string, fields ...string) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		count, e := tx.Hash().Delete(key, fields...)
		return int64(count), e
	})
	return newResult(n, err)
}

// HExists implements caches.HashCommand.
func (p *Provider) HExists(ctx context.Context, key, field string) caches.Result[bool] {
	key = p.prefix + key
	exists, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (bool, error) {
		return tx.Hash().Exists(key, field)
	})
	return newResult(exists, err)
}

// HIncrBy implements caches.HashCommand.
func (p *Provider) HIncrBy(ctx context.Context, key, field string, incr int64) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		val, e := tx.Hash().Incr(key, field, int(incr))
		return int64(val), e
	})
	return newResult(n, err)
}

// HIncrByFloat implements caches.HashCommand.
func (p *Provider) HIncrByFloat(ctx context.Context, key, field string, incr float64) caches.Result[float64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (float64, error) {
		return tx.Hash().IncrFloat(key, field, incr)
	})
	return newResult(n, err)
}

// HKeys implements caches.HashCommand.
func (p *Provider) HKeys(ctx context.Context, key string) caches.Result[[]string] {
	key = p.prefix + key
	keys, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]string, error) {
		return tx.Hash().Fields(key)
	})
	return newResult(keys, err)
}

// HLen implements caches.HashCommand.
func (p *Provider) HLen(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	n, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		l, e := tx.Hash().Len(key)
		return int64(l), e
	})
	return newResult(n, err)
}

// HMGet implements caches.HashCommand.
func (p *Provider) HMGet(ctx context.Context, key string, fields ...string) caches.Result[map[string][]byte] {
	key = p.prefix + key
	values, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (map[string][]byte, error) {
		vals, e := tx.Hash().GetMany(key, fields...)
		if e != nil {
			return nil, e
		}
		result := make(map[string][]byte, len(vals))
		for field, val := range vals {
			result[field] = val.Bytes()
		}
		return result, nil
	})
	return newResult(values, err)
}

// HMSet implements caches.HashCommand.
func (p *Provider) HMSet(ctx context.Context, key string, values map[string]any) caches.StatusResult {
	key = p.prefix + key
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]byte, error) {
		_, e := tx.Hash().SetMany(key, values)
		if e != nil {
			return nil, e
		}
		return []byte("OK"), nil
	})
	return newStatusResult(val, err)
}

// HSetNX implements caches.HashCommand.
func (p *Provider) HSetNX(ctx context.Context, key, field string, value any) caches.Result[bool] {
	key = p.prefix + key
	set, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (bool, error) {
		return tx.Hash().SetNotExists(key, field, value)
	})
	return newResult(set, err)
}

// HVals implements caches.HashCommand.
func (p *Provider) HVals(ctx context.Context, key string) caches.Result[[][]byte] {
	key = p.prefix + key
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		values, e := tx.Hash().Values(key)
		if e != nil {
			return nil, e
		}
		result := make([][]byte, len(values))
		for i, v := range values {
			result[i] = v.Bytes()
		}
		return result, nil
	})
	return newResult(vals, err)
}

// HScan implements caches.HashCommand.
func (p *Provider) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) caches.Result[caches.HScanResult] {
	key = p.prefix + key
	result, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (caches.HScanResult, error) {
		// Use Scan method for cursor-based scanning
		scanRes, e := tx.Hash().Scan(key, int(cursor), match, int(count))
		if e != nil {
			return caches.HScanResult{}, e
		}

		fields := make(map[string][]byte, len(scanRes.Items))
		for _, item := range scanRes.Items {
			fields[item.Field] = item.Value.Bytes()
		}

		return caches.HScanResult{
			Cursor: uint64(scanRes.Cursor),
			Fields: fields,
		}, nil
	})
	return newResult(result, err)
}
