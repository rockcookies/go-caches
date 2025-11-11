package redka

import (
	"context"
	"reflect"
	"time"
	"unsafe"

	rdk "github.com/nalgeon/redka"
)

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

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//nolint:staticcheck
func stringToBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
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
