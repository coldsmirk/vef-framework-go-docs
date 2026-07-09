---
sidebar_position: 8
---

# reflectx

`reflectx` 包提供反射与转换 helper：类型安全的 casting、类型兼容性检查、字符串字段
accessor、value predicate，以及一个 struct visitor/traversal engine。框架内部的
`mapx`、`copier`、`search` 都依赖它；需要相同反射原语的应用代码也可以直接使用它。

## API 参考

| API 组 | Public APIs |
| --- | --- |
| cast aliases | `reflectx.ToString`, `reflectx.ToStringE`, `reflectx.ToInt`, `reflectx.ToIntE`, `reflectx.ToInt8`, `reflectx.ToInt8E`, `reflectx.ToInt16`, `reflectx.ToInt16E`, `reflectx.ToInt32`, `reflectx.ToInt32E`, `reflectx.ToInt64`, `reflectx.ToInt64E`, `reflectx.ToUint`, `reflectx.ToUintE`, `reflectx.ToUint8`, `reflectx.ToUint8E`, `reflectx.ToUint16`, `reflectx.ToUint16E`, `reflectx.ToUint32`, `reflectx.ToUint32E`, `reflectx.ToUint64`, `reflectx.ToUint64E`, `reflectx.ToFloat32`, `reflectx.ToFloat32E`, `reflectx.ToFloat64`, `reflectx.ToFloat64E`, `reflectx.ToBool`, `reflectx.ToBoolE` |
| decimal conversion | `reflectx.ToDecimal`, `reflectx.ToDecimalE` |
| type compatibility and methods | `reflectx.Indirect`, `reflectx.IsPointerToStruct`, `reflectx.IsSimilarType`, `reflectx.IsTypeCompatible`, `reflectx.ConvertValue`, `reflectx.ErrCannotConvertType`, `reflectx.FindMethod`, `reflectx.CollectMethods` |
| string field helpers | `reflectx.IsStringType`, `reflectx.IsStringSliceType`, `reflectx.IsStringMapType`, `reflectx.GetStringValue`, `reflectx.SetStringValue`, `reflectx.GetStringSliceValue`, `reflectx.SetStringSliceValue`, `reflectx.GetStringMapValue`, `reflectx.SetStringMapValue` |
| value helpers | `reflectx.IsEmpty`, `reflectx.IsNotEmpty`, `reflectx.IsNumeric`, `reflectx.IsInteger`, `reflectx.IsSignedInt`, `reflectx.IsUnsignedInt`, `reflectx.IsFloat`, `reflectx.Equal`, `reflectx.Contains` |
| visitor actions, callbacks, traversal modes, and options | `reflectx.VisitAction`, `reflectx.Continue`, `reflectx.Stop`, `reflectx.SkipChildren`, `reflectx.TraversalMode`, `reflectx.DepthFirst`, `reflectx.BreadthFirst`, `reflectx.TagConfig`, `reflectx.VisitorConfig`, `reflectx.VisitorOption`, `reflectx.Visitor`, `reflectx.StructVisitor`, `reflectx.FieldVisitor`, `reflectx.MethodVisitor`, `reflectx.TypeVisitor`, `reflectx.StructTypeVisitor`, `reflectx.FieldTypeVisitor`, `reflectx.MethodTypeVisitor`, `reflectx.VisitOf`, `reflectx.Visit`, `reflectx.VisitType`, `reflectx.VisitFor`, `reflectx.WithTraversalMode`, `reflectx.WithDisableRecursive`, `reflectx.WithDiveTag`, `reflectx.WithMaxDepth` |

Exported fields：

| Type | Fields |
| --- | --- |
| `reflectx.TagConfig` | `TagConfig.Name`, `TagConfig.Value` |
| `reflectx.VisitorConfig` | `VisitorConfig.TraversalMode`, `VisitorConfig.Recursive`, `VisitorConfig.DiveTag`, `VisitorConfig.MaxDepth` |
| `reflectx.Visitor` | `Visitor.VisitStruct`, `Visitor.VisitField`, `Visitor.VisitMethod` |
| `reflectx.TypeVisitor` | `TypeVisitor.VisitStructType`, `TypeVisitor.VisitFieldType`, `TypeVisitor.VisitMethodType` |

## Cast Aliases 与 Decimal 转换

Cast aliases 是来自 `github.com/spf13/cast` 的 pass-through helpers：`To*E`
variants 会返回 conversion errors，non-E variants 在失败时返回目标类型的
zero value。

```go
import "github.com/coldsmirk/vef-framework-go/reflectx"

n, err := reflectx.ToInt64E("42")   // 42, nil
n = reflectx.ToInt64("not-a-number") // 0（error 被丢弃）
```

