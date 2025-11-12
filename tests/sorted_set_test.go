package tests

import (
	"context"
	"testing"

	"github.com/rockcookies/go-caches"
	"github.com/stretchr/testify/require"
)

// SortedSetCommandProvider defines the interface for testing SortedSetCommand implementations
type SortedSetCommandProvider interface {
	GetSortedSetCommand() caches.SortedSetCommand
	GetContext() context.Context
}

// RunSortedSetCommandTests runs all SortedSetCommand tests
func RunSortedSetCommandTests(t *testing.T, provider SortedSetCommandProvider) {
	t.Run("ZAdd_and_ZRange", func(t *testing.T) {
		testZAddAndZRange(t, provider)
	})
	t.Run("ZAdd_UpdateScore", func(t *testing.T) {
		testZAddUpdateScore(t, provider)
	})
	t.Run("ZAddArgs_NX", func(t *testing.T) {
		testZAddArgsNX(t, provider)
	})
	t.Run("ZAddArgs_XX", func(t *testing.T) {
		testZAddArgsXX(t, provider)
	})
	t.Run("ZAddArgs_GT", func(t *testing.T) {
		testZAddArgsGT(t, provider)
	})
	t.Run("ZAddArgs_LT", func(t *testing.T) {
		testZAddArgsLT(t, provider)
	})
	t.Run("ZAddArgs_CH", func(t *testing.T) {
		testZAddArgsCH(t, provider)
	})
	t.Run("ZCard", func(t *testing.T) {
		testZCard(t, provider)
	})
	t.Run("ZCard_NonExistent", func(t *testing.T) {
		testZCardNonExistent(t, provider)
	})
	t.Run("ZCount", func(t *testing.T) {
		testZCount(t, provider)
	})
	t.Run("ZIncrBy", func(t *testing.T) {
		testZIncrBy(t, provider)
	})
	t.Run("ZRank", func(t *testing.T) {
		testZRank(t, provider)
	})
	t.Run("ZRank_NonExistent", func(t *testing.T) {
		testZRankNonExistent(t, provider)
	})
	t.Run("ZRevRank", func(t *testing.T) {
		testZRevRank(t, provider)
	})
	t.Run("ZScore", func(t *testing.T) {
		testZScore(t, provider)
	})
	t.Run("ZScore_NonExistent", func(t *testing.T) {
		testZScoreNonExistent(t, provider)
	})
	t.Run("ZRem", func(t *testing.T) {
		testZRem(t, provider)
	})
	t.Run("ZRemRangeByRank", func(t *testing.T) {
		testZRemRangeByRank(t, provider)
	})
	t.Run("ZRemRangeByScore", func(t *testing.T) {
		testZRemRangeByScore(t, provider)
	})
	t.Run("ZRangeWithScores", func(t *testing.T) {
		testZRangeWithScores(t, provider)
	})
	t.Run("ZRevRange", func(t *testing.T) {
		testZRevRange(t, provider)
	})
	t.Run("ZRevRangeWithScores", func(t *testing.T) {
		testZRevRangeWithScores(t, provider)
	})
	t.Run("ZRangeByScore", func(t *testing.T) {
		testZRangeByScore(t, provider)
	})
	t.Run("ZRangeByScoreWithScores", func(t *testing.T) {
		testZRangeByScoreWithScores(t, provider)
	})
	t.Run("ZRevRangeByScore", func(t *testing.T) {
		testZRevRangeByScore(t, provider)
	})
	t.Run("ZInter", func(t *testing.T) {
		testZInter(t, provider)
	})
	t.Run("ZInterWithScores", func(t *testing.T) {
		testZInterWithScores(t, provider)
	})
	t.Run("ZInterStore", func(t *testing.T) {
		testZInterStore(t, provider)
	})
	t.Run("ZUnion", func(t *testing.T) {
		testZUnion(t, provider)
	})
	t.Run("ZUnionWithScores", func(t *testing.T) {
		testZUnionWithScores(t, provider)
	})
	t.Run("ZUnionStore", func(t *testing.T) {
		testZUnionStore(t, provider)
	})
	t.Run("ZScan_Basic", func(t *testing.T) {
		testZScanBasic(t, provider)
	})
	t.Run("ZScan_WithPattern", func(t *testing.T) {
		testZScanWithPattern(t, provider)
	})
	t.Run("ZScan_NonExistent", func(t *testing.T) {
		testZScanNonExistent(t, provider)
	})
}

