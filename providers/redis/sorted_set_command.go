package redis

import (
	"context"
	"fmt"

	rds "github.com/redis/go-redis/v9"
	"github.com/rockcookies/go-caches"
)

var _ caches.SortedSetCommand = (*Provider)(nil)

// ZAdd implements caches.SortedSetCommand.
func (p *Provider) ZAdd(ctx context.Context, key string, members ...caches.ZMember) caches.Result[int64] {
	key = p.prefix + key
	zmembers := make([]rds.Z, len(members))
	for i, m := range members {
		zmembers[i] = rds.Z{Score: m.Score, Member: m.Member}
	}
	res := p.db.ZAdd(ctx, key, zmembers...)
	res.SetErr(formatError(res.Err()))
	return res
}

// ZAddArgs implements caches.SortedSetCommand.
func (p *Provider) ZAddArgs(ctx context.Context, key string, mode string, ch bool, members ...caches.ZMember) caches.Result[int64] {
	key = p.prefix + key
	zmembers := make([]rds.Z, len(members))
	for i, m := range members {
		zmembers[i] = rds.Z{Score: m.Score, Member: m.Member}
	}

	args := rds.ZAddArgs{
		Members: zmembers,
		Ch:      ch,
	}

	switch mode {
	case "NX":
		args.NX = true
	case "XX":
		args.XX = true
	case "GT":
		args.GT = true
	case "LT":
		args.LT = true
	}

	res := p.db.ZAddArgs(ctx, key, args)
	res.SetErr(formatError(res.Err()))
	return res
}

// ZCard implements caches.SortedSetCommand.
func (p *Provider) ZCard(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.ZCard(ctx, key)
	res.SetErr(formatError(res.Err()))
	return res
}

// ZCount implements caches.SortedSetCommand.
func (p *Provider) ZCount(ctx context.Context, key string, min, max string) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.ZCount(ctx, key, min, max)
	res.SetErr(formatError(res.Err()))
	return res
}

// ZIncrBy implements caches.SortedSetCommand.
func (p *Provider) ZIncrBy(ctx context.Context, key string, increment float64, member string) caches.Result[float64] {
	key = p.prefix + key
	res := p.db.ZIncrBy(ctx, key, increment, member)
	res.SetErr(formatError(res.Err()))
	return res
}

