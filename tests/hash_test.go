package tests

import (
	"context"
	"testing"

	"github.com/rockcookies/go-caches"
	"github.com/stretchr/testify/require"
)

// HashCommandProvider defines the interface for testing HashCommand implementations
type HashCommandProvider interface {
	GetHashCommand() caches.HashCommand
	GetContext() context.Context
}

// RunHashCommandTests runs all HashCommand tests
func RunHashCommandTests(t *testing.T, provider HashCommandProvider) {
	t.Run("HSet_and_HGet", func(t *testing.T) {
		testHSetAndHGet(t, provider)
	})
	t.Run("HGet_NonExistent", func(t *testing.T) {
		testHGetNonExistent(t, provider)
	})
	t.Run("HSet_MultipleFields", func(t *testing.T) {
		testHSetMultipleFields(t, provider)
	})
	t.Run("HMSet_and_HMGet", func(t *testing.T) {
		testHMSetAndHMGet(t, provider)
	})
	t.Run("HMGet_PartialExist", func(t *testing.T) {
		testHMGetPartialExist(t, provider)
	})
	t.Run("HGetAll", func(t *testing.T) {
		testHGetAll(t, provider)
	})
	t.Run("HGetAll_NonExistent", func(t *testing.T) {
		testHGetAllNonExistent(t, provider)
	})
	t.Run("HDel_SingleField", func(t *testing.T) {
		testHDelSingleField(t, provider)
	})
	t.Run("HDel_MultipleFields", func(t *testing.T) {
		testHDelMultipleFields(t, provider)
	})
	t.Run("HDel_NonExistent", func(t *testing.T) {
		testHDelNonExistent(t, provider)
	})
	t.Run("HExists_ExistingField", func(t *testing.T) {
		testHExistsExistingField(t, provider)
	})
	t.Run("HExists_NonExistentField", func(t *testing.T) {
		testHExistsNonExistentField(t, provider)
	})
	t.Run("HSetNX_NewField", func(t *testing.T) {
		testHSetNXNewField(t, provider)
	})
	t.Run("HSetNX_ExistingField", func(t *testing.T) {
		testHSetNXExistingField(t, provider)
	})
	t.Run("HKeys", func(t *testing.T) {
		testHKeys(t, provider)
	})
	t.Run("HKeys_NonExistent", func(t *testing.T) {
		testHKeysNonExistent(t, provider)
	})
	t.Run("HVals", func(t *testing.T) {
		testHVals(t, provider)
	})
	t.Run("HVals_NonExistent", func(t *testing.T) {
		testHValsNonExistent(t, provider)
	})
	t.Run("HLen", func(t *testing.T) {
		testHLen(t, provider)
	})
	t.Run("HLen_NonExistent", func(t *testing.T) {
		testHLenNonExistent(t, provider)
	})
	t.Run("HIncrBy_NewField", func(t *testing.T) {
		testHIncrByNewField(t, provider)
	})
	t.Run("HIncrBy_ExistingField", func(t *testing.T) {
		testHIncrByExistingField(t, provider)
	})
	t.Run("HIncrByFloat", func(t *testing.T) {
		testHIncrByFloat(t, provider)
	})
	t.Run("HScan_BasicIteration", func(t *testing.T) {
		testHScanBasicIteration(t, provider)
	})
	t.Run("HScan_WithPattern", func(t *testing.T) {
		testHScanWithPattern(t, provider)
	})
	t.Run("HScan_WithCount", func(t *testing.T) {
		testHScanWithCount(t, provider)
	})
	t.Run("HScan_NonExistent", func(t *testing.T) {
		testHScanNonExistent(t, provider)
	})
}

// testHSetAndHGet tests basic HSet and HGet operations
func testHSetAndHGet(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:set_get"
	field := "field1"
	value := []byte("value1")

	// Set field
	values := map[string]any{field: value}
	result := cmd.HSet(ctx, key, values)
	require.NoError(t, result.Err())
	require.Equal(t, int64(1), result.Val())

	// Get field
	getResult := cmd.HGet(ctx, key, field)
	require.NoError(t, getResult.Err())
	require.Equal(t, value, getResult.Val())
}

// testHGetNonExistent tests HGet on a non-existent field
func testHGetNonExistent(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:get_nonexistent"

	result := cmd.HGet(ctx, key, "nonexistent")
	require.Equal(t, caches.Nil, result.Err())
	require.Nil(t, result.Val())
}

