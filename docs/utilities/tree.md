---
sidebar_position: 7
---

# Tree

The `tree` package provides a generic tree builder that converts flat slices into hierarchical tree structures.

## Quick Start

```go
import "github.com/coldsmirk/vef-framework-go/tree"

type Department struct {
    ID       string        `json:"id"`
    ParentID *string       `json:"parentId"`
    Name     string        `json:"name"`
    Children []Department  `json:"children"`
}

// Define the adapter
adapter := tree.Adapter[Department]{
    GetID:       func(d Department) string { return d.ID },
    GetParentID: func(d Department) *string { return d.ParentID },
    GetChildren: func(d Department) []Department { return d.Children },
    SetChildren: func(d *Department, children []Department) { d.Children = children },
}

// Build tree from flat slice
roots := tree.Build(flatDepartments, adapter)
```

## Adapter

The `Adapter[T]` struct defines how the tree builder accesses node properties:

```go
type Adapter[T any] struct {
    GetID       func(T) string      // Extract node ID
    GetParentID func(T) *string     // Extract parent ID (nil = root node)
    GetChildren func(T) []T         // Get children slice
    SetChildren func(*T, []T)       // Set children slice
}
```

Key rules:
- Nodes with `nil` parent ID are treated as roots
- Nodes whose parent ID references a non-existent node are also treated as roots
- Circular references are protected against via visited tracking

## Finding Nodes

### FindNode

Search for a node by ID in a tree:

```go
node, found := tree.FindNode(roots, "dept-123", adapter)
if found {
    fmt.Println(node.Name)
}
```

### FindNodePath

Get the full path from root to a target node:

```go
path, found := tree.FindNodePath(roots, "dept-456", adapter)
if found {
    for _, node := range path {
        fmt.Println(node.Name) // prints: "Root" → "Parent" → "dept-456"
    }
}
```

## Framework Integration

The `tree` package is used by the CRUD `FindTree` builder. `NewFindTree[T, S]` requires a builder of signature `func([]T) []T`, so you provide a thin wrapper that closes over the model's adapter:

```go
func buildDepartmentTree(flat []Department) []Department {
    adapter := tree.Adapter[Department]{
        GetID:       func(d Department) string { return d.ID },
        GetParentID: func(d Department) string { return d.ParentID },
        SetChildren: func(d *Department, children []Department) { d.Children = children },
        // GetChildren is only needed if you intend to call tree.FindNode /
        // tree.FindNodePath on the resulting tree; tree.Build itself doesn't use it.
        GetChildren: func(d Department) []Department { return d.Children },
    }
    return tree.Build(flat, adapter)
}

// Then plug the wrapper into the CRUD builder.
crud.NewFindTree[Department, DepartmentSearch](buildDepartmentTree)
```

> `tree.Build` has signature `Build[T any](nodes []T, adapter Adapter[T]) []T`, so it cannot be passed directly to `NewFindTree` — the wrapper bridges the two signatures.
> Nodes whose `GetID` returns `""` are skipped during indexing and won't have their children populated.
