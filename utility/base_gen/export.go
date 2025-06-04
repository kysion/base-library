// Package base_gen 提供基础生成器功能，包括账户号生成、UUID生成、随机数生成等
package base_gen

import (
	"github.com/kysion/base-library/utility/base_gen/internal/errors"
	"github.com/kysion/base-library/utility/base_gen/internal/fd_account"
	"github.com/kysion/base-library/utility/base_gen/internal/utils"
)

// 导出账户号生成器相关类型和函数

// Generator 定义账户号生成器接口
type Generator = fd_account.Generator

// GeneratorConfig 账户号生成器配置
type GeneratorConfig = fd_account.GeneratorConfig

// StrategyType 定义账户号生成策略类型
type StrategyType = fd_account.StrategyType

// 导出策略类型常量
const (
	StrategyRandom    = fd_account.StrategyRandom
	StrategyTimestamp = fd_account.StrategyTimestamp
	StrategyCounter   = fd_account.StrategyCounter
	StrategyUUID      = fd_account.StrategyUUID
)

// ConfigOption 配置选项函数类型
type ConfigOption = fd_account.ConfigOption

// WithLength 设置账户号总长度
func WithLength(length int) ConfigOption {
	return fd_account.WithLength(length)
}

// WithStrategy 设置账户号生成策略
func WithStrategy(strategy StrategyType) ConfigOption {
	return fd_account.WithStrategy(strategy)
}

// NewDefaultGenerator 创建默认账户号生成器
func NewDefaultGenerator(opts ...ConfigOption) (Generator, error) {
	return fd_account.NewDefaultGenerator(opts...)
}

// 导出校验器相关类型和函数

// Validator 账户号校验器接口
type Validator = fd_account.Validator

// LuhnValidator Luhn算法校验器
type LuhnValidator = fd_account.LuhnValidator

// LuhnOption Luhn校验器配置选项
type LuhnOption = fd_account.LuhnOption

// WithChecksumLength 设置校验码长度
func WithChecksumLength(length int) LuhnOption {
	return fd_account.WithChecksumLength(length)
}

// NewLuhnValidator 创建Luhn算法校验器
func NewLuhnValidator(opts ...LuhnOption) *LuhnValidator {
	return fd_account.NewLuhnValidator(opts...)
}

// 导出工具函数

// FormatAccountNumber 格式化账户号
func FormatAccountNumber(account string, separator string) string {
	return utils.FormatAccountNumber(account, separator)
}

// NextCounter 获取下一个线程安全的递增计数器值
func NextCounter() int64 {
	return utils.NextCounter()
}

// GenerateRandomDigits 生成指定长度的随机数字字符串
func GenerateRandomDigits(length int) (string, error) {
	return utils.GenerateRandomDigits(length)
}

// GenerateUUID 生成UUID
func GenerateUUID() (string, error) {
	return utils.GenerateUUID()
}

// CalculateLuhnChecksum 计算Luhn算法校验码
func CalculateLuhnChecksum(data string) int {
	return utils.CalculateLuhnChecksum(data)
}

// 导出错误处理相关函数

// CustomError 自定义错误类型
type CustomError = errors.CustomError

// New 创建新的自定义错误
func New(message string) error {
	return errors.New(message)
}

// Wrapf 包装错误并添加格式化上下文信息
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}
