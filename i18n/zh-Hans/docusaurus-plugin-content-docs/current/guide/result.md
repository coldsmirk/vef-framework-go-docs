---
sidebar_position: 9
---

# Result

`result` 包提供框架中使用的标准 API 响应信封。

## 响应结构

所有 API 响应都遵循以下格式：

```json
{
  "code": 0,
  "message": "成功",
  "data": {}
}
```

```go
type Result struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data"`
}
```

## 创建成功响应

```go
import "github.com/coldsmirk/vef-framework-go/result"

// 简单成功（无数据）
return result.Ok().Response(ctx)

// 带数据的成功
return result.Ok(user).Response(ctx)

// 自定义消息的成功
return result.Ok(result.WithMessage("创建成功")).Response(ctx)

// 带数据和自定义消息
return result.Ok(user, result.WithMessage("用户已创建")).Response(ctx)

// 自定义 HTTP 状态码
return result.Ok(user).Response(ctx, 201)
```

### Ok 选项

| 选项 | 效果 |
| --- | --- |
| `result.WithMessage(msg)` | 覆盖默认成功消息 |
| `result.WithCode(code)` | 覆盖默认成功码（0）|

## 创建错误响应

```go
// 简单业务错误
return result.Err("出错了").Response(ctx)

// 带自定义错误码
return result.Err("未找到", result.WithErrCode(result.ErrCodeRecordNotFound)).Response(ctx)

// 自定义 HTTP 状态码
return result.Err("未授权").Response(ctx, 401)
```

### Err 选项

| 选项 | 效果 |
| --- | --- |
| `result.WithErrCode(code)` | 设置特定错误码 |
| `result.WithErrData(data)` | 在错误响应中附加数据 |

## 检查结果

```go
r := result.Ok(data)
r.IsOk() // true（code == 0）

r = result.Err("fail")
r.IsOk() // false
```

## 错误码

### 认证错误（1000–1099）

| 错误码 | 常量 | 含义 |
| --- | --- | --- |
| 1000 | `ErrCodeUnauthenticated` | 未认证 |
| 1001 | `ErrCodeUnsupportedAuthenticationType` | 不支持的认证类型 |
| 1002 | `ErrCodeTokenExpired` | Token 过期 |
| 1003 | `ErrCodeTokenInvalid` | Token 无效 |
| 1004 | `ErrCodeTokenNotValidYet` | Token 尚未生效 |
| 1010 | `ErrCodePrincipalInvalid` | 主体无效 |
| 1011 | `ErrCodeCredentialsInvalid` | 凭证无效 |
| 1020 | `ErrCodeSignatureInvalid` | 签名无效 |

### 质询错误（1030–1039）

| 错误码 | 常量 | 含义 |
| --- | --- | --- |
| 1030 | `ErrCodeChallengeRequired` | 需要质询 |
| 1035 | `ErrCodeOTPCodeRequired` | 需要 OTP 验证码 |
| 1036 | `ErrCodeOTPCodeInvalid` | OTP 验证码无效 |
| 1037 | `ErrCodeNewPasswordRequired` | 需要新密码 |

### 授权错误（1100–1199）

| 错误码 | 常量 | 含义 |
| --- | --- | --- |
| 1100 | `ErrCodeAccessDenied` | 拒绝访问 |

### 请求错误（1200–1499）

| 错误码 | 常量 | 含义 |
| --- | --- | --- |
| 1200 | `ErrCodeNotFound` | 资源未找到 |
| 1300 | `ErrCodeUnsupportedMediaType` | 不支持的媒体类型 |
| 1400 | `ErrCodeBadRequest` | 错误请求 |
| 1401 | `ErrCodeTooManyRequests` | 请求频率限制 |
| 1402 | `ErrCodeRequestTimeout` | 请求超时 |

### 业务错误（2000+）

| 错误码 | 常量 | 含义 |
| --- | --- | --- |
| 2000 | `ErrCodeDefault` | 通用业务错误 |
| 2001 | `ErrCodeRecordNotFound` | 数据库记录未找到 |
| 2002 | `ErrCodeRecordAlreadyExists` | 记录已存在 |
| 2003 | `ErrCodeForeignKeyViolation` | 外键约束冲突 |

## I18n 集成

错误和成功消息通过 `i18n` 模块自动本地化。消息键（如 `"record_not_found"`、`"success"`）会在运行时从配置的语言包中查找。

## 预构建错误构造函数

框架提供常用的错误构造函数：

```go
result.ErrRecordNotFound()       // code: 2001
result.ErrRecordAlreadyExists()  // code: 2002
result.ErrForeignKeyViolation()  // code: 2003
result.ErrAccessDenied()         // code: 1100
result.ErrUnauthenticated()      // code: 1000
result.ErrBadRequest(msg)        // code: 1400
result.ErrNotFound()             // code: 1200
result.ErrTooManyRequests()      // code: 1401
```
