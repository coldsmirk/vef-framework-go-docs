---
sidebar_position: 8
---

# 验证

VEF 以 `go-playground/validator` 为基础验证引擎，并在其上叠加了框架自己的行为：

- 翻译后的错误消息
- `label` / `label_i18n` 支持
- 自定义验证规则
- 对框架 null 包装类型的支持

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
| `label_i18n:"..."` | 先经过 i18n 翻译，再回退到字段名 |
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

### `phone_number`

| 规则 | 说明 |
| --- | --- |
| 接受的值 | `1[3-9]\\d{9}` |
| 主要用途 | 手机号输入校验 |
| 常见报错语义 | “格式不正确” |

### `dec_min` / `dec_max`

| 规则 | 说明 |
| --- | --- |
| 支持的字段类型 | `decimal.Decimal` |
| 参数格式 | 十进制字符串，例如 `10.5` |
| 行为 | 把字段值与解析后的 decimal 阈值比较 |

### `alphanum_us`

| 规则 | 允许字符 |
| --- | --- |
| `alphanum_us` | 字母、数字、`_` |
| `alphanum_us_slash` | 字母、数字、`_`、`/` |
| `alphanum_us_dot` | 字母、数字、`_`、`.` |

常见用途：

| 规则 | 常见用途 |
| --- | --- |
| `alphanum_us` | action name、简单标识、代码值 |
| `alphanum_us_slash` | RPC resource name、斜杠分段标识 |
| `alphanum_us_dot` | 文件名、模块名、点分标识 |

## 支持的 Null 类型

VEF 已注册自定义 type func，因此以下 null 包装类型都能正确参与验证：

| 支持的 null 类型 |
| --- |
| `null.String` |
| `null.Int` |
| `null.Int16` |
| `null.Int32` |
| `null.Float` |
| `null.Bool` |
| `null.Byte` |
| `null.DateTime` |
| `null.Date` |
| `null.Time` |
| `null.Decimal` |

实际效果：

- 当 null 包装值有效时，校验作用于内部真实值
- 当 null 包装值无效时，该值会被视为 `nil`

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

### null 包装类型校验

```go
type UserParams struct {
	api.P

	Phone null.String `json:"phone" validate:"omitempty,phone_number" label:"Phone"`
}
```

## 实践建议

- 验证规则放在 typed params / meta 结构体上，不要写进 handler
- 尽量补 `label` 或 `label_i18n`，让错误消息保持用户可读
- 只要框架已有现成自定义规则，就优先复用
- 自定义规则应保持窄而明确，不要做过宽的“万能校验”
- 如果字段类型是框架 null wrapper，直接依赖内置 null 支持，不要手工拆值

## 下一步

继续阅读 [参数与元信息](./params-and-meta)，看验证是在请求解码链路中的哪个阶段触发的。
