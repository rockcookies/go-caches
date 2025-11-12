package redka

import (
	"context"

	rdk "github.com/nalgeon/redka"
	"github.com/rockcookies/go-caches"
)

var _ caches.SetCommand = (*Provider)(nil)

// SAdd implements caches.SetCommand.
func (p *Provider) SAdd(ctx context.Context, key string, members ...any) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		count, e := tx.Set().Add(key, members...)
		return int64(count), e
	})
	return newResult(n, err)
}

// SCard implements caches.SetCommand.
func (p *Provider) SCard(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	n, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		count, e := tx.Set().Len(key)
		return int64(count), e
	})
	return newResult(n, err)
}

// SDiff implements caches.SetCommand.
func (p *Provider) SDiff(ctx context.Context, keys ...string) caches.Result[[][]byte] {
	keys = prefixKeys(p.prefix, keys)
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		items, e := tx.Set().Diff(keys...)
		if e != nil {
			return nil, e
		}
		result := make([][]byte, len(items))
		for i, v := range items {
			result[i] = v.Bytes()
		}
		return result, nil
	})
	return newResult(vals, err)
}

// SDiffStore implements caches.SetCommand.
func (p *Provider) SDiffStore(ctx context.Context, destination string, keys ...string) caches.Result[int64] {
	destination = p.prefix + destination
	keys = prefixKeys(p.prefix, keys)
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		count, e := tx.Set().DiffStore(destination, keys...)
		return int64(count), e
	})
	return newResult(n, err)
}

// SInter implements caches.SetCommand.
func (p *Provider) SInter(ctx context.Context, keys ...string) caches.Result[[][]byte] {
	keys = prefixKeys(p.prefix, keys)
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		items, e := tx.Set().Inter(keys...)
		if e != nil {
			return nil, e
		}
		result := make([][]byte, len(items))
		for i, v := range items {
			result[i] = v.Bytes()
		}
		return result, nil
	})
	return newResult(vals, err)
}

// SInterStore implements caches.SetCommand.
func (p *Provider) SInterStore(ctx context.Context, destination string, keys ...string) caches.Result[int64] {
	destination = p.prefix + destination
	keys = prefixKeys(p.prefix, keys)
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		count, e := tx.Set().InterStore(destination, keys...)
		return int64(count), e
	})
	return newResult(n, err)
}

// SIsMember implements caches.SetCommand.
func (p *Provider) SIsMember(ctx context.Context, key string, member any) caches.Result[bool] {
	key = p.prefix + key
	exists, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (bool, error) {
		return tx.Set().Exists(key, member)
	})
	return newResult(exists, err)
}

// SMembers implements caches.SetCommand.
func (p *Provider) SMembers(ctx context.Context, key string) caches.Result[[][]byte] {
	key = p.prefix + key
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		items, e := tx.Set().Items(key)
		if e != nil {
			return nil, e
		}
		result := make([][]byte, len(items))
		for i, v := range items {
			result[i] = v.Bytes()
		}
		return result, nil
	})
	return newResult(vals, err)
}

// SMove implements caches.SetCommand.
func (p *Provider) SMove(ctx context.Context, source, destination string, member any) caches.Result[bool] {
	source = p.prefix + source
	destination = p.prefix + destination
	moved, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (bool, error) {
		e := tx.Set().Move(source, destination, member)
		if e == rdk.ErrNotFound {
			return false, nil
		}
		return e == nil, e
	})
	return newResult(moved, err)
}

// SPop implements caches.SetCommand.
func (p *Provider) SPop(ctx context.Context, key string) caches.Result[[]byte] {
	key = p.prefix + key
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]byte, error) {
		v, e := tx.Set().Pop(key)
		if e != nil {
			return nil, e
		}
		return v.Bytes(), nil
	})
	return newResult(val, err)
}

// SPopN implements caches.SetCommand.
func (p *Provider) SPopN(ctx context.Context, key string, count int64) caches.Result[[][]byte] {
	key = p.prefix + key
	vals, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		result := make([][]byte, 0, count)
		for i := int64(0); i < count; i++ {
			v, e := tx.Set().Pop(key)
			if e != nil {
				if e == rdk.ErrNotFound {
					break
				}
				return nil, e
			}
			result = append(result, v.Bytes())
		}
		return result, nil
	})
	return newResult(vals, err)
}

// SRandMember implements caches.SetCommand.
func (p *Provider) SRandMember(ctx context.Context, key string) caches.Result[[]byte] {
	key = p.prefix + key
	val, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]byte, error) {
		v, e := tx.Set().Random(key)
		if e != nil {
			return nil, e
		}
		return v.Bytes(), nil
	})
	return newResult(val, err)
}

// SRandMemberN implements caches.SetCommand.
func (p *Provider) SRandMemberN(ctx context.Context, key string, count int64) caches.Result[[][]byte] {
	key = p.prefix + key
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		result := make([][]byte, 0, count)
		for i := int64(0); i < count; i++ {
			v, e := tx.Set().Random(key)
			if e != nil {
				if e == rdk.ErrNotFound {
					break
				}
				return nil, e
			}
			result = append(result, v.Bytes())
		}
		return result, nil
	})
	return newResult(vals, err)
}

// SRem implements caches.SetCommand.
func (p *Provider) SRem(ctx context.Context, key string, members ...any) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		count, e := tx.Set().Delete(key, members...)
		return int64(count), e
	})
	return newResult(n, err)
}

// SScan implements caches.SetCommand.
func (p *Provider) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) caches.Result[caches.ScanResult] {
	key = p.prefix + key
	result, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (caches.ScanResult, error) {
		scanRes, e := tx.Set().Scan(key, int(cursor), match, int(count))
		if e != nil {
			return caches.ScanResult{}, e
		}

		elements := make([][]byte, len(scanRes.Items))
		for i, v := range scanRes.Items {
			elements[i] = v.Bytes()
		}

		return caches.ScanResult{
			Cursor:   uint64(scanRes.Cursor),
			Elements: elements,
		}, nil
	})
	return newResult(result, err)
}

// SUnion implements caches.SetCommand.
func (p *Provider) SUnion(ctx context.Context, keys ...string) caches.Result[[][]byte] {
	keys = prefixKeys(p.prefix, keys)
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		items, e := tx.Set().Union(keys...)
		if e != nil {
			return nil, e
		}
		result := make([][]byte, len(items))
		for i, v := range items {
			result[i] = v.Bytes()
		}
		return result, nil
	})
	return newResult(vals, err)
}

// SUnionStore implements caches.SetCommand.
func (p *Provider) SUnionStore(ctx context.Context, destination string, keys ...string) caches.Result[int64] {
	destination = p.prefix + destination
	keys = prefixKeys(p.prefix, keys)
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		count, e := tx.Set().UnionStore(destination, keys...)
		return int64(count), e
	})
	return newResult(n, err)
}
