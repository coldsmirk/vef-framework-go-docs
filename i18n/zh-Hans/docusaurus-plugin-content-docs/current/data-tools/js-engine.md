---
sidebar_position: 2
---

# JS 引擎

`js` 包内嵌了由 [goja](https://github.com/dop251/goja) 驱动的沙箱化
JavaScript 运行时。自 v0.39 起，它围绕 **Engine / Runtime / Lib** 三元
架构构建：不可变的 `Engine` 持有一组经过校验的库并批量产出一次性的
`Runtime`；具备副作用的能力（HTTP、SQL、缓存、事件……）都是按运行时安装的
`Lib`，脚本只能触及被安装的内容。

框架用该引擎运行集成适配器脚本、签名与验证脚本等脚本执行接缝；应用也可以
直接使用它满足自己的脚本化需求。

## 快速开始

```go
import "github.com/coldsmirk/vef-framework-go/js"

engine, err := js.NewEngine()
if err != nil {
    return err
}

// 一个引擎服务整个应用；每次执行创建一个运行时。
rt, err := engine.NewRuntime(js.WithRunTimeout(5 * time.Second))
if err != nil {
    return err
}

value, err := rt.RunString(ctx, `1 + 2`)
fmt.Println(value.Export()) // 3
```

## Engine

`js.NewEngine(opts ...EngineOption)` 构建引擎并急切校验库集合：nil 库、
空名称、名称冲突（跨标准库、常驻库与目录库）都会使构建失败。

| 选项 | 行为 |
| --- | --- |
| `js.WithBaseLibs(libs...)` | 注册常驻库：安装进每个运行时，无需 opt-in。仅用于安全、普适的工具 |
| `js.WithLibs(libs...)` | 注册目录（catalog）库：每个运行时经 `js.EnableLibs` 按需激活 |
| `js.WithoutStdLibs()` | 构建裸引擎，其运行时不带标准库 bundle |

错误：`js.ErrInvalidLib`（nil 库 / 空名）、`js.ErrDuplicateLib`（名称
冲突）、`js.ErrLibNotFound`（激活未注册的名称）。

引擎可并发使用。框架提供共享的 DI 引擎，内置能力库按两档预先接线：

- **常驻**（安装进每个运行时）：`console`、`crypto`、`cache`（内存存储，
  键前缀 `js:`）；
- **目录 opt-in**（运行时经 `EnableLibs` 按需激活）：`events`、`http`
  （默认 30s 超时）、`sql`（主数据库，只读）。

应用用 `vef.ProvideJSLib` 贡献或替换库：与内置同名的库在其所在档位内
替换，新名称加入 opt-in 目录。

```go
vef.ProvideJSLib(func(db orm.DB) js.Lib {
    return jssql.New(db, config.Postgres, jssql.WithExecute())
})
```

## Runtime

`Engine.NewRuntime(opts ...RuntimeOption)` 创建携带引擎基线加已激活目录库
的全新运行时。运行时**不**可并发使用，同一时间只允许一个 Run 调用在途，
用完即弃——每个 goroutine、每次执行创建一个。

| 选项 | 行为 |
| --- | --- |
| `js.EnableLibs(names...)` | 按名称激活目录库；安装顺序即参数顺序；重复激活为 no-op；未知名称使 `NewRuntime` 返回 `ErrLibNotFound` |
| `js.WithRunTimeout(d)` | 为每次 Run 调用设置上限；与调用方 context 组合——更早的期限生效 |
| `js.WithMaxCallStackSize(n)` | 约束 JS 调用栈深度，防失控递归 |

| 方法 | 契约 |
| --- | --- |
| `RunProgram(ctx, program)` / `RunString(ctx, source)` | 在 `ctx` 下执行：取消会中断运行中的脚本，并经 `Context()` 传导到宿主库的在途 IO；此时返回的错误即 context 的错误 |
| `Set(name, value)` | 绑定全局变量 |
| `Context()` | 在途 Run 调用的 context（空闲时为 `context.Background()`）；宿主库必须经它发起 IO |
| `AsFunction(value)` | 把值转换为可调用的 `js.Func` 句柄；调用必须发生在当前驱动运行时的 goroutine 上 |
| `VM()` | 底层 `*goja.Runtime`，用于高级库开发 |

每个运行时都配置了 JSON 字段名映射
（`goja.TagFieldNameMapper("json", true)`），经 `Set` 传入的 Go 结构体在
脚本中按 json 标签读取。

## 标准库 Bundle

除非引擎以 `WithoutStdLibs` 构建，每个运行时都自带一个 vendored esbuild
bundle，按各库生态原生的全局名安装：

| 全局名 | 库 | 用途 |
| --- | --- | --- |
| `BigNumber` | bignumber.js | 任意精度十进制 |
| `dayjs` | Day.js | 日期时间解析、格式化、运算 |
| `fxp` | fast-xml-parser | XML 解析与构建（`fxp.XMLParser`、`fxp.XMLBuilder`） |
| `radashi` | Radashi | 函数式工具集 |
| `z` | Zod | Schema 校验，内置 `en` 与 `zh-CN` 语言包（默认 `zh-CN`） |
| `URL` / `URLSearchParams` | core-js polyfill | WHATWG URL 处理 |

## 内置能力库

具备副作用的能力是 `js/*` 下的独立包，接入 DI 引擎的目录。每个库按包名
安装一个全局对象；运行时只能看到为它激活的库。所有宿主 IO 都经
`Runtime.Context()`，取消因此能到达阻塞中的 Go 调用。

### `jssql` —— 全局 `sql`

```js
sql.queryList('SELECT name FROM users WHERE age > ?', 18)  // → [{...}, ...]；无匹配行返回 []
sql.queryOne('SELECT ... WHERE id = ?', id)                // → {...} | null
sql.execute('UPDATE ...', args)                            // → { rowsAffected }
```

- `jssql.New(db, kind, opts...)` 在选定的数据源上构建库；由调用方决定脚本
  能触达什么。
- 只提供占位符绑定——刻意不提供字符串拼接辅助。
- 默认只读，由 AST 守卫（`sqlguard`）fail-closed 强制：写 CTE、堆叠语句、
  方言特有的副作用函数被拒绝；无法解析的 SQL 一律拒绝
  （`ErrQueryNotReadOnly`）。
- 未以 `jssql.WithExecute()` 构建时，`sql.execute` 抛出
  `ErrExecuteDisabled`。
- 结果集有上限（`jssql.WithMaxRows`，默认 1000；`ErrTooManyRows` 提示脚本
  加 `LIMIT`）。

### `jshttp` —— 全局 `http`

fetch 标准的同步化版本：

```js
http.fetch(url, { method, headers, query, body, redirect, timeout })
http.get(url, options?)          // fetch 语法糖，put/patch 同理
http.post(url, body, options?)
http.delete(url, options?)
```

- 每次调用返回 `{ status, statusText, ok, url, redirected, headers, body,
  text(), json(), arrayBuffer() }`；失败以可捕获异常抛出。响应头名小写化，
  多值以 `", "` 连接。
- 超出 fetch 的部分：`query` 追加 URL 参数、`timeout`（毫秒）取代
  AbortSignal、`redirect` 支持 `follow` / `error` / `manual`。
- 默认无任何限制——超时、响应体上限、主机白名单、私网守卫都通过选项
  （`jshttp.New(opts...)`）按需开启。

### `jscache` —— 全局 `cache`

脚本跨执行保存状态的唯一通道（运行时每次运行后即被丢弃）：

```js
cache.set('counter', { n: 1 })       // 默认 TTL
cache.set('token', value, 60000)     // TTL 毫秒
cache.get('counter')                 // → 值 | null
cache.has('counter')                 // → 布尔
cache.delete('counter')
```

`jscache.New(store, opts...)` 接受任意 `cache.Cache[any]` —— 单节点状态用
内存、共享状态用 Redis —— 外加 `WithKeyPrefix` 做键命名空间。

### `jsevents` —— 全局 `events`

```js
events.publish('report.generated', { reportId: id })
```

负载 JSON 编码后以 `event.RawPayload` 发布，Go 订阅方照常用
`event.SubscribeTyped` 解码。publish 刻意是唯一动词——订阅是长生命周期的，
属于宿主而不属于按执行创建的运行时。`jsevents.WithAllowedTypes(patterns...)`
可约束可发布的类型命名空间。

### `jscrypto` —— 全局 `crypto`

```js
crypto.md5(data)                    // 同理 sha1 / sha256 / sha512 / sm3
crypto.hmac('sha256', key, data)    // hex 摘要
crypto.base64Encode(data) / crypto.base64Decode(encoded)
crypto.hexEncode(data) / crypto.hexDecode(encoded)
crypto.uuid()
```

摘要为小写 hex；输入为 UTF-8 字符串。弱摘要（md5、sha1）只为兼容遗留 API
签名——密码存储属于安全模块的编码器。

### `jsconsole` —— 全局 `console`

```js
console.info('processed', count, payload)   // 亦有 warn / error
```

参数以空格连接；字符串原样、错误取其消息、其余 JSON 编码。由
`logx.Logger` 支撑（传入命名 logger 可区分脚本输出）。

## 编写自己的 Lib

```go
type Lib interface {
    Name() string          // Engine 内唯一键；按惯例即其安装的全局名
    Install(rt *js.Runtime) error
}
```

- 纯 JS 库：`js.SourceLib(name, source)`（急切编译）或
  `js.ProgramLib(name, program)`。
- 宿主能力：以 `rt.Set(...)` 实现 `Install`，只持有共享且 goroutine 安全的
  依赖——绝不持有按运行时状态——并经 `rt.Context()` 发起 IO，让取消得以
  传导。

```go
lib, err := js.SourceLib("fmtx", `var fmtx = { pad: (s, n) => String(s).padStart(n, '0') };`)
engine, err := js.NewEngine(js.WithLibs(lib))
rt, err := engine.NewRuntime(js.EnableLibs("fmtx"))
```

## 编译辅助

| API | 契约 |
| --- | --- |
| `js.Compile(name, source, strict)` / `js.MustCompile(...)` | 为重复执行预编译（goja 别名） |
| `js.Parse(name, source)` | 返回 `*js.AstProgram`（禁用 source map） |
| `js.Runtime`、`js.Value`、`js.Object`、`js.Program`、`js.AstProgram` | goja 类型别名 |
| `js.IsNaN`、`js.IsString`、`js.IsBigInt`、`js.IsNumber`、`js.IsInfinity`、`js.IsUndefined`、`js.IsNull` | goja 辅助函数别名 |

## 框架在哪些地方运行脚本

| 接缝 | 可见库 |
| --- | --- |
| [集成适配器脚本](../integration/outbound#适配器脚本环境) | 基线 + `errors`、`codes`、按系统的作用域 `http` / `sql` |
| [集成签名/验证脚本](../integration/outbound#出站认证-scheme) | 仅基线（零 IO）+ `request`、`params` 绑定 |
| [表达式引擎](./expression) | 自有求值器（非本引擎） |

## 线程安全

> **警告**：`Engine` 可并发使用；`Runtime` 不可。每个 goroutine、每次执行
> 创建一个运行时，用完即弃。
