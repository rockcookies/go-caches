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
