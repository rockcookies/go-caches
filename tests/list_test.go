package tests

import (
	"context"
	"testing"

	"github.com/rockcookies/go-caches"
	"github.com/stretchr/testify/require"
)

// ListCommandProvider defines the interface for testing ListCommand implementations
type ListCommandProvider interface {
	GetListCommand() caches.ListCommand
	GetContext() context.Context
}

// RunListCommandTests runs all ListCommand tests
func RunListCommandTests(t *testing.T, provider ListCommandProvider) {
	t.Run("LPush_and_LRange", func(t *testing.T) {
		testLPushAndLRange(t, provider)
	})
	t.Run("RPush_and_LRange", func(t *testing.T) {
		testRPushAndLRange(t, provider)
	})
	t.Run("LLen", func(t *testing.T) {
		testLLen(t, provider)
	})
	t.Run("LLen_NonExistent", func(t *testing.T) {
		testLLenNonExistent(t, provider)
	})
	t.Run("LIndex", func(t *testing.T) {
		testLIndex(t, provider)
	})
	t.Run("LIndex_Negative", func(t *testing.T) {
		testLIndexNegative(t, provider)
	})
	t.Run("LIndex_OutOfRange", func(t *testing.T) {
		testLIndexOutOfRange(t, provider)
	})
	t.Run("LPop", func(t *testing.T) {
		testLPop(t, provider)
	})
	t.Run("LPop_NonExistent", func(t *testing.T) {
		testLPopNonExistent(t, provider)
	})
	t.Run("LPopCount", func(t *testing.T) {
		testLPopCount(t, provider)
	})
	t.Run("RPop", func(t *testing.T) {
		testRPop(t, provider)
	})
	t.Run("RPop_NonExistent", func(t *testing.T) {
		testRPopNonExistent(t, provider)
	})
	t.Run("RPopCount", func(t *testing.T) {
		testRPopCount(t, provider)
	})
	t.Run("LInsert_Before", func(t *testing.T) {
		testLInsertBefore(t, provider)
	})
	t.Run("LInsert_After", func(t *testing.T) {
		testLInsertAfter(t, provider)
	})
	t.Run("LInsert_PivotNotFound", func(t *testing.T) {
		testLInsertPivotNotFound(t, provider)
	})
	t.Run("LSet", func(t *testing.T) {
		testLSet(t, provider)
	})
	t.Run("LSet_NegativeIndex", func(t *testing.T) {
		testLSetNegativeIndex(t, provider)
	})
	t.Run("LRem_Positive", func(t *testing.T) {
		testLRemPositive(t, provider)
	})
	t.Run("LRem_Negative", func(t *testing.T) {
		testLRemNegative(t, provider)
	})
	t.Run("LRem_Zero", func(t *testing.T) {
		testLRemZero(t, provider)
	})
	t.Run("LTrim", func(t *testing.T) {
		testLTrim(t, provider)
	})
	t.Run("RPopLPush", func(t *testing.T) {
		testRPopLPush(t, provider)
	})
	t.Run("RPopLPush_SameKey", func(t *testing.T) {
		testRPopLPushSameKey(t, provider)
	})
}

// testLPushAndLRange tests LPush and LRange operations
func testLPushAndLRange(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:lpush"

	// Push multiple elements
	result := cmd.LPush(ctx, key, "three", "two", "one")
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val())

	// Get all elements
	rangeResult := cmd.LRange(ctx, key, 0, -1)
	require.NoError(t, rangeResult.Err())
	require.Equal(t, 3, len(rangeResult.Val()))
	require.Equal(t, []byte("one"), rangeResult.Val()[0])
	require.Equal(t, []byte("two"), rangeResult.Val()[1])
	require.Equal(t, []byte("three"), rangeResult.Val()[2])
}

// testRPushAndLRange tests RPush and LRange operations
func testRPushAndLRange(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:rpush"

	// Push multiple elements
	result := cmd.RPush(ctx, key, "one", "two", "three")
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val())

	// Get all elements
	rangeResult := cmd.LRange(ctx, key, 0, -1)
	require.NoError(t, rangeResult.Err())
	require.Equal(t, 3, len(rangeResult.Val()))
	require.Equal(t, []byte("one"), rangeResult.Val()[0])
	require.Equal(t, []byte("two"), rangeResult.Val()[1])
	require.Equal(t, []byte("three"), rangeResult.Val()[2])
}

