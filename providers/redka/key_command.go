package redka

import (
	"context"
	"time"

	rdk "github.com/nalgeon/redka"
	"github.com/rockcookies/go-caches"
)

type expireType uint

const (
	expire expireType = iota
	expireNX
	expireXX
	expireGT
	expireLT
)

var _ caches.KeyCommand = (*Provider)(nil)

// DBSize implements caches.KeyCommand.
func (p *Provider) DBSize(ctx context.Context) caches.Result[int64] {
	val, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (int, error) {
		return tx.Key().Len()
	})
	return newResult(int64(val), err)
}

// Del implements caches.KeyCommand.
func (p *Provider) Del(ctx context.Context, keys ...string) caches.Result[int64] {
	keys = prefixKeys(p.prefix, keys)
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int, error) {
		return tx.Key().Delete(keys...)
	})
	return newResult(int64(val), err)
}

// Exists implements caches.KeyCommand.
func (p *Provider) Exists(ctx context.Context, keys ...string) caches.Result[int64] {
	keys = prefixKeys(p.prefix, keys)
	val, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (int, error) {
		return tx.Key().Count(keys...)
	})
	return newResult(int64(val), err)
}

func (p *Provider) expire(ctx context.Context, key string, exp time.Duration, expType expireType) caches.Result[bool] {
	secs := formatSec(exp)
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (ok bool, err error) {
		// 获取当前键信息
		keyInfo, err := tx.Key().Get(key)
		if err == rdk.ErrNotFound {
			return false, nil
		} else if err != nil {
			return false, err
		}

		// 计算新的过期时间戳（毫秒）
		newETime := time.Now().Add(time.Duration(secs) * time.Second).UnixMilli()
		currentETime := int64(-1)
		if keyInfo.ETime != nil {
			currentETime = *keyInfo.ETime
		}

		// 根据 expType 决定是否设置过期时间
		shouldExpire := false
		switch expType {
		case expire:
			// 无条件设置过期时间
			shouldExpire = true
		case expireNX:
			// 仅在键没有过期时间时设置
			shouldExpire = currentETime < 0
		case expireXX:
			// 仅在键已有过期时间时设置
			shouldExpire = currentETime >= 0
		case expireGT:
			// 仅在新过期时间大于当前过期时间时设置（键必须已有过期时间）
			shouldExpire = currentETime >= 0 && newETime > currentETime
		case expireLT:
			// 仅在新过期时间小于当前过期时间时设置（键必须已有过期时间）
			shouldExpire = currentETime >= 0 && newETime < currentETime
		}

		if !shouldExpire {
			return false, nil
		}

		// 设置过期时间
		expireTime := time.Now().Add(time.Duration(secs) * time.Second)
		err = tx.Key().ExpireAt(key, expireTime)
		return err == nil, err
	})

	return newResult(val, err)
}

// Expire implements caches.KeyCommand.
func (p *Provider) Expire(ctx context.Context, key string, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	return p.expire(ctx, key, expiration, expire)
}

// ExpireNX implements caches.KeyCommand.
func (p *Provider) ExpireNX(ctx context.Context, key string, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	return p.expire(ctx, key, expiration, expireNX)
}

// ExpireXX implements caches.KeyCommand.
func (p *Provider) ExpireXX(ctx context.Context, key string, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	return p.expire(ctx, key, expiration, expireXX)
}

// ExpireGT implements caches.KeyCommand.
func (p *Provider) ExpireGT(ctx context.Context, key string, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	return p.expire(ctx, key, expiration, expireGT)
}

// ExpireLT implements caches.KeyCommand.
func (p *Provider) ExpireLT(ctx context.Context, key string, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	return p.expire(ctx, key, expiration, expireLT)
}

