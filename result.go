package caches

type Result[T any] interface {
	SetErr(e error)
	Err() error
	SetVal(v T)
	Val() T
	Result() (T, error)
}

type StatusResult interface {
	Result[string]
	Bytes() ([]byte, error)
}

type genericResult[T any] struct {
	value T
	err   error
}

// NewResult creates a new genericResult with the given value and error.
// This is intended to be used by providers with their own error handling.
func NewResult[T any](result T, err error) *genericResult[T] {
	return &genericResult[T]{
		value: result,
		err:   err,
	}
}

// SetErr implements caches.Result.
func (b *genericResult[T]) SetErr(e error) {
	b.err = e
}

// Err implements caches.Result.
func (b *genericResult[T]) Err() error {
	return b.err
}

// SetVal implements caches.Result.
func (b *genericResult[T]) SetVal(v T) {
	b.value = v
}

// Val implements caches.Result.
func (b *genericResult[T]) Val() T {
	return b.value
}

// Result implements caches.Result.
func (b *genericResult[T]) Result() (T, error) {
	return b.value, b.err
}

// statusResult provides a concrete implementation of the StatusResult interface.
type statusResult struct {
	*genericResult[[]byte]
}

// NewStatusResult creates a new statusResult with the given byte slice value and error.
// This is intended to be used by providers with their own error handling.
func NewStatusResult(val []byte, err error) StatusResult {
	return &statusResult{
		NewResult(val, err),
	}
}

// Result implements caches.StatusResult.
func (s *statusResult) Result() (string, error) {
	return bytesToString(s.value), s.err
}

// SetVal implements caches.StatusResult.
func (s *statusResult) SetVal(v string) {
	s.value = stringToBytes(v)
}

// Val implements caches.StatusResult.
func (s *statusResult) Val() string {
	return bytesToString(s.value)
}

// Bytes implements caches.StatusResult.
func (s *statusResult) Bytes() ([]byte, error) {
	return s.value, s.err
}
