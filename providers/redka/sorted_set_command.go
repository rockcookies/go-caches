package redka

import (
	"context"
	"fmt"
	"strings"

	rdk "github.com/nalgeon/redka"
	"github.com/rockcookies/go-caches"
)

var _ caches.SortedSetCommand = (*Provider)(nil)

// ZAdd implements caches.SortedSetCommand.
func (p *Provider) ZAdd(ctx context.Context, key string, members ...caches.ZMember) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		items := make(map[any]float64, len(members))
		for _, m := range members {
			items[m.Member] = m.Score
		}
		count, e := tx.ZSet().AddMany(key, items)
		return int64(count), e
	})
	return newResult(n, err)
}

// ZAddArgs implements caches.SortedSetCommand.
func (p *Provider) ZAddArgs(ctx context.Context, key string, mode string, ch bool, members ...caches.ZMember) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		// Redka doesn't support mode flags directly, handle manually
		count := int64(0)
		for _, m := range members {
			exists := false
			if mode == "NX" || mode == "XX" {
				_, e := tx.ZSet().GetScore(key, m.Member)
				exists = (e == nil)
			}

			if mode == "NX" && exists {
				continue
			}
			if mode == "XX" && !exists {
				continue
			}

			created, e := tx.ZSet().Add(key, m.Member, m.Score)
			if e != nil {
				return 0, e
			}
			if ch || created {
				count++
			}
		}
		return count, nil
	})
	return newResult(n, err)
}

// ZCard implements caches.SortedSetCommand.
func (p *Provider) ZCard(ctx context.Context, key string) caches.Result[int64] {
	key = p.prefix + key
	n, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		count, e := tx.ZSet().Len(key)
		return int64(count), e
	})
	return newResult(n, err)
}

// ZCount implements caches.SortedSetCommand.
func (p *Provider) ZCount(ctx context.Context, key string, min, max string) caches.Result[int64] {
	key = p.prefix + key
	n, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		minScore, maxScore, e := parseScoreRange(min, max)
		if e != nil {
			return 0, e
		}
		count, e := tx.ZSet().Count(key, minScore, maxScore)
		return int64(count), e
	})
	return newResult(n, err)
}

// ZIncrBy implements caches.SortedSetCommand.
func (p *Provider) ZIncrBy(ctx context.Context, key string, increment float64, member string) caches.Result[float64] {
	key = p.prefix + key
	score, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (float64, error) {
		return tx.ZSet().Incr(key, member, increment)
	})
	return newResult(score, err)
}

// ZInter implements caches.SortedSetCommand.
func (p *Provider) ZInter(ctx context.Context, store caches.ZStore) caches.Result[[][]byte] {
	keys := prefixKeys(p.prefix, store.Keys)
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		cmd := tx.ZSet().InterWith(keys...)
		switch strings.ToUpper(store.Aggregate) {
		case "MIN":
			cmd = cmd.Min()
		case "MAX":
			cmd = cmd.Max()
		default:
			cmd = cmd.Sum()
		}

		items, e := cmd.Run()
		if e != nil {
			return nil, e
		}

		result := make([][]byte, len(items))
		for i, item := range items {
			result[i] = item.Elem.Bytes()
		}
		return result, nil
	})
	return newResult(vals, err)
}

// ZInterWithScores implements caches.SortedSetCommand.
func (p *Provider) ZInterWithScores(ctx context.Context, store caches.ZStore) caches.Result[[]caches.ZMember] {
	keys := prefixKeys(p.prefix, store.Keys)
	members, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]caches.ZMember, error) {
		cmd := tx.ZSet().InterWith(keys...)
		switch strings.ToUpper(store.Aggregate) {
		case "MIN":
			cmd = cmd.Min()
		case "MAX":
			cmd = cmd.Max()
		default:
			cmd = cmd.Sum()
		}

		items, e := cmd.Run()
		if e != nil {
			return nil, e
		}

		result := make([]caches.ZMember, len(items))
		for i, item := range items {
			result[i] = caches.ZMember{
				Member: item.Elem.Bytes(),
				Score:  item.Score,
			}
		}
		return result, nil
	})
	return newResult(members, err)
}

