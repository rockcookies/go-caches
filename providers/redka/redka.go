package redka

import (
	rdk "github.com/nalgeon/redka"
)

type Options struct {
	FormatKeyFunc func(key string) string
}

type RedkaCache struct {
	db        *rdk.DB
	formatKey func(key string) string
}

func New(db *rdk.DB) *RedkaCache {
	return NewWithOptions(db, nil)
}

func NewWithOptions(db *rdk.DB, opts *Options) *RedkaCache {
	if db == nil {
		panic("db is nil")
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

	return &RedkaCache{
		db:        db,
		formatKey: formatKey,
	}
}

func (r *RedkaCache) formatKeys(keys []string) []string {
	formatted := make([]string, len(keys))
	for i, key := range keys {
		formatted[i] = r.formatKey(key)
	}
	return formatted
}
