---
sidebar_position: 2
---

# 验证

VEF 以 `go-playground/validator` 为基础验证引擎，并在其上叠加了框架自己的行为：

- 翻译后的错误消息
- `label` / `label_i18n` 支持
- 自定义验证规则
- 基于指针和标准 `omitempty` 流程的可空字段支持

## API 参考

| API | Contract |
| --- | --- |
| `validator.Validate(value)` | 使用包级 validator 校验 `value`；成功返回 `nil`，校验失败时把第一个翻译后的错误包装成 bad-request `result.Error` |
| `validator.RegisterValidationRules(rules...)` | 将每个 `ValidationRule` 注册到共享 validator 和两个内置 translators；返回遇到的第一个注册错误 |
| `validator.RegisterTypeFunc(fn, types...)` | 通过 `RegisterCustomTypeFunc` 注册自定义类型提取函数 |
| `validator.CustomTypeFunc` | 与 `func(field reflect.Value) any` 兼容的自定义类型提取函数类型 |
| `validator.ValidationRule` | 用于定义自定义 validation tag 和翻译行为的 struct |
| `ValidationRule.RuleTag` | 传给 `RegisterValidation` 和 `RegisterTranslation` 的 validator tag 名 |
| `ValidationRule.ErrMessageTemplate` | 注册到 go-playground translators 的 fallback message template |
| `ValidationRule.ErrMessageI18nKey` | 可选的框架 i18n key；如果解析成真实消息，会优先于 fallback template |
| `ValidationRule.Validate` | 作用于 go-playground `FieldLevel` values 的实际规则 predicate |
| `ValidationRule.ParseParam` | 从 go-playground `FieldError` values 提取 placeholder values，用于渲染错误消息 |
| `ValidationRule.CallValidationEvenIfNull` | 透传给 go-playground `RegisterValidation` |

## 验证入口

typed params 和 typed meta 在解码完成后会自动验证。框架内部调用的是：

```go
validator.Validate(value)
```

这意味着只要 handler 参数使用 typed 请求结构体，就默认带有验证能力，不需要在 handler 里手工再调一次。

## 标准 Validator 标签

VEF 直接继承 `go-playground/validator` 的标准标签集合。常见标签例如：

| 标签 | 含义 | 示例 |
| --- | --- | --- |
| `required` | 字段必须存在且非零值 | `validate:"required"` |
| `email` | 字段必须是合法邮箱 | `validate:"required,email"` |
| `min` | 数值或长度下界 | `validate:"min=1"` |
| `max` | 数值或长度上界 | `validate:"max=32"` |
| `oneof` | 字段必须属于候选值集合 | `validate:"oneof=admin user guest"` |
| `len` | 精确长度 | `validate:"len=32"` |
| `omitempty` | 字段为空时跳过后续校验 | `validate:"omitempty,email"` |
| `dive` | 对 slice 或 map 的每个元素继续校验 | `validate:"required,dive"` |

VEF 不会重写这些上游规则，它做的是在其上增加框架自定义规则。

## 标签名解析

验证错误在可能的情况下会优先使用标签名，而不是原始 Go 字段名。

解析顺序如下：

| 来源 | 效果 |
| --- | --- |
| `label:"..."` | 直接使用显式标签文本 |
| `label_i18n:"..."` | 通过 `i18n.T` 解析；如果翻译缺失，错误消息里可能显示这个 i18n key |
| 都没有 | 使用 Go 字段名 |

示例：

```go
type UserParams struct {
	api.P

	Username string `json:"username" validate:"required" label:"Username"`
	Phone    string `json:"phone" validate:"phone_number" label_i18n:"user_phone"`
}
```

## 内置自定义规则

VEF 当前内置了以下自定义验证规则：

| 规则标签 | 期望字段类型 | 含义 | 示例 |
| --- | --- | --- | --- |
| `phone_number` | `string` | 校验中国大陆手机号 | `validate:"phone_number"` |
| `dec_min=<value>` | `decimal.Decimal` | 小数值必须大于等于给定阈值 | `validate:"dec_min=0"` |
| `dec_max=<value>` | `decimal.Decimal` | 小数值必须小于等于给定阈值 | `validate:"dec_max=999.99"` |
| `alphanum_us` | `string` | 只允许字母、数字、下划线 | `validate:"alphanum_us"` |
| `alphanum_us_slash` | `string` | 只允许字母、数字、下划线和斜线 | `validate:"alphanum_us_slash"` |
| `alphanum_us_dot` | `string` | 只允许字母、数字、下划线和点 | `validate:"alphanum_us_dot"` |

这些规则会在包初始化时注册。如果注册失败，启动会 panic，不会留下部分配置好的 validator。

### `phone_number`

| 规则 | 说明 |
| --- | --- |
| 接受的值 | `^1[3-9]\d{9}$` |
| 主要用途 | 手机号输入校验 |
| 常见报错语义 | “格式不正确” |

### `dec_min` / `dec_max`

| 规则 | 说明 |
| --- | --- |
| 支持的字段类型 | `decimal.Decimal` |
| 参数格式 | 十进制字符串，例如 `10.5` |
| 行为 | 把字段值与解析后的 decimal 阈值比较 |
| 失败模式 | 字段不是 `decimal.Decimal` 或阈值参数不是合法 decimal 时校验失败 |

### `alphanum_us`

| 规则 | 允许字符 |
| --- | --- |
| `alphanum_us` | 字母、数字、`_` |
| `alphanum_us_slash` | 字母、数字、`_`、`/` |
| `alphanum_us_dot` | 字母、数字、`_`、`.` |

