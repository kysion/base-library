# daoctl 数据访问控制工具包

`daoctl` 是一个基于 GoFrame (gf) 框架的数据访问控制工具包，提供了一套简洁高效的数据库操作接口，简化了数据库 CRUD 操作并提供了强大的缓存支持。

## 功能特点

- **统一的数据访问接口**：提供标准化的 DAO 接口，简化数据操作
- **ORM 缓存支持**：内置缓存机制，提高查询性能
- **扩展查询条件**：支持动态扩展的查询条件，实现复杂查询需求
- **泛型支持**：利用 Go 泛型特性，提供类型安全的数据访问方法
- **事务支持**：支持数据库事务操作
- **钩子机制**：提供数据库操作前后的钩子处理器

## 安装

```bash
go get github.com/kysion/base-library/utility/daoctl
```

## 核心组件

### 接口定义

#### IDao 接口

基础数据访问对象接口，定义了所有 DAO 对象必须实现的方法：

```go
type IDao interface {
    DB() gdb.DB                                                                               
    Table() string                                                                             
    Group() string                                                                            
    Ctx(ctx context.Context, cacheOption ...*gdb.CacheOption) *gdb.Model                      
    Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) 
    DaoConfig(ctx context.Context, cacheOption ...*gdb.CacheOption) *DaoConfig                
    GetExtWhereKeys() []string
    IsIgnoreCache() bool
}
```

#### TIDao 泛型接口

扩展了 IDao 接口，支持泛型操作：

```go
type TIDao[TColumns any] interface {
    IDao              
    Columns() TColumns
    IgnoreCache() IDao
    IgnoreExtModel(whereKey ...string) IDao
}
```

### 配置结构

```go
type DaoConfig struct {
    Dao         IDao                                                                               
    DB          gdb.DB                                                                            
    Table       string                                                                            
    Group       string                                                                             
    Model       *gdb.Model                                                                        
    CacheOption *gdb.CacheOption                                                                  
    HookHandler *gdb.HookHandler                                                                   
    ignoreCache bool                                                                              
    ExtWhere    map[string]func(model *gdb.Model, conf *DaoConfig, data ...interface{}) *gdb.Model
}
```

## 主要功能

### 数据查询

- `GetById`: 根据 ID 获取单个记录
- `GetByIdWithError`: 带错误返回的 ID 查询
- `Find`: 根据条件查找记录并排序
- `GetAll`: 获取所有符合条件的记录
- `Query`: 执行复杂的查询操作
- `Scan`: 扫描查询结果到结构体
- `ScanWithError`: 带错误返回的扫描操作

### 数据操作

- `Insert`: 插入新记录
- `Delete`: 删除记录
- `Save`: 保存记录（更新或插入）
- `Update`: 更新记录

### 缓存控制

- `EnableOrmCache`: 启用 ORM 缓存
- `DisabledOrmCache`: 禁用 ORM 缓存
- `MakeDaoCache`: 创建 DAO 缓存
- `SetTableIgnoreOrmCacheByCtx`: 设置上下文中忽略 ORM 缓存的表

### 扩展模型

- `MakeExtModelMap`: 注册外部模型的处理函数
- `ExecExWhere`: 执行扩展的条件查询
- `IgnoreExtModel`: 忽略特定的扩展模型字段
- `IsIgnoreExtModel`: 检查是否忽略某个扩展模型字段

## 使用示例

### 基本查询

```go
import (
    "context"
    "github.com/kysion/base-library/utility/daoctl"
)

// 根据 ID 查询
user := daoctl.GetById[User](dao.User.Ctx(ctx), 1)

// 条件查询
response, err := daoctl.Find[User](
    dao.User.Ctx(ctx), 
    []base_model.OrderBy{{Field: "id", Sort: "DESC"}},
    base_model.FilterInfo{Field: "status", Operator: "=", Value: 1},
)
```

### 带缓存的查询

```go
// 创建带缓存的 DAO 配置
config := daoctl.NewDaoConfig(ctx, dao.User)

// 使用配置进行查询
user := daoctl.GetById[User](config.Model, 1)
```

### 事务操作

```go
err := dao.User.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
    // 在事务中执行操作
    _, err := daoctl.Insert(dao.User.Ctx(ctx).TX(tx), &user)
    if err != nil {
        return err
    }
    
    // 更多事务操作...
    
    return nil
})
```

### 扩展查询条件

```go
// 注册扩展查询条件
daoctl.MakeExtModelMap("user", map[string]func(model *gdb.Model, conf *dao_interface.DaoConfig, data ...interface{}) *gdb.Model{
    "filter_active": func(model *gdb.Model, conf *dao_interface.DaoConfig, data ...interface{}) *gdb.Model {
        return model.Where("status", 1)
    },
})

// 查询时会自动应用扩展条件
users, err := daoctl.Find[User](dao.User.Ctx(ctx), nil)

// 忽略特定扩展条件
users, err := daoctl.Find[User](dao.User.IgnoreExtModel("filter_active").Ctx(ctx), nil)
```

## 最佳实践

1. **使用泛型接口**：优先使用 `TIDao` 泛型接口，获得类型安全的数据操作
2. **合理使用缓存**：根据业务需求配置缓存，提高查询性能
3. **扩展查询条件**：将常用查询条件封装为扩展条件，避免重复代码
4. **错误处理**：对于关键操作，使用带错误返回的方法，如 `GetByIdWithError`
5. **事务管理**：复杂操作使用事务保证数据一致性

## 配置说明

在配置文件中可以设置 ORM 缓存的相关参数：

```toml
[ormCache]
  [ormCache.ignore]
    # 不使用缓存的表列表
    tables = ["log_*", "cache_*"]
```

## 进阶主题

### 自定义 DAO 实现

创建自定义的 DAO 实现，需遵循 `IDao` 或 `TIDao` 接口：

```go
type UserDao struct {
    // 基础实现
}

func (d *UserDao) DB() gdb.DB {
    return g.DB()
}

func (d *UserDao) Table() string {
    return "user"
}

// 实现其他必要方法...
```

### 自定义钩子处理器

注册钩子处理器，在数据库操作前后执行自定义逻辑：

```go
hookHandler := &gdb.HookHandler{
    Before: func(ctx context.Context, in *gdb.HookBeforeInput) (result gdb.Result, err error) {
        // 操作前逻辑
        return
    },
    After: func(ctx context.Context, in *gdb.HookAfterInput) (err error) {
        // 操作后逻辑
        return
    },
}

// 注册钩子
model = daoctl.RegisterDaoHook(model, hookHandler)
```

## 注意事项

1. 合理配置缓存，避免缓存过多导致内存压力
2. 对于频繁变动的数据，考虑禁用缓存
3. 事务操作务必正确处理错误并回滚
4. 避免过于复杂的扩展查询条件，影响代码可读性
5. 使用泛型接口时，确保类型参数正确

## 贡献指南

欢迎提交 Pull Request 或 Issue 报告问题。在提交 PR 前，请确保：

1. 代码符合项目的编码规范
2. 添加了必要的测试用例
3. 文档已更新

## 许可证

本项目采用 MIT 许可证。
