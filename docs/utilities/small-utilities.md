---
sidebar_position: 8
---

# Small Utilities

This page documents lightweight utility packages that provide common helper functions.

## `ptr` — Pointer Helpers

Generic utility functions for safe pointer operations.

### `ptr.Of` — Value to Pointer

Returns a pointer to the value, or `nil` if the value is zero:

```go
import "github.com/coldsmirk/vef-framework-go/ptr"

p := ptr.Of("hello")   // *string → "hello"
p = ptr.Of("")          // *string → nil (zero value)
p = ptr.Of(42)          // *int → 42
p = ptr.Of(0)           // *int → nil
```

### `ptr.Value` — Pointer to Value

Dereferences a pointer with optional fallbacks:

```go
s := ptr.Value(p)                   // Dereference, or zero if nil
s = ptr.Value(p, fallback1, fallback2) // Try each fallback in order
```

### `ptr.ValueOrElse` — Lazy Fallback

Dereferences with a lazy fallback function:

```go
s := ptr.ValueOrElse(p, func() string {
    return computeDefault()
})
```

### `ptr.Zero` — Zero Value

Returns the zero value of any type:

```go
z := ptr.Zero[string]()  // ""
z := ptr.Zero[int]()     // 0
```

### `ptr.Equal` — Pointer Equality

Compares two pointers by value:

```go
ptr.Equal(a, b)  // Both nil → true, one nil → false, else compare values
```

### `ptr.Coalesce` — First Non-Nil

Returns the first non-nil pointer from a list:

```go
result := ptr.Coalesce(p1, p2, p3)  // First non-nil, or nil
```

## Practical Usage

The `ptr` package is commonly used in:

- **Search structs**: Optional filter fields use `*bool`, `*string`, etc.
- **Model fields**: Nullable database columns use pointer types
- **Params structs**: Optional update fields

```go
type UserSearch struct {
    api.P
    IsActive *bool   `json:"isActive" search:"eq,column=is_active"`
    DeptID   *string `json:"deptId" search:"eq,column=department_id"`
}

// Setting optional search filters
search := UserSearch{
    IsActive: ptr.Of(true),
    DeptID:   ptr.Of("dept-123"),
}
```
