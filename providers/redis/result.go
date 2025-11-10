package redis

import (
	rds "github.com/redis/go-redis/v9"
	"github.com/rockcookies/go-caches"
)

var _ caches.Result[any] = (*baseResult[any])(nil)

type baseResult[T any] struct {
	value T
	err   error
}

func newResult[T any](result T, err error) *baseResult[T] {
	if err == rds.Nil {
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