这三个 alphanum 规则都要求至少 1 个字符；空字符串不匹配对应 regex。

常见用途：

| 规则 | 常见用途 |
| --- | --- |
| `alphanum_us` | action name、简单标识、代码值 |
| `alphanum_us_slash` | RPC resource name、斜杠分段标识 |
| `alphanum_us_dot` | 文件名、模块名、点分标识 |

## 可空 / 可选字段

VEF 使用指针类型表达可空字段（老的 `null.*` 包装包在 v0.21 已经被移除，统一改用指针）。

```go
type UserParams struct {
    // 必填：校验直接作用在值上
    Username string `json:"username" validate:"required,min=3"`

    // 可选：指针为 nil 时跳过后续校验
    Phone *string `json:"phone" validate:"omitempty,phone_number"`

    // 可选 + "可为空" 语义
    Bio *string `json:"bio" validate:"omitempty,max=500"`
}
```

标准的 `omitempty` validator tag 已经能处理可空场景——指针为 `nil` 时跳过后续规则；指向具体值时则按链上规则继续校验。如果某个自定义类型需要参与校验，用 `validator.RegisterTypeFunc(...)` 搭配 `validator.CustomTypeFunc` 告诉校验器如何从该类型里取出可比较的值。

## 错误行为

验证失败时，框架会返回第一个翻译后的验证错误，并包装为框架 `result.Error`。

| 情况 | 结果 |
| --- | --- |
| 验证成功 | handler 正常继续执行 |
| 一条或多条验证规则失败 | 框架返回 bad-request 风格错误 |
| 验证过程中出现非验证类错误 | 框架也会包装成 bad-request 风格错误 |

HTTP 行为：

| 属性 | 值 |
| --- | --- |
| 业务码 | `result.ErrCodeBadRequest` |
| HTTP 状态 | `400 Bad Request` |

内置 go-playground 规则和 VEF 自定义规则都会在校验发生时读取当前
`i18n.CurrentLanguage()`。中文（`zh-CN`）使用中文 translator；其他支持语言使用英文 translator。

## 注册额外自定义规则

应用可以通过 `validator.RegisterValidationRules(...)` 注册自己的规则。

一个规则由 `validator.ValidationRule` 定义，主要字段包括：

| 字段 | 作用 |
| --- | --- |
| `RuleTag` | validator 标签名 |
| `ErrMessageTemplate` | 回退用的错误消息模板 |
| `ErrMessageI18nKey` | i18n 消息 key |
| `Validate` | 真正的校验函数 |
| `ParseParam` | 为错误消息占位符提取参数 |
| `CallValidationEvenIfNull` | 是否在 null 值上也执行校验 |

自定义规则消息会先检查 `ErrMessageI18nKey`。当 `i18n.T(key)` 返回不同于
key 的真实消息时，会用 `ParseParam` 的值替换 `{0}`、`{1}` 等占位符。
如果没有找到框架 i18n 消息，则由 go-playground translator 渲染
`ErrMessageTemplate`。

## 常见模式

### 简单必填字段

```go
type UserParams struct {
	api.P

	Username string `json:"username" validate:"required,alphanum,max=32" label:"Username"`
	Email    string `json:"email" validate:"omitempty,email,max=128" label:"Email"`
}
```

### decimal 范围校验

```go
type PriceParams struct {
	api.P

	Amount decimal.Decimal `json:"amount" validate:"dec_min=0,dec_max=999999.99" label:"Amount"`
}
```

### 可选指针字段校验

```go
type UserParams struct {
	api.P

	Phone *string `json:"phone" validate:"omitempty,phone_number" label:"Phone"`
}
```

## 自定义验证规则

前面的 struct-tag 校验（内置标签，加上框架自己的 `phone_number` /
`dec_min` / `dec_max` / `alphanum_us*` 规则）和自定义规则注册其实是同一个包的
两个视角，而不是两套独立系统：`validate:"..."` 标签用来选择在某个字段上运行
哪条规则，而 `validator.RegisterValidationRules` 就是让一个新标签变得"可选"
的方式。`validator.Validate` 会通过同一次调用同时运行内置规则和已注册的自定义
规则——并不存在一条单独的"自定义校验"代码路径。

| API | 作用 |
| --- | --- |
| `validator.Validate(value)` | 验证值并返回第一个框架验证错误 |
| `validator.RegisterValidationRules(rules...)` | 添加自定义 `ValidationRule` |
| `validator.RegisterTypeFunc(fn, types...)` | 为应用自定义 wrapper 注册类型提取函数 |
| `validator.CustomTypeFunc` | `RegisterTypeFunc` 接收的回调类型 |
| `validator.ValidationRule` | 自定义规则定义，包含 tag、消息、validator 回调、参数解析和 null-call 标记 |

完整的 `ValidationRule` 字段参考、消息解析顺序和实例，见上文的
[注册额外自定义规则](#注册额外自定义规则)。

## 实践建议

- 验证规则放在 typed params / meta 结构体上，不要写进 handler
- 尽量补 `label` 或 `label_i18n`，让错误消息保持用户可读
- 只要框架已有现成自定义规则，就优先复用
- 自定义规则应保持窄而明确，不要做过宽的“万能校验”
- 可空输入优先使用指针字段加 `omitempty`；只有应用自己的包装类型才需要额外注册类型提取函数

## 下一步

继续阅读 [参数与元信息](../building-apis/params-and-meta)，看验证是在请求解码链路中的哪个阶段触发的。
