package redis

import (
	rds "github.com/redis/go-redis/v9"
)

type Options struct {
	FormatKeyFunc func(key string) string
}

type RedisCache struct {
	db        rds.UniversalClient
	formatKey func(key string) string
}

func New(client rds.UniversalClient) *RedisCache {
	return NewWithOptions(client, nil)
}

func NewWithOptions(client rds.UniversalClient, opts *Options) *RedisCache {
	if client == nil {
		panic("client is nil")
	}

	if opts == nil {
		opts = &Options{}
	}

	formatKey := opts.FormatKeyFunc
	if formatKey == nil {
		formatKey = func(key string) string {
			return key
		}
	}

	return &RedisCache{
		db:        client,
		formatKey: formatKey,
	}
}
