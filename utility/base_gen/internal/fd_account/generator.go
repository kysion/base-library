package fd_account

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kysion/base-library/utility/base_gen/internal/utils"
)

// Generator 定义账户号生成器接口
// 提供生成和验证账户号的基础功能
type Generator interface {
	Generate(accountType string) (string, error) // 生成指定类型的账户号
	Validate(accountNumber string) bool          // 验证账户号有效性
}

// GeneratorConfig 账户号生成器配置
// 包含生成账户号所需的所有配置参数
type GeneratorConfig struct {
	Prefix         string       // 账户前缀（必填）
	Length         int          // 账户号总长度（必填）
	RandomLength   int          // 随机数部分长度
	ChecksumLength int          // 校验码长度
	TimeFormat     string       // 时间戳格式（默认YYMMDD）
	Separator      string       // 分隔符（默认"-"）
	Strategy       StrategyType // 生成策略（默认时间戳+随机数）
	Validator      Validator    // 校验器（默认Luhn算法）
}

// StrategyType 定义账户号生成策略类型
type StrategyType int

const (
	// StrategyRandom 纯随机策略：仅使用随机数生成账户号
	StrategyRandom StrategyType = iota
	// StrategyTimestamp 时间戳+随机数策略：使用当前时间戳作为基础
	StrategyTimestamp
	// StrategyCounter 计数器+随机数策略：使用单调递增计数器
	StrategyCounter
	// StrategyUUID UUID策略：基于UUID生成账户号
	StrategyUUID
)

type DefaultGenerator struct {
	config *GeneratorConfig
}

// ConfigOption 配置选项函数类型
// 用于通过选项模式配置生成器
type ConfigOption func(*GeneratorConfig) error

// WithLength 设置账户号总长度
// 参数：
//
//	length: 账户号总长度（必须大于0）
//
// 返回配置选项
func WithLength(length int) ConfigOption {
	return func(c *GeneratorConfig) error {
		if length <= 0 {
			return errors.New("length must be positive")
		}
		c.Length = length
		return nil
	}
}

// WithStrategy 设置账户号生成策略
// 参数：
//
//	strategy: 需要设置的策略类型
//
// 返回配置选项
func WithStrategy(strategy StrategyType) ConfigOption {
	return func(c *GeneratorConfig) error {
		if _, ok := strategies[strategy]; !ok {
			return fmt.Errorf("invalid strategy: %v", strategy)
		}
		c.Strategy = strategy
		return nil
	}
}

// NewDefaultGenerator 创建默认账户号生成器
func NewDefaultGenerator(opts ...ConfigOption) (*DefaultGenerator, error) {
	config := &GeneratorConfig{
		Prefix:         "KY",
		Length:         20,
		RandomLength:   10,
		ChecksumLength: 2,
		TimeFormat:     "060102", // YYMMDD
		Separator:      "-",
		Strategy:       StrategyTimestamp,
		Validator:      NewLuhnValidator(),
	}

	// 应用配置选项
	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, errors.Join(err, errors.New("failed to apply config option"))
		}
	}

	// 验证配置
	if err := validateConfig(config); err != nil {
		return nil, errors.Join(err, errors.New("invalid generator configuration"))
	}

	return &DefaultGenerator{config: config}, nil
}

// Generate 生成账户号
func (g *DefaultGenerator) Generate(accountType string) (string, error) {
	// 生成账户主体
	body, err := g.generateBody(accountType)
	if err != nil {
		return "", errors.Join(err, errors.New("failed to generate account body"))
	}

	// 计算校验码
	checksum, err := g.config.Validator.CalculateChecksum(body)
	if err != nil {
		return "", errors.Join(err, errors.New("failed to calculate checksum"))
	}

	// 组合完整账户号
	fullAccount := g.config.Prefix + g.config.Separator + body + g.config.Separator + checksum

	// 应用格式化规则
	return utils.FormatAccountNumber(fullAccount, g.config.Separator), nil
}

// Validate 验证账户号
func (g *DefaultGenerator) Validate(accountNumber string) bool {
	return g.config.Validator.Validate(accountNumber)
}

// generateBody 根据策略生成账户主体
func (g *DefaultGenerator) generateBody(accountType string) (string, error) {
	strategy, ok := strategies[g.config.Strategy]
	if !ok {
		return "", fmt.Errorf("unsupported strategy: %v", g.config.Strategy)
	}

	return strategy(g.config, accountType)
}

// validateConfig 验证配置
func validateConfig(config *GeneratorConfig) error {
	if len(config.Prefix) == 0 {
		return errors.New("prefix cannot be empty")
	}

	if config.Length <= len(config.Prefix)+config.ChecksumLength {
		return errors.New("total length must be greater than prefix length plus checksum length")
	}

	if config.RandomLength < 0 {
		return errors.New("random length cannot be negative")
	}

	if config.ChecksumLength < 0 {
		return errors.New("checksum length cannot be negative")
	}

	if _, ok := strategies[config.Strategy]; !ok {
		return fmt.Errorf("invalid strategy: %v", config.Strategy)
	}

	return nil
}

// 策略实现映射
var strategies = map[StrategyType]func(*GeneratorConfig, string) (string, error){
	StrategyRandom:    randomStrategy,
	StrategyTimestamp: timestampStrategy,
	StrategyCounter:   counterStrategy,
	StrategyUUID:      uuidStrategy,
}

// randomStrategy 纯随机策略
func randomStrategy(config *GeneratorConfig, _ string) (string, error) {
	bodyLength := config.Length - len(config.Prefix) - config.ChecksumLength
	return utils.GenerateRandomDigits(bodyLength)
}

// timestampStrategy 时间戳+随机数策略
func timestampStrategy(config *GeneratorConfig, _ string) (string, error) {
	timestamp := time.Now().Format(config.TimeFormat)

	// 计算剩余长度
	remaining := config.Length - len(config.Prefix) - config.ChecksumLength - len(timestamp)
	if remaining < 0 {
		return timestamp[:config.Length-len(config.Prefix)-config.ChecksumLength], nil
	}

	// 生成随机后缀
	randomPart, err := utils.GenerateRandomDigits(remaining)
	if err != nil {
		return "", err
	}

	return timestamp + randomPart, nil
}

// counterStrategy 计数器+随机数策略
func counterStrategy(config *GeneratorConfig, _ string) (string, error) {
	// 注意：这里使用简化的计数器实现，实际应用中应使用原子计数器或数据库序列
	counter := utils.NextCounter()
	counterStr := fmt.Sprintf("%d", counter)

	// 计算剩余长度
	remaining := config.Length - len(config.Prefix) - config.ChecksumLength - len(counterStr)
	if remaining < 0 {
		return counterStr[:config.Length-len(config.Prefix)-config.ChecksumLength], nil
	}

	// 生成随机后缀
	randomPart, err := utils.GenerateRandomDigits(remaining)
	if err != nil {
		return "", err
	}

	return counterStr + randomPart, nil
}

// uuidStrategy UUID策略
func uuidStrategy(config *GeneratorConfig, _ string) (string, error) {
	uuid, err := utils.GenerateUUID()
	if err != nil {
		return "", err
	}

	// 移除UUID中的分隔符
	uuid = strings.ReplaceAll(uuid, "-", "")

	// 截取需要的长度
	lengthNeeded := config.Length - len(config.Prefix) - config.ChecksumLength
	if len(uuid) > lengthNeeded {
		uuid = uuid[:lengthNeeded]
	}

	return uuid, nil
}