// ZInterStore implements caches.SortedSetCommand.
func (p *Provider) ZInterStore(ctx context.Context, destination string, store caches.ZStore) caches.Result[int64] {
	destination = p.prefix + destination
	keys := prefixKeys(p.prefix, store.Keys)
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		cmd := tx.ZSet().InterWith(keys...)
		switch strings.ToUpper(store.Aggregate) {
		case "MIN":
			cmd = cmd.Min()
		case "MAX":
			cmd = cmd.Max()
		default:
			cmd = cmd.Sum()
		}
		cmd = cmd.Dest(destination)
		count, e := cmd.Store()
		return int64(count), e
	})
	return newResult(n, err)
}

// ZRange implements caches.SortedSetCommand.
func (p *Provider) ZRange(ctx context.Context, key string, start, stop int64) caches.Result[[][]byte] {
	key = p.prefix + key
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		items, e := tx.ZSet().Range(key, int(start), int(stop))
		if e != nil {
			return nil, e
		}
		result := make([][]byte, len(items))
		for i, item := range items {
			result[i] = item.Elem.Bytes()
		}
		return result, nil
	})
	return newResult(vals, err)
}

// ZRangeWithScores implements caches.SortedSetCommand.
func (p *Provider) ZRangeWithScores(ctx context.Context, key string, start, stop int64) caches.Result[[]caches.ZMember] {
	key = p.prefix + key
	members, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]caches.ZMember, error) {
		items, e := tx.ZSet().Range(key, int(start), int(stop))
		if e != nil {
			return nil, e
		}
		result := make([]caches.ZMember, len(items))
		for i, item := range items {
			result[i] = caches.ZMember{
				Member: item.Elem.Bytes(),
				Score:  item.Score,
			}
		}
		return result, nil
	})
	return newResult(members, err)
}

// ZRangeArgs implements caches.SortedSetCommand.
func (p *Provider) ZRangeArgs(ctx context.Context, key string, args caches.ZRangeArgs) caches.Result[[][]byte] {
	key = p.prefix + key
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		cmd := tx.ZSet().RangeWith(key)

		if args.ByScore {
			minScore, maxScore, e := parseScoreRangeAny(args.Start, args.Stop)
			if e != nil {
				return nil, e
			}
			cmd = cmd.ByScore(minScore, maxScore)
		} else {
			start, stop := parseIndexRange(args.Start, args.Stop)
			cmd = cmd.ByRank(start, stop)
		}

		if args.Rev {
			cmd = cmd.Desc()
		}
		if args.Offset > 0 || args.Count > 0 {
			cmd = cmd.Offset(int(args.Offset)).Count(int(args.Count))
		}

		items, e := cmd.Run()
		if e != nil {
			return nil, e
		}

		result := make([][]byte, len(items))
		for i, item := range items {
			result[i] = item.Elem.Bytes()
		}
		return result, nil
	})
	return newResult(vals, err)
}

// ZRangeArgsWithScores implements caches.SortedSetCommand.
func (p *Provider) ZRangeArgsWithScores(ctx context.Context, key string, args caches.ZRangeArgs) caches.Result[[]caches.ZMember] {
	key = p.prefix + key
	members, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]caches.ZMember, error) {
		cmd := tx.ZSet().RangeWith(key)

		if args.ByScore {
			minScore, maxScore, e := parseScoreRangeAny(args.Start, args.Stop)
			if e != nil {
				return nil, e
			}
			cmd = cmd.ByScore(minScore, maxScore)
		} else {
			start, stop := parseIndexRange(args.Start, args.Stop)
			cmd = cmd.ByRank(start, stop)
		}

		if args.Rev {
			cmd = cmd.Desc()
		}
		if args.Offset > 0 || args.Count > 0 {
			cmd = cmd.Offset(int(args.Offset)).Count(int(args.Count))
		}

		items, e := cmd.Run()
		if e != nil {
			return nil, e
		}

		result := make([]caches.ZMember, len(items))
		for i, item := range items {
			result[i] = caches.ZMember{
				Member: item.Elem.Bytes(),
				Score:  item.Score,
			}
		}
		return result, nil
	})
	return newResult(members, err)
}