// testHSetMultipleFields tests HSet with multiple fields
func testHSetMultipleFields(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:set_multiple"
	values := map[string]any{
		"field1": []byte("value1"),
		"field2": []byte("value2"),
		"field3": []byte("value3"),
	}

	// Set multiple fields
	result := cmd.HSet(ctx, key, values)
	require.NoError(t, result.Err())
	require.Equal(t, int64(3), result.Val())

	// Verify all fields were set
	for field, expectedValue := range values {
		getResult := cmd.HGet(ctx, key, field)
		require.NoError(t, getResult.Err())
		require.Equal(t, expectedValue, getResult.Val())
	}
}

// testHMSetAndHMGet tests HMSet and HMGet operations
func testHMSetAndHMGet(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:mset_mget"
	values := map[string]any{
		"field1": []byte("value1"),
		"field2": []byte("value2"),
		"field3": []byte("value3"),
	}

	// HMSet multiple fields
	result := cmd.HMSet(ctx, key, values)
	require.NoError(t, result.Err())
	require.Equal(t, "OK", result.Val())

	// HMGet all fields
	fields := []string{"field1", "field2", "field3"}
	getResult := cmd.HMGet(ctx, key, fields...)
	require.NoError(t, getResult.Err())
	require.Equal(t, len(fields), len(getResult.Val()))

	// Verify all values
	resultMap := getResult.Val()
	for field, expectedValue := range values {
		require.Equal(t, expectedValue, resultMap[field])
	}
}

// testHMGetPartialExist tests HMGet when some fields exist
func testHMGetPartialExist(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:mget_partial"

	// Set some fields
	values := map[string]any{
		"field1": []byte("value1"),
		"field3": []byte("value3"),
	}
	cmd.HSet(ctx, key, values)

	// HMGet with mixed existing/non-existing fields
	fields := []string{"field1", "field2", "field3"}
	result := cmd.HMGet(ctx, key, fields...)
	require.NoError(t, result.Err())

	resultMap := result.Val()
	require.Equal(t, []byte("value1"), resultMap["field1"])
	require.Nil(t, resultMap["field2"])
	require.Equal(t, []byte("value3"), resultMap["field3"])
}

// testHGetAll tests HGetAll operation
func testHGetAll(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:getall"
	values := map[string]any{
		"field1": []byte("value1"),
		"field2": []byte("value2"),
		"field3": []byte("value3"),
	}

	// Set multiple fields
	cmd.HSet(ctx, key, values)

	// Get all fields and values
	result := cmd.HGetAll(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, len(values), len(result.Val()))

	// Verify all values
	for field, expectedValue := range values {
		require.Equal(t, expectedValue, result.Val()[field])
	}
}

// testHGetAllNonExistent tests HGetAll on a non-existent key
func testHGetAllNonExistent(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:getall_nonexistent"

	result := cmd.HGetAll(ctx, key)
	require.NoError(t, result.Err())
	require.Empty(t, result.Val())
}

// testHDelSingleField tests HDel on a single field
func testHDelSingleField(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:del_single"

	// Set a field
	values := map[string]any{"field1": []byte("value1")}
	cmd.HSet(ctx, key, values)

	// Delete the field
	result := cmd.HDel(ctx, key, "field1")
	require.NoError(t, result.Err())
	require.Equal(t, int64(1), result.Val())

	// Verify field is deleted
	exists := cmd.HExists(ctx, key, "field1")
	require.NoError(t, exists.Err())
	require.False(t, exists.Val())
}

// testHDelMultipleFields tests HDel on multiple fields
func testHDelMultipleFields(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:del_multiple"

	// Set multiple fields
	values := map[string]any{
		"field1": []byte("value1"),
		"field2": []byte("value2"),
		"field3": []byte("value3"),
	}
	cmd.HSet(ctx, key, values)

	// Delete multiple fields
	result := cmd.HDel(ctx, key, "field1", "field2")
	require.NoError(t, result.Err())
	require.Equal(t, int64(2), result.Val())

	// Verify fields are deleted
	exists1 := cmd.HExists(ctx, key, "field1")
	require.NoError(t, exists1.Err())
	require.False(t, exists1.Val())

	exists2 := cmd.HExists(ctx, key, "field2")
	require.NoError(t, exists2.Err())
	require.False(t, exists2.Val())

	// Verify field3 still exists
	exists3 := cmd.HExists(ctx, key, "field3")
	require.NoError(t, exists3.Err())
	require.True(t, exists3.Val())
}

// testHDelNonExistent tests HDel on a non-existent field
func testHDelNonExistent(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:del_nonexistent"

	// Delete non-existent field
	result := cmd.HDel(ctx, key, "nonexistent")
	require.NoError(t, result.Err())
	require.Equal(t, int64(0), result.Val())
}

