---
sidebar_position: 8
---

# reflectx

The `reflectx` package provides reflection and conversion helpers: type-safe
casting, type-compatibility checks, string-field accessors, value predicates,
and a struct visitor/traversal engine. It backs framework internals such as
`mapx`, `copier`, and `search`, and is available directly for application code
that needs the same reflection primitives.

## API Reference

| API group | Public APIs |
| --- | --- |
| cast aliases | `reflectx.ToString`, `reflectx.ToStringE`, `reflectx.ToInt`, `reflectx.ToIntE`, `reflectx.ToInt8`, `reflectx.ToInt8E`, `reflectx.ToInt16`, `reflectx.ToInt16E`, `reflectx.ToInt32`, `reflectx.ToInt32E`, `reflectx.ToInt64`, `reflectx.ToInt64E`, `reflectx.ToUint`, `reflectx.ToUintE`, `reflectx.ToUint8`, `reflectx.ToUint8E`, `reflectx.ToUint16`, `reflectx.ToUint16E`, `reflectx.ToUint32`, `reflectx.ToUint32E`, `reflectx.ToUint64`, `reflectx.ToUint64E`, `reflectx.ToFloat32`, `reflectx.ToFloat32E`, `reflectx.ToFloat64`, `reflectx.ToFloat64E`, `reflectx.ToBool`, `reflectx.ToBoolE` |
| decimal conversion | `reflectx.ToDecimal`, `reflectx.ToDecimalE` |
| type compatibility and methods | `reflectx.Indirect`, `reflectx.IsPointerToStruct`, `reflectx.IsSimilarType`, `reflectx.IsTypeCompatible`, `reflectx.ConvertValue`, `reflectx.ErrCannotConvertType`, `reflectx.FindMethod`, `reflectx.CollectMethods` |
| string field helpers | `reflectx.IsStringType`, `reflectx.IsStringSliceType`, `reflectx.IsStringMapType`, `reflectx.GetStringValue`, `reflectx.SetStringValue`, `reflectx.GetStringSliceValue`, `reflectx.SetStringSliceValue`, `reflectx.GetStringMapValue`, `reflectx.SetStringMapValue` |
| value helpers | `reflectx.IsEmpty`, `reflectx.IsNotEmpty`, `reflectx.IsNumeric`, `reflectx.IsInteger`, `reflectx.IsSignedInt`, `reflectx.IsUnsignedInt`, `reflectx.IsFloat`, `reflectx.Equal`, `reflectx.Contains` |
| visitor actions, callbacks, traversal modes, and options | `reflectx.VisitAction`, `reflectx.Continue`, `reflectx.Stop`, `reflectx.SkipChildren`, `reflectx.TraversalMode`, `reflectx.DepthFirst`, `reflectx.BreadthFirst`, `reflectx.TagConfig`, `reflectx.VisitorConfig`, `reflectx.VisitorOption`, `reflectx.Visitor`, `reflectx.StructVisitor`, `reflectx.FieldVisitor`, `reflectx.MethodVisitor`, `reflectx.TypeVisitor`, `reflectx.StructTypeVisitor`, `reflectx.FieldTypeVisitor`, `reflectx.MethodTypeVisitor`, `reflectx.VisitOf`, `reflectx.Visit`, `reflectx.VisitType`, `reflectx.VisitFor`, `reflectx.WithTraversalMode`, `reflectx.WithDisableRecursive`, `reflectx.WithDiveTag`, `reflectx.WithMaxDepth` |

Exported fields:

| Type | Fields |
| --- | --- |
| `reflectx.TagConfig` | `TagConfig.Name`, `TagConfig.Value` |
| `reflectx.VisitorConfig` | `VisitorConfig.TraversalMode`, `VisitorConfig.Recursive`, `VisitorConfig.DiveTag`, `VisitorConfig.MaxDepth` |
| `reflectx.Visitor` | `Visitor.VisitStruct`, `Visitor.VisitField`, `Visitor.VisitMethod` |
| `reflectx.TypeVisitor` | `TypeVisitor.VisitStructType`, `TypeVisitor.VisitFieldType`, `TypeVisitor.VisitMethodType` |

## Cast Aliases And Decimal Conversion

Cast aliases are pass-through helpers from `github.com/spf13/cast`: `To*E`
variants return conversion errors, while non-E variants return the
destination type's zero value on failure.

```go
import "github.com/coldsmirk/vef-framework-go/reflectx"

n, err := reflectx.ToInt64E("42")  // 42, nil
n = reflectx.ToInt64("not-a-number") // 0 (error discarded)
```

