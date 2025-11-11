package tests

import (
	"context"
	"testing"
	"time"

	"github.com/rockcookies/go-caches"
	"github.com/stretchr/testify/require"
)

// StringCommandProvider defines the interface for testing StringCommand implementations
type StringCommandProvider interface {
	GetStringCommand() caches.StringCommand
	GetContext() context.Context
}

// RunStringCommandTests runs all StringCommand tests
func RunStringCommandTests(t *testing.T, provider StringCommandProvider) {
	t.Run("Set_and_Get", func(t *testing.T) {
		testSetAndGet(t, provider)
	})
	t.Run("Get_NonExistent", func(t *testing.T) {
		testGetNonExistent(t, provider)
	})
	t.Run("Set_WithExpiration", func(t *testing.T) {
		testSetWithExpiration(t, provider)
	})
	t.Run("SetNX_NewKey", func(t *testing.T) {
		testSetNXNewKey(t, provider)
	})
	t.Run("SetNX_ExistingKey", func(t *testing.T) {
		testSetNXExistingKey(t, provider)
	})
	t.Run("SetXX_ExistingKey", func(t *testing.T) {
		testSetXXExistingKey(t, provider)
	})
	t.Run("SetXX_NonExistentKey", func(t *testing.T) {
		testSetXXNonExistentKey(t, provider)
	})
	t.Run("SetArgs_ModeNX", func(t *testing.T) {
		testSetArgsModeNX(t, provider)
	})
	t.Run("SetArgs_ModeXX", func(t *testing.T) {
		testSetArgsModeXX(t, provider)
	})
	t.Run("SetArgs_WithGet", func(t *testing.T) {
		testSetArgsWithGet(t, provider)
	})
	t.Run("SetArgs_WithTTL", func(t *testing.T) {
		testSetArgsWithTTL(t, provider)
	})
	t.Run("SetArgs_WithExpireAt", func(t *testing.T) {
		testSetArgsWithExpireAt(t, provider)
	})
	t.Run("SetArgs_KeepTTL", func(t *testing.T) {
		testSetArgsKeepTTL(t, provider)
	})
	t.Run("Incr_NewKey", func(t *testing.T) {
		testIncrNewKey(t, provider)
	})
	t.Run("Incr_ExistingKey", func(t *testing.T) {
		testIncrExistingKey(t, provider)
	})
	t.Run("IncrBy", func(t *testing.T) {
		testIncrBy(t, provider)
	})
	t.Run("Decr_NewKey", func(t *testing.T) {
		testDecrNewKey(t, provider)
	})
	t.Run("Decr_ExistingKey", func(t *testing.T) {
		testDecrExistingKey(t, provider)
	})
	t.Run("DecrBy", func(t *testing.T) {
		testDecrBy(t, provider)
	})
	t.Run("IncrByFloat", func(t *testing.T) {
		testIncrByFloat(t, provider)
	})
	t.Run("StrLen_ExistingKey", func(t *testing.T) {
		testStrLenExistingKey(t, provider)
	})
	t.Run("StrLen_NonExistentKey", func(t *testing.T) {
		testStrLenNonExistentKey(t, provider)
	})
	t.Run("MSet", func(t *testing.T) {
		testMSet(t, provider)
	})
	t.Run("MGet_AllExist", func(t *testing.T) {
		testMGetAllExist(t, provider)
	})
	t.Run("MGet_PartialExist", func(t *testing.T) {
		testMGetPartialExist(t, provider)
	})
	t.Run("MGet_NoneExist", func(t *testing.T) {
		testMGetNoneExist(t, provider)
	})
	t.Run("MSetNX_AllNew", func(t *testing.T) {
		testMSetNXAllNew(t, provider)
	})
	t.Run("MSetNX_SomeExist", func(t *testing.T) {
		testMSetNXSomeExist(t, provider)
	})
}

