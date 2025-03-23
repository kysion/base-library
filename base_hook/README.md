# 基础钩子 (base_hook)

基础钩子模块提供了一种灵活的钩子（Hook）机制，允许在特定事件或生命周期点注入自定义行为，实现功能扩展和解耦。

## 主要功能

基础钩子模块实现了以下主要功能：

1. **事件驱动钩子系统**：允许在关键操作点注册和触发自定义行为
2. **生命周期钩子**：支持程序或组件生命周期的各个阶段注入行为
3. **钩子优先级**：可以为钩子指定执行优先级，控制多个钩子的执行顺序
4. **钩子链**：支持钩子的链式调用和执行结果传递
5. **异步钩子**：支持异步执行钩子，提高性能

## 钩子类型

系统提供了多种类型的钩子：

### 1. 系统钩子 (SystemHook)

用于系统级别的事件处理：

```go
// 系统钩子接口
type SystemHook interface {
    // 启动前钩子
    BeforeStart(ctx context.Context) error
    // 启动后钩子
    AfterStart(ctx context.Context) error
    // 关闭前钩子
    BeforeStop(ctx context.Context) error
    // 关闭后钩子
    AfterStop(ctx context.Context) error
}
```

### 2. 数据钩子 (DataHook)

用于数据操作的各个阶段：

```go
// 数据钩子接口
type DataHook interface {
    // 数据创建前
    BeforeCreate(ctx context.Context, data interface{}) error
    // 数据创建后
    AfterCreate(ctx context.Context, data interface{}) error
    // 数据更新前
    BeforeUpdate(ctx context.Context, data interface{}) error
    // 数据更新后
    AfterUpdate(ctx context.Context, data interface{}) error
    // 数据删除前
    BeforeDelete(ctx context.Context, data interface{}) error
    // 数据删除后
    AfterDelete(ctx context.Context, data interface{}) error
}
```

### 3. 自定义钩子 (CustomHook)

可以根据需要实现自定义钩子接口：

```go
// 自定义钩子示例
type AuthHook interface {
    // 登录前钩子
    BeforeLogin(ctx context.Context, credentials interface{}) error
    // 登录后钩子
    AfterLogin(ctx context.Context, user interface{}, token string) error
    // 登出前钩子
    BeforeLogout(ctx context.Context, userId uint) error
    // 登出后钩子
    AfterLogout(ctx context.Context, userId uint) error
}
```

## 使用示例

### 注册和使用系统钩子

```go
package main

import (
    "context"
    "fmt"
    "github.com/kysion/base-library/base_hook"
)

// 实现系统钩子接口
type MySystemHook struct{}

func (h *MySystemHook) BeforeStart(ctx context.Context) error {
    fmt.Println("系统启动前执行...")
    return nil
}

func (h *MySystemHook) AfterStart(ctx context.Context) error {
    fmt.Println("系统启动后执行...")
    return nil
}

func (h *MySystemHook) BeforeStop(ctx context.Context) error {
    fmt.Println("系统停止前执行...")
    return nil
}

func (h *MySystemHook) AfterStop(ctx context.Context) error {
    fmt.Println("系统停止后执行...")
    return nil
}

func main() {
    // 注册系统钩子
    hookManager := base_hook.NewManager()
    hookManager.RegisterSystemHook(&MySystemHook{}, 100) // 100表示优先级
    
    // 触发系统钩子
    ctx := context.Background()
    hookManager.TriggerBeforeStart(ctx)
    
    // 应用程序主逻辑...
    fmt.Println("应用程序运行中...")
    
    // 触发停止钩子
    hookManager.TriggerBeforeStop(ctx)
    hookManager.TriggerAfterStop(ctx)
}
```

### 实现数据钩子

```go
// 实现数据钩子
type UserDataHook struct{}

func (h *UserDataHook) BeforeCreate(ctx context.Context, data interface{}) error {
    user, ok := data.(*User)
    if !ok {
        return fmt.Errorf("invalid data type")
    }
    
    // 数据验证
    if user.Username == "" {
        return fmt.Errorf("username cannot be empty")
    }
    
    // 密码加密
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)
    
    fmt.Println("用户创建前处理完成")
    return nil
}

func (h *UserDataHook) AfterCreate(ctx context.Context, data interface{}) error {
    user, ok := data.(*User)
    if !ok {
        return fmt.Errorf("invalid data type")
    }
    
    // 创建用户后的操作，如发送欢迎邮件
    fmt.Printf("用户 %s 创建成功，发送欢迎邮件\n", user.Username)
    return nil
}

// 其他方法实现...

// 注册数据钩子
func RegisterUserHooks() {
    hookManager := base_hook.GetManager()
    hookManager.RegisterDataHook("user", &UserDataHook{}, 100)
}

// 在用户服务中使用钩子
func CreateUser(ctx context.Context, user *User) error {
    // 触发创建前钩子
    if err := base_hook.GetManager().TriggerBeforeCreate(ctx, "user", user); err != nil {
        return err
    }
    
    // 保存用户数据
    if err := g.DB().Model("users").Data(user).Insert(); err != nil {
        return err
    }
    
    // 触发创建后钩子
    if err := base_hook.GetManager().TriggerAfterCreate(ctx, "user", user); err != nil {
        return err
    }
    
    return nil
}
```

### 自定义钩子示例

```go
// 定义自定义钩子接口
type PaymentHook interface {
    BeforePayment(ctx context.Context, order *Order) error
    AfterPayment(ctx context.Context, order *Order, paymentResult *PaymentResult) error
}

// 实现自定义钩子
type LogPaymentHook struct{}

func (h *LogPaymentHook) BeforePayment(ctx context.Context, order *Order) error {
    fmt.Printf("准备处理订单 #%d 的支付，金额: %.2f\n", order.Id, order.Amount)
    return nil
}

func (h *LogPaymentHook) AfterPayment(ctx context.Context, order *Order, result *PaymentResult) error {
    fmt.Printf("订单 #%d 支付%s，交易号: %s\n", 
        order.Id, 
        result.Success ? "成功" : "失败", 
        result.TransactionId)
    return nil
}

// 注册自定义钩子
func RegisterPaymentHooks() {
    hookManager := base_hook.GetManager()
    hookManager.RegisterCustomHook("payment", &LogPaymentHook{}, 100)
}

// 在支付服务中使用自定义钩子
func ProcessPayment(ctx context.Context, order *Order) (*PaymentResult, error) {
    // 获取钩子管理器
    hookManager := base_hook.GetManager()
    
    // 触发支付前钩子
    if err := hookManager.TriggerCustomHook(ctx, "payment", "BeforePayment", order); err != nil {
        return nil, err
    }
    
    // 处理支付逻辑
    result := &PaymentResult{
        Success:       true,
        TransactionId: generateTransactionId(),
        Amount:        order.Amount,
        Time:          time.Now(),
    }
    
    // 触发支付后钩子
    if err := hookManager.TriggerCustomHook(ctx, "payment", "AfterPayment", order, result); err != nil {
        return result, err
    }
    
    return result, nil
}
```

## 最佳实践

1. **适度使用钩子**：钩子机制虽然灵活，但过度使用会增加系统复杂性，应当合理使用
2. **错误处理**：钩子执行过程中的错误应当被妥善处理，避免影响主程序流程
3. **优先级设置**：为钩子设置合理的优先级，确保按照期望的顺序执行
4. **文档记录**：为自定义钩子提供清晰的文档，说明钩子的触发时机和用途
5. **避免副作用**：钩子应当尽量避免产生意外的副作用，保持功能的纯粹性
