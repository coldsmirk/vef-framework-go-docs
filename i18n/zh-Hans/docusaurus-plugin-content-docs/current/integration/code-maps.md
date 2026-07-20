---
sidebar_position: 4
---

# 码值映射

契约标准化的是*结构*；码值映射标准化的是*值*。一个码值映射（Code Map）在
宿主的标准码与某个外部系统的码之间翻译一个码值集（性别、婚姻状况、订单
状态……）的取值——它是标准数据模型理念在码值层面的实例。

## 模型

一条 `integration.CodeMap`（表 `itg_code_map`）把一个系统绑定到一个码值集：

| 字段 | 含义 |
| --- | --- |
| `systemId` | 所属系统 |
| `codeSet` | 被翻译码值集的标识（如 `"gender"`）。当宿主目录暴露同一集合（mold translate 标签、`mold.CodeSetInspector`）时，两边标识应一致 |
| `name` | 显示名 |
| `entries` | 映射对（见下） |
| `onUnmapped` | 无条目匹配时的默认行为：`reject`（默认，fail closed）、`passthrough` 或 `fallback` |
| `fallbackCanonical` / `fallbackExternal` | `fallback` 策略下两个方向各自返回的兜底值 |
| `isEnabled` | 禁用的映射视同不存在 |

每个 `CodeMapEntry` 是一个双向映射对：

```json
{
  "canonical": "F",
  "external": 2,
  "canonicalAliases": ["female"],
  "externalAliases": ["02"]
}
```

- 查找匹配主值或任一别名；翻译永远输出对侧的**主值**——别名只被匹配、
  从不输出。
- 值端到端保留 JSON 类型（字符串、数字、布尔）；查找按归一化字符串比较，
  因此 `1` 与 `"1"` 指向同一条目。
- 保存时校验拒绝任一侧（主值与别名合并后）的重复查找值
  （`ErrInvalidCodeMap`），保证每次查找是确定的。

## `codes` 脚本库

适配器脚本（两个方向）通过全局 `codes` 对象使用码值映射，作用域即执行中的
系统：

```js
codes.toExternal('gender', input.gender)          // 标准 → 外部
codes.toCanonical('gender', body.sex)             // 外部 → 标准
codes.toExternal('gender', v, { fallback: 'U' })  // 按调用覆盖未映射策略
codes.toCanonical('status', v, { passthrough: true })
codes.toCanonical('status', v, { reject: true })
codes.entries('gender')                           // 原始映射对
```

- `null` 与 `undefined` 原样透传——翻译"缺失"不算查找。
- 按调用的选项对象只能携带 `fallback`、`passthrough: true`、`reject: true`
  之一，覆盖该次调用的存储策略。选项对象格式错误会在每次调用时报错，
  而不是等到第一个未映射值。
- 对目标系统没有已启用映射的码值集执行查找抛出 `ErrMissingCodeMap`
  （分类为 `config`）；reject 策略下的未映射值抛出 `ErrUnmappedValue`。
- 一次运行看到一个快照（映射按执行记忆化）；保存对下一次调用即时生效。
  编译后的查找索引按内容哈希跨执行共享。

## 宿主码值目录

映射编辑器如果知道宿主定义了哪些码值集、每个集合有哪些取值，体验会好得多。
当宿主的 [mold 码值集注册](../data-tools/mold)同时实现
`mold.CodeSetInspector` 时，`integration/code_set` 资源就会暴露该目录：

```go
type CodeSetInspector interface {
    ListCodeSets(ctx context.Context) ([]mold.CodeSetInfo, error)
    ListCodes(ctx context.Context, codeSet string) ([]mold.CodeInfo, error)
}
```

- Inspector 优先从注册的 `mold.CodeSetLoader` 断言（常规路径），其次从
  `mold.CodeSetResolver`（整体替换 resolver 的宿主）。
- 存在 inspector 时，保存码值映射还会校验其 `codeSet` 标识确实在宿主目录中
  注册。
- 不存在时，目录操作返回 `supported: false`，编辑器退化为自由文本输入。

`integration/code_map` 与 `integration/code_set` 的逐字段请求/响应文档见
[RPC 资源](./resources)。

## 下一步

[RPC 资源](./resources) —— 完整的管理 API 参考。
