package redka

import (
	"context"
	"reflect"
	"unsafe"

	rdk "github.com/nalgeon/redka"
)

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