// ExpireAt implements caches.KeyCommand.
func (p *Provider) ExpireAt(ctx context.Context, key string, tm time.Time) caches.Result[bool] {
	key = p.prefix + key
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (ok bool, err error) {
		// 检查键是否存在
		_, err = tx.Key().Get(key)
		if err == rdk.ErrNotFound {
			return false, nil
		} else if err != nil {
			return false, err
		}

		err = tx.Key().ExpireAt(key, tm.Truncate(time.Millisecond))
		ok = err == nil
		return
	})
	return newResult(val, err)
}

func (p *Provider) expireTime(ctx context.Context, key string) (int64, error) {
	val, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		key, err := tx.Key().Get(key)

		if err == rdk.ErrNotFound {
			return -2, nil
		} else if err != nil {
			return 0, err
		} else if key.ETime == nil {
			return -1, nil
		} else {
			return *key.ETime, nil
		}
	})

	return val, err
}

// ExpireTime implements caches.KeyCommand.
func (p *Provider) ExpireTime(ctx context.Context, key string) caches.Result[time.Duration] {
	key = p.prefix + key
	exp, err := p.expireTime(ctx, key)
	if err != nil {
		return newResult(time.Duration(0), err)
	}

	if exp < 0 {
		// -2 或 -1，直接返回
		return newResult(time.Duration(exp), nil)
	}

	// exp 是 Unix 毫秒时间戳，转换为秒并返回
	expireTimeSec := exp / 1000
	return newResult(time.Duration(expireTimeSec)*time.Second, nil)
}

// PExpire implements caches.KeyCommand.
func (p *Provider) PExpire(ctx context.Context, key string, expiration time.Duration) caches.Result[bool] {
	key = p.prefix + key
	ms := formatMs(expiration)
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (ok bool, err error) {
		// 检查键是否存在
		_, err = tx.Key().Get(key)
		if err == rdk.ErrNotFound {
			return false, nil
		} else if err != nil {
			return false, err
		}

		err = tx.Key().ExpireAt(key, time.Now().Add(time.Duration(ms)*time.Millisecond))
		ok = err == nil
		return
	})
	return newResult(val, err)
}

// PExpireAt implements caches.KeyCommand.
func (p *Provider) PExpireAt(ctx context.Context, key string, tm time.Time) caches.Result[bool] {
	key = p.prefix + key
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (ok bool, err error) {
		// 检查键是否存在
		_, err = tx.Key().Get(key)
		if err == rdk.ErrNotFound {
			return false, nil
		} else if err != nil {
			return false, err
		}

		err = tx.Key().ExpireAt(key, tm)
		ok = err == nil
		return
	})
	return newResult(val, err)
}

// PExpireTime implements caches.KeyCommand.
func (p *Provider) PExpireTime(ctx context.Context, key string) caches.Result[time.Duration] {
	key = p.prefix + key
	exp, err := p.expireTime(ctx, key)
	if err != nil {
		return newResult(time.Duration(0), err)
	}

	if exp < 0 {
		// -2 或 -1，直接返回
		return newResult(time.Duration(exp), nil)
	}

	// exp 是 Unix 毫秒时间戳，直接作为毫秒返回
	return newResult(time.Duration(exp)*time.Millisecond, nil)
}

// FlushAll implements caches.KeyCommand.
func (p *Provider) FlushAll(ctx context.Context) caches.StatusResult {
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (res []byte, err error) {
		// 直接执行 tx.Key().DeleteAll().Error() 会报错： SQL logic error: cannot VACUUM from within a transaction (1)
		err = tx.Key().DeleteAll()
		if err == nil {
			res = []byte("OK")
		}
		return
	})

	return newStatusResult(val, err)
}

// Keys implements caches.KeyCommand.
func (p *Provider) Keys(ctx context.Context, pattern string) caches.Result[[]string] {
	pattern = p.prefix + pattern
	val, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (res []string, err error) {
		keys, err := tx.Key().Keys(pattern)
		if err == nil {
			res = make([]string, len(keys))
			prefixLen := len(p.prefix)

			for i, key := range keys {
				// 去除前缀，返回原始键名
				if prefixLen > 0 && len(key.Key) > prefixLen {
					res[i] = key.Key[prefixLen:]
				} else {
					res[i] = key.Key
				}
			}
		}

		return
	})

	return newResult(val, err)
}

