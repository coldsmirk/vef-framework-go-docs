---
sidebar_position: 5
---

# 数据库命名规范

本页定义 VEF 应用的数据库命名与 DDL 规范。

## 1. 总体结构

### 1.1 事务包裹

每个模块的 DDL 脚本必须包裹在事务中：

```sql
BEGIN;

-- DDL 语句 ...

COMMIT;
```

### 1.2 幂等 DDL

DDL 脚本必须具备幂等性，重复执行不得报错。应优先使用 `IF NOT EXISTS` / `IF EXISTS` 子句：

```sql
-- 建表
CREATE TABLE IF NOT EXISTS sys_user (...);

-- 删表（当需要时）
DROP TABLE IF EXISTS sys_user;

-- 建索引
CREATE INDEX IF NOT EXISTS idx_sys_user__email ON sys_user (email);
```

注意：

- `ALTER TABLE ADD CONSTRAINT` 不支持 `IF NOT EXISTS`
- 这种场景下，应通过查询系统目录判断是否已存在，或使用 `DROP ... IF EXISTS` 后重新创建的方式实现幂等

### 1.3 模块前缀

表名必须以模块缩写作为前缀，并使用下划线分隔。新模块前缀在团队达成共识前不得引入。

| 前缀 | 模块 | 说明 |
| --- | --- | --- |
| `sys_` | 系统基础模块 | 用户、角色、权限、配置等 |
| `md_` | 主数据模块 | 机构、部门、职员、编码等 |
| `hr_` | 人力资源模块 | 人事业务相关 |

### 1.4 脚本内注释

每张表之前必须使用单行中文注释标注表名：

```sql
-- 用户
CREATE TABLE sys_user (
    ...
);
```

## 2. 命名规范

### 2.1 通用规则

- 所有名称必须使用小写 `snake_case`
- 禁止使用数据库保留字作为裸名称；如不可避免，SQL 中必须使用双引号引用，例如 `"group"`、`"key"`
- 名称必须见名知意，阅读者不应依赖外部文档才能理解其用途
- 优先使用完整英文单词，不得随意缩写，例如 `organization` 而不是 `org`，`department` 而不是 `dept`
- 仅允许使用广泛认知的常见缩写；其他缩写必须先在团队内达成共识
- 同一概念在所有表中必须保持同名，例如统一使用 `remark`，不得混用 `note`、`comment`、`remark`
- 表名必须使用名词，列名必须使用名词
- 所有表名必须使用单数形式

允许的常见缩写：

| 缩写 | 全称 | 说明 |
| --- | --- | --- |
| `id` | identifier | 标识符 |
| `app` | application | 应用 |
| `config` | configuration | 配置 |
| `info` | information | 信息 |
| `stat` | statistics | 统计 |
| `log` | log | 日志 |

### 2.2 表命名

```
{模块前缀}_{实体名}
```

| 类型 | 示例 |
| --- | --- |
| 实体表 | `sys_user`、`md_organization` |
| 关联表 | `sys_user_role`、`md_department_staff` |
| 日志表 | `sys_login_log`、`sys_audit_log` |
| 规则/定义表 | `sys_config_definition`、`sys_sequence_rule` |

### 2.3 视图命名

```
vw_{模块前缀}_{视图名}
mv_{模块前缀}_{视图名}
```

| 类型 | 示例 |
| --- | --- |
| 普通视图 | `vw_sys_user_detail`、`vw_md_staff_summary` |
| 聚合视图 | `vw_hr_attendance_stat` |
| 物化视图 | `mv_sys_daily_login_stat`、`mv_hr_monthly_attendance` |

### 2.4 列命名

| 场景 | 命名规则 | 示例 |
| --- | --- | --- |
| 主键 | `id` | `id` |
| 外键引用 | `{被引用表短名}_id` | `role_id`、`organization_id`、`app_id` |
| 自引用（树形层级） | `parent_id` | `parent_id` |
| 布尔字段 | `is_{形容词/状态}` | `is_active`、`is_locked`、`is_default` |
| 时间戳字段 | `{动作}_at` | `created_at`、`updated_at`、`password_updated_at` |
| 排序字段 | `sort_order` | `sort_order` |
| 备注字段 | `remark` | `remark` |
| 元数据字段 | `meta` | `meta` |

