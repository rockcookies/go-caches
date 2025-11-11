# Tests Package

这个包聚合了 `go-caches` 所有 providers 的测试，提供统一的测试套件和测试执行环境。

## 目录结构

```
tests/
├── go.mod                        # 模块定义（包含所有 provider 依赖）
├── string_test.go                # StringCommand 接口的测试套件
├── redis_test.go                 # Redis provider 的测试
├── redka_test.go                 # Redka provider 的测试
├── helpers.go                    # 测试辅助函数
└── README.md                     # 本文档
```

## 设计理念

所有测试都集中在 `tests` 包中，使用 [testify](https://github.com/stretchr/testify) 框架组织，采用 **Provider Interface Pattern**：

- ✅ **统一管理**: 所有测试在一个地方，易于维护
- ✅ **避免循环依赖**: provider 不需要依赖测试工具
- ✅ **灵活扩展**: 轻松添加新的 Command 测试
- ✅ **自动组合**: provider 实现接口即可运行所有相关测试
- ✅ **规范化**: 使用 testify 提供清晰的测试结构

## 测试架构

### 核心概念

1. **Provider Interface**: 每个 Command 定义一个 Provider 接口
2. **Test Runner**: 提供运行该 Command 所有测试的函数
3. **Provider Suite**: 每个 provider (Redis/Redka) 实现所需接口

### 示例：StringCommand 测试

**`string_test.go`** - 定义接口和测试：

```go
// Provider 接口
type StringCommandProvider interface {
	suite.TestingSuite
	GetStringCommand() caches.StringCommand
	GetContext() context.Context
}

// 测试运行器
func RunStringCommandTests(s StringCommandProvider) {
	s.Run("SetAndGet", func() {
		testSetAndGet(s)
	})
	s.Run("SetNX", func() {
		testSetNX(s)
	})
	// ... 更多测试
}

// 具体测试函数
func testSetAndGet(s StringCommandProvider) {
	cmd := s.GetStringCommand()
	ctx := s.GetContext()
	// 测试实现...
}
```

**`redis_test.go`** - Redis Provider 实现：

```go
type RedisTestSuite struct {
	suite.Suite
	client *rds.Client
	cache  *redis.RedisCache
	ctx    context.Context
}

// 实现 StringCommandProvider 接口
func (s *RedisTestSuite) GetStringCommand() caches.StringCommand {
	return s.cache
}

func (s *RedisTestSuite) GetContext() context.Context {
	return s.ctx
}

// 运行 StringCommand 测试
func (s *RedisTestSuite) TestStringCommand() {
	RunStringCommandTests(s)
}

// 主测试入口
func TestRedis(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
```

## 测试覆盖

### StringCommand 测试

`TestStringCommand` 函数会测试以下所有方法：

- ✅ **Set/Get**: 基本的键值设置和获取
- ✅ **SetNX**: 仅当键不存在时设置
- ✅ **SetXX**: 仅当键存在时设置
- ✅ **SetArgs**: 使用高级参数设置（模式、TTL、Get、KeepTTL）
- ✅ **GetSet**: 设置新值并返回旧值
- ✅ **Incr/Decr**: 递增和递减整数值
- ✅ **IncrBy/DecrBy**: 按指定值递增和递减
- ✅ **IncrByFloat**: 按浮点数递增
- ✅ **StrLen**: 获取字符串长度
- ✅ **Expiration**: 测试过期时间功能

每个测试都会验证：
- 正常情况下的行为
- 边界情况（如键不存在时）
- 错误处理
- 返回值的正确性

## 最佳实践

### 1. 测试隔离
- 使用唯一的键前缀（如 `test:redis:`、`test:redka:`）避免冲突
- 使用 `SetupSuite` 和 `TearDownSuite` 管理资源
- 每个测试应该是独立的，不依赖其他测试的状态

### 2. Suite 生命周期

```go
SetupSuite()      // 在整个 suite 开始前运行一次
  TestXxx()       // 测试方法 1
  TestYyy()       // 测试方法 2
TearDownSuite()   // 在整个 suite 结束后运行一次
```

### 3. Testify 断言

```go
// 基本断言
s.NoError(err)                          // 断言无错误
s.Error(err)                            // 断言有错误
s.Equal(expected, actual)               // 断言相等
s.True(condition)                       // 断言为真
s.False(condition)                      // 断言为假

// 特殊断言
s.Nil(obj)                              // 断言为 nil
s.NotNil(obj)                           // 断言不为 nil
s.InDelta(expected, actual, delta)      // 浮点数比较

// Require vs Assert
s.Require().NoError(err)                // 失败时立即停止测试
s.Assert().NoError(err)                 // 失败时继续执行
```

### 4. 命名规范

- **Provider Interface**: `<Command>Provider`
  - 例如: `StringCommandProvider`, `HashCommandProvider`
- **Test Runner**: `Run<Command>Tests`
  - 例如: `RunStringCommandTests`, `RunHashCommandTests`
- **Test Functions**: `test<Operation>`
  - 例如: `testSetAndGet`, `testHSet`
- **Provider Suite**: `<Provider>TestSuite`
  - 例如: `RedisTestSuite`, `RedkaTestSuite`
- **Suite Test Methods**: `Test<Command>`
  - 例如: `TestStringCommand`, `TestHashCommand`## 运行测试

### 运行所有 Provider 测试

```bash
# 在 tests 目录下运行所有测试
cd tests
go test -v

# 或者从项目根目录运行
go test -v ./tests/

# 输出示例:
# === RUN   TestRedis
# === RUN   TestRedis/TestStringCommand
# === RUN   TestRedis/TestStringCommand/SetAndGet
# === RUN   TestRedis/TestStringCommand/SetNX
# ...
# === RUN   TestRedka
# === RUN   TestRedka/TestStringCommand
# ...
```

### 运行特定 Provider 的测试

```bash
# 仅运行 Redis 测试
go test -v -run TestRedis

# 仅运行 Redka 测试
go test -v -run TestRedka
```

### 运行特定的测试

```bash
# 运行所有 provider 的 StringCommand 测试
go test -v -run StringCommand

# 运行 Redis 的 StringCommand 测试
go test -v -run TestRedis/TestStringCommand

# 运行特定的子测试
go test -v -run TestRedis/TestStringCommand/SetNX
```

### 跳过集成测试

```bash
# 使用 -short 标志跳过需要外部服务的测试（如 Redis）
go test -v -short
```## 环境要求

### Redis Provider 测试
- 需要 Redis 服务器运行在 `localhost:6379`
- 使用 Docker 快速启动: `docker run -d -p 6379:6379 redis:latest`
- 如果 Redis 不可用，测试会自动跳过（使用 `-short` 标志）

### Redka Provider 测试
- 无需外部服务
- 使用内存数据库，开箱即用

## 添加新的 Command 测试

当你实现新的 Command 接口（如 `HashCommand`）时，按以下步骤添加测试：

### 步骤 1: 创建测试文件 (`hash_test.go`)

```go
package tests

import (
	"context"
	"github.com/rockcookies/go-caches"
	"github.com/stretchr/testify/suite"
)

// 1. 定义 Provider 接口
type HashCommandProvider interface {
	suite.TestingSuite
	GetHashCommand() caches.HashCommand
	GetContext() context.Context
}

// 2. 创建测试运行器
func RunHashCommandTests(s HashCommandProvider) {
	s.Run("HSet", func() {
		testHSet(s)
	})
	s.Run("HGet", func() {
		testHGet(s)
	})
	// 添加更多测试...
}

// 3. 实现具体测试函数
func testHSet(s HashCommandProvider) {
	cmd := s.GetHashCommand()
	ctx := s.GetContext()

	key := "test:hash:hset"
	field := "field1"
	value := "value1"

	result := cmd.HSet(ctx, key, field, value)
	s.NoError(result.Err())
	s.Equal(int64(1), result.Val())
}

func testHGet(s HashCommandProvider) {
	cmd := s.GetHashCommand()
	ctx := s.GetContext()

	key := "test:hash:hget"
	field := "field1"
	value := "value1"

	cmd.HSet(ctx, key, field, value)

	result := cmd.HGet(ctx, key, field)
	s.NoError(result.Err())
	s.Equal(value, string(result.Val()))
}
```

### 步骤 2: 更新 Redis Provider (`redis_test.go`)

```go
// 实现 HashCommandProvider 接口
func (s *RedisTestSuite) GetHashCommand() caches.HashCommand {
	return s.cache
}

// 添加测试方法
func (s *RedisTestSuite) TestHashCommand() {
	RunHashCommandTests(s)
}
```

### 步骤 3: 更新 Redka Provider (`redka_test.go`)

```go
// 实现 HashCommandProvider 接口
func (s *RedkaTestSuite) GetHashCommand() caches.HashCommand {
	return s.cache
}

// 添加测试方法
func (s *RedkaTestSuite) TestHashCommand() {
	RunHashCommandTests(s)
}
```

### 步骤 4: 运行测试

```bash
cd tests
go test -v

# 现在会运行：
# - TestRedis/TestStringCommand (所有 string 测试)
# - TestRedis/TestHashCommand (所有 hash 测试)
# - TestRedka/TestStringCommand (所有 string 测试)
# - TestRedka/TestHashCommand (所有 hash 测试)
```

完整示例请参考 `EXAMPLE_hash_test.go.txt` 文件。## 持续改进

欢迎贡献更多测试用例和改进建议！如果发现任何问题或有新的测试场景，请随时添加。