// testZAddAndZRange tests ZAdd and ZRange operations
func testZAddAndZRange(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zadd"

	// Add members with scores
	members := []caches.ZMember{
		{Member: []byte("one"), Score: 1},
		{Member: []byte("two"), Score: 2},
		{Member: []byte("three"), Score: 3},
	}
	result := cmd.ZAdd(ctx, key, members...)
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val())

	// Get range
	rangeResult := cmd.ZRange(ctx, key, 0, -1)
	require.NoError(t, rangeResult.Err())
	require.Equal(t, 3, len(rangeResult.Val()))
	require.Equal(t, []byte("one"), rangeResult.Val()[0])
	require.Equal(t, []byte("two"), rangeResult.Val()[1])
	require.Equal(t, []byte("three"), rangeResult.Val()[2])
}

// testZAddUpdateScore tests ZAdd updating existing member's score
func testZAddUpdateScore(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zadd_update"

	// Add member
	cmd.ZAdd(ctx, key, caches.ZMember{Member: []byte("one"), Score: 1})

	// Update score
	result := cmd.ZAdd(ctx, key, caches.ZMember{Member: []byte("one"), Score: 10})
	require.NoError(t, result.Err())
	require.Equal(t, int64(0), result.Val()) // No new members added

	// Verify score
	scoreResult := cmd.ZScore(ctx, key, "one")
	require.NoError(t, scoreResult.Err())
	require.Equal(t, 10.0, scoreResult.Val())
}

// testZAddArgsNX tests ZAddArgs with NX mode (only add new)
func testZAddArgsNX(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zadd_nx"

	// Add member
	cmd.ZAdd(ctx, key, caches.ZMember{Member: []byte("one"), Score: 1})

	// Try to update with NX (should fail)
	result := cmd.ZAddArgs(ctx, key, "NX", false, caches.ZMember{Member: []byte("one"), Score: 10})
	require.NoError(t, result.Err())
	require.Equal(t, int64(0), result.Val())

	// Verify score unchanged
	scoreResult := cmd.ZScore(ctx, key, "one")
	require.NoError(t, scoreResult.Err())
	require.Equal(t, 1.0, scoreResult.Val())

	// Add new member with NX (should succeed)
	result2 := cmd.ZAddArgs(ctx, key, "NX", false, caches.ZMember{Member: []byte("two"), Score: 2})
	require.NoError(t, result2.Err())
	require.Equal(t, int64(1), result2.Val())
}

// testZAddArgsXX tests ZAddArgs with XX mode (only update existing)
func testZAddArgsXX(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zadd_xx"

	// Add member
	cmd.ZAdd(ctx, key, caches.ZMember{Member: []byte("one"), Score: 1})

	// Update with XX (should succeed)
	result := cmd.ZAddArgs(ctx, key, "XX", false, caches.ZMember{Member: []byte("one"), Score: 10})
	require.NoError(t, result.Err())
	require.Equal(t, int64(0), result.Val())

	// Verify score updated
	scoreResult := cmd.ZScore(ctx, key, "one")
	require.NoError(t, scoreResult.Err())
	require.Equal(t, 10.0, scoreResult.Val())

	// Try to add new member with XX (should fail)
	result2 := cmd.ZAddArgs(ctx, key, "XX", false, caches.ZMember{Member: []byte("two"), Score: 2})
	require.NoError(t, result2.Err())
	require.Equal(t, int64(0), result2.Val())
}