### 2.5 约束命名

所有约束都必须显式命名，格式为 `{约束类型}_{表名}__{列名或语义}`。表名与列名或语义之间必须使用双下划线 `__` 分隔。

| 约束类型 | 前缀 | 命名格式 | 示例 |
| --- | --- | --- | --- |
| 主键 | `pk` | `pk_{表名}` | `pk_sys_user` |
| 唯一键 | `uk` | `uk_{表名}__{列名}` | `uk_sys_user__username`、`uk_sys_user__staff_id` |
| 外键 | `fk` | `fk_{表名}__{列名}` | `fk_sys_user__created_by` |
| 检查约束 | `ck` | `ck_{表名}__{列名}` | `ck_sys_user__gender` |

复合约束：

- 默认使用完整列名拼接，列名之间用单下划线连接
- 当完整命名超过 PostgreSQL 63 字节标识符长度限制时，允许使用语义化缩写

```sql
-- 默认：完整列名拼接
CONSTRAINT uk_sys_dictionary_item__dictionary_id_code UNIQUE (dictionary_id, code)
CONSTRAINT uk_sys_user_role__user_id_role_id UNIQUE (user_id, role_id)

-- 仅当完整列名拼接超长时：语义化缩写
CONSTRAINT uk_md_department__org_code UNIQUE (organization_id, code)
```

### 2.6 索引命名

索引命名必须区分索引类型。默认 B-tree 索引使用 `idx_` 前缀，其他索引类型使用对应前缀。

| 索引类型 | 前缀 | 命名格式 | 示例 |
| --- | --- | --- | --- |
| B-tree（默认） | `idx` | `idx_{表名}__{列名}` | `idx_sys_role__created_by` |
| GIN | `gin` | `gin_{表名}__{列名}` | `gin_md_medical_code__meta` |
| GiST | `gist` | `gist_{表名}__{列名}` | `gist_md_organization__location` |
| BRIN | `brin` | `brin_{表名}__{列名}` | `brin_sys_audit_log__created_at` |

复合索引中，多列默认使用完整列名并以单下划线连接。超长时允许使用语义化缩写，规则与复合约束一致。

```sql
CREATE INDEX idx_sys_audit_log__api_resource_api_action_api_version ON sys_audit_log (api_resource, api_action, api_version);

-- 超长时使用语义化缩写
CREATE INDEX idx_sys_audit_log__api ON sys_audit_log (api_resource, api_action, api_version);
```

部分索引：

- 索引名末尾必须追加 `__partial`
- 并在注释中说明过滤条件

```sql
-- 仅索引启用状态的用户
CREATE INDEX idx_sys_user__email__partial ON sys_user (email) WHERE is_active = TRUE;
```

覆盖索引（`INCLUDE`）：

- 索引名末尾必须追加 `__include`

```sql
-- 覆盖索引，避免回表查询
CREATE INDEX idx_sys_user__username__include ON sys_user (username) INCLUDE (name, email);
```

## 3. 数据类型规范

### 3.1 标准类型映射

| 用途 | 类型 | 说明 |
| --- | --- | --- |
| 主键 / 外键 | `VARCHAR(32)` | 存放去掉连字符的 UUID（32 字符） |
| 用户名称 | `VARCHAR(16)` | 人名等短名称 |
| 实体名称 | `VARCHAR(32)` | 角色名、部门名等 |
| 长名称 | `VARCHAR(64)` | 字典项名称、编码名称等 |
| 编码 | `VARCHAR(32)` | 各类业务编码 |
| URL / 文件路径 | `VARCHAR(128)` | 头像、链接、邮箱等 |
| 简介 / 长文本描述 | `VARCHAR(2048)` | 机构简介、失败原因等 |
| 备注 | `VARCHAR(256)` | 所有表的备注字段统一长度 |
| 布尔 | `BOOLEAN` | 搭配 `NOT NULL DEFAULT FALSE` |
| 排序 | `INTEGER` | 搭配 `NOT NULL DEFAULT 0` |
| 小范围整数 | `SMALLINT` | 如步长、宽度等 |
| 时间戳 | `TIMESTAMP` | 不带时区，使用 `LOCALTIMESTAMP` |
| 日期 | `DATE` | 仅日期场景，如出生日期 |
| 枚举 | `VARCHAR(n)` | 配合 `CHECK` 约束使用，长度根据实际枚举值决定 |
| 元数据 | `JSONB` | 可扩展的结构化附加信息 |
| 大文本 | `TEXT` | 无长度限制的文本，如配置值 |

