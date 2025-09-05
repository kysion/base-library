// Package base_timer 提供基础定时任务功能，包括间隔执行和延迟执行等特性
package base_timer

import (
	"context"
	"time"
)

// SetIntervalOptions 配置 SetInterval 函数的行为
// 用于定义间隔执行任务的各种参数和行为
type SetIntervalOptions struct {
	Interval  time.Duration                                                     // 执行间隔，两次执行开始时间的间隔
	Run       func(count int64, ticker *time.Ticker, ctx context.Context) error // 回调函数，count 表示当前执行次数，ticker 是定时器对象，ctx 是上下文
	Sync      bool                                                              // 是否同步执行（默认异步），同步模式下会等待当前任务执行完成后再开始下一次计时
	Immediate bool                                                              // 是否立即执行第一次，true 表示立即执行，false 表示等待一个间隔后再执行
	OnError   func(err error)                                                   // 错误处理函数，当任务执行出错时调用
	Stop      StopOptions                                                       // 停止相关配置，控制定时器停止时的行为
	Callback  CallbackOptions                                                   // 回调执行相关配置，控制回调函数的执行行为
}

// StopOptions 控制定时器停止时的行为
// 定义了当调用停止函数时，如何处理正在执行的任务
type StopOptions struct {
	Mode    StopWaitMode  // 停止时的等待模式，可选值见 StopWaitMode 常量
	Timeout time.Duration // 停止时的超时时间，仅在 Mode 为 StopWaitModeTimeout 时有效
}

// CallbackOptions 控制回调函数的执行行为
// 定义了回调函数执行时的超时控制和错误重试策略
type CallbackOptions struct {
	Timeout        time.Duration // 回调执行的超时时间，0 表示不设置超时
	TimeoutRetries int           // 超时重试次数（<0表示忽略错误不退出，>0表示重试次数，0表示即将停止定时器任务）
	ErrorRetries   int           // 错误重试次数（<0表示忽略错误不退出，>0表示重试次数，0表示即将停止定时器任务）
	RetryDelay     time.Duration // 重试间隔时间，两次重试之间的等待时间
}

// SetTimeoutOptions 配置 SetTimeout 函数的行为
// 用于定义延迟执行任务的各种参数和行为
type SetTimeoutOptions struct {
	Delay            time.Duration                   // 延迟执行时间，任务将在这个时间后执行
	Fn               func(ctx context.Context) error // 回调函数，ctx 是上下文
	OnError          func(err error)                 // 错误处理函数，当任务执行出错时调用
	ExecutionTimeout time.Duration                   // 执行超时时间，0 表示不设置超时
	TimeoutRetries   int                             // 超时重试次数（<0表示忽略错误不退出，>0表示重试次数，0表示即将停止定时器任务）
	ErrorRetries     int                             // 错误重试次数（<0表示忽略错误不退出，>0表示重试次数，0表示即将停止定时器任务）
	RetryDelay       time.Duration                   // 重试间隔时间，两次重试之间的等待时间
}

// 设置默认值
// 为 SetIntervalOptions 结构体的字段设置默认值
func (opts *SetIntervalOptions) setDefaults() *SetIntervalOptions {
	if opts.Interval <= 0 {
		opts.Interval = time.Second // 默认间隔1秒
	}
	if opts.Stop.Mode == 0 {
		opts.Stop.Mode = StopWaitModeBlocking // 默认使用阻塞等待模式
	}
	if opts.Stop.Timeout <= 0 && opts.Stop.Mode == StopWaitModeTimeout {
		opts.Stop.Timeout = 5 * time.Second // 默认超时等待时间为5秒
	}

	if opts.Callback.TimeoutRetries == 0 {
		opts.Callback.TimeoutRetries = -1 // 默认忽略超时错误
	}
	if opts.Callback.ErrorRetries == 0 {
		opts.Callback.ErrorRetries = -1 // 默认忽略普通错误
	}
	if opts.Callback.RetryDelay <= 0 {
		opts.Callback.RetryDelay = 500 * time.Millisecond // 默认重试间隔500ms
	}
	return opts
}

// 设置默认值
// 为 SetTimeoutOptions 结构体的字段设置默认值
func (opts *SetTimeoutOptions) setDefaults() *SetTimeoutOptions {
	if opts.Delay <= 0 {
		opts.Delay = time.Second // 默认延迟1秒
	}
	if opts.ExecutionTimeout <= 0 {
		opts.ExecutionTimeout = 0 // 默认不设置超时
	}
	if opts.TimeoutRetries == 0 {
		opts.TimeoutRetries = -1 // 默认忽略超时错误
	}
	if opts.ErrorRetries == 0 {
		opts.ErrorRetries = -1 // 默认忽略普通错误
	}
	if opts.RetryDelay <= 0 {
		opts.RetryDelay = 500 * time.Millisecond // 默认重试间隔500ms
	}

	return opts
}