// testSetAndGet tests basic Set and Get operations
func testSetAndGet(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:set_get"
	value := []byte("hello world")

	// Set value
	result := cmd.Set(ctx, key, value, 0)
	require.NoError(t, result.Err())
	require.Equal(t, "OK", result.Val())

	// Get value
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, value, getResult.Val())
}

// testGetNonExistent tests Get on a non-existent key
func testGetNonExistent(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:nonexistent"

	result := cmd.Get(ctx, key)
	require.Equal(t, caches.Nil, result.Err())
	require.Nil(t, result.Val())
}

// testSetWithExpiration tests Set with expiration
func testSetWithExpiration(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:expiration"
	value := []byte("expires soon")
	expiration := 100 * time.Millisecond

	// Set with expiration
	result := cmd.Set(ctx, key, value, expiration)
	require.NoError(t, result.Err())

	// Immediately verify the key exists
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, value, getResult.Val())

	// Wait for expiration
	time.Sleep(expiration + 50*time.Millisecond)

	// Key should be expired
	expiredResult := cmd.Get(ctx, key)
	require.Equal(t, caches.Nil, expiredResult.Err())
}

// testSetNXNewKey tests SetNX on a new key
func testSetNXNewKey(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:setnx_new"
	value := []byte("new value")

	result := cmd.SetNX(ctx, key, value, 0)
	require.NoError(t, result.Err())
	require.True(t, result.Val(), "SetNX should return true for new key")

	// Verify value was set
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, value, getResult.Val())
}

// testSetNXExistingKey tests SetNX on an existing key
func testSetNXExistingKey(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:setnx_existing"
	value1 := []byte("first value")
	value2 := []byte("second value")

	// Set initial value
	cmd.Set(ctx, key, value1, 0)

	// Try SetNX on existing key
	result := cmd.SetNX(ctx, key, value2, 0)
	require.NoError(t, result.Err())
	require.False(t, result.Val(), "SetNX should return false for existing key")

	// Verify original value is unchanged
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, value1, getResult.Val())
}

// testSetXXExistingKey tests SetXX on an existing key
func testSetXXExistingKey(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:setxx_existing"
	value1 := []byte("original value")
	value2 := []byte("updated value")

	// Set initial value
	cmd.Set(ctx, key, value1, 0)

	// Update with SetXX
	result := cmd.SetXX(ctx, key, value2, 0)
	require.NoError(t, result.Err())
	require.True(t, result.Val(), "SetXX should return true for existing key")

	// Verify value was updated
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, value2, getResult.Val())
}

// testSetXXNonExistentKey tests SetXX on a non-existent key
func testSetXXNonExistentKey(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:setxx_nonexistent"
	value := []byte("new value")

	result := cmd.SetXX(ctx, key, value, 0)
	require.NoError(t, result.Err())
	require.False(t, result.Val(), "SetXX should return false for non-existent key")

	// Verify key was not created
	getResult := cmd.Get(ctx, key)
	require.Equal(t, caches.Nil, getResult.Err())
}

// testSetArgsModeNX tests SetArgs with Mode="NX"
func testSetArgsModeNX(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:setargs_nx"
	value1 := []byte("first")
	value2 := []byte("second")

	// First SetArgs with NX should succeed
	args1 := caches.SetArgs{Mode: "NX"}
	result1 := cmd.SetArgs(ctx, key, value1, args1)
	require.NoError(t, result1.Err())
	require.Equal(t, "OK", result1.Val())

	// Second SetArgs with NX should fail (returns Nil error when key exists)
	args2 := caches.SetArgs{Mode: "NX"}
	result2 := cmd.SetArgs(ctx, key, value2, args2)
	require.Equal(t, caches.Nil, result2.Err(), "SetArgs with NX should return Nil error when key exists")

	// Verify original value is unchanged
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, value1, getResult.Val())
}

