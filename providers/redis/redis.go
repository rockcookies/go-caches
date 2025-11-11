package redis

import (
	"strings"

	rds "github.com/redis/go-redis/v9"
)

type Options struct {
	Prefix string
}

type Provider struct {
	db     rds.UniversalClient
	prefix string
}

func New(client rds.UniversalClient) *Provider {
	return NewWithOptions(client, nil)
}

func NewWithOptions(client rds.UniversalClient, opts *Options) *Provider {
	if client == nil {
		panic("client is nil")
	}

	if opts == nil {
		opts = &Options{}
	}

	return &Provider{
		db:     client,
		prefix: strings.TrimSpace(opts.Prefix),
	}
}
