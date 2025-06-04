package utils

import (
	"strings"
	"sync/atomic"
)

// FormatAccountNumber 使用指定分隔符格式化账户号码
// 将账户号按4位一组进行分割，自动处理短账户号
// 参数：
//
//	account: 需要格式化的原始账户号
//	separator: 使用的分隔符（如"-"）
//
// 返回格式化后的账户号
func FormatAccountNumber(account string, separator string) string {
	if separator == "" || len(account) <= 4 {
		return account
	}

	// 按每4个字符分组
	var parts []string
	for i := 0; i < len(account); i += 4 {
		end := i + 4
		if end > len(account) {
			end = len(account)
		}
		parts = append(parts, account[i:end])
	}

	return strings.Join(parts, separator)
}

// counter 是原子计数器变量，用于生成唯一序列号
var counter int64 = 0

// NextCounter 获取下一个线程安全的递增计数器值
// 使用原子操作保证并发安全性
// 返回单调递增的计数器值
func NextCounter() int64 {
	return atomic.AddInt64(&counter, 1)
}
