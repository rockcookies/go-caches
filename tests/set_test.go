package tests

import (
	"context"
	"testing"

	"github.com/rockcookies/go-caches"
	"github.com/stretchr/testify/require"
)

// SetCommandProvider defines the interface for testing SetCommand implementations
type SetCommandProvider interface {
	GetSetCommand() caches.SetCommand
	GetContext() context.Context
}

// RunSetCommandTests runs all SetCommand tests
func RunSetCommandTests(t *testing.T, provider SetCommandProvider) {
	t.Run("SAdd_and_SMembers", func(t *testing.T) {
		testSAddAndSMembers(t, provider)
	})
	t.Run("SAdd_Duplicates", func(t *testing.T) {
		testSAddDuplicates(t, provider)
	})
	t.Run("SCard", func(t *testing.T) {
		testSCard(t, provider)
	})
	t.Run("SCard_NonExistent", func(t *testing.T) {
		testSCardNonExistent(t, provider)
	})
	t.Run("SIsMember", func(t *testing.T) {
		testSIsMember(t, provider)
	})
	t.Run("SIsMember_NonExistent", func(t *testing.T) {
		testSIsMemberNonExistent(t, provider)
	})
	t.Run("SRem", func(t *testing.T) {
		testSRem(t, provider)
	})
	t.Run("SRem_NonExistent", func(t *testing.T) {
		testSRemNonExistent(t, provider)
	})
	t.Run("SPop", func(t *testing.T) {
		testSPop(t, provider)
	})
	t.Run("SPop_NonExistent", func(t *testing.T) {
		testSPopNonExistent(t, provider)
	})
	t.Run("SPopN", func(t *testing.T) {
		testSPopN(t, provider)
	})
	t.Run("SRandMember", func(t *testing.T) {
		testSRandMember(t, provider)
	})
	t.Run("SRandMemberN_Positive", func(t *testing.T) {
		testSRandMemberNPositive(t, provider)
	})
	t.Run("SRandMemberN_Negative", func(t *testing.T) {
		testSRandMemberNNegative(t, provider)
	})
	t.Run("SMove", func(t *testing.T) {
		testSMove(t, provider)
	})
	t.Run("SMove_NonExistent", func(t *testing.T) {
		testSMoveNonExistent(t, provider)
	})
	t.Run("SDiff", func(t *testing.T) {
		testSDiff(t, provider)
	})
	t.Run("SDiffStore", func(t *testing.T) {
		testSDiffStore(t, provider)
	})
	t.Run("SInter", func(t *testing.T) {
		testSInter(t, provider)
	})
	t.Run("SInterStore", func(t *testing.T) {
		testSInterStore(t, provider)
	})
	t.Run("SUnion", func(t *testing.T) {
		testSUnion(t, provider)
	})
	t.Run("SUnionStore", func(t *testing.T) {
		testSUnionStore(t, provider)
	})
	t.Run("SScan_Basic", func(t *testing.T) {
		testSScanBasic(t, provider)
	})
	t.Run("SScan_WithPattern", func(t *testing.T) {
		testSScanWithPattern(t, provider)
	})
	t.Run("SScan_NonExistent", func(t *testing.T) {
		testSScanNonExistent(t, provider)
	})
}

// testSAddAndSMembers tests SAdd and SMembers operations
func testSAddAndSMembers(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:sadd"

	// Add members
	result := cmd.SAdd(ctx, key, "one", "two", "three")
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val())

	// Get all members
	membersResult := cmd.SMembers(ctx, key)
	require.NoError(t, membersResult.Err())
	require.Equal(t, 3, len(membersResult.Val()))

	// Verify members (order may vary in sets)
	members := membersResult.Val()
	memberSet := make(map[string]bool)
	for _, m := range members {
		memberSet[string(m)] = true
	}
	require.True(t, memberSet["one"])
	require.True(t, memberSet["two"])
	require.True(t, memberSet["three"])
}

// testSAddDuplicates tests SAdd with duplicate members
func testSAddDuplicates(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:sadd_dup"

	// Add members
	result := cmd.SAdd(ctx, key, "one", "two")
	require.NoError(t, result.Err())
	require.Equal(t, int64(2), result.Val())

	// Add again with duplicates
	result2 := cmd.SAdd(ctx, key, "two", "three")
	require.NoError(t, result2.Err())
	require.Equal(t, int64(1), result2.Val()) // Only "three" is new

	// Verify total
	cardResult := cmd.SCard(ctx, key)
	require.NoError(t, cardResult.Err())
	require.Equal(t, int64(3), cardResult.Val())
}

