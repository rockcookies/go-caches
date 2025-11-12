package redka

import (
	"context"
	"time"

	rdk "github.com/nalgeon/redka"
	"github.com/rockcookies/go-caches"
)

// newResult creates a new BaseResult with Redka-specific error handling.
// It converts Redka ErrNotFound responses to caches.Nil for consistency.
func newResult[T any](result T, err error) caches.Result[T] {
	if err == rdk.ErrNotFound {
		err = caches.Nil
	}

	return caches.NewResult(result, err)
}

// newStatusResult creates a new statusResult with Redka-specific error handling.
// It converts Redka ErrNotFound responses to caches.Nil for consistency.
func newStatusResult(val []byte, err error) caches.StatusResult {
	if err == rdk.ErrNotFound {
		err = caches.Nil
	}

	return caches.NewStatusResult(val, err)
}

func prefixKeys(prefix string, keys []string) []string {
	if prefix == "" {
		return keys
	}

	prefixed := make([]string, len(keys))
	for i, key := range keys {
		prefixed[i] = prefix + key
	}

	return prefixed
}

func viewAndReturn[T any](ctx context.Context, db *rdk.DB, cb func(tx *rdk.Tx) (T, error)) (res T, err error) {
	err = db.ViewContext(ctx, func(tx *rdk.Tx) (e error) {
		res, e = cb(tx)
		return
	})
	return
}

func updateAndReturn[T any](ctx context.Context, db *rdk.DB, cb func(tx *rdk.Tx) (T, error)) (res T, err error) {
	err = db.UpdateContext(ctx, func(tx *rdk.Tx) (e error) {
		res, e = cb(tx)
		return
	})
	return
}

func formatMs(dur time.Duration) int64 {
	if dur > 0 && dur < time.Millisecond {
		return 1
	}
	return int64(dur / time.Millisecond)
}

func formatSec(dur time.Duration) int64 {
	if dur > 0 && dur < time.Second {
		return 1
	}
	return int64(dur / time.Second)
}
