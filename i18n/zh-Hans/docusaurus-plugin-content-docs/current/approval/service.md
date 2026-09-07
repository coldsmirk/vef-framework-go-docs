---
sidebar_position: 6
---

# 编程式流程控制

`approval.Service` 是审批引擎的编程式控制面：`approval/instance` 与 `approval/admin` 资源暴露的每一个运行时操作，都能在宿主代码里直接调用，完全不需要一次 HTTP 请求。只要启用了 `vef.ApprovalModule`，注入它即可。

```go
type OrderService struct {
    approval approval.Service
}
```

那些 API 资源本身就是这个接口的调用方，因此一个操作无论是由请求还是由宿主例程触发，行为都**完全一致**——同样的校验、同样顺序的领域事件、同样的审计行、生命周期钩子和业务回写。不存在需要同步维护的第二条代码路径。

不要绕过它去直接用 `cqrs.Bus` 加 `internal/approval/command` 里的类型。那些类型本来就无法从另一个模块引用；就算能，直接分发也会绕开负责发布事件、写操作日志的 collector behavior。

## 操作一览

| 方法 | 入参 | 说明 |
| --- | --- | --- |
| `StartInstance` | `approval.StartInstanceInput` | 在引擎遍历完开始节点后返回 `*Instance`——一个从开始节点一路走到结束节点都没有停下的流程，返回时就已经是完成状态。 |
| `WithdrawInstance` | `approval.WithdrawInstanceInput` | 操作人必须是申请人。 |
| `ResubmitInstance` | `approval.ResubmitInstanceInput` | 操作人必须是申请人；`FormData` 会整体替换实例的表单数据。 |
| `TerminateInstance` | `approval.TerminateInstanceInput` | 管理操作。运行中、已退回、已撤回的实例都可以终止——由实例状态机唯一裁定。 |
| `ApproveTask` | `approval.ApproveTaskInput` | 审批通过审批任务**或**办理完成办理任务；两者共用一条命令。 |
| `RejectTask` | `approval.RejectTaskInput` | 操作人必须持有该任务。 |
| `TransferTask` | `approval.TransferTaskInput` | 转办。节点必须允许转办。 |
| `RollbackTask` | `approval.RollbackTaskInput` | 退回。当回退类型是 `specified` 时，`TargetNodeID` 必须是已配置的回退目标之一。 |
| `ReassignTask` | `approval.ReassignTaskInput` | 管理操作。操作人不需要持有该任务。 |
| `AddAssignee` | `approval.AddAssigneeInput` | 加签。串行节点上会拒绝 `parallel`。 |
| `RemoveAssignee` | `approval.RemoveAssigneeInput` | 减签。`TaskID` 是被减掉那个人的任务，不是操作人自己的；节点上最后一个审批人不能被减掉。 |
| `AddCC` | `approval.AddCCInput` | 手动抄送。当前节点必须允许手动抄送。 |
| `MarkCCRead` | `approval.MarkCCReadInput` | 记录一次抄送已读回执。 |
| `UrgeTask` | `approval.UrgeTaskInput` | 催办，受节点上按 (任务, 催办人) 计的冷却时间约束。 |
| `RetryBusinessProjection` | `approval.RetryBusinessProjectionInput` | 管理操作。立即重放一条最终一致的业务投影，而不必等 worker 的下一个退避窗口。 |

流程定义管理——创建、部署、发布、更新、启停——刻意**不在**这份契约里。它属于设计期管理，信任姿态不同，由 `approval/flow` 资源提供。

## 句柄就是事务边界

每个方法都接收操作所运行的 `orm.DB` 句柄，而这个句柄决定了事务：

- 传入一个**已打开的 `RunInTx` 作用域**的句柄，操作就加入那个事务，于是业务写入与它触发的审批动作要么一起提交，要么一起回滚；
- 传入一个**普通句柄**，操作就自己开启并提交一个事务。

句柄是显式参数，而不是从 context 里读，原因在于 `orm.DB.RunInTx` 并不会把事务句柄挂到它交给回调的 context 上——一个只接收 context 的方法，会在调用方的事务里悄悄再开一个互不相干的事务。传入 nil 句柄会被 `approval.ErrDBRequired` 拒绝，而不是回落到框架自己的句柄；并且句柄必须属于**主数据源**：审批表在那里，领域事件也是用 `event.WithTx` 对着它发布的。