// ZInter implements caches.SortedSetCommand.
func (p *Provider) ZInter(ctx context.Context, store caches.ZStore) caches.Result[[][]byte] {
	zstore := convertZStore(p.prefix, store)
	res := p.db.ZInter(ctx, &zstore)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// ZInterWithScores implements caches.SortedSetCommand.
func (p *Provider) ZInterWithScores(ctx context.Context, store caches.ZStore) caches.Result[[]caches.ZMember] {
	zstore := convertZStore(p.prefix, store)
	res := p.db.ZInterWithScores(ctx, &zstore)

	if res.Err() != nil {
		return newResult([]caches.ZMember(nil), formatError(res.Err()))
	}

	result := make([]caches.ZMember, len(res.Val()))
	for i, z := range res.Val() {
		result[i] = caches.ZMember{
			Member: []byte(z.Member.(string)),
			Score:  z.Score,
		}
	}

	return newResult(result, nil)
}

// ZInterStore implements caches.SortedSetCommand.
func (p *Provider) ZInterStore(ctx context.Context, destination string, store caches.ZStore) caches.Result[int64] {
	destination = p.prefix + destination
	zstore := convertZStore(p.prefix, store)
	res := p.db.ZInterStore(ctx, destination, &zstore)
	res.SetErr(formatError(res.Err()))
	return res
}

// ZRange implements caches.SortedSetCommand.
func (p *Provider) ZRange(ctx context.Context, key string, start, stop int64) caches.Result[[][]byte] {
	key = p.prefix + key
	res := p.db.ZRange(ctx, key, start, stop)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// ZRangeWithScores implements caches.SortedSetCommand.
func (p *Provider) ZRangeWithScores(ctx context.Context, key string, start, stop int64) caches.Result[[]caches.ZMember] {
	key = p.prefix + key
	res := p.db.ZRangeWithScores(ctx, key, start, stop)

	if res.Err() != nil {
		return newResult([]caches.ZMember(nil), formatError(res.Err()))
	}

	result := make([]caches.ZMember, len(res.Val()))
	for i, z := range res.Val() {
		result[i] = caches.ZMember{
			Member: []byte(z.Member.(string)),
			Score:  z.Score,
		}
	}

	return newResult(result, nil)
}

// ZRangeArgs implements caches.SortedSetCommand.
func (p *Provider) ZRangeArgs(ctx context.Context, key string, args caches.ZRangeArgs) caches.Result[[][]byte] {
	key = p.prefix + key
	zargs := convertZRangeArgs(args)
	res := p.db.ZRangeArgs(ctx, zargs)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// ZRangeArgsWithScores implements caches.SortedSetCommand.
func (p *Provider) ZRangeArgsWithScores(ctx context.Context, key string, args caches.ZRangeArgs) caches.Result[[]caches.ZMember] {
	key = p.prefix + key
	zargs := convertZRangeArgs(args)
	res := p.db.ZRangeArgsWithScores(ctx, zargs)

	if res.Err() != nil {
		return newResult([]caches.ZMember(nil), formatError(res.Err()))
	}

	result := make([]caches.ZMember, len(res.Val()))
	for i, z := range res.Val() {
		result[i] = caches.ZMember{
			Member: []byte(z.Member.(string)),
			Score:  z.Score,
		}
	}

	return newResult(result, nil)
}

// ZRangeByScore implements caches.SortedSetCommand.
func (p *Provider) ZRangeByScore(ctx context.Context, key string, min, max string) caches.Result[[][]byte] {
	key = p.prefix + key
	res := p.db.ZRangeByScore(ctx, key, &rds.ZRangeBy{Min: min, Max: max})

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// ZRangeByScoreWithScores implements caches.SortedSetCommand.
func (p *Provider) ZRangeByScoreWithScores(ctx context.Context, key string, min, max string) caches.Result[[]caches.ZMember] {
	key = p.prefix + key
	res := p.db.ZRangeByScoreWithScores(ctx, key, &rds.ZRangeBy{Min: min, Max: max})

	if res.Err() != nil {
		return newResult([]caches.ZMember(nil), formatError(res.Err()))
	}

	result := make([]caches.ZMember, len(res.Val()))
	for i, z := range res.Val() {
		result[i] = caches.ZMember{
			Member: []byte(z.Member.(string)),
			Score:  z.Score,
		}
	}

	return newResult(result, nil)
}

// ZRank implements caches.SortedSetCommand.
func (p *Provider) ZRank(ctx context.Context, key string, member string) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.ZRank(ctx, key, member)
	res.SetErr(formatError(res.Err()))
	return res
}

// ZRankWithScore implements caches.SortedSetCommand.
func (p *Provider) ZRankWithScore(ctx context.Context, key string, member string) caches.Result[caches.ZRankScore] {
	key = p.prefix + key
	res := p.db.ZRankWithScore(ctx, key, member)

	if res.Err() != nil {
		return newResult(caches.ZRankScore{}, formatError(res.Err()))
	}

	return newResult(caches.ZRankScore{
		Rank:  res.Val().Rank,
		Score: res.Val().Score,
	}, nil)
}

// ZRem implements caches.SortedSetCommand.
func (p *Provider) ZRem(ctx context.Context, key string, members ...any) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.ZRem(ctx, key, members...)
	res.SetErr(formatError(res.Err()))
	return res
}

// ZRemRangeByRank implements caches.SortedSetCommand.
func (p *Provider) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.ZRemRangeByRank(ctx, key, start, stop)
	res.SetErr(formatError(res.Err()))
	return res
}

// ZRemRangeByScore implements caches.SortedSetCommand.
func (p *Provider) ZRemRangeByScore(ctx context.Context, key string, min, max string) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.ZRemRangeByScore(ctx, key, min, max)
	res.SetErr(formatError(res.Err()))
	return res
}

// ZRevRange implements caches.SortedSetCommand.
func (p *Provider) ZRevRange(ctx context.Context, key string, start, stop int64) caches.Result[[][]byte] {
	key = p.prefix + key
	res := p.db.ZRevRange(ctx, key, start, stop)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// ZRevRangeWithScores implements caches.SortedSetCommand.
func (p *Provider) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) caches.Result[[]caches.ZMember] {
	key = p.prefix + key
	res := p.db.ZRevRangeWithScores(ctx, key, start, stop)

	if res.Err() != nil {
		return newResult([]caches.ZMember(nil), formatError(res.Err()))
	}

	result := make([]caches.ZMember, len(res.Val()))
	for i, z := range res.Val() {
		result[i] = caches.ZMember{
			Member: []byte(z.Member.(string)),
			Score:  z.Score,
		}
	}

	return newResult(result, nil)
}