`reflectx.ToDecimalE` 对 nil values 和 nil pointer/interface values 返回
`decimal.Zero` 且没有 error；其他 pointer/interface values 会递归
dereference，然后委托 `decimal.NewFromAny`。`reflectx.ToDecimal` 会丢弃 error，
失败时返回 `decimal.Zero`。

## 类型兼容性与方法

`reflectx.Indirect` 只 dereference 一层 pointer type。
`reflectx.IsPointerToStruct` 要求传入非 nil `reflect.Type`，其 kind 是 pointer，
且 element kind 是 struct。`reflectx.IsSimilarType` 对相同 type 返回 true；
对 generic instantiations，只有 `PkgPath` 相同且 `[` 之前的 base type name
相同时才返回 true。`reflectx.IsTypeCompatible` 接受 exact/assignable types、
interface targets、pointer-to-pointer element compatibility、value-to-pointer
element assignment 和 pointer-to-value element assignment。
`reflectx.ConvertValue` 对这些 pointer/value conversions 做对应转换，对 nil
pointer inputs 返回 zero target values，对 value-to-pointer 和
pointer-to-pointer conversions 分配新 pointer，并用 `reflectx.ErrCannotConvertType`
包装不支持的转换。`reflectx.FindMethod` 先查 value，再对 non-pointer value
使用 addressable pointer copy。`reflectx.CollectMethods` 会 dereference pointers，
nil 或 non-struct values 返回 empty map，并从 pointer method set 按名称收集
promoted pointer/value methods。

## 字符串字段 Helper

String field helpers 要求调用方传入有效的 `reflect.Type` 或 `reflect.Value`。
`reflectx.IsStringType` 只接受 `string` 和 `*string`；`reflectx.IsStringSliceType`
只接受 `[]string`；`reflectx.IsStringMapType` 只接受 `map[string]string`。
getter helpers 对 incompatible types 和 nil string pointers/slices/maps 返回
`false`。setter helpers 对 incompatible values 是 no-op。
`reflectx.SetStringValue` 会把 `*string` 替换成 fresh `*string` pointer，而不是
修改已有 pointee。

## Value Helper

`reflectx.IsEmpty` 把 nil、invalid values、zero scalars、empty
strings/arrays/slices/maps、nil pointer/interface/channel/function values 和
zero structs 视为空。它对 `*string` 有特殊规则：指向 empty string 的非 nil
pointer 也是 empty，但其他非 nil pointers 不算 empty。`reflectx.Equal` 只在
signed integers、unsigned integers、floats 各自 category 内做 numeric 比较；
cross-category numeric comparisons 返回 false。exact-type comparable values
使用 `==`；same-type non-comparable nil-able values 只有在两边都 nil 时才相等。
`reflectx.Contains` 支持 string substring checks（element 必须是 string）、
slices/arrays 通过 `reflectx.Equal` 比较，以及 maps 通过 key lookup，且支持
convertible map keys。

## Visitor Traversal

Visitor traversal 默认使用 depth-first order。`reflectx.DepthFirst = 0`、
`reflectx.BreadthFirst = 1`、`reflectx.Continue = 0`、
`reflectx.Stop = 1`、`reflectx.SkipChildren = 2`。默认
`VisitorConfig.TraversalMode` 为 `DepthFirst`，`VisitorConfig.Recursive` 为 true，`VisitorConfig.DiveTag` 是
`TagConfig{Name: "visit", Value: "dive"}`，`VisitorConfig.MaxDepth == 0` 表示
unlimited depth。Anonymous struct fields 会自动 recurse；named struct fields
只有在 tag 匹配 `VisitorConfig.DiveTag` 时才 recurse。
`reflectx.WithTraversalMode` 选择 `DepthFirst` 或 `BreadthFirst`；
`reflectx.WithDisableRecursive` 设置 `VisitorConfig.Recursive = false`；
`reflectx.WithDiveTag` 替换 tag selector；`reflectx.WithMaxDepth` 设置 depth
cap，且 traversal 在 `depth >= MaxDepth` 时停止向下。`reflectx.SkipChildren`
跳过当前 field 的 recursion，`reflectx.Stop` 中止 traversal。`reflectx.Visit`
和 `reflectx.VisitOf` 使用 value callbacks，`reflectx.VisitType` 和
`reflectx.VisitFor` 使用 type-only callbacks；invalid、nil pointer 和
non-struct inputs 会被忽略。Field callbacks 收到的 `StructField.Index` 会被重写为
absolute index path，method callbacks 会通过 pointer method set 访问。

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

## 参见

- [Query Builder](../data-access/query-builder) —— `search` 包解析请求过滤条件
  时使用了 `reflectx` 的类型兼容性 helper。
