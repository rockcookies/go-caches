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