// testLLen tests LLen operation
func testLLen(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:llen"

	// Push elements
	cmd.RPush(ctx, key, "one", "two", "three")

	// Get length
	result := cmd.LLen(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val())
}

// testLLenNonExistent tests LLen on non-existent key
func testLLenNonExistent(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:llen_nonexistent"

	result := cmd.LLen(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(0), result.Val())
}

// testLIndex tests LIndex operation
func testLIndex(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:lindex"

	// Push elements
	cmd.RPush(ctx, key, "one", "two", "three")

	// Get by index
	result := cmd.LIndex(ctx, key, 0)
	require.NoError(t, result.Err())
	require.Equal(t, []byte("one"), result.Val())

	result2 := cmd.LIndex(ctx, key, 1)
	require.NoError(t, result2.Err())
	require.Equal(t, []byte("two"), result2.Val())
}

// testLIndexNegative tests LIndex with negative index
func testLIndexNegative(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:lindex_neg"

	// Push elements
	cmd.RPush(ctx, key, "one", "two", "three")

	// Get by negative index
	result := cmd.LIndex(ctx, key, -1)
	require.NoError(t, result.Err())
	require.Equal(t, []byte("three"), result.Val())

	result2 := cmd.LIndex(ctx, key, -2)
	require.NoError(t, result2.Err())
	require.Equal(t, []byte("two"), result2.Val())
}

// testLIndexOutOfRange tests LIndex with out of range index
func testLIndexOutOfRange(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:lindex_oor"

	// Push elements
	cmd.RPush(ctx, key, "one", "two")

	// Get out of range index
	result := cmd.LIndex(ctx, key, 10)
	require.Equal(t, caches.Nil, result.Err())
	require.Nil(t, result.Val())
}

// testLPop tests LPop operation
func testLPop(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:lpop"

	// Push elements
	cmd.RPush(ctx, key, "one", "two", "three")

	// Pop from left
	result := cmd.LPop(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, []byte("one"), result.Val())

	// Verify list
	lenResult := cmd.LLen(ctx, key)
	require.NoError(t, lenResult.Err())
	require.Equal(t, int64(2), lenResult.Val())
}

// testLPopNonExistent tests LPop on non-existent key
func testLPopNonExistent(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:lpop_nonexistent"

	result := cmd.LPop(ctx, key)
	require.Equal(t, caches.Nil, result.Err())
	require.Nil(t, result.Val())
}

// testLPopCount tests LPopCount operation
func testLPopCount(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:lpopcount"

	// Push elements
	cmd.RPush(ctx, key, "one", "two", "three", "four")

	// Pop multiple from left
	result := cmd.LPopCount(ctx, key, 2)
	require.NoError(t, result.Err())
	require.Equal(t, 2, len(result.Val()))
	require.Equal(t, []byte("one"), result.Val()[0])
	require.Equal(t, []byte("two"), result.Val()[1])

	// Verify remaining
	lenResult := cmd.LLen(ctx, key)
	require.NoError(t, lenResult.Err())
	require.Equal(t, int64(2), lenResult.Val())
}

// testRPop tests RPop operation
func testRPop(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:rpop"

	// Push elements
	cmd.RPush(ctx, key, "one", "two", "three")

	// Pop from right
	result := cmd.RPop(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, []byte("three"), result.Val())

	// Verify list
	lenResult := cmd.LLen(ctx, key)
	require.NoError(t, lenResult.Err())
	require.Equal(t, int64(2), lenResult.Val())
}

// testRPopNonExistent tests RPop on non-existent key
func testRPopNonExistent(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:rpop_nonexistent"

	result := cmd.RPop(ctx, key)
	require.Equal(t, caches.Nil, result.Err())
	require.Nil(t, result.Val())
}