### 3.2 枚举处理

不得使用 PostgreSQL `ENUM` 类型。统一使用 `VARCHAR(n)` 搭配 `CHECK` 约束，长度根据实际枚举值决定。

```sql
-- VARCHAR(1)：单字符枚举
gender    VARCHAR(1) NOT NULL DEFAULT 'U',
CONSTRAINT ck_sys_user__gender CHECK (gender = ANY (ARRAY['M', 'F', 'U']))
```

```sql
-- VARCHAR(8)：多值枚举
overflow_strategy  VARCHAR(8) NOT NULL DEFAULT 'error',
CONSTRAINT ck_sys_sequence_rule__overflow_strategy CHECK (overflow_strategy = ANY (ARRAY['error', 'reset', 'extend']))
```

## 4. 表结构模板

### 4.1 标准实体表

适用于需要完整增删改查的业务实体，例如用户、角色、机构、部门等：

```sql
CREATE TABLE {模块前缀}_{实体名} (
    -- ① 主键
    id                       VARCHAR(32) NOT NULL,

    -- ② 审计字段（固定 4 列，顺序不变）
    created_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    created_by               VARCHAR(32) NOT NULL DEFAULT 'system',
    updated_by               VARCHAR(32) NOT NULL DEFAULT 'system',

    -- ③ 业务字段
    ...

    -- ④ 通用可选字段（按需选用，顺序为：is_active → sort_order → remark → meta）
    is_active                BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order               INTEGER NOT NULL DEFAULT 0,
    remark                   VARCHAR(256),
    meta                     JSONB,

    -- ⑤ 约束（顺序：PK → UK → CK → FK）
    CONSTRAINT pk_{表名} PRIMARY KEY (id),
    CONSTRAINT uk_{表名}__{列名} UNIQUE (...),
    CONSTRAINT ck_{表名}__{列名} CHECK (...),
    CONSTRAINT fk_{表名}__created_by FOREIGN KEY (created_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_{表名}__updated_by FOREIGN KEY (updated_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
```

### 4.2 关联表（多对多）

适用于实体间多对多关系，例如用户角色、部门职员、角色权限等。这类表不需要 `updated_at` / `updated_by`。

```sql
CREATE TABLE {模块前缀}_{实体A}_{实体B} (
    -- ① 主键
    id                       VARCHAR(32) NOT NULL,

    -- ② 审计字段（仅 2 列）
    created_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    created_by               VARCHAR(32) NOT NULL DEFAULT 'system',

    -- ③ 关联外键
    {实体A}_id               VARCHAR(32) NOT NULL,
    {实体B}_id               VARCHAR(32) NOT NULL,

    -- ④ 附加业务字段（如有）
    ...

    -- ⑤ 约束
    CONSTRAINT pk_{表名} PRIMARY KEY (id),
    CONSTRAINT uk_{表名}__{实体A}_id_{实体B}_id UNIQUE ({实体A}_id, {实体B}_id),
    CONSTRAINT fk_{表名}__{实体A}_id FOREIGN KEY ({实体A}_id) REFERENCES {实体A表}(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_{表名}__{实体B}_id FOREIGN KEY ({实体B}_id) REFERENCES {实体B表}(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_{表名}__created_by FOREIGN KEY (created_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
```

### 4.3 日志表（只写不改）

适用于登录日志、审计日志等不可变记录。这类表不需要 `updated_at` / `updated_by`。