// testSCard tests SCard operation
func testSCard(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:scard"

	// Add members
	cmd.SAdd(ctx, key, "one", "two", "three")

	// Get cardinality
	result := cmd.SCard(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val())
}

// testSCardNonExistent tests SCard on non-existent key
func testSCardNonExistent(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:scard_nonexistent"

	result := cmd.SCard(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(0), result.Val())
}

// testSIsMember tests SIsMember operation
func testSIsMember(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:sismember"

	// Add members
	cmd.SAdd(ctx, key, "one", "two", "three")

	// Check membership
	result := cmd.SIsMember(ctx, key, "two")
	require.NoError(t, result.Err())
	require.True(t, result.Val())

	result2 := cmd.SIsMember(ctx, key, "four")
	require.NoError(t, result2.Err())
	require.False(t, result2.Val())
}

// testSIsMemberNonExistent tests SIsMember on non-existent key
func testSIsMemberNonExistent(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:sismember_nonexistent"

	result := cmd.SIsMember(ctx, key, "one")
	require.NoError(t, result.Err())
	require.False(t, result.Val())
}

// testSRem tests SRem operation
func testSRem(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:srem"

	// Add members
	cmd.SAdd(ctx, key, "one", "two", "three", "four")

	// Remove members
	result := cmd.SRem(ctx, key, "two", "four")
	require.NoError(t, result.Err())
	require.Equal(t, int64(2), result.Val())

	// Verify remaining
	cardResult := cmd.SCard(ctx, key)
	require.NoError(t, cardResult.Err())
	require.Equal(t, int64(2), cardResult.Val())
}

// testSRemNonExistent tests SRem on non-existent members
func testSRemNonExistent(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:srem_nonexistent"

	// Add members
	cmd.SAdd(ctx, key, "one", "two")

	// Remove non-existent members
	result := cmd.SRem(ctx, key, "three", "four")
	require.NoError(t, result.Err())
	require.Equal(t, int64(0), result.Val())
}

// testSPop tests SPop operation
func testSPop(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:spop"

	// Add members
	cmd.SAdd(ctx, key, "one", "two", "three")

	// Pop a member
	result := cmd.SPop(ctx, key)
	require.NoError(t, result.Err())
	require.NotNil(t, result.Val())

	// Verify cardinality decreased
	cardResult := cmd.SCard(ctx, key)
	require.NoError(t, cardResult.Err())
	require.Equal(t, int64(2), cardResult.Val())
}

// testSPopNonExistent tests SPop on non-existent key
func testSPopNonExistent(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:spop_nonexistent"

	result := cmd.SPop(ctx, key)
	require.Equal(t, caches.Nil, result.Err())
	require.Nil(t, result.Val())
}

// testSPopN tests SPopN operation
func testSPopN(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:spopn"

	// Add members
	cmd.SAdd(ctx, key, "one", "two", "three", "four", "five")

	// Pop multiple members
	result := cmd.SPopN(ctx, key, 3)
	require.NoError(t, result.Err())
	require.Equal(t, 3, len(result.Val()))

	// Verify remaining
	cardResult := cmd.SCard(ctx, key)
	require.NoError(t, cardResult.Err())
	require.Equal(t, int64(2), cardResult.Val())
}

// testSRandMember tests SRandMember operation
func testSRandMember(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:srandmember"

	// Add members
	cmd.SAdd(ctx, key, "one", "two", "three")

	// Get random member
	result := cmd.SRandMember(ctx, key)
	require.NoError(t, result.Err())
	require.NotNil(t, result.Val())

	// Verify set unchanged
	cardResult := cmd.SCard(ctx, key)
	require.NoError(t, cardResult.Err())
	require.Equal(t, int64(3), cardResult.Val())
}

