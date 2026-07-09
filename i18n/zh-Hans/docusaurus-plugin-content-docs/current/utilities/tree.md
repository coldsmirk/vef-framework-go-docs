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

公开 surface：

| API | Signature |
| --- | --- |
| `Adapter[T]` | `type Adapter[T any] struct` |
| `Adapter.GetID` | `func(T) string` |
| `Adapter.GetParentID` | `func(T) *string` |
| `Adapter.GetChildren` | `func(T) []T` |
| `Adapter.SetChildren` | `func(*T, []T)` |
| `Build` | `tree.Build[T any](nodes []T, adapter tree.Adapter[T]) []T` |
| `FindNode` | `tree.FindNode[T any](roots []T, targetID string, adapter tree.Adapter[T]) (T, bool)` |
| `FindNodePath` | `tree.FindNodePath[T any](roots []T, targetID string, adapter tree.Adapter[T]) ([]T, bool)` |

## Build 契约

`Build` 将扁平切片转换成嵌套根节点列表。

关键规则：

- `Build(nil, adapter)` 和 `Build([]T{}, adapter)` 返回非 nil 空切片
  (`[]T{}`)
- `GetID` 的值按原始 string key 使用；特殊字符和 Unicode 不会被 normalize
  或 escape
- `GetID` 应返回唯一且非空的 ID；空 ID 节点不会进入 parent lookup 的索引，
  其自身 children 也不会被填充
- 空 ID 节点仍可能因为自身 parent 关系出现在返回的 roots 中，或出现在某个
  parent 的 children 中
- `GetParentID(node) == nil` 时，该节点是 root
- 非 nil parent ID 如果不存在于已索引 node map，也会让该节点成为 root
- 如果一组节点形成 parent chain 永远到不了 root 的闭环，它们不会出现在返回的
  roots 中
- `Build` 在设置 children 时使用 visited tracking，避免 cyclic parent data
  无限递归
- `Build` 会对输入 slice 的元素调用 `SetChildren`，然后返回 root 元素的值拷贝；
  调用方应把输入 slice 元素视为可被修改
- `Build` 不会调用 `GetChildren`，所以只用于构建树的 wrapper 可以不提供它
- adapter callback 缺失时，代码执行到对应 callback 会自然 panic

## 查找节点

### FindNode

在树中按 ID 搜索节点：

```go
node, found := tree.FindNode(roots, "dept-123", adapter)
if found {
    fmt.Println(node.Name)
}
```

契约：

- 空 `targetID` 返回 `T` 的零值和 `false`
- 目标不存在时也返回 `T` 的零值和 `false`
- 遍历是 depth-first，并沿着 `GetChildren` 返回的 slice 继续
- duplicate IDs 不会被去重；返回第一个遍历命中的节点
- `FindNode` 不会围绕 `GetChildren` 额外加 cycle protection，因此应传入无环树

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

契约：

- 空 `targetID`、目标不存在或空树都会返回 `nil, false`
- 命中目标时返回完整 root-to-node path 和 `true`
- 遍历是 depth-first，并沿着 `GetChildren` 返回的 slice 继续
- `FindNodePath` 不会围绕 `GetChildren` 额外加 cycle protection，因此应传入无环树

## 框架集成

`tree` 包被 CRUD 的 `FindTree` 构建器使用。`NewFindTree[T, S]` 需要一个 `func([]T) []T` 形状的 builder，所以你写一个薄包装，把模型的 adapter 闭进去：

```go
func buildDepartmentTree(flat []Department) []Department {
    adapter := tree.Adapter[Department]{
        GetID:       func(d Department) string { return d.ID },
        GetParentID: func(d Department) *string { return d.ParentID },
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