```sql
CREATE TABLE {模块前缀}_{日志名}_log (
    -- ① 主键
    id                       VARCHAR(32) NOT NULL,

    -- ② 审计字段（仅 2 列）
    created_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    created_by               VARCHAR(32) NOT NULL DEFAULT 'system',

    -- ③ 业务字段
    ...

    -- ④ 约束
    CONSTRAINT pk_{表名} PRIMARY KEY (id),
    CONSTRAINT fk_{表名}__created_by FOREIGN KEY (created_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
```

## 5. 外键规则

### 5.1 `ON DELETE` / `ON UPDATE` 策略

| 场景 | ON DELETE | ON UPDATE |
| --- | --- | --- |
| `created_by` / `updated_by` -> `sys_user` | `RESTRICT` | `CASCADE` |
| 自引用层级 `parent_id` | `RESTRICT` | `CASCADE` |
| 业务实体间强依赖 | `RESTRICT` | `CASCADE` |
| 关联表引用主实体 | `CASCADE` | `CASCADE` |
| 配置键引用定义表 | `RESTRICT` | `CASCADE` |

### 5.2 审计外键

所有表中的 `created_by` 和 `updated_by`（如存在）都必须引用 `sys_user(id)`：

```sql
CONSTRAINT fk_{表名}__created_by FOREIGN KEY (created_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE,
CONSTRAINT fk_{表名}__updated_by FOREIGN KEY (updated_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE
```

## 6. 索引规则

### 6.1 自动创建的索引

PostgreSQL 会自动为 `PRIMARY KEY` 和 `UNIQUE` 约束创建索引，这些索引不得重复手工创建。

### 6.2 必建索引

| 列 | 排序 | 说明 |
| --- | --- | --- |
| `created_at` | `DESC` | 需要按创建时间倒序查询时建立 |
| `created_by` | 默认升序 | 所有表必建 |
| 非唯一外键列 | 默认升序 | 如 `parent_id`、`app_id`、`organization_id` 等 |
| `sort_order` | 默认升序 | 表确实使用排序字段时按需建立 |

### 6.3 特殊索引

| 场景 | 类型 | 示例 |
| --- | --- | --- |
| JSONB 查询 | GIN 索引 | `CREATE INDEX gin_md_medical_code__meta ON md_medical_code USING gin (meta);` |
| 复合条件查询 | 复合 B-tree 索引 | `CREATE INDEX idx_sys_audit_log__api ON sys_audit_log (api_resource, api_action, api_version);` |
| 树形路径查询 | B-tree 索引 | `CREATE INDEX idx_md_medical_code_category__tree_path ON md_medical_code_category (tree_path);` |

### 6.4 部分索引（Partial Index）

当查询长期针对某一固定过滤条件时，应使用部分索引以减少索引体积并提升查询性能：

```sql
-- 仅索引启用状态的用户邮箱
CREATE INDEX idx_sys_user__email__partial ON sys_user (email) WHERE is_active = TRUE;
```

适用场景：

- 状态过滤，例如 `WHERE is_active = TRUE`
- 仅索引非空值，例如 `WHERE deleted_reason IS NOT NULL`

### 6.5 覆盖索引（INCLUDE）

当查询返回列可以被索引完全覆盖时，应使用 `INCLUDE` 子句避免回表：

```sql
-- 按用户名查询时，同时覆盖 name 和 email，避免回表
CREATE INDEX idx_sys_user__username__include ON sys_user (username) INCLUDE (name, email);
```

注意：

- `INCLUDE` 列不参与索引排序和搜索
- 它们只作为附带数据存储在索引中

### 6.6 生产环境建索引

在已有数据的生产环境中创建索引时，必须使用 `CONCURRENTLY` 以避免阻塞写入：

```sql
-- CONCURRENTLY 不阻塞表的读写操作
CREATE INDEX CONCURRENTLY idx_sys_user__email ON sys_user (email);
```

注意事项：

- `CREATE INDEX CONCURRENTLY` 不能放在事务块中执行
- 它执行时间会比普通建索引更长
- 如果失败，PostgreSQL 可能留下 `INVALID` 状态的索引，重试前必须先手工删除
- 开发环境初始化 DDL 不需要使用 `CONCURRENTLY`，这一规则仅用于生产环境变更

## 7. COMMENT 规范

### 7.1 要求

