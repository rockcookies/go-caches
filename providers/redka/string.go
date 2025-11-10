package redka

import (
	"context"
	"strings"
	"time"

	rdk "github.com/nalgeon/redka"
	"github.com/rockcookies/go-caches"
)

var _ caches.StringCommand = (*RedkaCache)(nil)

func (r *RedkaCache) incr(ctx context.Context, key string, value int) caches.Result[int64] {
	key = r.formatKey(key)

	val, err := updateAndReturn(ctx, r.db, func(tx *rdk.Tx) (int, error) {
		return tx.Str().Incr(key, value)
	})

	return newResult(int64(val), err)
}

func (r *RedkaCache) incrFloat(ctx context.Context, key string, value float64) caches.Result[float64] {
	key = r.formatKey(key)

	val, err := updateAndReturn(ctx, r.db, func(tx *rdk.Tx) (float64, error) {
		return tx.Str().IncrFloat(key, value)
	})

	return newResult(val, err)
}

// Decr implements caches.StringCommand.
func (r *RedkaCache) Decr(ctx context.Context, key string) caches.Result[int64] {
	return r.incr(ctx, key, -1)
}

// DecrBy implements caches.StringCommand.
func (r *RedkaCache) DecrBy(ctx context.Context, key string, value int64) caches.Result[int64] {
	return r.incr(ctx, key, int(-value))
}

// Get implements caches.StringCommand.
func (r *RedkaCache) Get(ctx context.Context, key string) caches.Result[[]byte] {
	key = r.formatKey(key)

	val, err := viewAndReturn(ctx, r.db, func(tx *rdk.Tx) (rdk.Value, error) {
		return tx.Str().Get(key)
	})

	return newResult(val.Bytes(), err)
}

// GetSet implements caches.StringCommand.
func (r *RedkaCache) GetSet(ctx context.Context, key string, value any) caches.Result[[]byte] {
	key = r.formatKey(key)

	val, err := updateAndReturn(ctx, r.db, func(tx *rdk.Tx) ([]byte, error) {
		res, err := tx.Str().SetWith(key, value).Run()
		return res.Prev, err
	})

	return newResult(val, err)
}

// Incr implements caches.StringCommand.
func (r *RedkaCache) Incr(ctx context.Context, key string) caches.Result[int64] {
	return r.incr(ctx, key, 1)
}

// IncrBy implements caches.StringCommand.
func (r *RedkaCache) IncrBy(ctx context.Context, key string, value int64) caches.Result[int64] {
	return r.incr(ctx, key, int(value))
}

// IncrByFloat implements caches.StringCommand.
func (r *RedkaCache) IncrByFloat(ctx context.Context, key string, value float64) caches.Result[float64] {
	return r.incrFloat(ctx, key, value)
}

func (r *RedkaCache) setArgs(ctx context.Context, key string, value any, args *caches.SetArgs) ([]byte, bool, error) {
	key = r.formatKey(key)

	val, err := updateAndReturn(ctx, r.db, func(tx *rdk.Tx) (struct {
		Prev    rdk.Value
		Created bool
		Updated bool
	}, error,
	) {
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

		tx.Key().Len()

		return set.Run()
	})
	if err != nil {
		return nil, false, rdk.ErrNotFound
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

	return []byte{'O', 'K'}, true, nil
}

// Set implements caches.StringCommand.
func (r *RedkaCache) Set(ctx context.Context, key string, value any, expiration time.Duration) caches.StatusResult {
	val, _, err := r.setArgs(ctx, key, value, &caches.SetArgs{
		TTL:     expiration,
		KeepTTL: expiration == caches.KeepTTL,
	})
	return newStatusResult(val, err)
}

// SetArgs implements caches.StringCommand.
func (r *RedkaCache) SetArgs(ctx context.Context, key string, value any, args caches.SetArgs) caches.StatusResult {
	val, _, err := r.setArgs(ctx, key, value, &args)
	return newStatusResult(val, err)
}

// SetNX implements caches.StringCommand.
func (r *RedkaCache) SetNX(ctx context.Context, key string, value any, expiration time.Duration) caches.Result[bool] {
	_, ok, err := r.setArgs(ctx, key, value, &caches.SetArgs{
		Mode:    "NX",
		TTL:     expiration,
		KeepTTL: expiration == caches.KeepTTL,
	})
	return newResult(ok, err)
}

// SetXX implements caches.StringCommand.
func (r *RedkaCache) SetXX(ctx context.Context, key string, value any, expiration time.Duration) caches.Result[bool] {
	_, ok, err := r.setArgs(ctx, key, value, &caches.SetArgs{
		Mode:    "XX",
		TTL:     expiration,
		KeepTTL: expiration == caches.KeepTTL,
	})
	return newResult(ok, err)
}

func (r *RedkaCache) StrLen(ctx context.Context, key string) caches.Result[int64] {
	key = r.formatKey(key)

	val, err := viewAndReturn(ctx, r.db, func(tx *rdk.Tx) (rdk.Value, error) {
		return tx.Str().Get(key)
	})

	return newResult(int64(len(val.Bytes())), err)
}
