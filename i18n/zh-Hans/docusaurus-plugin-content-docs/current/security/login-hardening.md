---
sidebar_position: 6
---

# 登录加固（Login Hardening）

`security` 包为登录端点提供了五种相互独立、按需启用的加固能力。它们全部在
框架的安全模块内完成装配，都不会改变客户端调用 `security/auth.login` 的
方式；你可以只启用其中任意一部分——一个未做额外配置的新应用，默认不带任何一
项保护。

| 加固层 | 防护对象 | 配置节 | 应用需实现的接口 |
| --- | --- | --- | --- |
| 密码传输加密 | 客户端与服务端之间的凭据嗅探 | 无（仅 DI 装配） | `security.PasswordDecryptor` |
| 暴力破解锁定 | 撞库 / 密码猜测 | `vef.security.lockout` | 无（可将 `security.LoginGuard` 换成 Redis 版本） |
| 密码强度 | 弱密码 | `vef.security.password_policy` | 无 |
| 密码历史 | 密码重用 | `vef.security.password_policy.history_depth` | `security.PasswordHistoryStore` |
| 密码过期 | 密码陈旧未更换 | `vef.security.password_policy.max_age` | `security.PasswordMetadataLoader` |

本页按照由浅入深的顺序逐一介绍：先讲装配方式，再讲配置项，最后讲应用需要实
现的扩展接口。

## 密码传输加密

如果客户端在发送密码前先对其加密（一种常见做法："浏览器端加密、服务端再
哈希"，用于防御网络层的凭据嗅探），可以注册一个 `security.PasswordDecryptor`：

```go
type PasswordDecryptor interface {
	Decrypt(encryptedPassword string) (string, error)
}
```

`PasswordAuthenticator` 会在校验之前，先把收到的凭据解密为明文，再与存储的
哈希值比对。这样 `password.Encoder` 就能始终保持为一个纯粹的 KDF，存储和比
对都用同一套逻辑——解密是认证器这一层的职责，不属于 encoder。encoder 本身
请见 [Password](./password)；用来实现解密器的 `Cipher`/`CipherSigner` 请见
[Cryptox](./cryptox)。

`cryptox.NewRSA` 已经天然满足 `PasswordDecryptor`——它的
`Decrypt(ciphertext string) (string, error)` 方法与该接口逐字匹配，因此只需
要通过 DI 提供它即可：

```go
fx.Provide(func() (security.PasswordDecryptor, error) {
	return cryptox.NewRSA(privateKey, publicKey)
})
```

`PasswordDecryptor` 是 `NewPasswordAuthenticator` 的一个可选依赖（与
`UserLoader` 等一样）——不注册它，认证器就会把收到的凭据当作明文处理，这也
是零配置下的默认行为。

格式错乱的密文会被当作一次普通的密码错误来处理：认证器会执行一次哑
KDF 比对，让解密失败和真实的密码不匹配耗费相同的时间，从而堵住一个原本能
够区分"密文错误"与"密码错误"或"用户不存在"的时序侧信道。

## 暴力破解锁定

`security.LoginGuard` 在认证真正发生**之前**，根据某个身份已累积的失败次
数对登录端点进行限流：

```go
type LoginGuard interface {
	Check(ctx context.Context, attempt LoginAttempt) (LoginDecision, error)
	RecordFailure(ctx context.Context, attempt LoginAttempt) (LoginDecision, error)
	RecordSuccess(ctx context.Context, attempt LoginAttempt) error
}

type LoginAttempt struct {
	Identity string // 客户端提交的登录标识
	ClientIP string // 解析出的来源地址
}

type LoginDecision struct {
	Allowed    bool
	RetryAfter time.Duration // Allowed 为 true 时为零值
}
```

`AuthResource.Login` 会在认证前调用 `Check`，认证失败时调用
`RecordFailure`，凭据一旦通过验证（在任何第二因素挑战之前，因为撞库尝试所
用的凭据此时已经验证成功）就调用 `RecordSuccess`。失败次数按
`LockoutPolicy.Key` 维度累积，认证成功后清零。

从 v0.38 起，同一个 guard 也覆盖 `resolve_challenge`：第二因素猜测失败会
计入同一个锁定 key，锁定触发后两个端点同时被拦——攻击者即使走到挑战环节，
也无法在锁定预算之外暴力猜测。

