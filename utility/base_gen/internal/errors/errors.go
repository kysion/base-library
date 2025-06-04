package errors

import "fmt"

// CustomError 实现了标准error接口的自定义错误类型
// 包含错误信息和底层错误（可选）
type CustomError struct {
	Message string // 错误描述信息
	Err     error  // 底层错误（可为nil）
}

// Error 实现error接口方法
// 返回组合后的错误信息，包含当前消息和底层错误（如果存在）
func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 实现errors.Unwrap接口方法
// 返回底层错误用于错误链分析
func (e *CustomError) Unwrap() error {
	return e.Err
}

// Wrapf 包装错误并添加格式化上下文信息
// 参数：
//
//	err: 原始错误（可为nil）
//	format: 格式化字符串
//	args: 格式化参数
//
// 返回包装后的错误
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &CustomError{
		Message: fmt.Sprintf(format, args...),
		Err:     err,
	}
}

// New 创建新的自定义错误
// 参数：
//
//	message: 错误描述信息
//
// 返回新错误实例
func New(message string) error {
	return &CustomError{
		Message: message,
	}
}
