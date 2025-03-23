# Kysion Base Library 基础库

Kysion Base Library 是一个综合性的 Go 语言工具库，提供了一系列可复用的组件和工具函数，旨在简化应用开发流程，提高代码质量。

## 功能概览

基础库包含以下主要模块：

## 详细文档

- [**基础常量 (base_consts)**](base_consts/README.md) - 提供常用的系统常量定义
- [**基础模型 (base_model)**](base_model/README.md) - 提供数据模型的基础结构和方法
- [**基础钩子 (base_hook)**](base_hook/README.md) - 提供插件式扩展的钩子系统
- [**工具集 (utility)**](utility/README.md) - 提供各类实用工具函数和结构
  - [**枚举工具 (enum)**](utility/enum/README.md) - 提供类型安全的枚举实现

## 安装

通过 Go 模块引入库:

```bash
go get -u github.com/kysion/base-library
```

## 使用指南

### 基础常量 (base_consts)

基础常量模块提供标准化的常量定义，便于系统间的一致性和互操作性。

```go
import "github.com/kysion/base-library/base_consts"

func main() {
    // 使用预定义常量
    status := base_consts.StatusActive
}
```

### 基础模型 (base_model)

基础模型提供数据模型的通用结构和方法，简化数据处理流程。

```go
import "github.com/kysion/base-library/base_model"

// 使用基础模型创建自定义模型
type User struct {
    base_model.BaseModel
    Name string
    Age  int
}
```

### 基础钩子 (base_hook)

基础钩子模块实现了一个灵活的钩子系统，支持插件式扩展。

```go
import "github.com/kysion/base-library/base_hook"

func main() {
    // 创建钩子管理器
    hooks := base_hook.NewManager()
    
    // 注册钩子处理函数
    hooks.Register("before_save", func(data interface{}) error {
        // 处理逻辑
        return nil
    })
    
    // 触发钩子
    hooks.Trigger("before_save", userData)
}
```

### 工具集 (utility)

工具集模块提供各类实用工具函数，提高开发效率。

#### 枚举工具 (enum)

枚举工具提供了类型安全的枚举实现，支持整型、字符串等多种类型，并支持位运算操作。

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

func main() {
    // 创建组合权限
    userPerm := enum.New[Permission](0, "用户权限")
    userPerm.Add(ReadPerm, WritePerm)
    
    // 检查权限
    hasReadPerm := userPerm.Has(ReadPerm) // true
}
```

更多枚举工具用法请参考 [枚举工具文档](utility/enum/README.md)。

## 版本兼容性

- Go 1.18+ (需要泛型支持)

## 贡献指南

我们欢迎所有形式的贡献，包括但不限于功能增强、bug修复、文档改进等。

贡献流程：

1. Fork 项目仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 提交 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请查看 LICENSE 文件。
