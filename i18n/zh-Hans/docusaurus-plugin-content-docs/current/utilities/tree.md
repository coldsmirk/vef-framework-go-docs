---
sidebar_position: 7
---

# Tree

`tree` 包提供泛型树构建器，将扁平切片转换为层级树结构。

## 快速开始

```go
import "github.com/coldsmirk/vef-framework-go/tree"

type Department struct {
    ID       string        `json:"id"`
    ParentID *string       `json:"parentId"`
    Name     string        `json:"name"`
    Children []Department  `json:"children"`
}

// 定义适配器
adapter := tree.Adapter[Department]{
    GetID:       func(d Department) string { return d.ID },
    GetParentID: func(d Department) *string { return d.ParentID },
    GetChildren: func(d Department) []Department { return d.Children },
    SetChildren: func(d *Department, children []Department) { d.Children = children },
}

// 从扁平切片构建树
roots := tree.Build(flatDepartments, adapter)
```

## 适配器

`Adapter[T]` 结构体定义了树构建器如何访问节点属性：

```go
type Adapter[T any] struct {
    GetID       func(T) string      // 提取节点 ID
    GetParentID func(T) *string     // 提取父级 ID（nil = 根节点）
    GetChildren func(T) []T         // 获取子节点切片
    SetChildren func(*T, []T)       // 设置子节点切片
}
```

关键规则：
- `ParentID` 为 `nil` 的节点被视为根节点
- `ParentID` 引用不存在的节点时也被视为根节点
- 通过已访问标记防止循环引用

## 查找节点

### FindNode

在树中按 ID 搜索节点：

```go
node, found := tree.FindNode(roots, "dept-123", adapter)
if found {
    fmt.Println(node.Name)
}
```

### FindNodePath

获取从根节点到目标节点的完整路径：

```go
path, found := tree.FindNodePath(roots, "dept-456", adapter)
if found {
    for _, node := range path {
        fmt.Println(node.Name) // 输出：根节点 → 父节点 → dept-456
    }
}
```

## 框架集成

`tree` 包被 CRUD 的 `FindTree` 构建器使用。`NewFindTree[T, S]` 需要一个 `func([]T) []T` 形状的 builder，所以你写一个薄包装，把模型的 adapter 闭进去：

```go
func buildDepartmentTree(flat []Department) []Department {
    adapter := tree.Adapter[Department]{
        GetID:       func(d Department) string { return d.ID },
        GetParentID: func(d Department) string { return d.ParentID },
        SetChildren: func(d *Department, children []Department) { d.Children = children },
        // GetChildren 只在调用 tree.FindNode / tree.FindNodePath 时使用；
        // tree.Build 本身不会调它。
        GetChildren: func(d Department) []Department { return d.Children },
    }
    return tree.Build(flat, adapter)
}

// 然后把这个 wrapper 传给 CRUD builder。
crud.NewFindTree[Department, DepartmentSearch](buildDepartmentTree)
```

> `tree.Build` 的签名是 `Build[T any](nodes []T, adapter Adapter[T]) []T`，多了 `adapter` 参数，不能直接传给 `NewFindTree`，需要靠上面的 wrapper 桥接。
> 当 `GetID` 返回 `""` 时，该节点会被跳过索引；它的 children 也不会被填充。
