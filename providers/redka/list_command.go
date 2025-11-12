package redka

import (
	"context"

	rdk "github.com/nalgeon/redka"
	"github.com/rockcookies/go-caches"
)

var _ caches.ListCommand = (*Provider)(nil)

// LIndex implements caches.ListCommand.
func (p *Provider) LIndex(ctx context.Context, key string, index int64) caches.Result[[]byte] {
	key = p.prefix + key
	val, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]byte, error) {
		v, e := tx.List().Get(key, int(index))
		if e != nil {
			return nil, e
		}
		return v.Bytes(), nil
	})
	return newResult(val, err)
}

// LInsert implements caches.ListCommand.
func (p *Provider) LInsert(ctx context.Context, key string, position caches.LInsertPosition, pivot, element any) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		var length int
		var e error
		if position == caches.LInsertBefore {
			length, e = tx.List().InsertBefore(key, pivot, element)
		} else {
			length, e = tx.List().InsertAfter(key, pivot, element)
		}
		// Redis returns -1 when pivot is not found, not an error
		if e == rdk.ErrNotFound {
			return -1, nil
		}
		return int64(length), e
	})
	return newResult(n, err)
}

// LLen implements caches.ListCommand.
func (p *Provider) LLen(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	n, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		l, e := tx.List().Len(key)
		return int64(l), e
	})
	return newResult(n, err)
}

// LPop implements caches.ListCommand.
func (p *Provider) LPop(ctx context.Context, key string) caches.Result[[]byte] {
	key = p.prefix + key
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]byte, error) {
		v, e := tx.List().PopFront(key)
		if e != nil {
			return nil, e
		}
		return v.Bytes(), nil
	})
	return newResult(val, err)
}

// LPopCount implements caches.ListCommand.
func (p *Provider) LPopCount(ctx context.Context, key string, count int) caches.Result[[][]byte] {
	key = p.prefix + key
	vals, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		result := make([][]byte, 0, count)
		for range count {
			v, e := tx.List().PopFront(key)
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

// LPush implements caches.ListCommand.
func (p *Provider) LPush(ctx context.Context, key string, elements ...any) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		for _, elem := range elements {
			_, e := tx.List().PushFront(key, elem)
			if e != nil {
				return 0, e
			}
		}
		l, e := tx.List().Len(key)
		return int64(l), e
	})
	return newResult(n, err)
}

// LRange implements caches.ListCommand.
func (p *Provider) LRange(ctx context.Context, key string, start, stop int64) caches.Result[[][]byte] {
	key = p.prefix + key
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		items, e := tx.List().Range(key, int(start), int(stop))
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

// LRem implements caches.ListCommand.
func (p *Provider) LRem(ctx context.Context, key string, count int64, element any) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		var deleted int
		var e error
		if count == 0 {
			// Remove all occurrences
			deleted, e = tx.List().Delete(key, element)
		} else if count > 0 {
			// Remove from front
			deleted, e = tx.List().DeleteFront(key, element, int(count))
		} else {
			// Remove from back
			deleted, e = tx.List().DeleteBack(key, element, int(-count))
		}
		return int64(deleted), e
	})
	return newResult(n, err)
}

// LSet implements caches.ListCommand.
func (p *Provider) LSet(ctx context.Context, key string, index int64, element any) caches.StatusResult {
	key = p.prefix + key
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]byte, error) {
		e := tx.List().Set(key, int(index), element)
		if e != nil {
			return nil, e
		}
		return []byte("OK"), nil
	})
	return newStatusResult(val, err)
}

// LTrim implements caches.ListCommand.
func (p *Provider) LTrim(ctx context.Context, key string, start, stop int64) caches.StatusResult {
	key = p.prefix + key
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]byte, error) {
		_, e := tx.List().Trim(key, int(start), int(stop))
		if e != nil {
			return nil, e
		}
		return []byte("OK"), nil
	})
	return newStatusResult(val, err)
}

// RPop implements caches.ListCommand.
func (p *Provider) RPop(ctx context.Context, key string) caches.Result[[]byte] {
	key = p.prefix + key
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]byte, error) {
		v, e := tx.List().PopBack(key)
		if e != nil {
			return nil, e
		}
		return v.Bytes(), nil
	})
	return newResult(val, err)
}

// RPopCount implements caches.ListCommand.
func (p *Provider) RPopCount(ctx context.Context, key string, count int) caches.Result[[][]byte] {
	key = p.prefix + key
	vals, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		result := make([][]byte, 0, count)
		for i := 0; i < count; i++ {
			v, e := tx.List().PopBack(key)
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

// RPush implements caches.ListCommand.
func (p *Provider) RPush(ctx context.Context, key string, elements ...any) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		for _, elem := range elements {
			_, e := tx.List().PushBack(key, elem)
			if e != nil {
				return 0, e
			}
		}
		l, e := tx.List().Len(key)
		return int64(l), e
	})
	return newResult(n, err)
}

// RPopLPush implements caches.ListCommand.
func (p *Provider) RPopLPush(ctx context.Context, source, destination string) caches.Result[[]byte] {
	source = p.prefix + source
	destination = p.prefix + destination
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]byte, error) {
		v, e := tx.List().PopBackPushFront(source, destination)
		if e != nil {
			return nil, e
		}
		return v.Bytes(), nil
	})
	return newResult(val, err)
}