// testSetArgsModeXX tests SetArgs with Mode="XX"
func testSetArgsModeXX(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:setargs_xx"
	value1 := []byte("initial")
	value2 := []byte("updated")

	// SetArgs with XX on non-existent key should fail (returns Nil error when key does not exist)
	args1 := caches.SetArgs{Mode: "XX"}
	result1 := cmd.SetArgs(ctx, key, value1, args1)
	require.Equal(t, caches.Nil, result1.Err(), "SetArgs with XX should return Nil error when key does not exist")

	// Set initial value
	cmd.Set(ctx, key, value1, 0)

	// SetArgs with XX on existing key should succeed
	args2 := caches.SetArgs{Mode: "XX"}
	result2 := cmd.SetArgs(ctx, key, value2, args2)
	require.NoError(t, result2.Err())
	require.Equal(t, "OK", result2.Val())

	// Verify value was updated
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, value2, getResult.Val())
}

// testSetArgsWithGet tests SetArgs with Get=true
func testSetArgsWithGet(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:setargs_get"
	value1 := []byte("old value")
	value2 := []byte("new value")

	// Set initial value
	cmd.Set(ctx, key, value1, 0)

	// SetArgs with Get should return old value
	args := caches.SetArgs{Get: true}
	result := cmd.SetArgs(ctx, key, value2, args)
	require.NoError(t, result.Err())
	require.Equal(t, string(value1), result.Val())

	// Verify new value was set
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, value2, getResult.Val())
}

// testSetArgsWithTTL tests SetArgs with TTL
func testSetArgsWithTTL(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:setargs_ttl"
	value := []byte("temporary")
	ttl := 100 * time.Millisecond

	// Set with TTL
	args := caches.SetArgs{TTL: ttl}
	result := cmd.SetArgs(ctx, key, value, args)
	require.NoError(t, result.Err())

	// Immediately verify the key exists
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, value, getResult.Val())

	// Wait for expiration
	time.Sleep(ttl + 50*time.Millisecond)

	// Key should be expired
	expiredResult := cmd.Get(ctx, key)
	require.Equal(t, caches.Nil, expiredResult.Err())
}

// testSetArgsWithExpireAt tests SetArgs with ExpireAt
func testSetArgsWithExpireAt(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:setargs_expireat"
	value := []byte("expires at time")
	expireAt := time.Now().Add(1 * time.Second)

	// Set with ExpireAt
	args := caches.SetArgs{ExpireAt: expireAt}
	result := cmd.SetArgs(ctx, key, value, args)
	require.NoError(t, result.Err())

	// Immediately verify the key exists
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, value, getResult.Val())

	// Wait for expiration
	time.Sleep(time.Until(expireAt) + 50*time.Millisecond)

	// Key should be expired
	expiredResult := cmd.Get(ctx, key)
	require.Equal(t, caches.Nil, expiredResult.Err())
}

// testSetArgsKeepTTL tests SetArgs with KeepTTL
func testSetArgsKeepTTL(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:setargs_keepttl"
	value1 := []byte("first with ttl")
	value2 := []byte("second keeping ttl")
	ttl := 300 * time.Millisecond

	// Set initial value with TTL
	cmd.Set(ctx, key, value1, ttl)

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Update value with KeepTTL
	args := caches.SetArgs{KeepTTL: true}
	result := cmd.SetArgs(ctx, key, value2, args)
	require.NoError(t, result.Err())

	// Verify new value was set
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, value2, getResult.Val())

	// The key should still expire at the original time
	// Wait for remaining TTL
	time.Sleep(250*time.Millisecond + 50*time.Millisecond)

	// Key should be expired
	expiredResult := cmd.Get(ctx, key)
	require.Equal(t, caches.Nil, expiredResult.Err())
}

// testIncrNewKey tests Incr on a new key
func testIncrNewKey(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:incr_new"

	// Incr on new key should return 1
	result := cmd.Incr(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(1), result.Val())

	// Verify value
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, "1", string(getResult.Val()))
}