// testHExistsExistingField tests HExists on an existing field
func testHExistsExistingField(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:exists_existing"

	// Set a field
	values := map[string]any{"field1": []byte("value1")}
	cmd.HSet(ctx, key, values)

	// Check existence
	result := cmd.HExists(ctx, key, "field1")
	require.NoError(t, result.Err())
	require.True(t, result.Val())
}

// testHExistsNonExistentField tests HExists on a non-existent field
func testHExistsNonExistentField(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:exists_nonexistent"

	// Check non-existent field
	result := cmd.HExists(ctx, key, "nonexistent")
	require.NoError(t, result.Err())
	require.False(t, result.Val())
}

// testHSetNXNewField tests HSetNX on a new field
func testHSetNXNewField(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:setnx_new"
	field := "field1"
	value := []byte("value1")

	// HSetNX on new field should succeed
	result := cmd.HSetNX(ctx, key, field, value)
	require.NoError(t, result.Err())
	require.True(t, result.Val())

	// Verify value was set
	getResult := cmd.HGet(ctx, key, field)
	require.NoError(t, getResult.Err())
	require.Equal(t, value, getResult.Val())
}

// testHSetNXExistingField tests HSetNX on an existing field
func testHSetNXExistingField(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:setnx_existing"
	field := "field1"
	value1 := []byte("value1")
	value2 := []byte("value2")

	// Set initial value
	values := map[string]any{field: value1}
	cmd.HSet(ctx, key, values)

	// HSetNX on existing field should fail
	result := cmd.HSetNX(ctx, key, field, value2)
	require.NoError(t, result.Err())
	require.False(t, result.Val())

	// Verify original value is unchanged
	getResult := cmd.HGet(ctx, key, field)
	require.NoError(t, getResult.Err())
	require.Equal(t, value1, getResult.Val())
}

// testHKeys tests HKeys operation
func testHKeys(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:keys"
	values := map[string]any{
		"field1": []byte("value1"),
		"field2": []byte("value2"),
		"field3": []byte("value3"),
	}

	// Set multiple fields
	cmd.HSet(ctx, key, values)

	// Get all keys
	result := cmd.HKeys(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, len(values), len(result.Val()))

	// Verify all keys are present
	keys := result.Val()
	for field := range values {
		require.Contains(t, keys, field)
	}
}

// testHKeysNonExistent tests HKeys on a non-existent key
func testHKeysNonExistent(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:keys_nonexistent"

	result := cmd.HKeys(ctx, key)
	require.NoError(t, result.Err())
	require.Empty(t, result.Val())
}

// testHVals tests HVals operation
func testHVals(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:vals"
	values := map[string]any{
		"field1": []byte("value1"),
		"field2": []byte("value2"),
		"field3": []byte("value3"),
	}

	// Set multiple fields
	cmd.HSet(ctx, key, values)

	// Get all values
	result := cmd.HVals(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, len(values), len(result.Val()))

	// Verify all values are present
	vals := result.Val()
	for _, expectedValue := range values {
		found := false
		for _, val := range vals {
			if string(val) == string(expectedValue.([]byte)) {
				found = true
				break
			}
		}
		require.True(t, found, "Expected value not found in HVals result")
	}
}

// testHValsNonExistent tests HVals on a non-existent key
func testHValsNonExistent(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:vals_nonexistent"

	result := cmd.HVals(ctx, key)
	require.NoError(t, result.Err())
	require.Empty(t, result.Val())
}

// testHLen tests HLen operation
func testHLen(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:len"
	values := map[string]any{
		"field1": []byte("value1"),
		"field2": []byte("value2"),
		"field3": []byte("value3"),
	}

	// Set multiple fields
	cmd.HSet(ctx, key, values)

	// Get hash length
	result := cmd.HLen(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(len(values)), result.Val())
}

// testHLenNonExistent tests HLen on a non-existent key
func testHLenNonExistent(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:len_nonexistent"

	result := cmd.HLen(ctx, key)
	require.NoError(t, result.Err())
	require.Equal(t, int64(0), result.Val())
}

// testHIncrByNewField tests HIncrBy on a new field
func testHIncrByNewField(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:incrby_new"
	field := "counter"

	// HIncrBy on new field should start from 0
	result := cmd.HIncrBy(ctx, key, field, 5)
	require.NoError(t, result.Err())
	require.Equal(t, int64(5), result.Val())

	// Verify value
	getResult := cmd.HGet(ctx, key, field)
	require.NoError(t, getResult.Err())
	require.Equal(t, "5", string(getResult.Val()))
}