// ZRevRangeByScore implements caches.SortedSetCommand.
func (p *Provider) ZRevRangeByScore(ctx context.Context, key string, max, min string) caches.Result[[][]byte] {
	key = p.prefix + key
	res := p.db.ZRevRangeByScore(ctx, key, &rds.ZRangeBy{Min: min, Max: max})

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// ZRevRangeByScoreWithScores implements caches.SortedSetCommand.
func (p *Provider) ZRevRangeByScoreWithScores(ctx context.Context, key string, max, min string) caches.Result[[]caches.ZMember] {
	key = p.prefix + key
	res := p.db.ZRevRangeByScoreWithScores(ctx, key, &rds.ZRangeBy{Min: min, Max: max})

	if res.Err() != nil {
		return newResult([]caches.ZMember(nil), formatError(res.Err()))
	}

	result := make([]caches.ZMember, len(res.Val()))
	for i, z := range res.Val() {
		result[i] = caches.ZMember{
			Member: []byte(z.Member.(string)),
			Score:  z.Score,
		}
	}

	return newResult(result, nil)
}

// ZRevRank implements caches.SortedSetCommand.
func (p *Provider) ZRevRank(ctx context.Context, key string, member string) caches.Result[int64] {
	key = p.prefix + key
	res := p.db.ZRevRank(ctx, key, member)
	res.SetErr(formatError(res.Err()))
	return res
}

// ZRevRankWithScore implements caches.SortedSetCommand.
func (p *Provider) ZRevRankWithScore(ctx context.Context, key string, member string) caches.Result[caches.ZRankScore] {
	key = p.prefix + key
	res := p.db.ZRevRankWithScore(ctx, key, member)

	if res.Err() != nil {
		return newResult(caches.ZRankScore{}, formatError(res.Err()))
	}

	return newResult(caches.ZRankScore{
		Rank:  res.Val().Rank,
		Score: res.Val().Score,
	}, nil)
}

// ZScan implements caches.SortedSetCommand.
func (p *Provider) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) caches.Result[caches.ZScanResult] {
	key = p.prefix + key
	res := p.db.ZScan(ctx, key, cursor, match, count)

	if res.Err() != nil {
		return newResult(caches.ZScanResult{}, formatError(res.Err()))
	}

	keys, newCursor := res.Val()
	members := make([]caches.ZMember, 0, len(keys)/2)

	// ZScan returns member-score pairs as a flat slice
	for i := 0; i < len(keys); i += 2 {
		if i+1 < len(keys) {
			var score float64
			if _, err := fmt.Sscanf(keys[i+1], "%f", &score); err == nil {
				members = append(members, caches.ZMember{
					Member: []byte(keys[i]),
					Score:  score,
				})
			}
		}
	}

	return newResult(caches.ZScanResult{
		Cursor:  newCursor,
		Members: members,
	}, nil)
}

// ZScore implements caches.SortedSetCommand.
func (p *Provider) ZScore(ctx context.Context, key string, member string) caches.Result[float64] {
	key = p.prefix + key
	res := p.db.ZScore(ctx, key, member)
	res.SetErr(formatError(res.Err()))
	return res
}

// ZUnion implements caches.SortedSetCommand.
func (p *Provider) ZUnion(ctx context.Context, store caches.ZStore) caches.Result[[][]byte] {
	zstore := convertZStore(p.prefix, store)
	res := p.db.ZUnion(ctx, zstore)

	if res.Err() != nil {
		return newResult([][]byte(nil), formatError(res.Err()))
	}

	result := make([][]byte, len(res.Val()))
	for i, value := range res.Val() {
		result[i] = []byte(value)
	}

	return newResult(result, nil)
}

// ZUnionWithScores implements caches.SortedSetCommand.
func (p *Provider) ZUnionWithScores(ctx context.Context, store caches.ZStore) caches.Result[[]caches.ZMember] {
	zstore := convertZStore(p.prefix, store)
	res := p.db.ZUnionWithScores(ctx, zstore)

	if res.Err() != nil {
		return newResult([]caches.ZMember(nil), formatError(res.Err()))
	}

	result := make([]caches.ZMember, len(res.Val()))
	for i, z := range res.Val() {
		result[i] = caches.ZMember{
			Member: []byte(z.Member.(string)),
			Score:  z.Score,
		}
	}

	return newResult(result, nil)
}

// ZUnionStore implements caches.SortedSetCommand.
func (p *Provider) ZUnionStore(ctx context.Context, destination string, store caches.ZStore) caches.Result[int64] {
	destination = p.prefix + destination
	zstore := convertZStore(p.prefix, store)
	res := p.db.ZUnionStore(ctx, destination, &zstore)
	res.SetErr(formatError(res.Err()))
	return res
}

// Helper functions

func convertZStore(prefix string, store caches.ZStore) rds.ZStore {
	keys := prefixKeys(prefix, store.Keys)
	return rds.ZStore{
		Keys:      keys,
		Weights:   store.Weights,
		Aggregate: store.Aggregate,
	}
}

func convertZRangeArgs(args caches.ZRangeArgs) rds.ZRangeArgs {
	return rds.ZRangeArgs{
		ByScore: args.ByScore,
		ByLex:   args.ByLex,
		Rev:     args.Rev,
		Offset:  args.Offset,
		Count:   args.Count,
	}
}
