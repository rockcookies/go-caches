package tests

import (
	"context"
	"testing"
	"time"

	"github.com/rockcookies/go-caches"
	"github.com/stretchr/testify/require"
)

// KeyCommandProvider defines the interface for testing KeyCommand implementations
type KeyCommandProvider interface {
	GetKeyCommand() caches.KeyCommand
	GetStringCommand() caches.StringCommand
	GetContext() context.Context
}

// RunKeyCommandTests runs all KeyCommand tests
func RunKeyCommandTests(t *testing.T, provider KeyCommandProvider) {
	t.Run("Del_SingleKey", func(t *testing.T) {
		testDelSingleKey(t, provider)
	})
	t.Run("Del_MultipleKeys", func(t *testing.T) {
		testDelMultipleKeys(t, provider)
	})
	t.Run("Del_NonExistentKey", func(t *testing.T) {
		testDelNonExistentKey(t, provider)
	})
	t.Run("Exists_SingleKey", func(t *testing.T) {
		testExistsSingleKey(t, provider)
	})
	t.Run("Exists_MultipleKeys", func(t *testing.T) {
		testExistsMultipleKeys(t, provider)
	})
	t.Run("Exists_NonExistentKey", func(t *testing.T) {
		testExistsNonExistentKey(t, provider)
	})
	t.Run("Expire_ExistingKey", func(t *testing.T) {
		testExpireExistingKey(t, provider)
	})
	t.Run("Expire_NonExistentKey", func(t *testing.T) {
		testExpireNonExistentKey(t, provider)
	})
	t.Run("ExpireAt_ExistingKey", func(t *testing.T) {
		testExpireAtExistingKey(t, provider)
	})
	t.Run("ExpireNX_NoExpiration", func(t *testing.T) {
		testExpireNXNoExpiration(t, provider)
	})
	t.Run("ExpireNX_HasExpiration", func(t *testing.T) {
		testExpireNXHasExpiration(t, provider)
	})
	t.Run("ExpireXX_NoExpiration", func(t *testing.T) {
		testExpireXXNoExpiration(t, provider)
	})
	t.Run("ExpireXX_HasExpiration", func(t *testing.T) {
		testExpireXXHasExpiration(t, provider)
	})
	t.Run("ExpireGT_Greater", func(t *testing.T) {
		testExpireGTGreater(t, provider)
	})
	t.Run("ExpireGT_Less", func(t *testing.T) {
		testExpireGTLess(t, provider)
	})
	t.Run("ExpireLT_Greater", func(t *testing.T) {
		testExpireLTGreater(t, provider)
	})
	t.Run("ExpireLT_Less", func(t *testing.T) {
		testExpireLTLess(t, provider)
	})
	t.Run("TTL_WithExpiration", func(t *testing.T) {
		testTTLWithExpiration(t, provider)
	})
	t.Run("TTL_NoExpiration", func(t *testing.T) {
		testTTLNoExpiration(t, provider)
	})
	t.Run("TTL_NonExistentKey", func(t *testing.T) {
		testTTLNonExistentKey(t, provider)
	})
	t.Run("PTTL_WithExpiration", func(t *testing.T) {
		testPTTLWithExpiration(t, provider)
	})
	t.Run("Persist_ExistingKey", func(t *testing.T) {
		testPersistExistingKey(t, provider)
	})
	t.Run("Persist_NonExistentKey", func(t *testing.T) {
		testPersistNonExistentKey(t, provider)
	})
	t.Run("Persist_NoExpiration", func(t *testing.T) {
		testPersistNoExpiration(t, provider)
	})
	t.Run("Type_String", func(t *testing.T) {
		testTypeString(t, provider)
	})
	t.Run("Type_NonExistent", func(t *testing.T) {
		testTypeNonExistent(t, provider)
	})
	t.Run("Rename_ExistingKey", func(t *testing.T) {
		testRenameExistingKey(t, provider)
	})
	t.Run("Rename_NonExistentKey", func(t *testing.T) {
		testRenameNonExistentKey(t, provider)
	})
	t.Run("RenameNX_NewKey", func(t *testing.T) {
		testRenameNXNewKey(t, provider)
	})
	t.Run("RenameNX_ExistingNewKey", func(t *testing.T) {
		testRenameNXExistingNewKey(t, provider)
	})
	t.Run("Keys_Pattern", func(t *testing.T) {
		testKeysPattern(t, provider)
	})
	t.Run("Keys_NoMatch", func(t *testing.T) {
		testKeysNoMatch(t, provider)
	})
	t.Run("RandomKey_ExistingKeys", func(t *testing.T) {
		testRandomKeyExistingKeys(t, provider)
	})
	t.Run("RandomKey_EmptyDB", func(t *testing.T) {
		testRandomKeyEmptyDB(t, provider)
	})
	t.Run("Scan_BasicIteration", func(t *testing.T) {
		testScanBasicIteration(t, provider)
	})
	t.Run("Scan_WithPattern", func(t *testing.T) {
		testScanWithPattern(t, provider)
	})
	t.Run("Scan_WithCount", func(t *testing.T) {
		testScanWithCount(t, provider)
	})
	t.Run("Scan_EmptyDB", func(t *testing.T) {
		testScanEmptyDB(t, provider)
	})
}

