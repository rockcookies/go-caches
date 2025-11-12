package caches

import "context"

// SetCommand defines operations for Redis set data structure.
// Sets are unordered collections of unique strings.
type SetCommand interface {
	// SAdd adds one or more members to a set.
	// Creates the set if it does not exist.
	// Returns the number of members that were added to the set,
	// not including members already present.
	SAdd(ctx context.Context, key string, members ...any) Result[int64]

	// SCard returns the number of members in a set.
	// Returns 0 if the key does not exist.
	SCard(ctx context.Context, key string) Result[int64]

	// SDiff returns the difference of multiple sets.
	// The difference is the members of the first set that do not exist in any of the other sets.
	SDiff(ctx context.Context, keys ...string) Result[[][]byte]

	// SDiffStore stores the difference of multiple sets in a destination set.
	// If the destination set already exists, it is overwritten.
	// Returns the number of members in the resulting set.
	SDiffStore(ctx context.Context, destination string, keys ...string) Result[int64]

	// SInter returns the intersection of multiple sets.
	// The intersection is the members that exist in all given sets.
	SInter(ctx context.Context, keys ...string) Result[[][]byte]

	// SInterStore stores the intersection of multiple sets in a destination set.
	// If the destination set already exists, it is overwritten.
	// Returns the number of members in the resulting set.
	SInterStore(ctx context.Context, destination string, keys ...string) Result[int64]

	// SIsMember determines whether a member belongs to a set.
	// Returns true if the member is a member of the set, false otherwise.
	SIsMember(ctx context.Context, key string, member any) Result[bool]

	// SMembers returns all members of a set.
	// Returns an empty slice if the set does not exist.
	SMembers(ctx context.Context, key string) Result[[][]byte]

	// SMove moves a member from one set to another.
	// If the source set does not exist or does not contain the member, no operation is performed.
	// Returns true if the member was moved, false otherwise.
	SMove(ctx context.Context, source, destination string, member any) Result[bool]

	// SPop removes and returns a random member from a set.
	// Returns nil if the set does not exist or is empty.
	SPop(ctx context.Context, key string) Result[[]byte]

	// SPopN removes and returns up to count random members from a set.
	// Returns an empty slice if the set does not exist or is empty.
	SPopN(ctx context.Context, key string, count int64) Result[[][]byte]

	// SRandMember returns a random member from a set without removing it.
	// Returns nil if the set does not exist or is empty.
	SRandMember(ctx context.Context, key string) Result[[]byte]

	// SRandMemberN returns count random members from a set without removing them.
	// If count is positive, returns an array with distinct members.
	// If count is negative, returns an array with possibly repeated members.
	// Returns an empty slice if the set does not exist or is empty.
	SRandMemberN(ctx context.Context, key string, count int64) Result[[][]byte]

	// SRem removes one or more members from a set.
	// Members that are not members of the set are ignored.
	// Returns the number of members that were removed from the set,
	// not including non-existing members.
	SRem(ctx context.Context, key string, members ...any) Result[int64]

	// SScan iterates over members of a set.
	// cursor is the cursor to start iteration from (0 to start).
	// match is a glob-style pattern to filter members (empty string for no filter).
	// count is a hint for how many members to return per iteration.
	// Returns the next cursor and a slice of members.
	SScan(ctx context.Context, key string, cursor uint64, match string, count int64) Result[ScanResult]

	// SUnion returns the union of multiple sets.
	// The union is all members that exist in at least one of the given sets.
	SUnion(ctx context.Context, keys ...string) Result[[][]byte]

	// SUnionStore stores the union of multiple sets in a destination set.
	// If the destination set already exists, it is overwritten.
	// Returns the number of members in the resulting set.
	SUnionStore(ctx context.Context, destination string, keys ...string) Result[int64]
}

// ScanResult represents the result of a scan operation.
type ScanResult struct {
	// Cursor is the cursor to use in the next scan call.
	// A cursor value of 0 indicates the iteration is complete.
	Cursor uint64
	// Elements contains the scanned elements.
	Elements [][]byte
}