`fiber.Ctx` 会在绑定前被拆包。这一点很关键，因为 `contextx.SetDB` 是**就地**修改 `fiber.Ctx` 的 Locals：把事务绑上去会让整个请求都指向一个在调用方提交时就死掉的事务。

### 要传 `RunInTx` 交给回调的那个 context

而不是你调用 `RunInTx` 时用的那个。`RunInTx` 确实会往回调 context 上放一样东西：`orm.OnCommit` 读取的提交钩子注册表，`tx_memory` 事件传输正是通过它延迟投递的。因此在使用 `tx_memory` 的路由下，用外层 context 会让发布以 `orm.ErrNoCommitScope` 失败，并把整个动作回滚掉。

```go
db.RunInTx(ctx, func(txCtx context.Context, tx orm.DB) error {
    if err := writeBusinessRow(txCtx, tx); err != nil {
        return err
    }

    return svc.ApproveTask(txCtx, tx, approval.ApproveTaskInput{
        TaskID:   taskID,
        Operator: operator,
        Opinion:  "approved by the settlement routine",
        Caller:   approval.SystemCaller,
    })
})
```

### 绑定审计操作人

审计列（`created_by` / `updated_by`）写的是绑定在句柄上的操作人，而 `orm.DB.WithNamedArg` 只在连接池级别有效——对事务句柄调用它会 **panic**。要在开启事务之前把人绑好，再使用该句柄产出的 tx：

```go
acting := db.WithNamedArg(orm.PlaceholderKeyOperator, userID)
acting.RunInTx(ctx, func(txCtx context.Context, tx orm.DB) error { … })
```

没有绑定操作人的句柄——也就是注入进来的主 `orm.DB`——会往这些列里写 `orm.OperatorSystem`。无论哪种情况，操作日志都会记录入参里的 `Operator`，所以不绑定损失的是行级归属，而不是审计轨迹。

## Service 不做鉴权

:::caution[鉴权是调用方的责任]
把守 HTTP 端点的 RBAC 令牌（`approval.instance.terminate`、`approval.task.reassign`、`approval.binding.retry`）挂在 API **操作**上，而不在其背后的命令里。`TerminateInstance`、`ReassignTask`、`RetryBusinessProjection` 只强制 `Caller` 携带的租户范围。宿主把它们中的任何一个包装到自己的接口上时，就自己承担了这个检查——调用点和框架自己的写法一模一样，因此漏掉时没有任何东西会提醒你。
:::

每个入参都带一个 `Caller`（`approval.CallerContext`），而且是必填的：零值会失败关闭——`StartInstanceInput` 报的是 `ErrFlowNotFound` 而不是权限错误，所以调用方无法借此探测别的租户有哪些流程。可信自动化用 `approval.SystemCaller`，需要钉死单一租户则用 `approval.CallerContext{TenantID: …}`。

`StartInstanceInput.Globals` 同样需要谨慎。在请求路径上，框架之所以坚持通过 `InstanceGlobalsResolver` 在服务端解析 globals，正是因为它们会左右条件分支；在这里由宿主直接传入，因此请传你自己信得过的值，而不要把客户端提交的任何东西原样转发进来。

## 错误

错误就是 `approval` 包里那些公开哨兵（`ErrFlowNotActive`、`ErrNotAllowedInitiate`、`ErrTaskNotPending` 等），可以用 `errors.Is` 匹配——`result.Error` 按业务码比较，因此不管消息渲染成哪种语言，哨兵都能匹配上。

```go
if errors.Is(err, approval.ErrTaskNotPending) {
    // 已经有别人处理过了
}
```

## 可重入

`approval.Service` 是第一条能让宿主**从正在执行的命令内部**重新进入命令管线的路径：`InstanceLifecycleHook` 在事务中触发，而它自身可能又去操作另一条任务。CQRS collector 已经处理了这种情况——内层分发会复用 context 上已有的 collector 并放弃所有权，因此仍由最外层分发负责 flush，事件顺序也仍与发生顺序一致。内层分发失败时只回退它自己追加的那部分，于是外层命令即便吞掉了内层错误，也仍保留自己的条目，只丢掉那些描述“并未真正发生的写入”的条目。

---

下一步：阅读 [事件与集成](./integration.md) 了解这些操作会触发的扩展点，或阅读 [RPC 资源](./resources.md) 了解调用同一接口的 HTTP 接口面。
