// Package base_timer 提供基础定时任务功能，包括间隔执行和延迟执行等特性
package base_timer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// SetTimeout 设置延迟执行的函数，返回停止通道
// 该函数创建一个定时器，在指定的延迟时间后执行回调函数
// 参数:
//   - ctx: 上下文，用于取消操作
//   - options: 配置选项，详见 SetTimeoutOptions 结构体
//
// 返回值:
//   - chan struct{}: 停止通道，向该通道发送信号可以停止定时器，或者等待该通道关闭表示任务已完成
func SetTimeout(ctx context.Context, options SetTimeoutOptions) chan struct{} {
	options.setDefaults() // 设置默认值

	// 参数校验
	if options.Fn == nil {
		panic("callback function is required")
	}

	stop := make(chan struct{})
	timer := time.NewTimer(options.Delay)
	var wg sync.WaitGroup
	var executed bool

	// 错误处理函数
	handleError := func(err error) {
		if options.OnError != nil {
			options.OnError(err)
		} else {
			log.Printf("Timeout callback error: %v", err)
		}
	}

	go func() {
		defer timer.Stop()

		select {
		case <-timer.C:
			// 延迟时间到达，执行回调
			executed = true
			wg.Add(1)
			go func() {
				defer wg.Done()
				runTimeoutCallback(ctx, options, handleError)
			}()

		case <-ctx.Done():
			// 上下文取消
			handleError(ctx.Err())
			if !timer.Stop() {
				// 尝试停止定时器，如果已经触发则等待任务完成
				if !executed {
					<-timer.C // 排空定时器通道
				}
			}
			wg.Wait() // 等待可能已启动的任务完成
			close(stop)
			return

		case <-stop:
			// 手动停止
			if !timer.Stop() {
				// 尝试停止定时器，如果已经触发则等待任务完成
				if !executed {
					<-timer.C // 排空定时器通道
				}
			}
			wg.Wait() // 等待可能已启动的任务完成
			return
		}

		// 任务执行完成后关闭通道
		wg.Wait()
		close(stop)
	}()

	return stop
}

// runTimeoutCallback 执行超时任务的回调函数，处理重试逻辑
// 该函数负责执行回调函数并处理可能的错误和重试
// 参数:
//   - ctx: 上下文，用于取消操作
//   - options: 配置选项
//   - handleError: 错误处理函数
func runTimeoutCallback(
	ctx context.Context,
	options SetTimeoutOptions,
	handleError func(error),
) {
	// 执行回调并处理重试
	shouldExit, err := executeWithRetry(
		ctx,
		func(execCtx context.Context) error {
			return options.Fn(execCtx)
		},
		options.ExecutionTimeout,
		options.TimeoutRetries,
		options.ErrorRetries,
		options.RetryDelay,
		handleError,
	)

	// 如果需要退出，记录错误
	if shouldExit && err != nil {
		handleError(fmt.Errorf("timeout task failed permanently: %v", err))
	}
}