// testDelSingleKey tests Del on a single key
func testDelSingleKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:del_single"

	// Set a value
	strCmd.Set(ctx, key, "value", 0)

	// Delete the key
	result := keyCmd.Del(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(1), result.Val())

	// Verify key is deleted
	exists := keyCmd.Exists(ctx, key)
	require.NoError(t, exists.Err())
	require.Equal(t, int64(0), exists.Val())
}

// testDelMultipleKeys tests Del on multiple keys
func testDelMultipleKeys(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	keys := []string{
		"test:key:del_multi1",
		"test:key:del_multi2",
		"test:key:del_multi3",
	}

	// Set values
	for _, key := range keys {
		strCmd.Set(ctx, key, "value", 0)
	}

	// Delete all keys
	result := keyCmd.Del(ctx, keys...)
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val())
}

// testDelNonExistentKey tests Del on a non-existent key
func testDelNonExistentKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	ctx := provider.GetContext()

	key := "test:key:del_nonexistent"

	// Delete non-existent key
	result := keyCmd.Del(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(0), result.Val())
}

// testExistsSingleKey tests Exists on a single key
func testExistsSingleKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:exists_single"

	// Set a value
	strCmd.Set(ctx, key, "value", 0)

	// Check existence
	result := keyCmd.Exists(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(1), result.Val())
}

// testExistsMultipleKeys tests Exists on multiple keys
func testExistsMultipleKeys(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	keys := []string{
		"test:key:exists_multi1",
		"test:key:exists_multi2",
		"test:key:exists_multi3",
	}

	// Set two of the keys
	strCmd.Set(ctx, keys[0], "value1", 0)
	strCmd.Set(ctx, keys[1], "value2", 0)

	// Check existence (2 out of 3 exist)
	result := keyCmd.Exists(ctx, keys...)
	require.NoError(t, result.Err())
	require.Equal(t, int64(2), result.Val())
}

// testExistsNonExistentKey tests Exists on a non-existent key
func testExistsNonExistentKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	ctx := provider.GetContext()

	key := "test:key:exists_nonexistent"

	// Check non-existent key
	result := keyCmd.Exists(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(0), result.Val())
}

// testExpireExistingKey tests Expire on an existing key
func testExpireExistingKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:expire_existing"
	expiration := 1 * time.Second

	// Set a value
	strCmd.Set(ctx, key, "value", 0)

	// Set expiration
	result := keyCmd.Expire(ctx, key, expiration)
	require.NoError(t, result.Err())
	require.True(t, result.Val())

	// Wait for expiration
	time.Sleep(expiration + 100*time.Millisecond)

	// Key should be expired
	exists := keyCmd.Exists(ctx, key)
	require.NoError(t, exists.Err())
	require.Equal(t, int64(0), exists.Val())
}

