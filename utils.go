package caches

import (
	"reflect"
	"unsafe"
)

// bytesToString performs an unsafe conversion from []byte to string without allocation.
// This is used for performance-critical paths where the underlying data is not modified.
func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// stringToBytes performs an unsafe conversion from string to []byte without allocation.
// This is used for performance-critical paths where the underlying data is not modified.
//nolint:staticcheck
func stringToBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}