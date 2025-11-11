package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/rockcookies/go-caches"
)

func formatError(err error) error {
	if err == redis.Nil {
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
