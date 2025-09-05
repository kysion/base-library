# Base Generator 基础生成器库

Base Generator 是一个用于生成和验证各类标识符的基础工具库，主要提供账户号生成、UUID生成、随机数生成等功能。该库采用模块化设计，支持多种生成策略和验证方式。

## 目录结构

```
base_gen/
├── internal/       # 内部实现模块
│   ├── errors/     # 错误处理模块
│   ├── fd_account/ # 账户号生成模块
│   └── utils/      # 通用工具模块
├── export.go       # 统一导出接口
└── README.md       # 项目文档
```

## 快速开始

```go
import "github.com/kysion/base-library/utility/base_gen"

// 创建默认生成器
generator, err := base_gen.NewDefaultGenerator(
    base_gen.WithLength(20),
    base_gen.WithStrategy(base_gen.StrategyTimestamp),
)

// 生成账户号
accountNumber, err := generator.Generate("SAVINGS")

// 验证账户号
isValid := generator.Validate(accountNumber)
```

## 功能特性

### 1. 账户号生成

提供灵活的账户号生成和验证功能。

#### 主要特性

- 支持多种生成策略：
  - 纯随机策略 (StrategyRandom)
  - 时间戳+随机数策略 (StrategyTimestamp)
  - 计数器+随机数策略 (StrategyCounter)
  - UUID策略 (StrategyUUID)
- 内置 Luhn 算法校验
- 可配置的账户号格式（前缀、长度、分隔符等）
- 线程安全的计数器实现

#### 使用示例

```go
// 创建自定义生成器
generator, err := base_gen.NewDefaultGenerator(
    base_gen.WithLength(20),
    base_gen.WithStrategy(base_gen.StrategyTimestamp),
)

// 创建 Luhn 校验器
validator := base_gen.NewLuhnValidator(
    base_gen.WithChecksumLength(2),
)

// 生成并验证账户号
accountNumber, err := generator.Generate("SAVINGS")
isValid := validator.Validate(accountNumber)
```

### 2. 工具函数

提供各种通用工具函数。

#### 主要功能

1. 字符串处理

   ```go
   // 格式化账户号
   formatted := base_gen.FormatAccountNumber("1234567890", "-")
   
   // 获取递增计数器值
   counter := base_gen.NextCounter()
   ```

2. 随机数生成

   ```go
   // 生成随机数字
   random, err := base_gen.GenerateRandomDigits(10)
   
   // 生成 UUID
   uuid, err := base_gen.GenerateUUID()
   ```

3. 校验码计算

   ```go
   // 计算 Luhn 校验码
   checksum := base_gen.CalculateLuhnChecksum("123456789")
   ```

### 3. 错误处理

提供统一的错误处理机制。

```go
// 创建新错误
err := base_gen.New("操作失败")

// 包装错误
wrappedErr := base_gen.Wrapf(err, "处理用户 %s 时", userId)
```

## 配置选项

### 账户号生成器配置

```go
type GeneratorConfig struct {
    Prefix         string       // 账户前缀
    Length         int          // 账户号总长度
    RandomLength   int          // 随机数部分长度
    ChecksumLength int          // 校验码长度
    TimeFormat     string       // 时间戳格式
    Separator      string       // 分隔符
    Strategy       StrategyType // 生成策略
    Validator      Validator    // 校验器
}
```

### Luhn 校验器配置

```go
type LuhnValidator struct {
    checksumLength int // 校验码长度
}
```

## 最佳实践

1. 账户号生成
   - 建议使用 `StrategyTimestamp` 策略生成业务账户号
   - 对于需要唯一性的场景，可以使用 `StrategyUUID` 策略
   - 合理设置账户号长度，确保包含足够的随机性

2. 错误处理
   - 使用 `Wrapf` 包装错误，保留错误上下文
   - 在关键业务逻辑处进行错误检查

3. 随机数生成
   - 使用 `GenerateRandomDigits` 生成安全的随机数字
   - 对于需要全局唯一性的场景，使用 `GenerateUUID`

## 注意事项

1. 账户号生成
   - 确保配置的账户号长度足够容纳所有必要信息
   - 注意校验码长度对总长度的影响
   - 合理选择生成策略，避免冲突

2. 安全性
   - 使用 `crypto/rand` 包生成随机数，确保安全性
   - 避免在日志中打印完整的账户号信息

3. 性能考虑
   - 生成器实例可以重复使用
   - 大量生成时注意内存使用

## 贡献指南

1. 遵循 Go 标准代码规范
2. 添加必要的单元测试
3. 更新文档说明
4. 提交前进行代码审查

## 许可证

MIT License
