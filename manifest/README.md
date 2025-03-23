# 清单模块 (manifest)

清单模块提供了用于描述和管理应用程序组件、服务和依赖关系的工具，帮助构建可扩展的模块化应用程序。

## 主要功能

清单模块实现了以下主要功能：

1. **应用程序组件注册**：提供机制注册和发现应用程序组件
2. **依赖关系管理**：管理组件之间的依赖关系，确保正确的初始化顺序
3. **服务发现**：允许应用程序组件相互发现和通信
4. **版本控制**：跟踪组件版本，确保兼容性
5. **配置管理**：管理组件的配置信息

## 核心组件

### 1. 组件清单 (ComponentManifest)

描述单个组件的信息：

```go
// 组件清单
type ComponentManifest struct {
    Name        string            // 组件名称
    Version     string            // 组件版本
    Description string            // 组件描述
    Author      string            // 作者信息
    Dependencies []string         // 依赖的其他组件
    Provides    []string          // 提供的服务或功能
    Requires    []string          // 依赖的服务或功能
    Config      map[string]string // 配置信息
}
```

### 2. 应用程序清单 (AppManifest)

管理整个应用程序的组件集合：

```go
// 应用程序清单
type AppManifest struct {
    AppName     string                      // 应用程序名称
    Version     string                      // 应用程序版本
    Components  map[string]ComponentManifest // 组件集合
    EntryPoints []string                    // 入口点组件
}
```

### 3. 清单管理器 (ManifestManager)

负责注册、发现和管理组件：

```go
// 清单管理器接口
type ManifestManager interface {
    // 注册组件
    RegisterComponent(manifest ComponentManifest) error
    
    // 获取组件
    GetComponent(name string) (ComponentManifest, error)
    
    // 获取所有组件
    GetAllComponents() map[string]ComponentManifest
    
    // 解析依赖关系
    ResolveDependencies() ([]string, error)
    
    // 验证组件兼容性
    ValidateCompatibility() error
    
    // 初始化所有组件
    InitializeComponents(ctx context.Context) error
}
```

## 使用示例

### 创建和注册组件

```go
package main

import (
    "context"
    "fmt"
    "github.com/kysion/base-library/manifest"
)

// 创建用户组件清单
func createUserComponentManifest() manifest.ComponentManifest {
    return manifest.ComponentManifest{
        Name:        "user-service",
        Version:     "1.0.0",
        Description: "用户管理服务",
        Author:      "Kysion Team",
        Dependencies: []string{"db-service"},
        Provides:    []string{"user-management"},
        Requires:    []string{"database-connection"},
        Config: map[string]string{
            "user.table": "users",
            "user.cache": "true",
        },
    }
}

// 创建数据库组件清单
func createDbComponentManifest() manifest.ComponentManifest {
    return manifest.ComponentManifest{
        Name:        "db-service",
        Version:     "1.0.0",
        Description: "数据库服务",
        Author:      "Kysion Team",
        Dependencies: []string{},
        Provides:    []string{"database-connection"},
        Requires:    []string{},
        Config: map[string]string{
            "db.host": "localhost",
            "db.port": "3306",
        },
    }
}

func main() {
    // 创建清单管理器
    manager := manifest.NewManager()
    
    // 注册组件
    userManifest := createUserComponentManifest()
    dbManifest := createDbComponentManifest()
    
    manager.RegisterComponent(dbManifest)  // 先注册数据库组件
    manager.RegisterComponent(userManifest) // 再注册用户组件
    
    // 解析依赖关系，获取初始化顺序
    initOrder, err := manager.ResolveDependencies()
    if err != nil {
        fmt.Printf("解析依赖关系失败: %v\n", err)
        return
    }
    
    // 显示初始化顺序
    fmt.Println("组件初始化顺序:")
    for i, componentName := range initOrder {
        fmt.Printf("%d. %s\n", i+1, componentName)
    }
    
    // 验证组件兼容性
    if err := manager.ValidateCompatibility(); err != nil {
        fmt.Printf("组件兼容性校验失败: %v\n", err)
        return
    }
    
    // 初始化所有组件
    ctx := context.Background()
    if err := manager.InitializeComponents(ctx); err != nil {
        fmt.Printf("组件初始化失败: %v\n", err)
        return
    }
    
    fmt.Println("所有组件初始化成功!")
}
```

### 创建应用程序清单

```go
// 创建应用程序清单
func createAppManifest() manifest.AppManifest {
    // 创建组件清单
    userManifest := createUserComponentManifest()
    dbManifest := createDbComponentManifest()
    
    // 创建应用程序清单
    appManifest := manifest.AppManifest{
        AppName: "my-application",
        Version: "1.0.0",
        Components: map[string]manifest.ComponentManifest{
            userManifest.Name: userManifest,
            dbManifest.Name:   dbManifest,
        },
        EntryPoints: []string{"user-service"},
    }
    
    return appManifest
}

// 使用应用程序清单启动应用
func StartApplication() error {
    appManifest := createAppManifest()
    
    // 创建清单管理器
    manager := manifest.NewManagerFromApp(appManifest)
    
    // 解析依赖关系
    initOrder, err := manager.ResolveDependencies()
    if err != nil {
        return fmt.Errorf("解析依赖关系失败: %v", err)
    }
    
    // 初始化组件
    ctx := context.Background()
    for _, componentName := range initOrder {
        component, _ := manager.GetComponent(componentName)
        fmt.Printf("初始化组件: %s (版本 %s)\n", component.Name, component.Version)
        
        // 这里调用组件的初始化方法
        // ...
    }
    
    // 启动入口点组件
    for _, entryPoint := range appManifest.EntryPoints {
        fmt.Printf("启动入口点组件: %s\n", entryPoint)
        
        // 这里调用入口点组件的启动方法
        // ...
    }
    
    return nil
}
```

### 组件配置管理

```go
// 获取组件配置
func GetComponentConfig(manager manifest.ManifestManager, componentName, configKey string) (string, error) {
    component, err := manager.GetComponent(componentName)
    if err != nil {
        return "", fmt.Errorf("获取组件失败: %v", err)
    }
    
    value, exists := component.Config[configKey]
    if !exists {
        return "", fmt.Errorf("配置项 %s 不存在", configKey)
    }
    
    return value, nil
}

// 使用组件配置
func UseComponentConfig() {
    manager := manifest.NewManager()
    // 注册组件...
    
    // 获取数据库主机配置
    dbHost, err := GetComponentConfig(manager, "db-service", "db.host")
    if err != nil {
        fmt.Printf("获取数据库主机配置失败: %v\n", err)
        return
    }
    
    fmt.Printf("数据库主机: %s\n", dbHost)
    
    // 连接数据库
    // db, err := sql.Open("mysql", fmt.Sprintf("root:password@tcp(%s:3306)/database", dbHost))
    // ...
}
```

## 最佳实践

1. **组件粒度**：选择合适的组件粒度，既不要过大导致不灵活，也不要过小导致过度复杂
2. **依赖声明**：明确声明组件的依赖关系，避免隐式依赖
3. **版本管理**：遵循语义化版本规范，方便管理组件之间的兼容性
4. **配置管理**：通过配置实现组件的可配置性，避免硬编码
5. **文档记录**：为每个组件提供清晰的文档，包括功能描述、依赖关系和配置选项