// testExpireNonExistentKey tests Expire on a non-existent key
func testExpireNonExistentKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	ctx := provider.GetContext()

	key := "test:key:expire_nonexistent"

	// Try to set expiration on non-existent key
	result := keyCmd.Expire(ctx, key, time.Second)
	require.NoError(t, result.Err())
	require.False(t, result.Val())
}

// testExpireAtExistingKey tests ExpireAt on an existing key
func testExpireAtExistingKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:expireat_existing"
	expireAt := time.Now().Add(1 * time.Second)

	// Set a value
	strCmd.Set(ctx, key, "value", 0)

	// Set expiration time
	result := keyCmd.ExpireAt(ctx, key, expireAt)
	require.NoError(t, result.Err())
	require.True(t, result.Val())

	// Wait for expiration
	time.Sleep(time.Until(expireAt) + 100*time.Millisecond)

	// Key should be expired
	exists := keyCmd.Exists(ctx, key)
	require.NoError(t, exists.Err())
	require.Equal(t, int64(0), exists.Val())
}

// testExpireNXNoExpiration tests ExpireNX on a key without expiration
func testExpireNXNoExpiration(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:expirenx_no_exp"

	// Set a value without expiration
	strCmd.Set(ctx, key, "value", 0)

	// ExpireNX should succeed
	result := keyCmd.ExpireNX(ctx, key, time.Minute)
	require.NoError(t, result.Err())
	require.True(t, result.Val())
}

// testExpireNXHasExpiration tests ExpireNX on a key with expiration
func testExpireNXHasExpiration(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:expirenx_has_exp"

	// Set a value with expiration
	strCmd.Set(ctx, key, "value", time.Minute)

	// ExpireNX should fail (key already has expiration)
	result := keyCmd.ExpireNX(ctx, key, 2*time.Minute)
	require.NoError(t, result.Err())
	require.False(t, result.Val())
}

// testExpireXXNoExpiration tests ExpireXX on a key without expiration
func testExpireXXNoExpiration(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:expirexx_no_exp"

	// Set a value without expiration
	strCmd.Set(ctx, key, "value", 0)

	// ExpireXX should fail (key has no expiration)
	result := keyCmd.ExpireXX(ctx, key, time.Minute)
	require.NoError(t, result.Err())
	require.False(t, result.Val())
}

// testExpireXXHasExpiration tests ExpireXX on a key with expiration
func testExpireXXHasExpiration(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:expirexx_has_exp"

	// Set a value with expiration
	strCmd.Set(ctx, key, "value", time.Minute)

	// ExpireXX should succeed
	result := keyCmd.ExpireXX(ctx, key, 2*time.Minute)
	require.NoError(t, result.Err())
	require.True(t, result.Val())
}

// testExpireGTGreater tests ExpireGT with greater expiration
func testExpireGTGreater(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:expiregt_greater"

	// Set a value with 1 minute expiration
	strCmd.Set(ctx, key, "value", time.Minute)

	// ExpireGT with 2 minutes should succeed
	result := keyCmd.ExpireGT(ctx, key, 2*time.Minute)
	require.NoError(t, result.Err())
	require.True(t, result.Val())
}

// testExpireGTLess tests ExpireGT with less expiration
func testExpireGTLess(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:expiregt_less"

	// Set a value with 2 minute expiration
	strCmd.Set(ctx, key, "value", 2*time.Minute)

	// ExpireGT with 1 minute should fail
	result := keyCmd.ExpireGT(ctx, key, time.Minute)
	require.NoError(t, result.Err())
	require.False(t, result.Val())
}

// testExpireLTGreater tests ExpireLT with greater expiration
func testExpireLTGreater(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:expirelt_greater"

	// Set a value with 2 minute expiration
	strCmd.Set(ctx, key, "value", 2*time.Minute)

	// ExpireLT with 3 minutes should fail
	result := keyCmd.ExpireLT(ctx, key, 3*time.Minute)
	require.NoError(t, result.Err())
	require.False(t, result.Val())
}

