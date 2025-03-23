# 工具集 (Utility)

工具集模块提供了一系列实用工具函数和结构，用于简化常见的开发任务，提高代码质量和开发效率。

## 模块概览

工具集包含以下主要组件：

- **枚举工具 (enum)** - 提供类型安全的枚举实现，支持整型、字符串等多种类型的枚举值

## 枚举工具 (enum)

枚举工具提供了一个灵活的枚举类型实现，支持整型、字符串型等多种类型的枚举值，以及枚举之间的位运算操作。

### 主要特性

- 支持多种数据类型的枚举值（整型、字符串型等）
- 支持枚举值之间的位运算（检查、添加、移除）
- 支持带有附加数据的枚举类型
- 高性能实现，基本操作无内存分配

### 使用示例

#### 整型枚举

```go
import "github.com/kysion/base-library/utility/enum"

// 定义权限枚举
type Permission interface {
    enum.IEnumCode[int]
}

// 创建枚举值
var (
    ReadPerm = enum.New[Permission](1, "读取权限")
    WritePerm = enum.New[Permission](2, "写入权限")
)

// 创建组合权限
userPerm := enum.New[Permission](0, "用户权限")
userPerm.Add(ReadPerm, WritePerm)

// 检查权限
hasReadPerm := userPerm.Has(ReadPerm) // true
```

#### 字符串枚举

```go
import "github.com/kysion/base-library/utility/enum"

// 定义国家枚举
type Country interface {
    enum.IEnumCode[string]
}

// 创建枚举值
var (
    China = enum.New[Country]("CN", "中国")
    USA = enum.New[Country]("US", "美国")
)

// 创建枚举实例
country := enum.New[Country]("CN", "中国")

// 检查是否匹配
isChina := country.Has(China) // true
```

更多详细信息，请查看 [枚举工具文档](enum/README.md)。
