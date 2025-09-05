package fd_account

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kysion/base-library/utility/base_gen/internal/utils"
)

// Validator 账户号校验器接口
type Validator interface {
	CalculateChecksum(data string) (string, error)
	Validate(accountNumber string) bool
}

// LuhnValidator 实现Luhn算法校验器
// 用于计算和验证账户号的校验码
type LuhnValidator struct {
	checksumLength int // 校验码长度（默认2位）
}

// NewLuhnValidator 创建Luhn算法校验器
func NewLuhnValidator(opts ...LuhnOption) *LuhnValidator {
	v := &LuhnValidator{
		checksumLength: 2,
	}

	for _, opt := range opts {
		opt(v)
	}

	return v
}

// LuhnOption Luhn校验器配置选项
type LuhnOption func(*LuhnValidator)

// WithChecksumLength 设置校验码长度
func WithChecksumLength(length int) LuhnOption {
	return func(v *LuhnValidator) {
		v.checksumLength = length
	}
}

// CalculateChecksum 计算Luhn算法校验码
// 参数：
//
//	data: 需要计算校验码的数据部分
//
// 返回：
//
//	校验码字符串
//	错误信息（如果存在）
func (v *LuhnValidator) CalculateChecksum(data string) (string, error) {
	if len(data) == 0 {
		return "", errors.New("data cannot be empty")
	}

	checksum := utils.CalculateLuhnChecksum(data)
	return fmt.Sprintf("%0*d", v.checksumLength, checksum), nil
}

// Validate 验证账户号是否符合Luhn算法
// 参数：
//
//	accountNumber: 需要验证的账户号
//
// 返回：
//
//	验证结果（true为有效）
func (v *LuhnValidator) Validate(accountNumber string) bool {
	// 移除所有分隔符
	cleaned := strings.ReplaceAll(accountNumber, "-", "")

	// 检查最小长度要求
	if len(cleaned) <= v.checksumLength {
		return false
	}

	// 分离数据部分和校验码
	data := cleaned[:len(cleaned)-v.checksumLength]
	checksumStr := cleaned[len(cleaned)-v.checksumLength:]

	// 计算预期校验码
	expectedChecksum, err := v.CalculateChecksum(data)
	if err != nil {
		return false
	}

	return checksumStr == expectedChecksum
}