// testZAddArgsGT tests ZAddArgs with GT mode (only if new score is greater)
func testZAddArgsGT(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zadd_gt"

	// Add member
	cmd.ZAdd(ctx, key, caches.ZMember{Member: []byte("one"), Score: 5})

	// Update with greater score (should succeed)
	result := cmd.ZAddArgs(ctx, key, "GT", false, caches.ZMember{Member: []byte("one"), Score: 10})
	require.NoError(t, result.Err())

	// Verify score updated
	scoreResult := cmd.ZScore(ctx, key, "one")
	require.NoError(t, scoreResult.Err())
	require.Equal(t, 10.0, scoreResult.Val())

	// Try to update with lesser score (should not update)
	result2 := cmd.ZAddArgs(ctx, key, "GT", false, caches.ZMember{Member: []byte("one"), Score: 3})
	require.NoError(t, result2.Err())

	// Verify score unchanged
	scoreResult2 := cmd.ZScore(ctx, key, "one")
	require.NoError(t, scoreResult2.Err())
	require.Equal(t, 10.0, scoreResult2.Val())
}

// testZAddArgsLT tests ZAddArgs with LT mode (only if new score is less)
func testZAddArgsLT(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zadd_lt"

	// Add member
	cmd.ZAdd(ctx, key, caches.ZMember{Member: []byte("one"), Score: 5})

	// Update with lesser score (should succeed)
	result := cmd.ZAddArgs(ctx, key, "LT", false, caches.ZMember{Member: []byte("one"), Score: 3})
	require.NoError(t, result.Err())

	// Verify score updated
	scoreResult := cmd.ZScore(ctx, key, "one")
	require.NoError(t, scoreResult.Err())
	require.Equal(t, 3.0, scoreResult.Val())

	// Try to update with greater score (should not update)
	result2 := cmd.ZAddArgs(ctx, key, "LT", false, caches.ZMember{Member: []byte("one"), Score: 10})
	require.NoError(t, result2.Err())

	// Verify score unchanged
	scoreResult2 := cmd.ZScore(ctx, key, "one")
	require.NoError(t, scoreResult2.Err())
	require.Equal(t, 3.0, scoreResult2.Val())
}

// testZAddArgsCH tests ZAddArgs with CH flag (return changed count)
func testZAddArgsCH(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zadd_ch"

	// Add member
	cmd.ZAdd(ctx, key, caches.ZMember{Member: []byte("one"), Score: 1})

	// Update with CH flag
	result := cmd.ZAddArgs(ctx, key, "", true, caches.ZMember{Member: []byte("one"), Score: 10})
	require.NoError(t, result.Err())
	require.Equal(t, int64(1), result.Val()) // Changed count = 1
}

// testZCard tests ZCard operation
func testZCard(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zcard"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("one"), Score: 1},
		caches.ZMember{Member: []byte("two"), Score: 2},
		caches.ZMember{Member: []byte("three"), Score: 3},
	)

	// Get cardinality
	result := cmd.ZCard(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val())
}

// testZCardNonExistent tests ZCard on non-existent key
func testZCardNonExistent(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zcard_nonexistent"

	result := cmd.ZCard(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(0), result.Val())
}

// testZCount tests ZCount operation
func testZCount(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zcount"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("one"), Score: 1},
		caches.ZMember{Member: []byte("two"), Score: 2},
		caches.ZMember{Member: []byte("three"), Score: 3},
		caches.ZMember{Member: []byte("four"), Score: 4},
		caches.ZMember{Member: []byte("five"), Score: 5},
	)

	// Count members in range
	result := cmd.ZCount(ctx, key, "2", "4")
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val()) // 2, 3, 4
}

// testZIncrBy tests ZIncrBy operation
func testZIncrBy(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zincrby"

	// Add member
	cmd.ZAdd(ctx, key, caches.ZMember{Member: []byte("one"), Score: 1})

	// Increment score
	result := cmd.ZIncrBy(ctx, key, 5, "one")
	require.NoError(t, result.Err())
	require.Equal(t, 6.0, result.Val())

	// Verify score
	scoreResult := cmd.ZScore(ctx, key, "one")
	require.NoError(t, scoreResult.Err())
	require.Equal(t, 6.0, scoreResult.Val())
}

// testZRank tests ZRank operation
func testZRank(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zrank"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("one"), Score: 1},
		caches.ZMember{Member: []byte("two"), Score: 2},
		caches.ZMember{Member: []byte("three"), Score: 3},
	)

	// Get rank
	result := cmd.ZRank(ctx, key, "two")
	require.NoError(t, result.Err())
	require.Equal(t, int64(1), result.Val()) // 0-based index
}