// ZRangeByScore implements caches.SortedSetCommand.
func (p *Provider) ZRangeByScore(ctx context.Context, key string, min, max string) caches.Result[[][]byte] {
	key = p.prefix + key
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		minScore, maxScore, e := parseScoreRange(min, max)
		if e != nil {
			return nil, e
		}

		items, e := tx.ZSet().RangeWith(key).ByScore(minScore, maxScore).Run()
		if e != nil {
			return nil, e
		}

		result := make([][]byte, len(items))
		for i, item := range items {
			result[i] = item.Elem.Bytes()
		}
		return result, nil
	})
	return newResult(vals, err)
}

// ZRangeByScoreWithScores implements caches.SortedSetCommand.
func (p *Provider) ZRangeByScoreWithScores(ctx context.Context, key string, min, max string) caches.Result[[]caches.ZMember] {
	key = p.prefix + key
	members, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]caches.ZMember, error) {
		minScore, maxScore, e := parseScoreRange(min, max)
		if e != nil {
			return nil, e
		}

		items, e := tx.ZSet().RangeWith(key).ByScore(minScore, maxScore).Run()
		if e != nil {
			return nil, e
		}

		result := make([]caches.ZMember, len(items))
		for i, item := range items {
			result[i] = caches.ZMember{
				Member: item.Elem.Bytes(),
				Score:  item.Score,
			}
		}
		return result, nil
	})
	return newResult(members, err)
}

// ZRank implements caches.SortedSetCommand.
func (p *Provider) ZRank(ctx context.Context, key string, member string) caches.Result[int64] {
	key = p.prefix + key
	rank, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		r, _, e := tx.ZSet().GetRank(key, member)
		return int64(r), e
	})
	return newResult(rank, err)
}

// ZRankWithScore implements caches.SortedSetCommand.
func (p *Provider) ZRankWithScore(ctx context.Context, key string, member string) caches.Result[caches.ZRankScore] {
	key = p.prefix + key
	rankScore, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (caches.ZRankScore, error) {
		rank, score, e := tx.ZSet().GetRank(key, member)
		if e != nil {
			return caches.ZRankScore{}, e
		}
		return caches.ZRankScore{
			Rank:  int64(rank),
			Score: score,
		}, nil
	})
	return newResult(rankScore, err)
}

// ZRem implements caches.SortedSetCommand.
func (p *Provider) ZRem(ctx context.Context, key string, members ...any) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		count, e := tx.ZSet().Delete(key, members...)
		return int64(count), e
	})
	return newResult(n, err)
}

// ZRemRangeByRank implements caches.SortedSetCommand.
func (p *Provider) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		count, e := tx.ZSet().DeleteWith(key).ByRank(int(start), int(stop)).Run()
		return int64(count), e
	})
	return newResult(n, err)
}

// ZRemRangeByScore implements caches.SortedSetCommand.
func (p *Provider) ZRemRangeByScore(ctx context.Context, key string, min, max string) caches.Result[int64] {
	key = p.prefix + key
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		minScore, maxScore, e := parseScoreRange(min, max)
		if e != nil {
			return 0, e
		}
		count, e := tx.ZSet().DeleteWith(key).ByScore(minScore, maxScore).Run()
		return int64(count), e
	})
	return newResult(n, err)
}

// ZRevRange implements caches.SortedSetCommand.
func (p *Provider) ZRevRange(ctx context.Context, key string, start, stop int64) caches.Result[[][]byte] {
	key = p.prefix + key
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		items, e := tx.ZSet().RangeWith(key).ByRank(int(start), int(stop)).Desc().Run()
		if e != nil {
			return nil, e
		}
		result := make([][]byte, len(items))
		for i, item := range items {
			result[i] = item.Elem.Bytes()
		}
		return result, nil
	})
	return newResult(vals, err)
}

