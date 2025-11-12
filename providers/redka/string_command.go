package redka

import (
	"context"
	"strings"
	"time"

	rdk "github.com/nalgeon/redka"
	"github.com/rockcookies/go-caches"
)

var _ caches.StringCommand = (*Provider)(nil)

func (p *Provider) incr(ctx context.Context, key string, value int) caches.Result[int64] {
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int, error) {
		return tx.Str().Incr(key, value)
	})

	return newResult(int64(val), err)
}

// Decr implements caches.StringCommand.
func (p *Provider) Decr(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	return p.incr(ctx, key, -1)
}

// DecrBy implements caches.StringCommand.
func (p *Provider) DecrBy(ctx context.Context, key string, value int64) caches.Result[int64] {
	key = p.prefix + key
	return p.incr(ctx, key, int(-value))
}

// Get implements caches.StringCommand.
func (p *Provider) Get(ctx context.Context, key string) caches.Result[[]byte] {
	key = p.prefix + key
	val, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]byte, error) {
		res, err := tx.Str().Get(key)
		if err != nil {
			return nil, err
		}
		return res.Bytes(), nil
	})

	return newResult(val, err)
}

// Incr implements caches.StringCommand.
func (p *Provider) Incr(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	return p.incr(ctx, key, 1)
}

// IncrBy implements caches.StringCommand.
func (p *Provider) IncrBy(ctx context.Context, key string, value int64) caches.Result[int64] {
	key = p.prefix + key
	return p.incr(ctx, key, int(value))
}

func (p *Provider) incrFloat(ctx context.Context, key string, value float64) caches.Result[float64] {
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (float64, error) {
		return tx.Str().IncrFloat(key, value)
	})

	return newResult(val, err)
}

// IncrByFloat implements caches.StringCommand.
func (p *Provider) IncrByFloat(ctx context.Context, key string, value float64) caches.Result[float64] {
	key = p.prefix + key
	return p.incrFloat(ctx, key, value)
}

// Set implements caches.StringCommand.
func (p *Provider) Set(ctx context.Context, key string, value any, expiration time.Duration) caches.StatusResult {
	key = p.prefix + key
	val, _, err := p.setArgs(ctx, key, value, &caches.SetArgs{
		TTL:     expiration,
		KeepTTL: expiration == caches.KeepTTL,
	})
	return newStatusResult(val, err)
}

// SetArgs implements caches.StringCommand.
func (p *Provider) SetArgs(ctx context.Context, key string, value any, args caches.SetArgs) caches.StatusResult {
	key = p.prefix + key
	val, _, err := p.setArgs(ctx, key, value, &args)
	return newStatusResult(val, err)
}

// SetNX implements caches.StringCommand.
func (p *Provider) SetNX(ctx context.Context, key string, value any, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	_, ok, err := p.setArgs(ctx, key, value, &caches.SetArgs{
		Mode:    "NX",
		TTL:     expiration,
		KeepTTL: expiration == caches.KeepTTL,
	})
	// SetNX 在失败时不应该返回错误，而是返回 false
	if err == rdk.ErrNotFound {
		return newResult(false, nil)
	}
	return newResult(ok, err)
}

// SetXX implements caches.StringCommand.
func (p *Provider) SetXX(ctx context.Context, key string, value any, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	_, ok, err := p.setArgs(ctx, key, value, &caches.SetArgs{
		Mode:    "XX",
		TTL:     expiration,
		KeepTTL: expiration == caches.KeepTTL,
	})
	// SetXX 在失败时不应该返回错误，而是返回 false
	if err == rdk.ErrNotFound {
		return newResult(false, nil)
	}
	return newResult(ok, err)
}

func (p *Provider) StrLen(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	val, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (int, error) {
		res, err := tx.Str().Get(key)
		if err == rdk.ErrNotFound {
			// 键不存在时返回 0
			return 0, nil
		} else if err != nil {
			return 0, err
		}
		return len(res.Bytes()), nil
	})

	return newResult(int64(val), err)
}

func (p *Provider) setArgs(ctx context.Context, key string, value any, args *caches.SetArgs) ([]byte, bool, error) {
	type Res struct {
		Prev    rdk.Value
		Created bool
		Updated bool
	}

	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (*Res, error) {
		set := tx.Str().SetWith(key, value)

		if strings.ToUpper(args.Mode) == "NX" {
			set = set.IfNotExists()
		} else if strings.ToUpper(args.Mode) == "XX" {
			set = set.IfExists()
		}

		if args.KeepTTL {
			set = set.KeepTTL()
		}

		if !args.ExpireAt.IsZero() {
			set = set.At(args.ExpireAt)
		}

		if args.TTL > 0 {
			set = set.TTL(args.TTL)
		}

		res, err := set.Run()
		if err != nil {
			return nil, err
		} else {
			return &Res{
				Prev:    res.Prev,
				Created: res.Created,
				Updated: res.Updated,
			}, nil
		}
	})
	if err != nil {
		return nil, false, err
	}

	if args.Get {
		if val.Prev.IsZero() {
			return nil, false, rdk.ErrNotFound
		} else {
			return val.Prev, true, nil
		}
	}

	if !val.Created && !val.Updated {
		return nil, false, rdk.ErrNotFound
	}

	return []byte("OK"), true, nil
}

// MGet implements caches.StringCommand.
func (p *Provider) MGet(ctx context.Context, keys ...string) caches.Result[map[string][]byte] {
	keys = prefixKeys(p.prefix, keys)
	val, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (map[string][]byte, error) {
		values, err := tx.Str().GetMany(keys...)
		if err != nil {
			return nil, err
		}

		result := make(map[string][]byte, len(values))
		prefixLen := len(p.prefix)

		for key, value := range values {
			// 去除前缀，返回原始键名
			originalKey := key
			if prefixLen > 0 && len(key) > prefixLen {
				originalKey = key[prefixLen:]
			}
			result[originalKey] = value.Bytes()
		}

		return result, nil
	})

	return newResult(val, err)
}

// MSet implements caches.StringCommand.
func (p *Provider) MSet(ctx context.Context, values map[string]any) caches.StatusResult {
	prefixedValues := make(map[string]any, len(values))
	for key, value := range values {
		prefixedValues[p.prefix+key] = value
	}

	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]byte, error) {
		err := tx.Str().SetMany(prefixedValues)
		if err != nil {
			return nil, err
		}
		return []byte("OK"), nil
	})
	return newStatusResult(val, err)
}

// MSetNX implements caches.StringCommand.
func (p *Provider) MSetNX(ctx context.Context, values map[string]any) caches.Result[bool] {
	prefixedValues := make(map[string]any, len(values))
	for key, value := range values {
		prefixedValues[p.prefix+key] = value
	}

	err := p.db.UpdateContext(ctx, func(tx *rdk.Tx) error {
		// 使用 SetArgs 设置所有键，如果任何键已存在则失败
		for key, value := range prefixedValues {
			set := tx.Str().SetWith(key, value).IfNotExists()
			res, err := set.Run()
			if err != nil {
				return err
			}
			// 如果没有创建成功（键已存在），返回错误以回滚事务
			if !res.Created {
				return rdk.ErrNotFound
			}
		}
		return nil
	})

	if err == rdk.ErrNotFound {
		// 有键已存在，返回 false
		return newResult(false, nil)
	} else if err != nil {
		// 其他错误
		return newResult(false, err)
	}

	return newResult(true, nil)
}