### 启用与配置

锁定功能**默认开启**（`max_failures = 10`）。在 `vef.security.lockout` 下配
置：

```toml
[vef.security.lockout]
enabled = true          # 默认值：true
max_failures = 10       # 默认值：10
window = "15m"           # 默认值：15m —— 连续这么久没有新失败，计数器即清零
lock_duration = "15m"    # 默认值：15m —— "lock" 策略下的封锁时长
strategy = "lock"        # "lock" 或 "backoff"，默认值："lock"
backoff_base = "1s"      # 默认值：1s —— "backoff" 策略下第一次的延迟
backoff_max = "15m"      # 默认值：15m —— backoff 延迟的上限
key = "user_ip"          # "user"、"ip" 或 "user_ip"，默认值："user_ip"
```

将 `enabled` 设为 `false` 可彻底关闭锁定功能。其余字段留空或为零值时都会解
析为各自的默认值——如果自己组装策略，应通过 `config.LockoutConfig` 的
`Effective*` 访问器读取，而不是直接读原始字段。

- **`strategy = "lock"`**（对应 Go 常量 `security.LockoutStrategyLock`）：失
  败次数**达到** `max_failures` 即在 `lock_duration` 时长内封锁所有尝试——
  攻击者恰好获得 `max_failures` 次机会，触及阈值的那次失败即触发锁定。
- **`strategy = "backoff"`**（对应 Go 常量 `security.LockoutStrategyBackoff`）：
  改为施加逐步升高的延迟——触及阈值的那次失败开始等待 `backoff_base`，此后
  每多失败一次延迟翻倍，直到 `backoff_max` 封顶。合法用户会被拖慢，但永远
  不会被彻底锁死；攻击者也无法借此把受害者无限期地锁在门外。
- **`key`** 选择失败次数按哪个身份维度计数：`"user"`
  （`security.LockoutKeyUser`——按登录标识计数，跨所有来源 IP）、`"ip"`
  （`security.LockoutKeyIP`——按来源地址计数，跨所有登录标识）、或
  `"user_ip"`（`security.LockoutKeyUserIP`——默认值，按"标识 + 来源"的组合
  计数，既能限制凭据猜测，又不会让攻击者仅凭猜出一个账号就把该账号从所有
  IP 上锁死）。

### 存储后端

默认的守卫实现是 `security.MemoryLoginGuard`，由
`security.NewMemoryLoginGuard(policy)` 构造——适用于单实例部署。多节点部
署应通过 `fx.Decorate` 把它换成 `security.NewRedisLoginGuard`，该函数返回
一个由 Redis 共享计数器支撑的 `security.RedisLoginGuard`，使失败计数器在
各节点间共享：

```go
vef.Run(
	// ...
	fx.Decorate(func(client *redis.Client, cfg *config.SecurityConfig) security.LoginGuard {
		l := cfg.Lockout
		return security.NewRedisLoginGuard(client, security.LockoutPolicy{
			MaxFailures:  l.EffectiveMaxFailures(),
			Window:       l.EffectiveWindow(),
			LockDuration: l.EffectiveLockDuration(),
			Strategy:     security.LockoutStrategy(l.EffectiveStrategy()),
			BackoffBase:  l.EffectiveBackoffBase(),
			BackoffMax:   l.EffectiveBackoffMax(),
			Key:          security.LockoutKey(l.EffectiveKey()),
		})
	}),
)
```

这需要 `vef.redis.enabled = true`，以便 DI 容器中能拿到 `*redis.Client`
（[会话管理](./session-management) 页中，opaque token 的会话存储也是同样
的换法）。

### 故障处理与返回的错误

`LoginGuard` 后端出错（例如 Redis 不可达）时会**失败开放（fail open）**：
守卫记录一条警告日志，放行本次尝试，而不是因为计数器存储不可用就拒绝所有
登录。守卫是纵深防御的一环，不是认证结果的权威来源。

一旦触发锁定，会返回 `security.ErrAccountLocked(retryAfter)`——HTTP
429，业务码为 `security.ErrCodeAccountLocked`（`1023`），响应消息中的重试
等待时长会向上取整到整分钟（最少为一分钟）。该消息由 i18n key
`security.ErrMessageAccountLocked`（`"security_account_locked"`）渲染而
成。

