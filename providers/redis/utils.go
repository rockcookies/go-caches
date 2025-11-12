package redis

import (
	rds "github.com/redis/go-redis/v9"
	"github.com/rockcookies/go-caches"
)

// newResult creates a new BaseResult with Redis-specific error handling.
// It converts Redis nil responses to caches.Nil for consistency.
func newResult[T any](result T, err error) caches.Result[T] {
	if err == rds.Nil {
		err = caches.Nil
	}

	return caches.NewResult(result, err)
}

// newStatusResult creates a new statusResult with Redis-specific error handling.
// It converts Redis nil responses to caches.Nil for consistency.
func newStatusResult(val []byte, err error) caches.StatusResult {
	if err == rds.Nil {
		err = caches.Nil
	}

	return caches.NewStatusResult(val, err)
}

func formatError(err error) error {
	if err == rds.Nil {
		return caches.Nil
	}
	return err
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
