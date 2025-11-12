package tests

import (
	"context"
	"testing"

	rdk "github.com/nalgeon/redka"
	"github.com/rockcookies/go-caches"
	"github.com/rockcookies/go-caches/providers/redka"
	"github.com/stretchr/testify/suite"

	// Import SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

// RedkaTestSuite is the base test suite for Redka provider
type RedkaTestSuite struct {
	suite.Suite
	db       *rdk.DB
	provider *redka.Provider
	ctx      context.Context
}

// SetupSuite runs once before all tests
func (s *RedkaTestSuite) SetupSuite() {
	// Create in-memory Redka database
	db, err := rdk.Open(":memory:", nil)
	s.Require().NoError(err, "Failed to open Redka database")
	s.db = db

	// Create cache instance with key prefix to avoid conflicts
	s.provider = redka.NewWithOptions(db, &redka.Options{
		Prefix: "test:redka:",
	})

	s.ctx = context.Background()
}

// TearDownSuite runs once after all tests
func (s *RedkaTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

// GetStringCommand implements StringCommandProvider interface
func (s *RedkaTestSuite) GetStringCommand() caches.StringCommand {
	return s.provider
}

// GetKeyCommand implements KeyCommandProvider interface
func (s *RedkaTestSuite) GetKeyCommand() caches.KeyCommand {
	return s.provider
}

// GetHashCommand implements HashCommandProvider interface
func (s *RedkaTestSuite) GetHashCommand() caches.HashCommand {
	return s.provider
}

// GetListCommand implements ListCommandProvider interface
func (s *RedkaTestSuite) GetListCommand() caches.ListCommand {
	return s.provider
}

// GetSetCommand implements SetCommandProvider interface
func (s *RedkaTestSuite) GetSetCommand() caches.SetCommand {
	return s.provider
}

// GetSortedSetCommand implements SortedSetCommandProvider interface
func (s *RedkaTestSuite) GetSortedSetCommand() caches.SortedSetCommand {
	return s.provider
}

// GetContext implements StringCommandProvider interface
func (s *RedkaTestSuite) GetContext() context.Context {
	return s.ctx
}

// TestStringCommand runs all StringCommand tests
func (s *RedkaTestSuite) TestStringCommand() {
	RunStringCommandTests(s.T(), s)
}

// TestKeyCommand runs all KeyCommand tests
func (s *RedkaTestSuite) TestKeyCommand() {
	RunKeyCommandTests(s.T(), s)
}

// TestHashCommand runs all HashCommand tests
func (s *RedkaTestSuite) TestHashCommand() {
	RunHashCommandTests(s.T(), s)
}

// TestListCommand runs all ListCommand tests
func (s *RedkaTestSuite) TestListCommand() {
	RunListCommandTests(s.T(), s)
}

// TestSetCommand runs all SetCommand tests
func (s *RedkaTestSuite) TestSetCommand() {
	RunSetCommandTests(s.T(), s)
}

// TestSortedSetCommand runs all SortedSetCommand tests
func (s *RedkaTestSuite) TestSortedSetCommand() {
	RunSortedSetCommandTests(s.T(), s)
}

// TestRedka runs all Redka provider tests
func TestRedka(t *testing.T) {
	suite.Run(t, new(RedkaTestSuite))
}