// testSRandMemberNPositive tests SRandMemberN with positive count
func testSRandMemberNPositive(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:srandmembern_pos"

	// Add members
	cmd.SAdd(ctx, key, "one", "two", "three", "four", "five")

	// Get random members (distinct)
	result := cmd.SRandMemberN(ctx, key, 3)
	require.NoError(t, result.Err())
	require.Equal(t, 3, len(result.Val()))

	// Verify all distinct
	seen := make(map[string]bool)
	for _, m := range result.Val() {
		require.False(t, seen[string(m)], "Expected distinct members")
		seen[string(m)] = true
	}
}

// testSRandMemberNNegative tests SRandMemberN with negative count
func testSRandMemberNNegative(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:srandmembern_neg"

	// Add members
	cmd.SAdd(ctx, key, "one", "two", "three")

	// Get random members (possibly repeated)
	result := cmd.SRandMemberN(ctx, key, -5)
	require.NoError(t, result.Err())
	require.Equal(t, 5, len(result.Val()))
}

// testSMove tests SMove operation
func testSMove(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	source := "test:set:smove_src"
	dest := "test:set:smove_dest"

	// Setup source
	cmd.SAdd(ctx, source, "one", "two", "three")

	// Setup destination
	cmd.SAdd(ctx, dest, "a", "b")

	// Move member
	result := cmd.SMove(ctx, source, dest, "two")
	require.NoError(t, result.Err())
	require.True(t, result.Val())

	// Verify source
	isMemberResult := cmd.SIsMember(ctx, source, "two")
	require.NoError(t, isMemberResult.Err())
	require.False(t, isMemberResult.Val())

	// Verify destination
	isMemberResult2 := cmd.SIsMember(ctx, dest, "two")
	require.NoError(t, isMemberResult2.Err())
	require.True(t, isMemberResult2.Val())
}

// testSMoveNonExistent tests SMove with non-existent member
func testSMoveNonExistent(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	source := "test:set:smove_src_ne"
	dest := "test:set:smove_dest_ne"

	// Setup source
	cmd.SAdd(ctx, source, "one", "two")

	// Move non-existent member
	result := cmd.SMove(ctx, source, dest, "three")
	require.NoError(t, result.Err())
	require.False(t, result.Val())
}

// testSDiff tests SDiff operation
func testSDiff(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key1 := "test:set:sdiff1"
	key2 := "test:set:sdiff2"
	key3 := "test:set:sdiff3"

	// Setup sets
	cmd.SAdd(ctx, key1, "a", "b", "c", "d")
	cmd.SAdd(ctx, key2, "c")
	cmd.SAdd(ctx, key3, "a", "c", "e")

	// Get difference
	result := cmd.SDiff(ctx, key1, key2, key3)
	require.NoError(t, result.Err())

	// Expected: key1 - key2 - key3 = {b, d}
	diff := result.Val()
	require.Equal(t, 2, len(diff))

	diffSet := make(map[string]bool)
	for _, m := range diff {
		diffSet[string(m)] = true
	}
	require.True(t, diffSet["b"])
	require.True(t, diffSet["d"])
}

// testSDiffStore tests SDiffStore operation
func testSDiffStore(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key1 := "test:set:sdiffstore1"
	key2 := "test:set:sdiffstore2"
	dest := "test:set:sdiffstore_dest"

	// Setup sets
	cmd.SAdd(ctx, key1, "a", "b", "c")
	cmd.SAdd(ctx, key2, "c", "d")

	// Store difference
	result := cmd.SDiffStore(ctx, dest, key1, key2)
	require.NoError(t, result.Err())
	require.Equal(t, int64(2), result.Val())

	// Verify destination
	membersResult := cmd.SMembers(ctx, dest)
	require.NoError(t, membersResult.Err())
	require.Equal(t, 2, len(membersResult.Val()))
}

// testSInter tests SInter operation
func testSInter(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key1 := "test:set:sinter1"
	key2 := "test:set:sinter2"
	key3 := "test:set:sinter3"

	// Setup sets
	cmd.SAdd(ctx, key1, "a", "b", "c", "d")
	cmd.SAdd(ctx, key2, "c", "d", "e")
	cmd.SAdd(ctx, key3, "a", "c", "e", "f")

	// Get intersection
	result := cmd.SInter(ctx, key1, key2, key3)
	require.NoError(t, result.Err())

	// Expected: {c}
	inter := result.Val()
	require.Equal(t, 1, len(inter))
	require.Equal(t, []byte("c"), inter[0])
}

