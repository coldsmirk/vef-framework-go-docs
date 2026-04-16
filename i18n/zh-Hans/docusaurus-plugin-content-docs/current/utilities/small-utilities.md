---
sidebar_position: 8
---

# 小工具

本页文档记录了提供常用辅助函数的轻量级工具包。

## `ptr` — 指针助手

用于安全指针操作的泛型工具函数。

### `ptr.Of` — 值转指针

返回值的指针，如果是零值则返回 `nil`：

```go
import "github.com/coldsmirk/vef-framework-go/ptr"

p := ptr.Of("hello")   // *string → "hello"
p = ptr.Of("")          // *string → nil（零值）
p = ptr.Of(42)          // *int → 42
p = ptr.Of(0)           // *int → nil
```

### `ptr.Value` — 指针取值

解引用指针，支持可选回退：

```go
s := ptr.Value(p)                        // 解引用，nil 时返回零值
s = ptr.Value(p, fallback1, fallback2)   // 依次尝试每个回退
```

### `ptr.ValueOrElse` — 懒回退

使用懒加载回退函数解引用：

```go
s := ptr.ValueOrElse(p, func() string {
    return computeDefault()
})
```

### `ptr.Zero` — 零值

返回任意类型的零值：

```go
z := ptr.Zero[string]()  // ""
z := ptr.Zero[int]()     // 0
```

### `ptr.Equal` — 指针相等

按值比较两个指针：

```go
ptr.Equal(a, b)  // 都为 nil → true，一个 nil → false，否则比较值
```

### `ptr.Coalesce` — 第一个非 nil

返回列表中第一个非 nil 的指针：

```go
result := ptr.Coalesce(p1, p2, p3)  // 第一个非 nil，或 nil
```

## 实际用法

`ptr` 包常用于：

- **搜索结构体**：可选筛选字段使用 `*bool`、`*string` 等
- **模型字段**：可空数据库列使用指针类型
- **参数结构体**：可选更新字段

```go
type UserSearch struct {
    api.P
    IsActive *bool   `json:"isActive" search:"eq,column=is_active"`
    DeptID   *string `json:"deptId" search:"eq,column=department_id"`
}

// 设置可选搜索筛选条件
search := UserSearch{
    IsActive: ptr.Of(true),
    DeptID:   ptr.Of("dept-123"),
}
```
