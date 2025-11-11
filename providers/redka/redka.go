package redka

import (
	"strings"

	rdk "github.com/nalgeon/redka"
)

type Options struct {
	Prefix string
}

type Provider struct {
	db     *rdk.DB
	prefix string
}

func New(db *rdk.DB) *Provider {
	return NewWithOptions(db, nil)
}

func NewWithOptions(db *rdk.DB, opts *Options) *Provider {
	if db == nil {
		panic("db is nil")
	}

	if opts == nil {
		opts = &Options{}
	}

	return &Provider{
		db:     db,
		prefix: strings.TrimSpace(opts.Prefix),
	}
}

func (p *Provider) Prefix() string {
	return p.prefix
}