// testZRankNonExistent tests ZRank on non-existent member
func testZRankNonExistent(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zrank_nonexistent"

	// Add members
	cmd.ZAdd(ctx, key, caches.ZMember{Member: []byte("one"), Score: 1})

	// Get rank of non-existent member
	result := cmd.ZRank(ctx, key, "two")
	require.Equal(t, caches.Nil, result.Err())
}

// testZRevRank tests ZRevRank operation
func testZRevRank(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zrevrank"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("one"), Score: 1},
		caches.ZMember{Member: []byte("two"), Score: 2},
		caches.ZMember{Member: []byte("three"), Score: 3},
	)

	// Get reverse rank
	result := cmd.ZRevRank(ctx, key, "two")
	require.NoError(t, result.Err())
	require.Equal(t, int64(1), result.Val()) // Reversed: three(0), two(1), one(2)
}

// testZScore tests ZScore operation
func testZScore(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zscore"

	// Add member
	cmd.ZAdd(ctx, key, caches.ZMember{Member: []byte("one"), Score: 1.5})

	// Get score
	result := cmd.ZScore(ctx, key, "one")
	require.NoError(t, result.Err())
	require.Equal(t, 1.5, result.Val())
}

// testZScoreNonExistent tests ZScore on non-existent member
func testZScoreNonExistent(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zscore_nonexistent"

	result := cmd.ZScore(ctx, key, "one")
	require.Equal(t, caches.Nil, result.Err())
}

// testZRem tests ZRem operation
func testZRem(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zrem"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("one"), Score: 1},
		caches.ZMember{Member: []byte("two"), Score: 2},
		caches.ZMember{Member: []byte("three"), Score: 3},
	)

	// Remove members
	result := cmd.ZRem(ctx, key, "one", "three")
	require.NoError(t, result.Err())
	require.Equal(t, int64(2), result.Val())

	// Verify remaining
	cardResult := cmd.ZCard(ctx, key)
	require.NoError(t, cardResult.Err())
	require.Equal(t, int64(1), cardResult.Val())
}

// testZRemRangeByRank tests ZRemRangeByRank operation
func testZRemRangeByRank(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zremrangebyrank"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("one"), Score: 1},
		caches.ZMember{Member: []byte("two"), Score: 2},
		caches.ZMember{Member: []byte("three"), Score: 3},
		caches.ZMember{Member: []byte("four"), Score: 4},
		caches.ZMember{Member: []byte("five"), Score: 5},
	)

	// Remove by rank range
	result := cmd.ZRemRangeByRank(ctx, key, 1, 3)
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val()) // Removes two, three, four

	// Verify remaining
	cardResult := cmd.ZCard(ctx, key)
	require.NoError(t, cardResult.Err())
	require.Equal(t, int64(2), cardResult.Val())
}

// testZRemRangeByScore tests ZRemRangeByScore operation
func testZRemRangeByScore(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zremrangebyscore"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("one"), Score: 1},
		caches.ZMember{Member: []byte("two"), Score: 2},
		caches.ZMember{Member: []byte("three"), Score: 3},
		caches.ZMember{Member: []byte("four"), Score: 4},
		caches.ZMember{Member: []byte("five"), Score: 5},
	)

	// Remove by score range
	result := cmd.ZRemRangeByScore(ctx, key, "2", "4")
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val()) // Removes 2, 3, 4

	// Verify remaining
	cardResult := cmd.ZCard(ctx, key)
	require.NoError(t, cardResult.Err())
	require.Equal(t, int64(2), cardResult.Val())
}

