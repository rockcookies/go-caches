package redis

import (
	"context"

	"github.com/rockcookies/go-caches"
)

var _ caches.SetCommand = (*Provider)(nil)

// SAdd implements caches.SetCommand.
func (p *Provider) SAdd(ctx context.Context, key string, members ...any) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.SAdd(ctx, key, members...)
	res.SetErr(formatError(res.Err()))
	return res
}

// SCard implements caches.SetCommand.
func (p *Provider) SCard(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.SCard(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// SDiff implements caches.SetCommand.
func (p *Provider) SDiff(ctx context.Context, keys ...string) caches.Result[[][]byte] {
	keys = prefixKeys(p.prefix, keys)
	res := p.db.SDiff(ctx, keys...)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// SDiffStore implements caches.SetCommand.
func (p *Provider) SDiffStore(ctx context.Context, destination string, keys ...string) caches.Result[int64] {
	destination = p.prefix + destination
	keys = prefixKeys(p.prefix, keys)
	res := p.db.SDiffStore(ctx, destination, keys...)
	res.SetErr(formatError(res.Err()))
	return res
}

// SInter implements caches.SetCommand.
func (p *Provider) SInter(ctx context.Context, keys ...string) caches.Result[[][]byte] {
	keys = prefixKeys(p.prefix, keys)
	res := p.db.SInter(ctx, keys...)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// SInterStore implements caches.SetCommand.
func (p *Provider) SInterStore(ctx context.Context, destination string, keys ...string) caches.Result[int64] {
	destination = p.prefix + destination
	keys = prefixKeys(p.prefix, keys)
	res := p.db.SInterStore(ctx, destination, keys...)
	res.SetErr(formatError(res.Err()))
	return res
}

// SIsMember implements caches.SetCommand.
func (p *Provider) SIsMember(ctx context.Context, key string, member any) caches.Result[bool] {
	key = p.prefix + key
	res := p.db.SIsMember(ctx, key, member)
	res.SetErr(formatError(res.Err()))
	return res
}

// SMembers implements caches.SetCommand.
func (p *Provider) SMembers(ctx context.Context, key string) caches.Result[[][]byte] {
	key = p.prefix + key
	res := p.db.SMembers(ctx, key)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// SMove implements caches.SetCommand.
func (p *Provider) SMove(ctx context.Context, source, destination string, member any) caches.Result[bool] {
	source = p.prefix + source
	destination = p.prefix + destination
	res := p.db.SMove(ctx, source, destination, member)
	res.SetErr(formatError(res.Err()))
	return res
}

// SPop implements caches.SetCommand.
func (p *Provider) SPop(ctx context.Context, key string) caches.Result[[]byte] {
	key = p.prefix + key
	res := p.db.SPop(ctx, key)
	return newResult(res.Bytes())
}

// SPopN implements caches.SetCommand.
func (p *Provider) SPopN(ctx context.Context, key string, count int64) caches.Result[[][]byte] {
	key = p.prefix + key
	res := p.db.SPopN(ctx, key, count)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// SRandMember implements caches.SetCommand.
func (p *Provider) SRandMember(ctx context.Context, key string) caches.Result[[]byte] {
	key = p.prefix + key
	res := p.db.SRandMember(ctx, key)
	return newResult(res.Bytes())
}

// SRandMemberN implements caches.SetCommand.
func (p *Provider) SRandMemberN(ctx context.Context, key string, count int64) caches.Result[[][]byte] {
	key = p.prefix + key
	res := p.db.SRandMemberN(ctx, key, count)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// SRem implements caches.SetCommand.
func (p *Provider) SRem(ctx context.Context, key string, members ...any) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.SRem(ctx, key, members...)
	res.SetErr(formatError(res.Err()))
	return res
}

// SScan implements caches.SetCommand.
func (p *Provider) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) caches.Result[caches.ScanResult] {
	key = p.prefix + key
	res := p.db.SScan(ctx, key, cursor, match, count)

	if res.Err() != nil {
		return newResult(caches.ScanResult{}, formatError(res.Err()))
	}

	keys, newCursor := res.Val()
	elements := make([][]byte, len(keys))
	for i, value := range keys {
		elements[i] = []byte(value)
	}

	return newResult(caches.ScanResult{
		Cursor:   newCursor,
		Elements: elements,
	}, nil)
}

// SUnion implements caches.SetCommand.
func (p *Provider) SUnion(ctx context.Context, keys ...string) caches.Result[[][]byte] {
	keys = prefixKeys(p.prefix, keys)
	res := p.db.SUnion(ctx, keys...)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// SUnionStore implements caches.SetCommand.
func (p *Provider) SUnionStore(ctx context.Context, destination string, keys ...string) caches.Result[int64] {
	destination = p.prefix + destination
	keys = prefixKeys(p.prefix, keys)
	res := p.db.SUnionStore(ctx, destination, keys...)
	res.SetErr(formatError(res.Err()))
	return res
}