// testRPopCount tests RPopCount operation
func testRPopCount(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:rpopcount"

	// Push elements
	cmd.RPush(ctx, key, "one", "two", "three", "four")

	// Pop multiple from right
	result := cmd.RPopCount(ctx, key, 2)
	require.NoError(t, result.Err())
	require.Equal(t, 2, len(result.Val()))
	require.Equal(t, []byte("four"), result.Val()[0])
	require.Equal(t, []byte("three"), result.Val()[1])

	// Verify remaining
	lenResult := cmd.LLen(ctx, key)
	require.NoError(t, lenResult.Err())
	require.Equal(t, int64(2), lenResult.Val())
}

// testLInsertBefore tests LInsert with BEFORE position
func testLInsertBefore(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:linsert_before"

	// Push elements
	cmd.RPush(ctx, key, "one", "three")

	// Insert before "three"
	result := cmd.LInsert(ctx, key, caches.LInsertBefore, "three", "two")
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val())

	// Verify list
	rangeResult := cmd.LRange(ctx, key, 0, -1)
	require.NoError(t, rangeResult.Err())
	require.Equal(t, []byte("one"), rangeResult.Val()[0])
	require.Equal(t, []byte("two"), rangeResult.Val()[1])
	require.Equal(t, []byte("three"), rangeResult.Val()[2])
}

// testLInsertAfter tests LInsert with AFTER position
func testLInsertAfter(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:linsert_after"

	// Push elements
	cmd.RPush(ctx, key, "one", "three")

	// Insert after "one"
	result := cmd.LInsert(ctx, key, caches.LInsertAfter, "one", "two")
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val())

	// Verify list
	rangeResult := cmd.LRange(ctx, key, 0, -1)
	require.NoError(t, rangeResult.Err())
	require.Equal(t, []byte("one"), rangeResult.Val()[0])
	require.Equal(t, []byte("two"), rangeResult.Val()[1])
	require.Equal(t, []byte("three"), rangeResult.Val()[2])
}

// testLInsertPivotNotFound tests LInsert with non-existent pivot
func testLInsertPivotNotFound(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:linsert_notfound"

	// Push elements
	cmd.RPush(ctx, key, "one", "two")

	// Insert with non-existent pivot
	result := cmd.LInsert(ctx, key, caches.LInsertBefore, "three", "value")
	require.NoError(t, result.Err())
	require.Equal(t, int64(-1), result.Val())
}

// testLSet tests LSet operation
func testLSet(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:lset"

	// Push elements
	cmd.RPush(ctx, key, "one", "two", "three")

	// Set element at index
	result := cmd.LSet(ctx, key, 1, "TWO")
	require.NoError(t, result.Err())
	require.Equal(t, "OK", result.Val())

	// Verify
	indexResult := cmd.LIndex(ctx, key, 1)
	require.NoError(t, indexResult.Err())
	require.Equal(t, []byte("TWO"), indexResult.Val())
}

// testLSetNegativeIndex tests LSet with negative index
func testLSetNegativeIndex(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:lset_neg"

	// Push elements
	cmd.RPush(ctx, key, "one", "two", "three")

	// Set element at negative index
	result := cmd.LSet(ctx, key, -1, "THREE")
	require.NoError(t, result.Err())
	require.Equal(t, "OK", result.Val())

	// Verify
	indexResult := cmd.LIndex(ctx, key, -1)
	require.NoError(t, indexResult.Err())
	require.Equal(t, []byte("THREE"), indexResult.Val())
}

// testLRemPositive tests LRem with positive count
func testLRemPositive(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:lrem_pos"

	// Push elements with duplicates
	cmd.RPush(ctx, key, "a", "b", "a", "c", "a")

	// Remove first 2 occurrences of "a"
	result := cmd.LRem(ctx, key, 2, "a")
	require.NoError(t, result.Err())
	require.Equal(t, int64(2), result.Val())

	// Verify
	rangeResult := cmd.LRange(ctx, key, 0, -1)
	require.NoError(t, rangeResult.Err())
	require.Equal(t, 3, len(rangeResult.Val()))
	require.Equal(t, []byte("b"), rangeResult.Val()[0])
	require.Equal(t, []byte("c"), rangeResult.Val()[1])
	require.Equal(t, []byte("a"), rangeResult.Val()[2])
}

