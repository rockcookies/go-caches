package redis

import (
	"context"

	"github.com/rockcookies/go-caches"
)

var _ caches.ListCommand = (*Provider)(nil)

// LIndex implements caches.ListCommand.
func (p *Provider) LIndex(ctx context.Context, key string, index int64) caches.Result[[]byte] {
	key = p.prefix + key
	res := p.db.LIndex(ctx, key, index)
	return newResult(res.Bytes())
}

// LInsert implements caches.ListCommand.
func (p *Provider) LInsert(ctx context.Context, key string, position caches.LInsertPosition, pivot, element any) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.LInsert(ctx, key, string(position), pivot, element)
	res.SetErr(formatError(res.Err()))
	return res
}

// LLen implements caches.ListCommand.
func (p *Provider) LLen(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.LLen(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// LPop implements caches.ListCommand.
func (p *Provider) LPop(ctx context.Context, key string) caches.Result[[]byte] {
	key = p.prefix + key
	res := p.db.LPop(ctx, key)
	return newResult(res.Bytes())
}

// LPopCount implements caches.ListCommand.
func (p *Provider) LPopCount(ctx context.Context, key string, count int) caches.Result[[][]byte] {
	key = p.prefix + key
	res := p.db.LPopCount(ctx, key, count)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// LPush implements caches.ListCommand.
func (p *Provider) LPush(ctx context.Context, key string, elements ...any) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.LPush(ctx, key, elements...)
	res.SetErr(formatError(res.Err()))
	return res
}

// LRange implements caches.ListCommand.
func (p *Provider) LRange(ctx context.Context, key string, start, stop int64) caches.Result[[][]byte] {
	key = p.prefix + key
	res := p.db.LRange(ctx, key, start, stop)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// LRem implements caches.ListCommand.
func (p *Provider) LRem(ctx context.Context, key string, count int64, element any) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.LRem(ctx, key, count, element)
	res.SetErr(formatError(res.Err()))
	return res
}

// LSet implements caches.ListCommand.
func (p *Provider) LSet(ctx context.Context, key string, index int64, element any) caches.StatusResult {
	key = p.prefix + key
	res := p.db.LSet(ctx, key, index, element)
	res.SetErr(formatError(res.Err()))
	return res
}

// LTrim implements caches.ListCommand.
func (p *Provider) LTrim(ctx context.Context, key string, start, stop int64) caches.StatusResult {
	key = p.prefix + key
	res := p.db.LTrim(ctx, key, start, stop)
	res.SetErr(formatError(res.Err()))
	return res
}

// RPop implements caches.ListCommand.
func (p *Provider) RPop(ctx context.Context, key string) caches.Result[[]byte] {
	key = p.prefix + key
	res := p.db.RPop(ctx, key)
	return newResult(res.Bytes())
}

// RPopCount implements caches.ListCommand.
func (p *Provider) RPopCount(ctx context.Context, key string, count int) caches.Result[[][]byte] {
	key = p.prefix + key
	res := p.db.RPopCount(ctx, key, count)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// RPopLPush implements caches.ListCommand.
func (p *Provider) RPopLPush(ctx context.Context, source, destination string) caches.Result[[]byte] {
	source = p.prefix + source
	destination = p.prefix + destination
	res := p.db.RPopLPush(ctx, source, destination)
	return newResult(res.Bytes())
}

// RPush implements caches.ListCommand.
func (p *Provider) RPush(ctx context.Context, key string, elements ...any) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.RPush(ctx, key, elements...)
	res.SetErr(formatError(res.Err()))
	return res
}