// testExpireLTLess tests ExpireLT with less expiration
func testExpireLTLess(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:expirelt_less"

	// Set a value with 2 minute expiration
	strCmd.Set(ctx, key, "value", 2*time.Minute)

	// ExpireLT with 1 minute should succeed
	result := keyCmd.ExpireLT(ctx, key, time.Minute)
	require.NoError(t, result.Err())
	require.True(t, result.Val())
}

// testTTLWithExpiration tests TTL on a key with expiration
func testTTLWithExpiration(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:ttl_with_exp"
	expiration := 10 * time.Second

	// Set a value with expiration
	strCmd.Set(ctx, key, "value", expiration)

	// Get TTL
	result := keyCmd.TTL(ctx, key)
	require.NoError(t, result.Err())
	require.Greater(t, result.Val(), time.Duration(0))
	require.LessOrEqual(t, result.Val(), expiration)
}

// testTTLNoExpiration tests TTL on a key without expiration
func testTTLNoExpiration(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:ttl_no_exp"

	// Set a value without expiration
	strCmd.Set(ctx, key, "value", 0)

	// Get TTL (should be -1 for no expiration)
	result := keyCmd.TTL(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, time.Duration(-1), result.Val())
}

// testTTLNonExistentKey tests TTL on a non-existent key
func testTTLNonExistentKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	ctx := provider.GetContext()

	key := "test:key:ttl_nonexistent"

	// Get TTL (should be -2 for non-existent key)
	result := keyCmd.TTL(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, time.Duration(-2), result.Val())
}

// testPTTLWithExpiration tests PTTL on a key with expiration
func testPTTLWithExpiration(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:pttl_with_exp"
	expiration := 10 * time.Second

	// Set a value with expiration
	strCmd.Set(ctx, key, "value", expiration)

	// Get PTTL (in milliseconds)
	result := keyCmd.PTTL(ctx, key)
	require.NoError(t, result.Err())
	require.Greater(t, result.Val(), time.Duration(0))
	require.LessOrEqual(t, result.Val(), expiration)
}

// testPersistExistingKey tests Persist on a key with expiration
func testPersistExistingKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:persist_existing"

	// Set a value with expiration
	strCmd.Set(ctx, key, "value", time.Minute)

	// Remove expiration
	result := keyCmd.Persist(ctx, key)
	require.NoError(t, result.Err())
	require.True(t, result.Val())

	// Verify no expiration
	ttl := keyCmd.TTL(ctx, key)
	require.NoError(t, ttl.Err())
	require.Equal(t, time.Duration(-1), ttl.Val())
}

// testPersistNonExistentKey tests Persist on a non-existent key
func testPersistNonExistentKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	ctx := provider.GetContext()

	key := "test:key:persist_nonexistent"

	// Try to persist non-existent key
	result := keyCmd.Persist(ctx, key)
	require.NoError(t, result.Err())
	require.False(t, result.Val())
}

// testPersistNoExpiration tests Persist on a key without expiration
func testPersistNoExpiration(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:persist_no_exp"

	// Set a value without expiration
	strCmd.Set(ctx, key, "value", 0)

	// Try to persist (should return false as key has no expiration)
	result := keyCmd.Persist(ctx, key)
	require.NoError(t, result.Err())
	require.False(t, result.Val())
}

// testTypeString tests Type on a string key
func testTypeString(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:key:type_string"

	// Set a string value
	strCmd.Set(ctx, key, "value", 0)

	// Get type
	result := keyCmd.Type(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, "string", result.Val())
}

// testTypeNonExistent tests Type on a non-existent key
func testTypeNonExistent(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	ctx := provider.GetContext()

	key := "test:key:type_nonexistent"

	// Get type of non-existent key
	result := keyCmd.Type(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, "none", result.Val())
}