// testZRangeWithScores tests ZRangeWithScores operation
func testZRangeWithScores(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zrangewithscores"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("one"), Score: 1},
		caches.ZMember{Member: []byte("two"), Score: 2},
		caches.ZMember{Member: []byte("three"), Score: 3},
	)

	// Get range with scores
	result := cmd.ZRangeWithScores(ctx, key, 0, -1)
	require.NoError(t, result.Err())
	require.Equal(t, 3, len(result.Val()))

	members := result.Val()
	require.Equal(t, []byte("one"), members[0].Member)
	require.Equal(t, 1.0, members[0].Score)
	require.Equal(t, []byte("two"), members[1].Member)
	require.Equal(t, 2.0, members[1].Score)
	require.Equal(t, []byte("three"), members[2].Member)
	require.Equal(t, 3.0, members[2].Score)
}

// testZRevRange tests ZRevRange operation
func testZRevRange(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zrevrange"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("one"), Score: 1},
		caches.ZMember{Member: []byte("two"), Score: 2},
		caches.ZMember{Member: []byte("three"), Score: 3},
	)

	// Get reverse range
	result := cmd.ZRevRange(ctx, key, 0, -1)
	require.NoError(t, result.Err())
	require.Equal(t, 3, len(result.Val()))
	require.Equal(t, []byte("three"), result.Val()[0])
	require.Equal(t, []byte("two"), result.Val()[1])
	require.Equal(t, []byte("one"), result.Val()[2])
}

// testZRevRangeWithScores tests ZRevRangeWithScores operation
func testZRevRangeWithScores(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zrevrangewithscores"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("one"), Score: 1},
		caches.ZMember{Member: []byte("two"), Score: 2},
		caches.ZMember{Member: []byte("three"), Score: 3},
	)

	// Get reverse range with scores
	result := cmd.ZRevRangeWithScores(ctx, key, 0, -1)
	require.NoError(t, result.Err())
	require.Equal(t, 3, len(result.Val()))

	members := result.Val()
	require.Equal(t, []byte("three"), members[0].Member)
	require.Equal(t, 3.0, members[0].Score)
}

// testZRangeByScore tests ZRangeByScore operation
func testZRangeByScore(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zrangebyscore"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("one"), Score: 1},
		caches.ZMember{Member: []byte("two"), Score: 2},
		caches.ZMember{Member: []byte("three"), Score: 3},
		caches.ZMember{Member: []byte("four"), Score: 4},
		caches.ZMember{Member: []byte("five"), Score: 5},
	)

	// Get range by score
	result := cmd.ZRangeByScore(ctx, key, "2", "4")
	require.NoError(t, result.Err())
	require.Equal(t, 3, len(result.Val()))
	require.Equal(t, []byte("two"), result.Val()[0])
	require.Equal(t, []byte("three"), result.Val()[1])
	require.Equal(t, []byte("four"), result.Val()[2])
}

// testZRangeByScoreWithScores tests ZRangeByScoreWithScores operation
func testZRangeByScoreWithScores(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zrangebyscorewithscores"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("one"), Score: 1},
		caches.ZMember{Member: []byte("two"), Score: 2},
		caches.ZMember{Member: []byte("three"), Score: 3},
	)

	// Get range by score with scores
	result := cmd.ZRangeByScoreWithScores(ctx, key, "1", "2")
	require.NoError(t, result.Err())
	require.Equal(t, 2, len(result.Val()))

	members := result.Val()
	require.Equal(t, []byte("one"), members[0].Member)
	require.Equal(t, 1.0, members[0].Score)
	require.Equal(t, []byte("two"), members[1].Member)
	require.Equal(t, 2.0, members[1].Score)
}

// testZRevRangeByScore tests ZRevRangeByScore operation
func testZRevRangeByScore(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zrevrangebyscore"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("one"), Score: 1},
		caches.ZMember{Member: []byte("two"), Score: 2},
		caches.ZMember{Member: []byte("three"), Score: 3},
		caches.ZMember{Member: []byte("four"), Score: 4},
		caches.ZMember{Member: []byte("five"), Score: 5},
	)

	// Get reverse range by score
	result := cmd.ZRevRangeByScore(ctx, key, "4", "2")
	require.NoError(t, result.Err())
	require.Equal(t, 3, len(result.Val()))
	require.Equal(t, []byte("four"), result.Val()[0])
	require.Equal(t, []byte("three"), result.Val()[1])
	require.Equal(t, []byte("two"), result.Val()[2])
}