- 每张表都必须有 `COMMENT ON TABLE`
- 每个列都必须有 `COMMENT ON COLUMN`
- 注释必须使用简体中文
- `COMMENT` 语句必须紧跟在 `CREATE TABLE` 之后、索引语句之前

### 7.2 格式

```sql
CREATE TABLE sys_user (
    ...
);
COMMENT ON TABLE sys_user IS '用户';
COMMENT ON COLUMN sys_user.id IS '主键';
COMMENT ON COLUMN sys_user.created_at IS '创建时间';
COMMENT ON COLUMN sys_user.updated_at IS '更新时间';
COMMENT ON COLUMN sys_user.created_by IS '创建人';
COMMENT ON COLUMN sys_user.updated_by IS '更新人';
-- 业务列注释 ...
CREATE INDEX ...
```

### 7.3 审计字段固定注释

| 列 | 注释 |
| --- | --- |
| `id` | `主键` |
| `created_at` | `创建时间` |
| `updated_at` | `更新时间` |
| `created_by` | `创建人` |
| `updated_by` | `更新人` |
| `remark` | `备注` |
| `meta` | `元数据` |
| `sort_order` | `排序` |
| `is_active` | `是否启用` |

### 7.4 外键列注释

外键列的注释必须使用 `{被引用实体名称}ID` 格式，不得写成“主键”。

| 列 | 正确注释 | 错误注释 |
| --- | --- | --- |
| `user_id` | `用户ID` | ~~用户主键~~ |
| `role_id` | `角色ID` | ~~角色主键~~ |
| `organization_id` | `机构ID` | ~~机构主键~~ |
| `parent_id` | `上级ID` | ~~上级主键~~ |

## 8. 列定义格式

### 8.1 对齐

列名与数据类型之间必须对齐，保证同一张表内视觉整齐：

```sql
CREATE TABLE sys_user (
    id                       VARCHAR(32) NOT NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    username                 VARCHAR(32) NOT NULL,
    is_active                BOOLEAN NOT NULL DEFAULT FALSE,
    remark                   VARCHAR(256),
    meta                     JSONB,
    ...
);
```

格式要求：

- 列名左对齐
- 数据类型起始位置对齐，建议列名区域宽度为 25 个字符
- `NOT NULL`、`DEFAULT` 等修饰符紧跟在类型之后

### 8.2 列顺序

表内列必须按以下顺序排列：

1. `id`
2. `created_at`、`updated_at`
3. `created_by`、`updated_by`
4. 外键引用列，例如 `parent_id`、`{entity}_id`
5. 核心业务字段
6. 状态或开关字段，例如 `is_active`、`is_locked`、`status`
7. `sort_order`（如有）
8. `remark`（如有）
9. `meta`（如有）

### 8.3 约束顺序

`CONSTRAINT` 子句必须按以下顺序排列：

1. `PRIMARY KEY`
2. `UNIQUE`
3. `CHECK`
4. `FOREIGN KEY`

其中，`created_by` 和 `updated_by` 的外键必须放在最后。

## 9. 默认值规范

| 类型 / 用途 | 默认值 |
| --- | --- |
| `created_at` | `LOCALTIMESTAMP` |
| `updated_at` | `LOCALTIMESTAMP` |
| `created_by` | `'system'` |
| `updated_by` | `'system'` |
| `is_active` 等布尔字段 | `FALSE` |
| `sort_order` | `0` |
| `gender` | `'U'` |
| 业务布尔字段在确有必要默认启用时 | `TRUE`，但必须明确说明原因 |

### 9.1 审计字段更新机制

- `created_at` 和 `created_by` 只在插入时由数据库默认值填充，后续不得修改
- `updated_at` 和 `updated_by` 必须由应用程序在每次更新时赋值
- 不得使用数据库触发器承担这一职责

### 9.2 系统用户初始化前提

`created_by` 和 `updated_by` 的默认值 `'system'` 通过外键引用 `sys_user(id)`。因此，数据库初始化时必须先插入一条 `id = 'system'` 的系统用户记录，否则依赖这些默认值的插入操作会因外键校验失败：