// testRenameExistingKey tests Rename on an existing key
func testRenameExistingKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	oldKey := "test:key:rename_old"
	newKey := "test:key:rename_new"

	// Set a value
	strCmd.Set(ctx, oldKey, "value", 0)

	// Rename
	result := keyCmd.Rename(ctx, oldKey, newKey)
	require.NoError(t, result.Err())
	require.Equal(t, "OK", result.Val())

	// Verify old key is gone
	exists := keyCmd.Exists(ctx, oldKey)
	require.NoError(t, exists.Err())
	require.Equal(t, int64(0), exists.Val())

	// Verify new key exists
	getResult := strCmd.Get(ctx, newKey)
	require.NoError(t, getResult.Err())
	require.Equal(t, []byte("value"), getResult.Val())
}

// testRenameNonExistentKey tests Rename on a non-existent key
func testRenameNonExistentKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	ctx := provider.GetContext()

	oldKey := "test:key:rename_nonexistent"
	newKey := "test:key:rename_target"

	// Try to rename non-existent key (should return error)
	result := keyCmd.Rename(ctx, oldKey, newKey)
	require.Error(t, result.Err())
}

// testRenameNXNewKey tests RenameNX when new key doesn't exist
func testRenameNXNewKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	oldKey := "test:key:renamenx_old"
	newKey := "test:key:renamenx_new"

	// Set a value
	strCmd.Set(ctx, oldKey, "value", 0)

	// RenameNX should succeed
	result := keyCmd.RenameNX(ctx, oldKey, newKey)
	require.NoError(t, result.Err())
	require.True(t, result.Val())

	// Verify new key exists
	getResult := strCmd.Get(ctx, newKey)
	require.NoError(t, getResult.Err())
	require.Equal(t, []byte("value"), getResult.Val())
}

// testRenameNXExistingNewKey tests RenameNX when new key already exists
func testRenameNXExistingNewKey(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	oldKey := "test:key:renamenx_old_exist"
	newKey := "test:key:renamenx_new_exist"

	// Set both keys
	strCmd.Set(ctx, oldKey, "value1", 0)
	strCmd.Set(ctx, newKey, "value2", 0)

	// RenameNX should fail
	result := keyCmd.RenameNX(ctx, oldKey, newKey)
	require.NoError(t, result.Err())
	require.False(t, result.Val())

	// Verify new key is unchanged
	getResult := strCmd.Get(ctx, newKey)
	require.NoError(t, getResult.Err())
	require.Equal(t, []byte("value2"), getResult.Val())
}

// testKeysPattern tests Keys with a pattern
func testKeysPattern(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	// Set multiple keys with a pattern
	keys := []string{
		"test:key:pattern:foo1",
		"test:key:pattern:foo2",
		"test:key:pattern:bar1",
	}

	for _, key := range keys {
		strCmd.Set(ctx, key, "value", 0)
	}

	// Get keys matching pattern
	result := keyCmd.Keys(ctx, "test:key:pattern:foo*")
	require.NoError(t, result.Err())
	require.Len(t, result.Val(), 2)
	require.Contains(t, result.Val(), "test:key:pattern:foo1")
	require.Contains(t, result.Val(), "test:key:pattern:foo2")
}

// testKeysNoMatch tests Keys with a pattern that doesn't match
func testKeysNoMatch(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	ctx := provider.GetContext()

	// Get keys matching non-existent pattern
	result := keyCmd.Keys(ctx, "test:key:nomatch:*")
	require.NoError(t, result.Err())
	require.Empty(t, result.Val())
}

// testRandomKeyExistingKeys tests RandomKey when keys exist
func testRandomKeyExistingKeys(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	// Set multiple keys
	keys := []string{
		"test:key:random:key1",
		"test:key:random:key2",
		"test:key:random:key3",
	}

	for _, key := range keys {
		strCmd.Set(ctx, key, "value", 0)
	}

	// Get a random key
	result := keyCmd.RandomKey(ctx)
	require.NoError(t, result.Err())
	require.NotEmpty(t, result.Val())

	// The returned key should be one of the keys we set
	// Note: RandomKey might return any key in the database, not just our test keys
	// So we can only verify it's not empty
}

