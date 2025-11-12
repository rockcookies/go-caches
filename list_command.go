package caches

import "context"

// LInsertPosition represents the position for list insert operation.
type LInsertPosition string

const (
	// LInsertBefore inserts an element before the pivot.
	LInsertBefore LInsertPosition = "BEFORE"
	// LInsertAfter inserts an element after the pivot.
	LInsertAfter LInsertPosition = "AFTER"
)

// ListCommand defines operations for Redis list data structure.
// Lists are sequences of strings sorted by insertion order.
type ListCommand interface {
	// LIndex returns the element at index in the list stored at key.
	// The index is zero-based, so 0 means the first element, 1 the second element and so on.
	// Negative indices can be used to designate elements starting at the tail of the list.
	LIndex(ctx context.Context, key string, index int64) Result[[]byte]

	// LInsert inserts element in the list stored at key either before or after the reference value pivot.
	// Returns the length of the list after the insert operation, or -1 when the value pivot was not found.
	LInsert(ctx context.Context, key string, position LInsertPosition, pivot, element any) Result[int64]

	// LLen returns the length of the list stored at key.
	// If key does not exist, it is interpreted as an empty list and 0 is returned.
	LLen(ctx context.Context, key string) Result[int64]

	// LPop removes and returns the first element of the list stored at key.
	LPop(ctx context.Context, key string) Result[[]byte]

	// LPopCount removes and returns the first count elements of the list stored at key.
	LPopCount(ctx context.Context, key string, count int) Result[[][]byte]

	// LPush inserts all the specified values at the head of the list stored at key.
	// If key does not exist, it is created as empty list before performing the push operations.
	// Returns the length of the list after the push operations.
	LPush(ctx context.Context, key string, elements ...any) Result[int64]

	// LRange returns the specified elements of the list stored at key.
	// The offsets start and stop are zero-based indexes.
	// These offsets can be negative numbers indicating offsets starting at the end of the list.
	LRange(ctx context.Context, key string, start, stop int64) Result[[][]byte]

	// LRem removes the first count occurrences of elements equal to element from the list stored at key.
	// The count argument influences the operation in the following ways:
	// count > 0: Remove elements equal to element moving from head to tail.
	// count < 0: Remove elements equal to element moving from tail to head.
	// count = 0: Remove all elements equal to element.
	// Returns the number of removed elements.
	LRem(ctx context.Context, key string, count int64, element any) Result[int64]

	// LSet sets the list element at index to element.
	LSet(ctx context.Context, key string, index int64, element any) StatusResult

	// LTrim trims an existing list so that it will contain only the specified range of elements.
	// Both start and stop are zero-based indexes.
	LTrim(ctx context.Context, key string, start, stop int64) StatusResult

	// RPop removes and returns the last element of the list stored at key.
	RPop(ctx context.Context, key string) Result[[]byte]

	// RPopCount removes and returns the last count elements of the list stored at key.
	RPopCount(ctx context.Context, key string, count int) Result[[][]byte]

	// RPopLPush atomically returns and removes the last element (tail) of the list stored at source,
	// and pushes the element at the first element (head) of the list stored at destination.
	RPopLPush(ctx context.Context, source, destination string) Result[[]byte]

	// RPush inserts all the specified values at the tail of the list stored at key.
	// If key does not exist, it is created as empty list before performing the push operations.
	// Returns the length of the list after the push operations.
	RPush(ctx context.Context, key string, elements ...any) Result[int64]
}