// testIncrExistingKey tests Incr on an existing key
func testIncrExistingKey(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:incr_existing"

	// Set initial value
	cmd.Set(ctx, key, "10", 0)

	// Incr should increment the value
	result := cmd.Incr(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(11), result.Val())

	// Verify value
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, "11", string(getResult.Val()))
}

// testIncrBy tests IncrBy operation
func testIncrBy(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:incrby"

	// Set initial value
	cmd.Set(ctx, key, "100", 0)

	// IncrBy 50
	result := cmd.IncrBy(ctx, key, 50)
	require.NoError(t, result.Err())
	require.Equal(t, int64(150), result.Val())

	// IncrBy -30 (decrement)
	result2 := cmd.IncrBy(ctx, key, -30)
	require.NoError(t, result2.Err())
	require.Equal(t, int64(120), result2.Val())
}

// testDecrNewKey tests Decr on a new key
func testDecrNewKey(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:decr_new"

	// Decr on new key should return -1
	result := cmd.Decr(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(-1), result.Val())

	// Verify value
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, "-1", string(getResult.Val()))
}

// testDecrExistingKey tests Decr on an existing key
func testDecrExistingKey(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:decr_existing"

	// Set initial value
	cmd.Set(ctx, key, "10", 0)

	// Decr should decrement the value
	result := cmd.Decr(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(9), result.Val())

	// Verify value
	getResult := cmd.Get(ctx, key)
	require.NoError(t, getResult.Err())
	require.Equal(t, "9", string(getResult.Val()))
}

// testDecrBy tests DecrBy operation
func testDecrBy(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:decrby"

	// Set initial value
	cmd.Set(ctx, key, "100", 0)

	// DecrBy 30
	result := cmd.DecrBy(ctx, key, 30)
	require.NoError(t, result.Err())
	require.Equal(t, int64(70), result.Val())

	// DecrBy -20 (increment)
	result2 := cmd.DecrBy(ctx, key, -20)
	require.NoError(t, result2.Err())
	require.Equal(t, int64(90), result2.Val())
}

// testIncrByFloat tests IncrByFloat operation
func testIncrByFloat(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:incrbyfloat"

	// Set initial value
	cmd.Set(ctx, key, "10.5", 0)

	// IncrByFloat 2.1
	result := cmd.IncrByFloat(ctx, key, 2.1)
	require.NoError(t, result.Err())
	require.InDelta(t, 12.6, result.Val(), 0.0001)

	// IncrByFloat -0.1
	result2 := cmd.IncrByFloat(ctx, key, -0.1)
	require.NoError(t, result2.Err())
	require.InDelta(t, 12.5, result2.Val(), 0.0001)
}

// testStrLenExistingKey tests StrLen on an existing key
func testStrLenExistingKey(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:strlen_existing"
	value := []byte("hello world")

	// Set value
	cmd.Set(ctx, key, value, 0)

	// Get string length
	result := cmd.StrLen(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(len(value)), result.Val())
}

// testStrLenNonExistentKey tests StrLen on a non-existent key
func testStrLenNonExistentKey(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	key := "test:string:strlen_nonexistent"

	// StrLen on non-existent key should return 0
	result := cmd.StrLen(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(0), result.Val())
}

// testMSet tests MSet operation
func testMSet(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	values := map[string]any{
		"test:string:mset1": []byte("value1"),
		"test:string:mset2": []byte("value2"),
		"test:string:mset3": []byte("value3"),
	}

	// MSet multiple key-value pairs
	result := cmd.MSet(ctx, values)
	require.NoError(t, result.Err())
	require.Equal(t, "OK", result.Val())

	// Verify all values were set
	for key, expectedValue := range values {
		getResult := cmd.Get(ctx, key)
		require.NoError(t, getResult.Err())
		require.Equal(t, expectedValue, getResult.Val())
	}
}