// testHIncrByExistingField tests HIncrBy on an existing field
func testHIncrByExistingField(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:incrby_existing"
	field := "counter"

	// Set initial value
	values := map[string]any{field: "10"}
	cmd.HSet(ctx, key, values)

	// Increment
	result := cmd.HIncrBy(ctx, key, field, 5)
	require.NoError(t, result.Err())
	require.Equal(t, int64(15), result.Val())

	// Decrement
	result2 := cmd.HIncrBy(ctx, key, field, -3)
	require.NoError(t, result2.Err())
	require.Equal(t, int64(12), result2.Val())
}

// testHIncrByFloat tests HIncrByFloat operation
func testHIncrByFloat(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:incrbyfloat"
	field := "price"

	// Set initial value
	values := map[string]any{field: "10.5"}
	cmd.HSet(ctx, key, values)

	// Increment by float
	result := cmd.HIncrByFloat(ctx, key, field, 2.3)
	require.NoError(t, result.Err())
	require.InDelta(t, 12.8, result.Val(), 0.0001)

	// Decrement by float
	result2 := cmd.HIncrByFloat(ctx, key, field, -0.3)
	require.NoError(t, result2.Err())
	require.InDelta(t, 12.5, result2.Val(), 0.0001)
}

// testHScanBasicIteration tests HScan basic iteration
func testHScanBasicIteration(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:scan"
	values := map[string]any{
		"field1": []byte("value1"),
		"field2": []byte("value2"),
		"field3": []byte("value3"),
		"field4": []byte("value4"),
		"field5": []byte("value5"),
	}

	// Set multiple fields
	cmd.HSet(ctx, key, values)

	// Scan all fields
	var allFields map[string][]byte = make(map[string][]byte)
	cursor := uint64(0)
	iterations := 0
	maxIterations := 10 // Prevent infinite loop

	for iterations < maxIterations {
		result := cmd.HScan(ctx, key, cursor, "", 10)
		require.NoError(t, result.Err())

		scanResult := result.Val()
		for field, value := range scanResult.Fields {
			allFields[field] = value
		}
		cursor = scanResult.Cursor

		iterations++
		if cursor == 0 {
			break
		}
	}

	// Should have found all fields
	require.Equal(t, len(values), len(allFields))
	for field, expectedValue := range values {
		require.Equal(t, expectedValue, allFields[field])
	}
}

// testHScanWithPattern tests HScan with pattern matching
func testHScanWithPattern(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:scanpat"
	values := map[string]any{
		"foo1": []byte("value1"),
		"foo2": []byte("value2"),
		"bar1": []byte("value3"),
		"bar2": []byte("value4"),
	}

	// Set multiple fields
	cmd.HSet(ctx, key, values)

	// Scan only "foo" fields
	var allFields map[string][]byte = make(map[string][]byte)
	cursor := uint64(0)
	iterations := 0
	maxIterations := 10

	for iterations < maxIterations {
		result := cmd.HScan(ctx, key, cursor, "foo*", 10)
		require.NoError(t, result.Err())

		scanResult := result.Val()
		for field, value := range scanResult.Fields {
			allFields[field] = value
		}
		cursor = scanResult.Cursor

		iterations++
		if cursor == 0 {
			break
		}
	}

	// Should have found only foo fields
	require.Contains(t, allFields, "foo1")
	require.Contains(t, allFields, "foo2")
	require.NotContains(t, allFields, "bar1")
	require.NotContains(t, allFields, "bar2")
}

// testHScanWithCount tests HScan with count parameter
func testHScanWithCount(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:scancount"
	values := map[string]any{}
	for i := 1; i <= 20; i++ {
		field := string(rune('a' + i - 1))
		values[field] = []byte("value")
	}

	// Set multiple fields
	cmd.HSet(ctx, key, values)

	// Scan with small count
	result := cmd.HScan(ctx, key, 0, "", 5)
	require.NoError(t, result.Err())

	scanResult := result.Val()
	// Count is a hint, so we just verify we got some fields
	require.NotNil(t, scanResult.Fields)
}

// testHScanNonExistent tests HScan on a non-existent key
func testHScanNonExistent(t *testing.T, provider HashCommandProvider) {
	cmd := provider.GetHashCommand()
	ctx := provider.GetContext()

	key := "test:hash:scan_nonexistent"

	result := cmd.HScan(ctx, key, 0, "", 10)
	require.NoError(t, result.Err())

	scanResult := result.Val()
	require.Empty(t, scanResult.Fields)
	require.Equal(t, uint64(0), scanResult.Cursor)
}
