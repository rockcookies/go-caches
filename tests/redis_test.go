package tests

import (
	"context"
	"testing"

	rds "github.com/redis/go-redis/v9"
	"github.com/rockcookies/go-caches"
	"github.com/rockcookies/go-caches/providers/redis"
	"github.com/stretchr/testify/suite"
)

// RedisTestSuite is the base test suite for Redis provider
type RedisTestSuite struct {
	suite.Suite
	client  *rds.Client
	provder *redis.Provider
	ctx     context.Context
}

// SetupSuite runs once before all tests
func (s *RedisTestSuite) SetupSuite() {
	// Skip if running in short mode
	if testing.Short() {
		s.T().Skip("Skipping integration test in short mode")
	}

	// Create Redis client
	s.client = rds.NewClient(&rds.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	// Ping to verify connection
	s.ctx = context.Background()
	if err := s.client.Ping(s.ctx).Err(); err != nil {
		s.T().Skipf("Redis is not available: %v", err)
	}

	// Create cache instance with key prefix to avoid conflicts
	s.provder = redis.NewWithOptions(s.client, &redis.Options{
		Prefix: "test:redis:",
	})
}

// TearDownSuite runs once after all tests
func (s *RedisTestSuite) TearDownSuite() {
	if s.client != nil {
		// Delete all test keys
		iter := s.client.Scan(s.ctx, 0, "test:redis:*", 0).Iterator()
		for iter.Next(s.ctx) {
			s.client.Del(s.ctx, iter.Val())
		}
		if err := iter.Err(); err != nil {
			s.T().Logf("Cleanup error: %v", err)
		}

		// Close client
		s.client.Close()
	}
}

// GetStringCommand implements StringCommandProvider interface
func (s *RedisTestSuite) GetStringCommand() caches.StringCommand {
	return s.provder
}

// GetKeyCommand implements KeyCommandProvider interface
func (s *RedisTestSuite) GetKeyCommand() caches.KeyCommand {
	return s.provder
}

// GetContext implements StringCommandProvider interface
func (s *RedisTestSuite) GetContext() context.Context {
	return s.ctx
}

// TestStringCommand runs all StringCommand tests
func (s *RedisTestSuite) TestStringCommand() {
	RunStringCommandTests(s.T(), s)
}

// TestKeyCommand runs all KeyCommand tests
func (s *RedisTestSuite) TestKeyCommand() {
	RunKeyCommandTests(s.T(), s)
}

// TestRedis runs all Redis provider tests
func TestRedis(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