// testLRemNegative tests LRem with negative count
func testLRemNegative(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:lrem_neg"

	// Push elements with duplicates
	cmd.RPush(ctx, key, "a", "b", "a", "c", "a")

	// Remove last 2 occurrences of "a"
	result := cmd.LRem(ctx, key, -2, "a")
	require.NoError(t, result.Err())
	require.Equal(t, int64(2), result.Val())

	// Verify
	rangeResult := cmd.LRange(ctx, key, 0, -1)
	require.NoError(t, rangeResult.Err())
	require.Equal(t, 3, len(rangeResult.Val()))
	require.Equal(t, []byte("a"), rangeResult.Val()[0])
	require.Equal(t, []byte("b"), rangeResult.Val()[1])
	require.Equal(t, []byte("c"), rangeResult.Val()[2])
}

// testLRemZero tests LRem with count = 0 (remove all)
func testLRemZero(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:lrem_zero"

	// Push elements with duplicates
	cmd.RPush(ctx, key, "a", "b", "a", "c", "a")

	// Remove all occurrences of "a"
	result := cmd.LRem(ctx, key, 0, "a")
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val())

	// Verify
	rangeResult := cmd.LRange(ctx, key, 0, -1)
	require.NoError(t, rangeResult.Err())
	require.Equal(t, 2, len(rangeResult.Val()))
	require.Equal(t, []byte("b"), rangeResult.Val()[0])
	require.Equal(t, []byte("c"), rangeResult.Val()[1])
}

// testLTrim tests LTrim operation
func testLTrim(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:ltrim"

	// Push elements
	cmd.RPush(ctx, key, "one", "two", "three", "four", "five")

	// Trim to keep only middle elements
	result := cmd.LTrim(ctx, key, 1, 3)
	require.NoError(t, result.Err())
	require.Equal(t, "OK", result.Val())

	// Verify
	rangeResult := cmd.LRange(ctx, key, 0, -1)
	require.NoError(t, rangeResult.Err())
	require.Equal(t, 3, len(rangeResult.Val()))
	require.Equal(t, []byte("two"), rangeResult.Val()[0])
	require.Equal(t, []byte("three"), rangeResult.Val()[1])
	require.Equal(t, []byte("four"), rangeResult.Val()[2])
}

// testRPopLPush tests RPopLPush operation
func testRPopLPush(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	source := "test:list:rpoplpush_src"
	dest := "test:list:rpoplpush_dest"

	// Setup source list
	cmd.RPush(ctx, source, "one", "two", "three")

	// Setup destination list
	cmd.RPush(ctx, dest, "a", "b")

	// Move element
	result := cmd.RPopLPush(ctx, source, dest)
	require.NoError(t, result.Err())
	require.Equal(t, []byte("three"), result.Val())

	// Verify source
	sourceResult := cmd.LRange(ctx, source, 0, -1)
	require.NoError(t, sourceResult.Err())
	require.Equal(t, 2, len(sourceResult.Val()))

	// Verify destination
	destResult := cmd.LRange(ctx, dest, 0, -1)
	require.NoError(t, destResult.Err())
	require.Equal(t, 3, len(destResult.Val()))
	require.Equal(t, []byte("three"), destResult.Val()[0])
	require.Equal(t, []byte("a"), destResult.Val()[1])
	require.Equal(t, []byte("b"), destResult.Val()[2])
}

// testRPopLPushSameKey tests RPopLPush on same key (rotation)
func testRPopLPushSameKey(t *testing.T, provider ListCommandProvider) {
	cmd := provider.GetListCommand()
	ctx := provider.GetContext()

	key := "test:list:rpoplpush_same"

	// Setup list
	cmd.RPush(ctx, key, "one", "two", "three")

	// Rotate
	result := cmd.RPopLPush(ctx, key, key)
	require.NoError(t, result.Err())
	require.Equal(t, []byte("three"), result.Val())

	// Verify rotation
	rangeResult := cmd.LRange(ctx, key, 0, -1)
	require.NoError(t, rangeResult.Err())
	require.Equal(t, []byte("three"), rangeResult.Val()[0])
	require.Equal(t, []byte("one"), rangeResult.Val()[1])
	require.Equal(t, []byte("two"), rangeResult.Val()[2])
}
