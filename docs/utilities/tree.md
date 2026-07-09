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

Public surface:

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

## Build Contract

`Build` converts a flat slice into nested roots.

Key rules:

- `Build(nil, adapter)` and `Build([]T{}, adapter)` return a non-nil empty
  slice (`[]T{}`)
- `GetID` values are raw string keys; special characters and Unicode are not
  normalized or escaped
- `GetID` is expected to return a unique non-empty ID; empty-ID nodes are not
  indexed for parent lookup and their own children are not populated
- empty-ID nodes can still appear in the returned roots or in a parent's
  children when their parent relationship puts them there
- `GetParentID(node) == nil` makes the node a root
- a non-nil parent ID that does not exist in the indexed node map also makes the
  node a root
- closed cycles whose parent chain never reaches a root are omitted from the
  returned roots
- `Build` uses visited tracking while assigning children so cyclic parent data
  does not recurse forever
- `Build` calls `SetChildren` on elements of the input slice and returns value
  copies of the root elements; treat the input slice elements as mutable
- `GetChildren` is not called by `Build`, so a wrapper that only builds a tree
  can omit it
- missing adapter callbacks panic naturally when the operation reaches them

## Finding Nodes

### FindNode

Search for a node by ID in a tree:

```go
node, found := tree.FindNode(roots, "dept-123", adapter)
if found {
    fmt.Println(node.Name)
}
```

Contract:

- an empty `targetID` returns the zero value of `T` and `false`
- a missing target also returns the zero value of `T` and `false`
- traversal is depth-first and follows the slices returned by `GetChildren`
- duplicate IDs are not de-duplicated; the first traversal match wins
- `FindNode` does not add cycle protection around `GetChildren`, so pass an
  acyclic tree

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

Contract:

- an empty `targetID`, a missing target, or an empty tree returns `nil, false`
- a found target returns the full root-to-node path and `true`
- traversal is depth-first and follows the slices returned by `GetChildren`
- `FindNodePath` does not add cycle protection around `GetChildren`, so pass an
  acyclic tree

## Framework Integration

The `tree` package is used by the CRUD `FindTree` builder. `NewFindTree[T, S]` requires a builder of signature `func([]T) []T`, so you provide a thin wrapper that closes over the model's adapter:

```go
func buildDepartmentTree(flat []Department) []Department {
    adapter := tree.Adapter[Department]{
        GetID:       func(d Department) string { return d.ID },
        GetParentID: func(d Department) *string { return d.ParentID },
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