// testZInter tests ZInter operation
func testZInter(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key1 := "test:zset:zinter1"
	key2 := "test:zset:zinter2"

	// Setup sorted sets
	cmd.ZAdd(ctx, key1,
		caches.ZMember{Member: []byte("a"), Score: 1},
		caches.ZMember{Member: []byte("b"), Score: 2},
		caches.ZMember{Member: []byte("c"), Score: 3},
	)
	cmd.ZAdd(ctx, key2,
		caches.ZMember{Member: []byte("b"), Score: 4},
		caches.ZMember{Member: []byte("c"), Score: 5},
		caches.ZMember{Member: []byte("d"), Score: 6},
	)

	// Get intersection
	store := caches.ZStore{Keys: []string{key1, key2}}
	result := cmd.ZInter(ctx, store)
	require.NoError(t, result.Err())

	// Expected: {b, c}
	inter := result.Val()
	require.Equal(t, 2, len(inter))
}

// testZInterWithScores tests ZInterWithScores operation
func testZInterWithScores(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key1 := "test:zset:zinterwithscores1"
	key2 := "test:zset:zinterwithscores2"

	// Setup sorted sets
	cmd.ZAdd(ctx, key1,
		caches.ZMember{Member: []byte("a"), Score: 1},
		caches.ZMember{Member: []byte("b"), Score: 2},
	)
	cmd.ZAdd(ctx, key2,
		caches.ZMember{Member: []byte("b"), Score: 3},
		caches.ZMember{Member: []byte("c"), Score: 4},
	)

	// Get intersection with scores
	store := caches.ZStore{Keys: []string{key1, key2}}
	result := cmd.ZInterWithScores(ctx, store)
	require.NoError(t, result.Err())

	// Expected: {b with score 5 (2+3)}
	members := result.Val()
	require.Equal(t, 1, len(members))
	require.Equal(t, []byte("b"), members[0].Member)
	require.Equal(t, 5.0, members[0].Score) // Default aggregate is SUM
}

// testZInterStore tests ZInterStore operation
func testZInterStore(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key1 := "test:zset:zinterstore1"
	key2 := "test:zset:zinterstore2"
	dest := "test:zset:zinterstore_dest"

	// Setup sorted sets
	cmd.ZAdd(ctx, key1,
		caches.ZMember{Member: []byte("a"), Score: 1},
		caches.ZMember{Member: []byte("b"), Score: 2},
	)
	cmd.ZAdd(ctx, key2,
		caches.ZMember{Member: []byte("b"), Score: 3},
		caches.ZMember{Member: []byte("c"), Score: 4},
	)

	// Store intersection
	store := caches.ZStore{Keys: []string{key1, key2}}
	result := cmd.ZInterStore(ctx, dest, store)
	require.NoError(t, result.Err())
	require.Equal(t, int64(1), result.Val())

	// Verify destination
	cardResult := cmd.ZCard(ctx, dest)
	require.NoError(t, cardResult.Err())
	require.Equal(t, int64(1), cardResult.Val())
}

// testZUnion tests ZUnion operation
func testZUnion(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key1 := "test:zset:zunion1"
	key2 := "test:zset:zunion2"

	// Setup sorted sets
	cmd.ZAdd(ctx, key1,
		caches.ZMember{Member: []byte("a"), Score: 1},
		caches.ZMember{Member: []byte("b"), Score: 2},
	)
	cmd.ZAdd(ctx, key2,
		caches.ZMember{Member: []byte("c"), Score: 3},
		caches.ZMember{Member: []byte("d"), Score: 4},
	)

	// Get union
	store := caches.ZStore{Keys: []string{key1, key2}}
	result := cmd.ZUnion(ctx, store)
	require.NoError(t, result.Err())

	// Expected: {a, b, c, d}
	union := result.Val()
	require.Equal(t, 4, len(union))
}

