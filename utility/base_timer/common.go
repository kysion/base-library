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

// executeWithRetry 执行任务并处理超时和错误重试逻辑
// 该函数实现了任务执行的核心逻辑，包括超时控制和错误重试
// 参数:
//   - ctx: 上下文，用于取消操作
//   - execute: 要执行的函数
//   - timeout: 执行超时时间，0表示不设置超时
//   - timeoutRetries: 超时重试次数，<0表示忽略超时错误，0表示立即退出，>0表示重试次数
//   - errorRetries: 错误重试次数，<0表示忽略普通错误，0表示立即退出，>0表示重试次数
//   - retryDelay: 重试间隔时间
//   - handleError: 错误处理函数
//
// 返回值:
//   - shouldExit: 表示是否应该退出定时器
//   - lastError: 最后一次错误（如果有）
func executeWithRetry(
	ctx context.Context,
	execute func(ctx context.Context) error,
	timeout time.Duration,
	timeoutRetries int,
	errorRetries int,
	retryDelay time.Duration,
	handleError func(error),
) (bool, error) {
	var lastError error
	shouldExit := false

	// 计算最大尝试次数（包括初始尝试和重试）
	maxTimeoutAttempts := 1
	if timeoutRetries >= 0 {
		maxTimeoutAttempts = timeoutRetries + 1
	}

	maxErrorAttempts := 1
	if errorRetries >= 0 {
		maxErrorAttempts = errorRetries + 1
	}

	// 外层循环控制超时重试
	for timeoutAttempt := 0; timeoutAttempt < maxTimeoutAttempts; timeoutAttempt++ {
		if timeoutAttempt > 0 {
			// 非首次尝试，等待重试延迟
			time.Sleep(retryDelay)
		}

		var cancel context.CancelFunc
		executionCtx := ctx

		// 设置超时上下文
		if timeout > 0 {
			executionCtx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		// 执行回调
		errCh := make(chan error, 1)
		go func() {
			errCh <- execute(executionCtx)
		}()

		// 等待执行结果或超时
		var err error
		select {
		case err = <-errCh:
			// 执行完成，检查是否有错误
			if err != nil {
				lastError = err
				handleError(fmt.Errorf("attempt %d/%d failed: %v", timeoutAttempt+1, maxTimeoutAttempts, err))

				// 检查错误重试次数
				if errorRetries >= 0 {
					for errorAttempt := 1; errorAttempt < maxErrorAttempts; errorAttempt++ {
						time.Sleep(retryDelay)

						go func() {
							errCh <- execute(executionCtx)
						}()

						err = <-errCh
						if err == nil {
							return false, nil // 错误重试成功
						}

						lastError = err
						handleError(fmt.Errorf("error retry attempt %d/%d failed: %v",
							errorAttempt+1, maxErrorAttempts, err))
					}

					// 错误重试次数耗尽
					if errorRetries >= 0 {
						shouldExit = true
					}
				}
			} else {
				return false, nil // 成功执行
			}

		case <-executionCtx.Done():
			// 执行超时
			lastError = fmt.Errorf("execution timeout after %v", timeout)
			handleError(fmt.Errorf("timeout attempt %d/%d: %v", timeoutAttempt+1, maxTimeoutAttempts, lastError))

			// 如果超时重试次数为0或已耗尽，退出
			if timeoutRetries == 0 || timeoutAttempt >= timeoutRetries {
				shouldExit = true
				break
			}
		}

		// 如果应该退出，不再进行超时重试
		if shouldExit {
			break
		}
	}

	return shouldExit, lastError
}

// stopTimerAndWait 停止定时器并等待任务完成
// 根据不同的停止模式，采取不同的等待策略
// 参数:
//   - ticker: 定时器对象
//   - wg: 等待组，用于等待所有任务完成
//   - running: 正在运行的任务计数
//   - stopOpts: 停止选项
//   - handleError: 错误处理函数
func stopTimerAndWait(
	ticker *time.Ticker,
	wg *sync.WaitGroup,
	running *int64,
	stopOpts StopOptions,
	handleError func(error),
) {
	ticker.Stop() // 停止定时器，防止新任务启动

	switch stopOpts.Mode {
	case StopWaitModeBlocking:
		// 阻塞等待所有任务完成
		wg.Wait()

	case StopWaitModeNonBlocking:
		// 非阻塞，立即返回
		log.Printf("Stopping without waiting for %d running tasks", atomic.LoadInt64(running))

	case StopWaitModeTimeout:
		// 带超时的等待
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			log.Println("All tasks completed")

		case <-time.After(stopOpts.Timeout):
			handleError(fmt.Errorf("timed out waiting for tasks, %d still running", atomic.LoadInt64(running)))
		}
	}
}