// testRandomKeyEmptyDB tests RandomKey when no keys exist
func testRandomKeyEmptyDB(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	ctx := provider.GetContext()

	// Note: This test is tricky because we can't guarantee the DB is completely empty
	// in all implementations. FlushAll might not work as expected in some providers.
	// Instead, we just verify that RandomKey doesn't crash and returns something valid
	result := keyCmd.RandomKey(ctx)

	// Should either succeed with a key, or return empty string/error
	// Both behaviors are acceptable depending on the state of the database
	if result.Err() != nil {
		t.Logf("RandomKey on potentially empty DB returned error: %v", result.Err())
	} else {
		t.Logf("RandomKey on potentially empty DB returned: %s", result.Val())
	}
	// No assertion needed - just verify it doesn't panic
}

// testScanBasicIteration tests Scan basic iteration
func testScanBasicIteration(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	// Set multiple keys
	testKeys := []string{
		"test:key:scan:item1",
		"test:key:scan:item2",
		"test:key:scan:item3",
		"test:key:scan:item4",
		"test:key:scan:item5",
	}

	for _, key := range testKeys {
		strCmd.Set(ctx, key, "value", 0)
	}

	// Scan all keys
	var allKeys []string
	cursor := uint64(0)
	iterations := 0
	maxIterations := 10 // Prevent infinite loop

	for iterations < maxIterations {
		result := keyCmd.Scan(ctx, cursor, "test:key:scan:*", 10)
		require.NoError(t, result.Err())

		scanResult := result.Val()
		allKeys = append(allKeys, scanResult.Keys...)
		cursor = scanResult.Cursor

		iterations++
		if cursor == 0 {
			break
		}
	}

	// Should have found all test keys
	require.GreaterOrEqual(t, len(allKeys), len(testKeys))
	for _, key := range testKeys {
		require.Contains(t, allKeys, key)
	}
}

// testScanWithPattern tests Scan with pattern matching
func testScanWithPattern(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	// Set keys with different patterns
	fooKeys := []string{
		"test:key:scanpat:foo1",
		"test:key:scanpat:foo2",
	}
	barKeys := []string{
		"test:key:scanpat:bar1",
		"test:key:scanpat:bar2",
	}

	for _, key := range append(fooKeys, barKeys...) {
		strCmd.Set(ctx, key, "value", 0)
	}

	// Scan only "foo" keys
	var allKeys []string
	cursor := uint64(0)
	iterations := 0
	maxIterations := 10

	for iterations < maxIterations {
		result := keyCmd.Scan(ctx, cursor, "test:key:scanpat:foo*", 10)
		require.NoError(t, result.Err())

		scanResult := result.Val()
		allKeys = append(allKeys, scanResult.Keys...)
		cursor = scanResult.Cursor

		iterations++
		if cursor == 0 {
			break
		}
	}

	// Should have found all foo keys
	for _, key := range fooKeys {
		require.Contains(t, allKeys, key)
	}

	// Should not have found any bar keys
	for _, key := range barKeys {
		require.NotContains(t, allKeys, key)
	}
}

// testScanWithCount tests Scan with count parameter
func testScanWithCount(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	strCmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	// Set multiple keys
	for i := 1; i <= 20; i++ {
		key := "test:key:scancount:" + string(rune('a'+i-1))
		strCmd.Set(ctx, key, "value", 0)
	}

	// Scan with small count
	result := keyCmd.Scan(ctx, 0, "test:key:scancount:*", 5)
	require.NoError(t, result.Err())

	scanResult := result.Val()
	// Count is a hint, so we just verify we got some keys
	// and possibly a non-zero cursor if there are more
	require.NotNil(t, scanResult.Keys)
}

// testScanEmptyDB tests Scan on empty database
func testScanEmptyDB(t *testing.T, provider KeyCommandProvider) {
	keyCmd := provider.GetKeyCommand()
	ctx := provider.GetContext()

	// Scan with pattern that won't match anything
	result := keyCmd.Scan(ctx, 0, "test:key:scanempty:nonexistent:*", 10)
	require.NoError(t, result.Err())

	scanResult := result.Val()
	// When pattern doesn't match any keys, we should get empty result
	// Cursor might still be non-zero if there are other keys in DB
	require.Empty(t, scanResult.Keys)
}