```sql
INSERT INTO sys_user (id, username, name) VALUES ('system', 'system', '系统');
```

## 10. ALTER TABLE 规范

当跨模块存在循环依赖时，应在目标表创建完成后使用 `ALTER TABLE` 补充外键：

```sql
-- 在 md_staff 表创建之后补充 sys_user.staff_id 的外键
ALTER TABLE sys_user
    ADD CONSTRAINT fk_sys_user__staff_id FOREIGN KEY (staff_id) REFERENCES md_staff(id) ON DELETE RESTRICT ON UPDATE CASCADE;
```

此类语句必须放在相关 DDL 文件末尾、`COMMIT;` 之前。

## 11. 完整示例

以下是标准实体表的完整 DDL 示例：

```sql
-- 角色
CREATE TABLE sys_role (
    id                       VARCHAR(32) NOT NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    created_by               VARCHAR(32) NOT NULL DEFAULT 'system',
    updated_by               VARCHAR(32) NOT NULL DEFAULT 'system',
    name                     VARCHAR(32) NOT NULL,
    is_active                BOOLEAN NOT NULL DEFAULT FALSE,
    remark                   VARCHAR(256),
    meta                     JSONB,

    CONSTRAINT pk_sys_role PRIMARY KEY (id),
    CONSTRAINT fk_sys_role__created_by FOREIGN KEY (created_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_sys_role__updated_by FOREIGN KEY (updated_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
COMMENT ON TABLE sys_role IS '角色';
COMMENT ON COLUMN sys_role.id IS '主键';
COMMENT ON COLUMN sys_role.created_at IS '创建时间';
COMMENT ON COLUMN sys_role.updated_at IS '更新时间';
COMMENT ON COLUMN sys_role.created_by IS '创建人';
COMMENT ON COLUMN sys_role.updated_by IS '更新人';
COMMENT ON COLUMN sys_role.name IS '角色名称';
COMMENT ON COLUMN sys_role.is_active IS '是否启用';
COMMENT ON COLUMN sys_role.remark IS '备注';
COMMENT ON COLUMN sys_role.meta IS '元数据';
CREATE INDEX idx_sys_role__created_at ON sys_role (created_at DESC);
CREATE INDEX idx_sys_role__created_by ON sys_role (created_by);
```

以下是关联表的完整 DDL 示例：

```sql
-- 用户角色关系
CREATE TABLE sys_user_role (
    id                       VARCHAR(32) NOT NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    created_by               VARCHAR(32) NOT NULL DEFAULT 'system',
    user_id                  VARCHAR(32) NOT NULL,
    role_id                  VARCHAR(32) NOT NULL,
    effective_from           TIMESTAMP,
    effective_to             TIMESTAMP,

    CONSTRAINT pk_sys_user_role PRIMARY KEY (id),
    CONSTRAINT uk_sys_user_role__user_id_role_id UNIQUE (user_id, role_id),
    CONSTRAINT fk_sys_user_role__role_id FOREIGN KEY (role_id) REFERENCES sys_role(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sys_user_role__user_id FOREIGN KEY (user_id) REFERENCES sys_user(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sys_user_role__created_by FOREIGN KEY (created_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
COMMENT ON TABLE sys_user_role IS '用户角色关系';
COMMENT ON COLUMN sys_user_role.id IS '主键';
COMMENT ON COLUMN sys_user_role.created_at IS '创建时间';
COMMENT ON COLUMN sys_user_role.created_by IS '创建人';
COMMENT ON COLUMN sys_user_role.user_id IS '用户ID';
COMMENT ON COLUMN sys_user_role.role_id IS '角色ID';
COMMENT ON COLUMN sys_user_role.effective_from IS '生效时间';
COMMENT ON COLUMN sys_user_role.effective_to IS '失效时间';
CREATE INDEX idx_sys_user_role__created_by ON sys_user_role (created_by);
CREATE INDEX idx_sys_user_role__role_id ON sys_user_role (role_id);
```

## 延伸阅读

- [应用代码命名规范](./application-naming-conventions)：Go 代码与 API 命名规则
- [模型](../guide/models)：Go 字段名、JSON tag 与数据库列名之间的关系
