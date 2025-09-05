// Package base_timer 提供基础定时任务功能，包括间隔执行和延迟执行等特性
package base_timer

// StopWaitMode 停止时的等待模式
// 定义了定时器停止时如何处理正在执行的任务
type StopWaitMode int

const (
	// StopWaitModeBlocking 停止时阻塞等待所有任务完成
	// 当定时器停止时，会等待所有正在执行的任务完成后才返回
	StopWaitModeBlocking StopWaitMode = iota

	// StopWaitModeNonBlocking 停止时非阻塞，立即返回
	// 当定时器停止时，会立即返回，不等待正在执行的任务完成
	StopWaitModeNonBlocking

	// StopWaitModeTimeout 停止时超时等待
	// 当定时器停止时，会等待正在执行的任务完成，但如果超过指定的超时时间，则会立即返回
	StopWaitModeTimeout
)
