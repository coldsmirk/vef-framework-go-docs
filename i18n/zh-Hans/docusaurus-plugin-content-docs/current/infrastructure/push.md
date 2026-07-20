---
sidebar_position: 6
---

# 服务端推送

`push` 模块（v0.39）是基于 WebSocket 的服务端推送通道：业务代码向用户、
角色或全体发送类型化消息，已连接的客户端以 JSON 文本帧接收。

投递在契约上是**尽力而为**的：离线、已断开或来不及消费队列的接收方会错过
消息。可靠通知应落在业务存储里（供客户端拉取的通知表），推送只作为实时
提示。

## 从业务代码发送

`push.Notifier` 可从 DI 注入。端点关闭时它依然可用——投递被静默丢弃
（不存在任何连接），业务代码无需对特性开关做分支。

```go
type OrderService struct {
    notifier push.Notifier
}

func (s *OrderService) NotifyApprovers(ctx context.Context, order *Order, approverIDs []string) error {
    return s.notifier.Push(ctx,
        push.NewMessage("order.pending", map[string]any{"orderId": order.ID}),
        push.ToUsers(approverIDs...),
    )
}
```

| API | 契约 |
| --- | --- |
| `Push(ctx, message, targets...)` | 向所有目标选中的接收方投递（取并集，每个连接至多一次）。消息 ID 或时间为零值时自动填充。仅当消息或目标集非法时返回错误，从不因接收方错过而报错 |
| `push.NewMessage(type, payload)` | 构建带生成 ID 与当前时间的消息 |
| `push.ToUsers(userIDs...)` | 指向给定用户 ID 的在线连接 |
| `push.ToRoles(roles...)` | 指向持有任一给定角色的连接（角色在握手时快照） |
| `push.Broadcast()` | 指向所有在线连接 |

错误：`push.ErrNoTarget`（目标集为空或 users/roles 选择器为空——投递给
空集合永远是调用方缺陷）、`push.ErrTypeRequired`（客户端按 type 分发）、
`push.ErrUnknownTargetKind`（手工构造的目标超出词汇表）。

### 消息信封

每次推送是一个 JSON 文本帧：

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `id` | `string` | 唯一消息 ID（零值时生成） |
| `type` | `string` | 业务定义的判别符，客户端按其分发 |
| `payload` | 任意 JSON 值 | 任意可 JSON 序列化的负载；为空省略 |
| `time` | 时间戳 | 发送时间（零值时填充） |

## WebSocket 端点

端点为 opt-in：

```toml
[vef.push]
enabled = true
path = "/ws"                      # 默认
allowed_origins = []              # 空表示允许所有来源（握手本身已由令牌认证）
ping_interval = "30s"             # 连续错过两次 pong 即断开
write_timeout = "10s"             # 单个出站帧上限
send_buffer = 32                  # 每连接出站队列；消费过慢的客户端被断开
max_connections_per_user = 0     # 每节点每用户上限；0 为不限
session_recheck_interval = "60s"  # 不透明令牌会话复核节律
```

### 客户端握手

握手使用与 API 相同的令牌机制（`vef.security.token_type`）认证，且发生在
任何 socket 建立之前。浏览器 WebSocket API 无法设置 `Authorization` 头，
因此令牌可从两个通道之一提供：

- `Authorization: Bearer <token>` 头（非浏览器客户端），或
- 标准访问令牌查询参数 `__accessToken`。

```js
const ws = new WebSocket(`wss://host/ws?__accessToken=${accessToken}`);
ws.onmessage = (e) => {
    const msg = JSON.parse(e.data); // { id, type, payload, time }
    dispatch(msg.type, msg.payload);
};
```

### 会话集成（不透明令牌）

在 `vef.security.token_type = "opaque_token"` 下，连接与其登录会话绑定：

- 握手时检查会话，注册完成后立即复核一次，堵住握手窗口内的吊销竞态；复核
  通过前连接被隔离在投递之外；
- 会话吊销（登出、并发登录驱逐、管理员踢出）通过
  `security.SessionRevocationListener` 接缝立即关闭连接；
- 周期清扫（`session_recheck_interval`）对所有连接复核会话，捕获过期。

终止性关闭码是客户端协议契约的一部分——收到后不得自动重连（相反，传输层
故障应带退避重连）：

| 关闭码 | 常量 | 含义 |
| --- | --- | --- |
| `4401` | `push.CloseSessionInvalid` | 登录会话被吊销或过期；进入登出流程 |
| `4429` | `push.CloseTooManyConnections` | 达到每用户连接上限；在其他连接关闭前不要重试 |

无状态 JWT 令牌下没有可吊销的会话；连接存活到自行关闭或心跳失败。

## 多节点中继

单节点时 Hub 在节点内投递。多节点部署下，启用 Redis
（`vef.redis.enabled = true`）即自动经 Redis pub/sub 在节点间中继推送与
吊销踢出：

- 中继频道按 Redis 库号和应用名命名空间化
  （`vef:push:relay:<db>:<app-name>`），共享同一 Redis 的无关部署绝不会
  互相投递；
- 启用中继而未设置 `vef.app.name` 会拒绝启动——否则命名空间会冲突；
- 发布经有界工作协程执行，与吊销调用路径解耦。

除 Redis 本身外无需额外配置。

## 设计说明

- 每个 socket 的写由单一 goroutine 拥有；每连接队列有界（`send_buffer`），
  消费过慢的客户端被断开，而不是让其对 Hub 施加背压。
- 角色定向匹配的是连接握手时快照的角色；角色变更在下一次连接生效。
- 没有客户端到服务端的消息契约：通道刻意单向（客户端只发送 pong）。

## 下一步

将推送与持久化状态变更配合使用：经[事件总线](./event-bus)发布领域事件、
落库通知记录，再推送实时提示。
