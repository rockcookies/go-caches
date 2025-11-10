package redka

import (
	rdk "github.com/nalgeon/redka"
	"github.com/rockcookies/go-caches"
)

var (
	_ caches.Result[any]  = (*baseResult[any])(nil)
	_ caches.StatusResult = (*statusResult)(nil)
)

type baseResult[T any] struct {
	value T
	err   error
}

func newResult[T any](result T, err error) *baseResult[T] {
	if err == rdk.ErrNotFound {
		err = caches.Nil
	}

	return &baseResult[T]{
		value: result,
		err:   err,
	}
}

// SetErr implements caches.BaseCommand.
func (b *baseResult[T]) SetErr(e error) {
	b.err = e
}

// Err implements caches.BaseCommand.
func (b *baseResult[T]) Err() error {
	return b.err
}

// SetVal implements caches.BaseCommand.
func (b *baseResult[T]) SetVal(v T) {
	b.value = v
}

// Val implements caches.BaseCommand.
func (b *baseResult[T]) Val() T {
	return b.value
}

// Result implements caches.BaseCommand.
func (b *baseResult[T]) Result() (T, error) {
	return b.value, b.err
}

type statusResult struct {
	*baseResult[[]byte]
}

func newStatusResult(val []byte, err error) *statusResult {
	return &statusResult{
		newResult(val, err),
	}
}

// Result implements caches.StringCommand.
func (s *statusResult) Result() (string, error) {
	return bytesToString(s.value), s.err
}

// SetVal implements caches.StringCommand.
func (s *statusResult) SetVal(v string) {
	s.value = stringToBytes(v)
}

// Val implements caches.StringCommand.
func (s *statusResult) Val() string {
	return bytesToString(s.value)
}

// Bytes implements caches.StringCommand.
func (s *statusResult) Bytes() ([]byte, error) {
	return s.value, s.err
}
