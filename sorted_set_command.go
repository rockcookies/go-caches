package caches

import "context"

// ZMember represents a member with its score in a sorted set.
type ZMember struct {
	// Member is the member value.
	Member []byte
	// Score is the score associated with the member.
	Score float64
}

// ZRangeArgs provides arguments for range operations on sorted sets.
type ZRangeArgs struct {
	// ByScore determines whether the range is by score or by rank (index).
	ByScore bool
	// ByLex determines whether the range is by lexicographical order.
	ByLex bool
	// Rev determines whether to return members in reverse order.
	Rev bool
	// Offset is the number of elements to skip (used with Limit).
	Offset int64
	// Count is the maximum number of elements to return (used with Limit, -1 for no limit).
	Count int64
}

// ZStore specifies aggregation method for set operations.
type ZStore struct {
	// Keys are the keys of the sorted sets to operate on.
	Keys []string
	// Weights are the multiplication factors for scores from each sorted set.
	// If not specified, defaults to 1 for each set.
	Weights []float64
	// Aggregate specifies how to aggregate scores: "SUM", "MIN", or "MAX".
	// Defaults to "SUM".
	Aggregate string
}

// SortedSetCommand defines operations for Redis sorted set data structure.
// Sorted sets (zsets) are collections of unique strings ordered by each string's associated score.
type SortedSetCommand interface {
	// ZAdd adds one or more members with scores to a sorted set.
	// If a member already exists, its score is updated.
	// Returns the number of members added (not including members whose score was updated).
	ZAdd(ctx context.Context, key string, members ...ZMember) Result[int64]

	// ZAddArgs adds members with additional options.
	// mode can be "NX" (only add new members), "XX" (only update existing members),
	// "GT" (only update if new score is greater), "LT" (only update if new score is less).
	// ch returns the number of members changed (added or updated) instead of just added.
	ZAddArgs(ctx context.Context, key string, mode string, ch bool, members ...ZMember) Result[int64]

	// ZCard returns the number of members in a sorted set.
	// Returns 0 if the key does not exist.
	ZCard(ctx context.Context, key string) Result[int64]

	// ZCount returns the number of members in a sorted set within a range of scores.
	// min and max can be inclusive or exclusive (use "(" prefix for exclusive).
	ZCount(ctx context.Context, key string, min, max string) Result[int64]

	// ZIncrBy increments the score of a member in a sorted set by increment.
	// If the member does not exist, it is added with increment as its score.
	// Returns the new score of the member.
	ZIncrBy(ctx context.Context, key string, increment float64, member string) Result[float64]

	// ZInter returns the intersection of multiple sorted sets.
	// The intersection contains members that exist in all given sets.
	ZInter(ctx context.Context, store ZStore) Result[[][]byte]

	// ZInterWithScores returns the intersection with scores.
	ZInterWithScores(ctx context.Context, store ZStore) Result[[]ZMember]

	// ZInterStore stores the intersection of multiple sorted sets in a destination key.
	// Returns the number of members in the resulting sorted set.
	ZInterStore(ctx context.Context, destination string, store ZStore) Result[int64]

	// ZRange returns members in a sorted set within a range of indexes.
	// start and stop are zero-based indexes (can be negative to indicate offsets from the end).
	// Returns members in ascending order by score.
	ZRange(ctx context.Context, key string, start, stop int64) Result[[][]byte]

	// ZRangeWithScores returns members with their scores.
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) Result[[]ZMember]

	// ZRangeArgs returns members based on custom range arguments.
	ZRangeArgs(ctx context.Context, key string, args ZRangeArgs) Result[[][]byte]

	// ZRangeArgsWithScores returns members with scores based on custom range arguments.
	ZRangeArgsWithScores(ctx context.Context, key string, args ZRangeArgs) Result[[]ZMember]

	// ZRangeByScore returns members in a sorted set within a range of scores.
	// min and max can be inclusive or exclusive (use "(" prefix for exclusive).
	ZRangeByScore(ctx context.Context, key string, min, max string) Result[[][]byte]

	// ZRangeByScoreWithScores returns members with scores within a range of scores.
	ZRangeByScoreWithScores(ctx context.Context, key string, min, max string) Result[[]ZMember]

	// ZRank returns the rank (index) of a member in a sorted set ordered by ascending scores.
	// Ranks start at 0 for the member with the lowest score.
	// Returns -1 if the member does not exist.
	ZRank(ctx context.Context, key string, member string) Result[int64]

	// ZRankWithScore returns the rank and score of a member.
	ZRankWithScore(ctx context.Context, key string, member string) Result[ZRankScore]

	// ZRem removes one or more members from a sorted set.
	// Returns the number of members removed.
	ZRem(ctx context.Context, key string, members ...any) Result[int64]

	// ZRemRangeByRank removes members in a sorted set within a range of indexes.
	// start and stop are zero-based indexes (can be negative).
	// Returns the number of members removed.
	ZRemRangeByRank(ctx context.Context, key string, start, stop int64) Result[int64]

	// ZRemRangeByScore removes members in a sorted set within a range of scores.
	// min and max can be inclusive or exclusive (use "(" prefix for exclusive).
	// Returns the number of members removed.
	ZRemRangeByScore(ctx context.Context, key string, min, max string) Result[int64]

	// ZRevRange returns members in a sorted set within a range of indexes in reverse order.
	// Members are ordered by descending scores.
	ZRevRange(ctx context.Context, key string, start, stop int64) Result[[][]byte]

	// ZRevRangeWithScores returns members with scores in reverse order.
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) Result[[]ZMember]

	// ZRevRangeByScore returns members in a sorted set within a range of scores in reverse order.
	ZRevRangeByScore(ctx context.Context, key string, max, min string) Result[[][]byte]

	// ZRevRangeByScoreWithScores returns members with scores within a range of scores in reverse order.
	ZRevRangeByScoreWithScores(ctx context.Context, key string, max, min string) Result[[]ZMember]

	// ZRevRank returns the rank (index) of a member in a sorted set ordered by descending scores.
	// Ranks start at 0 for the member with the highest score.
	// Returns -1 if the member does not exist.
	ZRevRank(ctx context.Context, key string, member string) Result[int64]

	// ZRevRankWithScore returns the rank and score of a member in reverse order.
	ZRevRankWithScore(ctx context.Context, key string, member string) Result[ZRankScore]

	// ZScan iterates over members and scores of a sorted set.
	// cursor is the cursor to start iteration from (0 to start).
	// match is a glob-style pattern to filter members (empty string for no filter).
	// count is a hint for how many members to return per iteration.
	// Returns the next cursor and a slice of members with scores.
	ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) Result[ZScanResult]

	// ZScore returns the score of a member in a sorted set.
	// Returns an error if the member does not exist.
	ZScore(ctx context.Context, key string, member string) Result[float64]

	// ZUnion returns the union of multiple sorted sets.
	// The union contains all members that exist in at least one of the given sets.
	ZUnion(ctx context.Context, store ZStore) Result[[][]byte]

	// ZUnionWithScores returns the union with scores.
	ZUnionWithScores(ctx context.Context, store ZStore) Result[[]ZMember]

	// ZUnionStore stores the union of multiple sorted sets in a destination key.
	// Returns the number of members in the resulting sorted set.
	ZUnionStore(ctx context.Context, destination string, store ZStore) Result[int64]
}

// ZRankScore represents the rank and score of a member.
type ZRankScore struct {
	// Rank is the rank (index) of the member.
	Rank int64
	// Score is the score of the member.
	Score float64
}

// ZScanResult represents the result of a zscan operation.
type ZScanResult struct {
	// Cursor is the cursor to use in the next scan call.
	// A cursor value of 0 indicates the iteration is complete.
	Cursor uint64
	// Members contains the scanned members with their scores.
	Members []ZMember
}