// testSInterStore tests SInterStore operation
func testSInterStore(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key1 := "test:set:sinterstore1"
	key2 := "test:set:sinterstore2"
	dest := "test:set:sinterstore_dest"

	// Setup sets
	cmd.SAdd(ctx, key1, "a", "b", "c")
	cmd.SAdd(ctx, key2, "b", "c", "d")

	// Store intersection
	result := cmd.SInterStore(ctx, dest, key1, key2)
	require.NoError(t, result.Err())
	require.Equal(t, int64(2), result.Val())

	// Verify destination
	membersResult := cmd.SMembers(ctx, dest)
	require.NoError(t, membersResult.Err())
	require.Equal(t, 2, len(membersResult.Val()))
}

// testSUnion tests SUnion operation
func testSUnion(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key1 := "test:set:sunion1"
	key2 := "test:set:sunion2"
	key3 := "test:set:sunion3"

	// Setup sets
	cmd.SAdd(ctx, key1, "a", "b")
	cmd.SAdd(ctx, key2, "c", "d")
	cmd.SAdd(ctx, key3, "e")

	// Get union
	result := cmd.SUnion(ctx, key1, key2, key3)
	require.NoError(t, result.Err())

	// Expected: {a, b, c, d, e}
	union := result.Val()
	require.Equal(t, 5, len(union))
}

// testSUnionStore tests SUnionStore operation
func testSUnionStore(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key1 := "test:set:sunionstore1"
	key2 := "test:set:sunionstore2"
	dest := "test:set:sunionstore_dest"

	// Setup sets
	cmd.SAdd(ctx, key1, "a", "b")
	cmd.SAdd(ctx, key2, "c", "d")

	// Store union
	result := cmd.SUnionStore(ctx, dest, key1, key2)
	require.NoError(t, result.Err())
	require.Equal(t, int64(4), result.Val())

	// Verify destination
	membersResult := cmd.SMembers(ctx, dest)
	require.NoError(t, membersResult.Err())
	require.Equal(t, 4, len(membersResult.Val()))
}

// testSScanBasic tests SScan basic iteration
func testSScanBasic(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:sscan"

	// Add members
	cmd.SAdd(ctx, key, "member1", "member2", "member3", "member4", "member5")

	// Scan all members
	var allMembers [][]byte
	cursor := uint64(0)
	iterations := 0
	maxIterations := 10

	for iterations < maxIterations {
		result := cmd.SScan(ctx, key, cursor, "", 10)
		require.NoError(t, result.Err())

		scanResult := result.Val()
		allMembers = append(allMembers, scanResult.Elements...)
		cursor = scanResult.Cursor

		iterations++
		if cursor == 0 {
			break
		}
	}

	// Should have found all members
	require.GreaterOrEqual(t, len(allMembers), 5)
}

// testSScanWithPattern tests SScan with pattern matching
func testSScanWithPattern(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:sscan_pattern"

	// Add members with pattern
	cmd.SAdd(ctx, key, "foo1", "foo2", "bar1", "bar2")

	// Scan with pattern
	var allMembers [][]byte
	cursor := uint64(0)
	iterations := 0
	maxIterations := 10

	for iterations < maxIterations {
		result := cmd.SScan(ctx, key, cursor, "foo*", 10)
		require.NoError(t, result.Err())

		scanResult := result.Val()
		allMembers = append(allMembers, scanResult.Elements...)
		cursor = scanResult.Cursor

		iterations++
		if cursor == 0 {
			break
		}
	}

	// Should have found foo members
	for _, m := range allMembers {
		member := string(m)
		require.True(t, member == "foo1" || member == "foo2", "Expected foo* pattern match")
	}
}

// testSScanNonExistent tests SScan on non-existent key
func testSScanNonExistent(t *testing.T, provider SetCommandProvider) {
	cmd := provider.GetSetCommand()
	ctx := provider.GetContext()

	key := "test:set:sscan_nonexistent"

	result := cmd.SScan(ctx, key, 0, "", 10)
	require.NoError(t, result.Err())

	scanResult := result.Val()
	require.Empty(t, scanResult.Elements)
	require.Equal(t, uint64(0), scanResult.Cursor)
}
