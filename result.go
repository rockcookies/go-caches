package caches

// Result provides a generic container for operation results with error handling.
// This interface allows providers to return values along with potential errors in a consistent way.
// The generic type T represents the type of the value returned by the operation.
type Result[T any] interface {
	// SetErr sets the error on the result.
	// This is typically used by providers during error handling.
	SetErr(e error)

	// Err returns the error associated with the result.
	// Returns nil if the operation was successful.
	Err() error

	// SetVal sets the value on the result.
	// This is typically used by providers during successful operations.
	SetVal(v T)

	// Val returns the value stored in the result.
	// If an error is present, this may return the zero value for T.
	Val() T

	// Result returns both the value and error.
	// This is the most convenient way to handle the result in a single call.
	Result() (T, error)
}

// StatusResult specializes Result[string] for operations that return status messages.
// This interface is used for operations where the result is typically a status string
// (like "OK") but also provides access to the underlying byte representation.
type StatusResult interface {
	Result[string]

	// Bytes returns the status value as a byte slice along with any error.
	// This is useful when you need the raw byte representation of the status.
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
	return string(s.value), s.err
}

// SetVal implements caches.StatusResult.
func (s *statusResult) SetVal(v string) {
	s.value = []byte(v)
}

// Val implements caches.StatusResult.
func (s *statusResult) Val() string {
	return string(s.value)
}

// Bytes implements caches.StatusResult.
func (s *statusResult) Bytes() ([]byte, error) {
	return s.value, s.err
}