// testZUnionWithScores tests ZUnionWithScores operation
func testZUnionWithScores(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key1 := "test:zset:zunionwithscores1"
	key2 := "test:zset:zunionwithscores2"

	// Setup sorted sets
	cmd.ZAdd(ctx, key1,
		caches.ZMember{Member: []byte("a"), Score: 1},
		caches.ZMember{Member: []byte("b"), Score: 2},
	)
	cmd.ZAdd(ctx, key2,
		caches.ZMember{Member: []byte("b"), Score: 3},
		caches.ZMember{Member: []byte("c"), Score: 4},
	)

	// Get union with scores
	store := caches.ZStore{Keys: []string{key1, key2}}
	result := cmd.ZUnionWithScores(ctx, store)
	require.NoError(t, result.Err())

	// Expected: {a:1, b:5, c:4}
	members := result.Val()
	require.Equal(t, 3, len(members))
}

// testZUnionStore tests ZUnionStore operation
func testZUnionStore(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key1 := "test:zset:zunionstore1"
	key2 := "test:zset:zunionstore2"
	dest := "test:zset:zunionstore_dest"

	// Setup sorted sets
	cmd.ZAdd(ctx, key1,
		caches.ZMember{Member: []byte("a"), Score: 1},
		caches.ZMember{Member: []byte("b"), Score: 2},
	)
	cmd.ZAdd(ctx, key2,
		caches.ZMember{Member: []byte("c"), Score: 3},
	)

	// Store union
	store := caches.ZStore{Keys: []string{key1, key2}}
	result := cmd.ZUnionStore(ctx, dest, store)
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val())

	// Verify destination
	cardResult := cmd.ZCard(ctx, dest)
	require.NoError(t, cardResult.Err())
	require.Equal(t, int64(3), cardResult.Val())
}

// testZScanBasic tests ZScan basic iteration
func testZScanBasic(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zscan"

	// Add members
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("member1"), Score: 1},
		caches.ZMember{Member: []byte("member2"), Score: 2},
		caches.ZMember{Member: []byte("member3"), Score: 3},
		caches.ZMember{Member: []byte("member4"), Score: 4},
		caches.ZMember{Member: []byte("member5"), Score: 5},
	)

	// Scan all members
	var allMembers []caches.ZMember
	cursor := uint64(0)
	iterations := 0
	maxIterations := 10

	for iterations < maxIterations {
		result := cmd.ZScan(ctx, key, cursor, "", 10)
		require.NoError(t, result.Err())

		scanResult := result.Val()
		allMembers = append(allMembers, scanResult.Members...)
		cursor = scanResult.Cursor

		iterations++
		if cursor == 0 {
			break
		}
	}

	// Should have found all members
	require.GreaterOrEqual(t, len(allMembers), 5)
}

// testZScanWithPattern tests ZScan with pattern matching
func testZScanWithPattern(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zscan_pattern"

	// Add members with pattern
	cmd.ZAdd(ctx, key,
		caches.ZMember{Member: []byte("foo1"), Score: 1},
		caches.ZMember{Member: []byte("foo2"), Score: 2},
		caches.ZMember{Member: []byte("bar1"), Score: 3},
		caches.ZMember{Member: []byte("bar2"), Score: 4},
	)

	// Scan with pattern
	var allMembers []caches.ZMember
	cursor := uint64(0)
	iterations := 0
	maxIterations := 10

	for iterations < maxIterations {
		result := cmd.ZScan(ctx, key, cursor, "foo*", 10)
		require.NoError(t, result.Err())

		scanResult := result.Val()
		allMembers = append(allMembers, scanResult.Members...)
		cursor = scanResult.Cursor

		iterations++
		if cursor == 0 {
			break
		}
	}

	// Should have found foo members
	for _, m := range allMembers {
		member := string(m.Member)
		require.True(t, member == "foo1" || member == "foo2", "Expected foo* pattern match")
	}
}

// testZScanNonExistent tests ZScan on non-existent key
func testZScanNonExistent(t *testing.T, provider SortedSetCommandProvider) {
	cmd := provider.GetSortedSetCommand()
	ctx := provider.GetContext()

	key := "test:zset:zscan_nonexistent"

	result := cmd.ZScan(ctx, key, 0, "", 10)
	require.NoError(t, result.Err())

	scanResult := result.Val()
	require.Empty(t, scanResult.Members)
	require.Equal(t, uint64(0), scanResult.Cursor)
}
