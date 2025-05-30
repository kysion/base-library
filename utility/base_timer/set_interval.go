// Package base_timer 提供基础定时任务功能，包括间隔执行和延迟执行等特性
package base_timer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// SetInterval 按指定间隔重复执行函数，返回停止通道
// 该函数创建一个定时器，按照指定的时间间隔重复执行回调函数
// 参数:
//   - ctx: 上下文，用于取消操作
//   - options: 配置选项，详见 SetIntervalOptions 结构体
//
// 返回值:
//   - chan struct{}: 停止通道，向该通道发送信号可以停止定时器
func SetInterval(ctx context.Context, options SetIntervalOptions) chan struct{} {
	options.setDefaults() // 设置默认值

	// 参数校验
	if options.Interval <= 0 {
		panic("interval must be positive")
	}
	if options.Fn == nil {
		panic("callback function is required")
	}

	stop := make(chan struct{})
	ticker := time.NewTicker(options.Interval)
	var wg sync.WaitGroup
	var running int64 // 记录正在运行的任务数

	// 错误处理函数
	handleError := func(err error) {
		if options.OnError != nil {
			options.OnError(err)
		} else {
			log.Printf("Interval callback error: %v", err)
		}
	}

	// 执行回调函数的包装
	executeCallback := func(currentCount int64) {
		if options.Sync {
			// 同步执行
			runIntervalCallback(ctx, currentCount, ticker, options, handleError)
		} else {
			// 异步执行
			wg.Add(1)
			atomic.AddInt64(&running, 1)
			go func() {
				defer wg.Done()
				defer atomic.AddInt64(&running, -1)
				runIntervalCallback(ctx, currentCount, ticker, options, handleError)
			}()
		}
	}

	// 立即执行第一次
	if options.Immediate {
		executeCallback(0)
	}

	go func() {
		defer ticker.Stop()
		count := int64(0)

		for {
			select {
			case <-ticker.C:
				count++
				executeCallback(count)

			case <-ctx.Done():
				handleError(ctx.Err())
				stopTimerAndWait(ticker, &wg, &running, options.Stop, handleError)
				close(stop)
				return

			case <-stop:
				stopTimerAndWait(ticker, &wg, &running, options.Stop, handleError)
				return
			}
		}
	}()

	return stop
}

// runIntervalCallback 执行间隔任务的回调函数，处理重试逻辑
// 该函数负责执行回调函数并处理可能的错误和重试
// 参数:
//   - ctx: 上下文，用于取消操作
//   - count: 当前执行次数
//   - ticker: 定时器对象
//   - options: 配置选项
//   - handleError: 错误处理函数
func runIntervalCallback(
	ctx context.Context,
	count int64,
	ticker *time.Ticker,
	options SetIntervalOptions,
	handleError func(error),
) {
	// 执行回调并处理重试
	shouldExit, err := executeWithRetry(
		ctx,
		func(execCtx context.Context) error {
			return options.Fn(count, ticker, execCtx)
		},
		options.Callback.Timeout,
		options.Callback.TimeoutRetries,
		options.Callback.ErrorRetries,
		options.Callback.RetryDelay,
		handleError,
	)

	// 如果需要退出，记录错误
	if shouldExit && err != nil {
		handleError(fmt.Errorf("interval task failed permanently: %v", err))
	}
}