// ZRevRangeWithScores implements caches.SortedSetCommand.
func (p *Provider) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) caches.Result[[]caches.ZMember] {
	key = p.prefix + key
	members, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]caches.ZMember, error) {
		items, e := tx.ZSet().RangeWith(key).ByRank(int(start), int(stop)).Desc().Run()
		if e != nil {
			return nil, e
		}
		result := make([]caches.ZMember, len(items))
		for i, item := range items {
			result[i] = caches.ZMember{
				Member: item.Elem.Bytes(),
				Score:  item.Score,
			}
		}
		return result, nil
	})
	return newResult(members, err)
}

// ZRevRangeByScore implements caches.SortedSetCommand.
func (p *Provider) ZRevRangeByScore(ctx context.Context, key string, max, min string) caches.Result[[][]byte] {
	key = p.prefix + key
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		minScore, maxScore, e := parseScoreRange(min, max)
		if e != nil {
			return nil, e
		}

		items, e := tx.ZSet().RangeWith(key).ByScore(minScore, maxScore).Desc().Run()
		if e != nil {
			return nil, e
		}

		result := make([][]byte, len(items))
		for i, item := range items {
			result[i] = item.Elem.Bytes()
		}
		return result, nil
	})
	return newResult(vals, err)
}

// ZRevRangeByScoreWithScores implements caches.SortedSetCommand.
func (p *Provider) ZRevRangeByScoreWithScores(ctx context.Context, key string, max, min string) caches.Result[[]caches.ZMember] {
	key = p.prefix + key
	members, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]caches.ZMember, error) {
		minScore, maxScore, e := parseScoreRange(min, max)
		if e != nil {
			return nil, e
		}

		items, e := tx.ZSet().RangeWith(key).ByScore(minScore, maxScore).Desc().Run()
		if e != nil {
			return nil, e
		}

		result := make([]caches.ZMember, len(items))
		for i, item := range items {
			result[i] = caches.ZMember{
				Member: item.Elem.Bytes(),
				Score:  item.Score,
			}
		}
		return result, nil
	})
	return newResult(members, err)
}

// ZRevRank implements caches.SortedSetCommand.
func (p *Provider) ZRevRank(ctx context.Context, key string, member string) caches.Result[int64] {
	key = p.prefix + key
	rank, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		r, _, e := tx.ZSet().GetRankRev(key, member)
		return int64(r), e
	})
	return newResult(rank, err)
}

// ZRevRankWithScore implements caches.SortedSetCommand.
func (p *Provider) ZRevRankWithScore(ctx context.Context, key string, member string) caches.Result[caches.ZRankScore] {
	key = p.prefix + key
	rankScore, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (caches.ZRankScore, error) {
		rank, score, e := tx.ZSet().GetRankRev(key, member)
		if e != nil {
			return caches.ZRankScore{}, e
		}
		return caches.ZRankScore{
			Rank:  int64(rank),
			Score: score,
		}, nil
	})
	return newResult(rankScore, err)
}

// ZScan implements caches.SortedSetCommand.
func (p *Provider) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) caches.Result[caches.ZScanResult] {
	key = p.prefix + key
	result, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (caches.ZScanResult, error) {
		scanRes, e := tx.ZSet().Scan(key, int(cursor), match, int(count))
		if e != nil {
			return caches.ZScanResult{}, e
		}

		members := make([]caches.ZMember, len(scanRes.Items))
		for i, item := range scanRes.Items {
			members[i] = caches.ZMember{
				Member: item.Elem.Bytes(),
				Score:  item.Score,
			}
		}

		return caches.ZScanResult{
			Cursor:  uint64(scanRes.Cursor),
			Members: members,
		}, nil
	})
	return newResult(result, err)
}

// ZScore implements caches.SortedSetCommand.
func (p *Provider) ZScore(ctx context.Context, key string, member string) caches.Result[float64] {
	key = p.prefix + key
	score, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) (float64, error) {
		return tx.ZSet().GetScore(key, member)
	})
	return newResult(score, err)
}