## 密码强度

`security.PasswordValidator` 校验候选明文密码是否满足策略：

```go
type PasswordValidator interface {
	Validate(ctx context.Context, principal *Principal, plaintext string) error
}
```

它由多条可组合的 `PasswordRule` 通过 `NewRuleBasedValidator` 组装而成；不注
册任何规则时会接受所有密码（零配置下的默认行为）：

```go
type PasswordRule interface {
	Check(principal *Principal, plaintext string) error
}
```

内置规则：

| 构造函数 | 规则 |
| --- | --- |
| `NewMinLengthRule(minLength)` | 至少 `minLength` 个 rune |
| `NewMaxLengthRule(maxLength)` | 至多 `maxLength` 个 rune（同时防御慢 KDF 拒绝服务和 bcrypt 静默截断） |
| `NewCharacterClassRule(requireUpper, requireLower, requireDigit, requireSymbol, minClasses)` | 要求的字符类别，及/或要求同时出现的不同类别的最少数量。符号类是任何非字母、非数字、非空白的 rune；无大小写的字母（如中文）不计入任何类别（v0.38 修正） |
| `NewDisallowIdentityRule()` | 拒绝包含 principal 的 `ID` 或 `Name` 的密码（大小写不敏感；短于 3 个 rune 的片段会被忽略——按 rune 计数，两个汉字的名字不会拒绝掉大部分密码） |
| `NewBlocklistRule(entries)` | 拒绝匹配黑名单条目的密码（大小写不敏感，比对前会去除首尾空白） |

### 配置

框架会根据 `vef.security.password_policy` 自动构建一个 `PasswordValidator`，
并注入到框架自身需要它的地方（下文的强制改密挑战）。每个字段都是按需启用
——留空或为零值即关闭对应规则：

```toml
[vef.security.password_policy]
min_length = 12
max_length = 128
require_upper = true
require_lower = true
require_digit = true
require_symbol = false
min_char_classes = 3
disallow_username = true
blocklist = ["password", "123456", "qwerty"]
```

也可以把同一个 `security.PasswordValidator` 注入到你自己的注册或重置流程
中，复用同一份已配置好的策略，而不必重新声明一遍规则。

### 违规错误

每种强度违规都携带业务码 `security.ErrCodePasswordPolicyViolation`
（`1050`），HTTP 400；i18n 消息会说明具体触发了哪条规则：

| 错误 | 触发条件 |
| --- | --- |
| `ErrPasswordTooShort(minLength)` | 低于 `min_length` |
| `ErrPasswordTooLong(maxLength)` | 高于 `max_length` |
| `ErrPasswordMissingUppercase` / `ErrPasswordMissingLowercase` / `ErrPasswordMissingDigit` / `ErrPasswordMissingSymbol` | 缺少某个必须的字符类别 |
| `ErrPasswordTooFewCharClasses(minClasses)` | 出现的不同字符类别少于 `min_char_classes` |
| `ErrPasswordContainsIdentity` | 密码包含账号的 `ID` 或 `Name` |
| `ErrPasswordBlocked` | 密码匹配黑名单条目 |

上面三条带模板参数的消息各自对应一个具名 i18n key 常量，供只需要原始 key
而非构造好的 `result.Error` 的调用方使用：`security.ErrMessagePasswordTooShort`、
`security.ErrMessagePasswordTooLong`、`security.ErrMessagePasswordTooFewCharClasses`。

## 密码历史（防重用）

密码历史会拒绝重复该主体最近使用过的某个密码的新密码。框架只负责读取历史
记录来判断是否重用，哈希比对也由框架自己完成——存储历史记录的职责在应用
一侧，因为用户数据库归应用所有：

```go
type PasswordHistoryStore interface {
	// Recent 返回该主体最近的若干条已编码密码，按时间从新到旧排列，最多 limit 条。
	Recent(ctx context.Context, principalID string, limit int) ([]string, error)
	// Add 将 encodedPassword 记录为该主体最新的一条历史记录。
	Add(ctx context.Context, principalID, encodedPassword string) error
}
```

将你自己的实现注册为一个普通的 DI 值：

```go
fx.Provide(func(db orm.DB) security.PasswordHistoryStore {
	return myapp.NewPasswordHistoryStore(db)
})
```

