# base_timer 包

`base_timer` 是一个提供基础定时任务功能的 Go 语言包，它封装了 Go 标准库中的定时器功能，提供了更加灵活和强大的定时任务管理能力。

## 功能特点

- **间隔执行**：通过 `SetInterval` 函数，可以按照指定的时间间隔重复执行任务
- **延迟执行**：通过 `SetTimeout` 函数，可以在指定的延迟时间后执行任务
- **错误处理**：提供了完善的错误处理机制，包括错误回调和重试策略
- **超时控制**：可以为任务执行设置超时时间，防止任务执行时间过长
- **灵活停止**：支持多种停止模式，包括阻塞等待、非阻塞立即返回和超时等待
- **同步/异步执行**：支持同步和异步执行模式

## 安装

```bash
go get github.com/kysion/base-library/utility/base_timer
```

## 使用方法

### SetInterval - 间隔执行

`SetInterval` 函数用于按照指定的时间间隔重复执行任务。

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kysion/base-library/utility/base_timer"
)

func main() {
	ctx := context.Background()

	// 创建一个每秒执行一次的定时任务
	stop := base_timer.SetInterval(ctx, base_timer.SetIntervalOptions{
		Interval: time.Second,
		Fn: func(count int64, ticker *time.Ticker, ctx context.Context) error {
			fmt.Printf("执行次数: %d\n", count)
			return nil
		},
		Immediate: true, // 立即执行第一次
		OnError: func(err error) {
			fmt.Printf("错误: %v\n", err)
		},
	})

	// 10秒后停止
	time.Sleep(10 * time.Second)
	stop <- struct{}{}
}
```

### SetTimeout - 延迟执行

`SetTimeout` 函数用于在指定的延迟时间后执行任务。

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kysion/base-library/utility/base_timer"
)

func main() {
	ctx := context.Background()

	// 创建一个3秒后执行的延迟任务
	stop := base_timer.SetTimeout(ctx, base_timer.SetTimeoutOptions{
		Delay: 3 * time.Second,
		Fn: func(ctx context.Context) error {
			fmt.Println("延迟任务执行")
			return nil
		},
		OnError: func(err error) {
			fmt.Printf("错误: %v\n", err)
		},
	})

	// 等待任务完成
	<-stop
	fmt.Println("任务已完成")
}
```

## 高级用法

### 设置超时和重试

```go
base_timer.SetInterval(ctx, base_timer.SetIntervalOptions{
	Interval: time.Second,
	Fn: func(count int64, ticker *time.Ticker, ctx context.Context) error {
		// 模拟耗时操作
		time.Sleep(2 * time.Second)
		return nil
	},
	Callback: base_timer.CallbackOptions{
		Timeout:        1 * time.Second,  // 设置执行超时为1秒
		TimeoutRetries: 3,               // 超时后重试3次
		ErrorRetries:   2,               // 发生错误后重试2次
		RetryDelay:     time.Second,     // 重试间隔为1秒
	},
})
```

### 不同的停止模式

```go
// 阻塞等待所有任务完成
stop := base_timer.SetInterval(ctx, base_timer.SetIntervalOptions{
	Interval: time.Second,
	Fn: func(count int64, ticker *time.Ticker, ctx context.Context) error {
		return nil
	},
	Stop: base_timer.StopOptions{
		Mode: base_timer.StopWaitModeBlocking,
	},
})

// 非阻塞立即返回
stop := base_timer.SetInterval(ctx, base_timer.SetIntervalOptions{
	Interval: time.Second,
	Fn: func(count int64, ticker *time.Ticker, ctx context.Context) error {
		return nil
	},
	Stop: base_timer.StopOptions{
		Mode: base_timer.StopWaitModeNonBlocking,
	},
})

// 超时等待
stop := base_timer.SetInterval(ctx, base_timer.SetIntervalOptions{
	Interval: time.Second,
	Fn: func(count int64, ticker *time.Ticker, ctx context.Context) error {
		return nil
	},
	Stop: base_timer.StopOptions{
		Mode:    base_timer.StopWaitModeTimeout,
		Timeout: 5 * time.Second,
	},
})
```

### 同步执行

```go
base_timer.SetInterval(ctx, base_timer.SetIntervalOptions{
	Interval: time.Second,
	Fn: func(count int64, ticker *time.Ticker, ctx context.Context) error {
		return nil
	},
	Sync: true, // 同步执行，每次任务执行完成后才会开始下一次计时
})
```

## 配置选项

### SetIntervalOptions

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| Interval | time.Duration | 执行间隔 | 1秒 |
| Fn | func(count int64, ticker *time.Ticker, ctx context.Context) error | 回调函数 | 必填 |
| Sync | bool | 是否同步执行 | false |
| Immediate | bool | 是否立即执行第一次 | false |
| OnError | func(err error) | 错误处理函数 | nil |
| Stop | StopOptions | 停止相关配置 | - |
| Callback | CallbackOptions | 回调执行相关配置 | - |

### SetTimeoutOptions

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| Delay | time.Duration | 延迟执行时间 | 1秒 |
| Fn | func(ctx context.Context) error | 回调函数 | 必填 |
| OnError | func(err error) | 错误处理函数 | nil |
| ExecutionTimeout | time.Duration | 执行超时时间 | 0（不设置超时）|
| TimeoutRetries | int | 超时重试次数 | -1（忽略超时错误）|
| ErrorRetries | int | 错误重试次数 | -1（忽略普通错误）|
| RetryDelay | time.Duration | 重试间隔时间 | 500ms |

### StopOptions

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| Mode | StopWaitMode | 停止时的等待模式 | StopWaitModeBlocking（支持StopWaitModeBlocking、StopWaitModeNonBlocking、StopWaitModeTimeout三种模式） |
| Timeout | time.Duration | 停止时的超时时间 | 5秒（仅在 Mode 为 StopWaitModeTimeout 时有效）|

### CallbackOptions

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| Timeout | time.Duration | 回调执行的超时时间 | 0（不设置超时）|
| TimeoutRetries | int | 超时重试次数 | -1（忽略超时错误）|
| ErrorRetries | int | 错误重试次数 | -1（忽略普通错误）|
| RetryDelay | time.Duration | 重试间隔时间 | 500ms |

## 注意事项

1. 当使用 `SetInterval` 时，如果回调函数执行时间超过了间隔时间，下一次执行会立即开始，不会等待完整的间隔时间。如果需要确保每次执行之间有固定的间隔，可以使用 `Sync: true` 选项。

2. 重试次数设置说明：
   - 值为 `-1`：忽略错误，继续执行
   - 值为 `0`：遇到错误立即退出
   - 值大于 `0`：表示最大重试次数

3. 停止定时器时，根据不同的停止模式有不同的行为：
   - `StopWaitModeBlocking`：阻塞等待所有任务完成
   - `StopWaitModeNonBlocking`：立即返回，不等待任务完成
   - `StopWaitModeTimeout`：等待指定时间，超时后返回

4. 使用 `context.Context` 可以在外部取消定时任务的执行。