// testMGetAllExist tests MGet when all keys exist
func testMGetAllExist(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	// Set multiple values
	keys := []string{
		"test:string:mget1",
		"test:string:mget2",
		"test:string:mget3",
	}
	values := map[string]any{
		keys[0]: []byte("value1"),
		keys[1]: []byte("value2"),
		keys[2]: []byte("value3"),
	}
	cmd.MSet(ctx, values)

	// MGet all keys
	result := cmd.MGet(ctx, keys...)
	require.NoError(t, result.Err())
	require.Equal(t, len(keys), len(result.Val()))

	// Verify all values
	resultMap := result.Val()
	for _, key := range keys {
		require.Equal(t, values[key], resultMap[key])
	}
}

// testMGetPartialExist tests MGet when some keys exist
func testMGetPartialExist(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	// Set some values
	cmd.Set(ctx, "test:string:mget_partial1", []byte("value1"), 0)
	cmd.Set(ctx, "test:string:mget_partial3", []byte("value3"), 0)

	keys := []string{
		"test:string:mget_partial1",
		"test:string:mget_partial2", // non-existent
		"test:string:mget_partial3",
	}

	// MGet with mixed existing/non-existing keys
	result := cmd.MGet(ctx, keys...)
	require.NoError(t, result.Err())

	resultMap := result.Val()
	require.Equal(t, []byte("value1"), resultMap["test:string:mget_partial1"])
	require.Nil(t, resultMap["test:string:mget_partial2"])
	require.Equal(t, []byte("value3"), resultMap["test:string:mget_partial3"])
}

// testMGetNoneExist tests MGet when no keys exist
func testMGetNoneExist(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	keys := []string{
		"test:string:mget_none1",
		"test:string:mget_none2",
		"test:string:mget_none3",
	}

	// MGet non-existent keys
	result := cmd.MGet(ctx, keys...)
	require.NoError(t, result.Err())

	resultMap := result.Val()
	for _, key := range keys {
		require.Nil(t, resultMap[key])
	}
}

// testMSetNXAllNew tests MSetNX when all keys are new
func testMSetNXAllNew(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	values := map[string]any{
		"test:string:msetnx_new1": []byte("value1"),
		"test:string:msetnx_new2": []byte("value2"),
		"test:string:msetnx_new3": []byte("value3"),
	}

	// MSetNX should succeed when all keys are new
	result := cmd.MSetNX(ctx, values)
	require.NoError(t, result.Err())
	require.True(t, result.Val(), "MSetNX should return true when all keys are new")

	// Verify all values were set
	for key, expectedValue := range values {
		getResult := cmd.Get(ctx, key)
		require.NoError(t, getResult.Err())
		require.Equal(t, expectedValue, getResult.Val())
	}
}

// testMSetNXSomeExist tests MSetNX when some keys already exist
func testMSetNXSomeExist(t *testing.T, provider StringCommandProvider) {
	cmd := provider.GetStringCommand()
	ctx := provider.GetContext()

	// Set one existing key
	cmd.Set(ctx, "test:string:msetnx_exist2", []byte("existing"), 0)

	values := map[string]any{
		"test:string:msetnx_exist1": []byte("value1"),
		"test:string:msetnx_exist2": []byte("value2"), // already exists
		"test:string:msetnx_exist3": []byte("value3"),
	}

	// MSetNX should fail when any key exists
	result := cmd.MSetNX(ctx, values)
	require.NoError(t, result.Err())
	require.False(t, result.Val(), "MSetNX should return false when any key exists")

	// Verify existing key is unchanged
	getResult := cmd.Get(ctx, "test:string:msetnx_exist2")
	require.NoError(t, getResult.Err())
	require.Equal(t, []byte("existing"), getResult.Val())

	// Verify new keys were not set
	getResult1 := cmd.Get(ctx, "test:string:msetnx_exist1")
	require.Equal(t, caches.Nil, getResult1.Err())

	getResult3 := cmd.Get(ctx, "test:string:msetnx_exist3")
	require.Equal(t, caches.Nil, getResult3.Err())
}