// ZUnion implements caches.SortedSetCommand.
func (p *Provider) ZUnion(ctx context.Context, store caches.ZStore) caches.Result[[][]byte] {
	keys := prefixKeys(p.prefix, store.Keys)
	vals, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([][]byte, error) {
		cmd := tx.ZSet().UnionWith(keys...)
		switch strings.ToUpper(store.Aggregate) {
		case "MIN":
			cmd = cmd.Min()
		case "MAX":
			cmd = cmd.Max()
		default:
			cmd = cmd.Sum()
		}

		items, e := cmd.Run()
		if e != nil {
			return nil, e
		}

		result := make([][]byte, len(items))
		for i, item := range items {
			result[i] = item.Elem.Bytes()
		}
		return result, nil
	})
	return newResult(vals, err)
}

// ZUnionWithScores implements caches.SortedSetCommand.
func (p *Provider) ZUnionWithScores(ctx context.Context, store caches.ZStore) caches.Result[[]caches.ZMember] {
	keys := prefixKeys(p.prefix, store.Keys)
	members, err := viewAndReturn(ctx, p.db, func(tx *rdk.Tx) ([]caches.ZMember, error) {
		cmd := tx.ZSet().UnionWith(keys...)
		switch strings.ToUpper(store.Aggregate) {
		case "MIN":
			cmd = cmd.Min()
		case "MAX":
			cmd = cmd.Max()
		default:
			cmd = cmd.Sum()
		}

		items, e := cmd.Run()
		if e != nil {
			return nil, e
		}

		result := make([]caches.ZMember, len(items))
		for i, item := range items {
			result[i] = caches.ZMember{
				Member: item.Elem.Bytes(),
				Score:  item.Score,
			}
		}
		return result, nil
	})
	return newResult(members, err)
}

// ZUnionStore implements caches.SortedSetCommand.
func (p *Provider) ZUnionStore(ctx context.Context, destination string, store caches.ZStore) caches.Result[int64] {
	destination = p.prefix + destination
	keys := prefixKeys(p.prefix, store.Keys)
	n, err := updateAndReturn(ctx, p.db, func(tx *rdk.Tx) (int64, error) {
		cmd := tx.ZSet().UnionWith(keys...)
		switch strings.ToUpper(store.Aggregate) {
		case "MIN":
			cmd = cmd.Min()
		case "MAX":
			cmd = cmd.Max()
		default:
			cmd = cmd.Sum()
		}
		cmd = cmd.Dest(destination)
		count, e := cmd.Store()
		return int64(count), e
	})
	return newResult(n, err)
}

// Helper functions

func parseScoreRange(min, max string) (float64, float64, error) {
	return parseScoreRangeAny(min, max)
}

func parseScoreRangeAny(start, stop any) (float64, float64, error) {
	minScore, err := parseScore(start)
	if err != nil {
		return 0, 0, err
	}
	maxScore, err := parseScore(stop)
	if err != nil {
		return 0, 0, err
	}
	return minScore, maxScore, nil
}

func parseScore(val any) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return parseScoreString(v)
	default:
		return 0, fmt.Errorf("invalid score type")
	}
}

func parseScoreString(s string) (float64, error) {
	if s == "-inf" {
		return -1e308, nil
	}
	if s == "+inf" || s == "inf" {
		return 1e308, nil
	}

	// Handle exclusive markers
	if len(s) > 0 && s[0] == '(' {
		s = s[1:]
	}

	var score float64
	_, err := fmt.Sscanf(s, "%f", &score)
	return score, err
}

func parseIndexRange(start, stop any) (int, int) {
	startIdx := 0
	stopIdx := -1

	switch v := start.(type) {
	case int:
		startIdx = v
	case int64:
		startIdx = int(v)
	case string:
		fmt.Sscanf(v, "%d", &startIdx)
	}

	switch v := stop.(type) {
	case int:
		stopIdx = v
	case int64:
		stopIdx = int(v)
	case string:
		fmt.Sscanf(v, "%d", &stopIdx)
	}

	return startIdx, stopIdx
}