将 `vef.security.password_policy.history_depth` 设为一个正数：

```toml
[vef.security.password_policy]
history_depth = 5
```

当同时满足"注册了 `PasswordHistoryStore`"和"`history_depth > 0`"这两个条
件时，框架会通过 `NewChainValidator` 把重用检查组合进注入的
`PasswordValidator`——强度规则先执行，然后才是重用检查。密码匹配最近
`history_depth` 条记录中的任意一条时，会以
`security.ErrPasswordReused` 失败（同样是 `ErrCodePasswordPolicyViolation` /
400）。

请在你自己实现的 `PasswordChanger.ChangePassword` 中、持久化新哈希之后立即
调用 `Add`——框架只负责读取和比对历史记录，写入由应用完成。如果你自己组装
校验器链而不依赖配置驱动的默认实现，也可以直接使用
`NewHistoryValidator`（`NewHistoryValidator(store, encoder, depth)`）。

## 密码过期

密码过期会在密码使用时长超过配置的最大年龄后强制要求改密。框架需要知道密
码最后一次设置的时间，但这个数据不归框架所有，因此由应用实现一个加载器：

```go
type PasswordMetadataLoader interface {
	// PasswordChangedAt 返回该 principal 密码最后一次设置的时间。零值表示"未知"，
	// 会被当作"尚未过期"处理，而不是在数据不完整时强行要求改密。
	PasswordChangedAt(ctx context.Context, principal *Principal) (time.Time, error)
}
```

将其包装为一个 `ExpiryPasswordChangeChecker`：

```go
checker := security.NewExpiryPasswordChangeChecker(myMetadataLoader, 90*24*time.Hour)
```

按约定，年龄上限声明在 `vef.security.password_policy.max_age`（零值表示
关闭过期检查）——但注意框架**不会**自行消费这个键。与 `history_depth`
（注册 `PasswordHistoryStore` 后自动组链）不同，`max_age` 是纯声明性字段：
由你的装配代码从 `config.SecurityConfig` 读出并传给
`NewExpiryPasswordChangeChecker`，如上例所示。

```toml
[vef.security.password_policy]
max_age = "2160h" # 90 天
```

`ExpiryPasswordChangeChecker` 实现了 `security.PasswordChangeChecker`，这个
接口同样用于其他强制改密场景（比如首次登录）。可以用
`NewCompositePasswordChangeChecker` 把多个原因组合在一起，它会返回第一个
命中的原因：

```go
checker := security.NewCompositePasswordChangeChecker(
	firstLoginChecker,
	security.NewExpiryPasswordChangeChecker(myMetadataLoader, 90*24*time.Hour),
)
```

把组合后的 checker、你的 `PasswordChanger`，以及（可选的）
`PasswordValidator` 一起传给 `NewPasswordChangeChallengeProvider`，再将其注
册为一个登录挑战提供者：

```go
vef.ProvideChallengeProvider(func(
	checker security.PasswordChangeChecker,
	changer security.PasswordChanger,
	validator security.PasswordValidator,
) security.ChallengeProvider {
	return security.NewPasswordChangeChallengeProvider(checker, changer, validator)
})
```

一旦 checker 命中，`security/auth.login` 就会返回一个 `password_change`
挑战（`Reason: "expired"`）而不是令牌；客户端通过 `resolve_challenge` 把新
密码作为响应提交上来完成解决。该 provider 会先用传入的 `PasswordValidator`
校验新密码——因此强度和历史规则在这里同样适用——校验通过后再调用
`PasswordChanger.ChangePassword` 持久化。完整的登录/挑战请求结构请见
[认证](./authentication)。

## 组合使用

一个启用全部加固层的部署，会把上面各节的配置合并到同一个 `[vef.security]`
表下：

```toml
[vef.security.lockout]
max_failures = 10
strategy = "lock"
key = "user_ip"

[vef.security.password_policy]
min_length = 12
require_upper = true
require_lower = true
require_digit = true
min_char_classes = 3
disallow_username = true
history_depth = 5
max_age = "2160h"
```

……再加上 Go 侧按需实现的 `PasswordDecryptor`、`PasswordHistoryStore`、
`PasswordMetadataLoader`——具体实现哪些取决于你的威胁模型；它们相互独立，
少实现某一个也不影响其余部分正常工作。