`reflectx.ToDecimalE` returns `decimal.Zero` with no error for nil values and
nil pointer/interface values; otherwise it recursively dereferences
pointer/interface values and delegates to `decimal.NewFromAny`.
`reflectx.ToDecimal` discards the error and returns `decimal.Zero` on
failure.

## Type Compatibility And Methods

`reflectx.Indirect` dereferences one pointer type. `reflectx.IsPointerToStruct`
requires a non-nil `reflect.Type` whose kind is pointer and whose element kind
is struct. `reflectx.IsSimilarType` is true for identical types, or for
generic instantiations with the same `PkgPath` and the same base type name
before `[`. `reflectx.IsTypeCompatible` accepts exact/assignable types,
interface targets, pointer-to-pointer element compatibility,
value-to-pointer element assignment, and pointer-to-value element assignment.
`reflectx.ConvertValue` mirrors those pointer/value conversions, returns zero
target values for nil pointer inputs, allocates new pointers for
value-to-pointer and pointer-to-pointer conversions, and wraps unsupported
conversions with `reflectx.ErrCannotConvertType`. `reflectx.FindMethod` checks
the value first, then an addressable pointer copy for non-pointer values.
`reflectx.CollectMethods` dereferences pointers, returns an empty map for nil
or non-struct values, and collects promoted pointer/value methods by name
from the pointer method set.

## String Field Helpers

String field helpers require callers to pass a valid `reflect.Type` or
`reflect.Value`. `reflectx.IsStringType` accepts `string` and `*string` only;
`reflectx.IsStringSliceType` accepts `[]string`; `reflectx.IsStringMapType`
accepts `map[string]string`. Getter helpers return `false` for incompatible
types and nil string pointers/slices/maps. Setter helpers are no-ops for
incompatible values. `reflectx.SetStringValue` replaces a `*string` with a
fresh `*string` pointer instead of mutating an existing pointee.

## Value Helpers

`reflectx.IsEmpty` treats nil, invalid values, zero scalars, empty
strings/arrays/slices/maps, nil pointer/interface/channel/function values, and
zero structs as empty. It has a special `*string` rule: a non-nil pointer to
an empty string is empty, while other non-nil pointers are not empty.
`reflectx.Equal` compares signed integers only with signed integers, unsigned
integers only with unsigned integers, and floats only with floats;
cross-category numeric comparisons return false. Exact-type comparable values
use `==`, and same-type non-comparable nil-able values are equal only when
both are nil. `reflectx.Contains` supports string substring checks with
string elements, slices/arrays via `reflectx.Equal`, and maps by key lookup
with convertible map keys.

## Visitor Traversal

Visitor traversal defaults to depth-first order. `reflectx.DepthFirst = 0`,
`reflectx.BreadthFirst = 1`, `reflectx.Continue = 0`, `reflectx.Stop = 1`, and
`reflectx.SkipChildren = 2`. By default `VisitorConfig.TraversalMode` is
`DepthFirst`, `VisitorConfig.Recursive` is true, `VisitorConfig.DiveTag` is
`TagConfig{Name: "visit", Value: "dive"}`, and `VisitorConfig.MaxDepth == 0`
means unlimited depth. Anonymous struct fields recurse automatically; named
struct fields recurse only when their tag matches `VisitorConfig.DiveTag`.
`reflectx.WithTraversalMode` selects `DepthFirst` or `BreadthFirst`;
`reflectx.WithDisableRecursive` sets `VisitorConfig.Recursive = false`;
`reflectx.WithDiveTag` replaces the tag selector; `reflectx.WithMaxDepth` sets
the depth cap, and traversal stops descending when `depth >= MaxDepth`.
`reflectx.SkipChildren` skips recursion for the current field and
`reflectx.Stop` aborts traversal. `reflectx.Visit` and `reflectx.VisitOf` use
value callbacks, `reflectx.VisitType` and `reflectx.VisitFor` use type-only
callbacks, and invalid, nil pointer, and non-struct inputs are ignored. Field
callbacks receive `StructField.Index` rewritten to the absolute index path,
and method callbacks are visited through the pointer method set.

```go
import (
	"fmt"
	"reflect"

	"github.com/coldsmirk/vef-framework-go/reflectx"
)

type Profile struct {
	Name string `visit:"dive"`
	Bio  string
}

reflectx.VisitOf(myProfile, reflectx.Visitor{
	VisitField: func(field reflect.StructField, value reflect.Value, depth int) reflectx.VisitAction {
		fmt.Println(field.Name, value.Interface())
		return reflectx.Continue
	},
})
```

## See Also

- [Query Builder](../data-access/query-builder) — the `search` package parses
  request filters using `reflectx` type-compatibility helpers.