// TTL implements caches.KeyCommand.
func (p *Provider) TTL(ctx context.Context, key string) caches.Result[time.Duration] {
	key = p.prefix + key
	exp, err := p.expireTime(ctx, key)
	if err != nil {
		return newResult(time.Duration(0), err)
	}

	if exp < 0 {
		// -2 或 -1，直接返回
		return newResult(time.Duration(exp), nil)
	}

	// 计算剩余时间（秒）
	ttl := max(time.Until(time.UnixMilli(exp)).Truncate(time.Second), 0)
	return newResult(ttl, nil)
}

// PTTL implements caches.KeyCommand.
func (p *Provider) PTTL(ctx context.Context, key string) caches.Result[time.Duration] {
	key = p.prefix + key
	exp, err := p.expireTime(ctx, key)
	if err != nil {
		return newResult(time.Duration(0), err)
	}

	if exp < 0 {
		// -2 或 -1，直接返回
		return newResult(time.Duration(exp), nil)
	}

	// 计算剩余时间（毫秒）
	ttl := max(time.Until(time.UnixMilli(exp)), 0)
	return newResult(ttl, nil)
}

// Persist implements caches.KeyCommand.
func (p *Provider) Persist(ctx context.Context, key string) caches.Result[bool] {
	key = p.prefix + key
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (bool, error) {
		// 先检查键是否存在
		keyInfo, err := tx.Key().Get(key)
		if err == rdk.ErrNotFound {
			return false, nil
		} else if err != nil {
			return false, err
		}

		// 如果键没有过期时间，返回 false
		if keyInfo.ETime == nil {
			return false, nil
		}

		// 移除过期时间
		err = tx.Key().Persist(key)
		return err == nil, err
	})
	return newResult(val, err)
}

// Rename implements caches.KeyCommand.
func (p *Provider) Rename(ctx context.Context, key string, newKey string) caches.StatusResult {
	key = p.prefix + key
	newKey = p.prefix + newKey
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]byte, error) {
		err := tx.Key().Rename(key, newKey)
		if err != nil {
			return nil, err
		}
		return []byte("OK"), nil
	})
	return newStatusResult(val, err)
}

// RenameNX implements caches.KeyCommand.
func (p *Provider) RenameNX(ctx context.Context, key string, newKey string) caches.Result[bool] {
	key = p.prefix + key
	newKey = p.prefix + newKey
	val, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (bool, error) {
		// 检查新键是否已存在
		_, err := tx.Key().Get(newKey)
		if err == nil {
			// 新键已存在，返回 false
			return false, nil
		} else if err != rdk.ErrNotFound {
			// 其他错误
			return false, err
		}

		// 新键不存在，执行重命名
		err = tx.Key().Rename(key, newKey)
		if err != nil {
			return false, err
		}
		return true, nil
	})
	return newResult(val, err)
}

// Type implements caches.KeyCommand.
func (p *Provider) Type(ctx context.Context, key string) caches.Result[string] {
	key = p.prefix + key
	val, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (string, error) {
		keyInfo, err := tx.Key().Get(key)
		if err == rdk.ErrNotFound {
			return "none", nil
		} else if err != nil {
			return "", err
		}

		// 将 TypeID 转换为字符串
		var typeStr string
		switch keyInfo.Type {
		case 1: // String
			typeStr = "string"
		case 2: // List
			typeStr = "list"
		case 3: // Set
			typeStr = "set"
		case 4: // Hash
			typeStr = "hash"
		case 5: // ZSet
			typeStr = "zset"
		default:
			typeStr = "unknown"
		}

		return typeStr, nil
	})
	return newResult(val, err)
}
