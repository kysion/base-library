# 基础常量 (base_consts)

基础常量模块提供了一系列标准化的常量定义，便于系统间的一致性和互操作性。这些常量涵盖了常见的状态码、错误类型、配置项等。

## 主要常量组

基础常量模块包含以下主要常量组：

### 状态常量

系统常用的状态标识：

```go
import "github.com/kysion/base-library/base_consts"

const (
    StatusActive   = 1  // 激活状态
    StatusInactive = 2  // 非激活状态
    StatusDeleted  = 3  // 已删除状态
    // 其他状态常量...
)
```

### 错误码常量

标准化的错误码定义：

```go
import "github.com/kysion/base-library/base_consts"

const (
    ErrCodeSuccess        = 0      // 成功
    ErrCodeUnknown        = 10000  // 未知错误
    ErrCodeNotFound       = 10001  // 资源不存在
    ErrCodeUnauthorized   = 10002  // 未授权访问
    // 其他错误码常量...
)
```

### 配置键常量

系统配置的标准键名：

```go
import "github.com/kysion/base-library/base_consts"

const (
    ConfigKeyDatabase = "database"  // 数据库配置键
    ConfigKeyRedis    = "redis"     // Redis配置键
    ConfigKeyServer   = "server"    // 服务器配置键
    // 其他配置键常量...
)
```

### 时间常量

常用的时间单位和格式：

```go
import "github.com/kysion/base-library/base_consts"

const (
    TimeFormatDefault       = "2006-01-02 15:04:05"  // 默认时间格式
    TimeFormatDate          = "2006-01-02"           // 日期格式
    TimeFormatTime          = "15:04:05"             // 时间格式
    TimeFormatYearMonth     = "2006-01"              // 年月格式
    // 其他时间常量...
)
```

## 使用示例

基础常量的使用非常简单，直接导入包并使用相应的常量即可：

```go
package main

import (
    "fmt"
    "github.com/kysion/base-library/base_consts"
)

func main() {
    // 使用状态常量
    if userStatus == base_consts.StatusActive {
        fmt.Println("用户处于激活状态")
    }
    
    // 使用错误码常量
    if errorCode == base_consts.ErrCodeNotFound {
        fmt.Println("资源不存在")
    }
    
    // 使用时间格式常量
    timeStr := time.Now().Format(base_consts.TimeFormatDefault)
    fmt.Println("当前时间:", timeStr)
}
```

## 扩展常量

如果需要在项目中扩展或自定义常量，建议创建专门的常量包，并遵循类似的命名和组织方式：

```go
package project_consts

import "github.com/kysion/base-library/base_consts"

// 扩展状态常量
const (
    StatusPending  = 4  // 在基础常量之后继续编号
    StatusApproved = 5
    StatusRejected = 6
)

// 扩展错误码常量
const (
    ErrCodeCustomBase     = 20000  // 自定义错误码基准值
    ErrCodeInvalidFormat  = 20001
    ErrCodeDuplicate      = 20002
)
```

这样可以保持项目中常量的一致性和可维护性。
