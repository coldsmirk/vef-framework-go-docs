---
sidebar_position: 6
title: Database Conventions
---

# Database Conventions

This page defines the mandatory database naming and DDL conventions for VEF application projects.

## 1. Overall Structure

### 1.1 Transaction Wrapper

Each module's DDL script must be wrapped in a transaction:

```sql
BEGIN;

-- DDL statements ...

COMMIT;
```

### 1.2 Idempotent DDL

DDL scripts must be idempotent. Re-running the same script must not fail. Use `IF NOT EXISTS` and `IF EXISTS` whenever the database supports them:

```sql
-- Create table
CREATE TABLE IF NOT EXISTS sys_user (...);

-- Drop table when needed
DROP TABLE IF EXISTS sys_user;

-- Create index
CREATE INDEX IF NOT EXISTS idx_sys_user__email ON sys_user (email);
```

Note:

- `ALTER TABLE ADD CONSTRAINT` does not support `IF NOT EXISTS`
- for that case, idempotency must be implemented by checking system catalogs first, or by using `DROP ... IF EXISTS` followed by recreation

### 1.3 Module Prefixes

Table names must start with a module prefix followed by an underscore. A new module prefix must not be introduced until the team agrees on it.

| Prefix | Module | Description |
| --- | --- | --- |
| `sys_` | system foundation module | users, roles, permissions, configuration |
| `md_` | master data module | organizations, departments, staff, codes |
| `hr_` | human resources module | HR-related business data |

## 2. Naming Rules

### 2.1 General Rules

- all names must use lowercase `snake_case`
- database reserved words should be avoided as bare names whenever possible, but not at the cost of semantic accuracy or naming consistency; when a reserved word is still the best standard-compliant name, it must be quoted with double quotes in SQL, for example `"group"` or `"key"`
- names must be self-explanatory and convey meaning without requiring external documentation
- full words must be preferred over arbitrary abbreviations, for example `organization` instead of `org`, and `department` instead of `dept`
- only widely understood abbreviations may be used by default; any other abbreviation must first be agreed on by the team
- the same concept must use the same name everywhere, for example always use `remark` instead of mixing `note`, `comment`, and `remark`
- table names must use nouns, and column names must use nouns
- table names must use singular form

Allowed common abbreviations:

| Abbreviation | Full form | Description |
| --- | --- | --- |
| `id` | identifier | identifier |
| `app` | application | application |
| `config` | configuration | configuration |
| `info` | information | information |
| `stat` | statistics | statistics |
| `log` | log | log |

### 2.2 Table Naming

```
{module_prefix}_{entity_name}
```

| Type | Examples |
| --- | --- |
| entity tables | `sys_user`, `md_organization` |
| relation tables | `sys_user_role`, `md_department_staff` |
| log tables | `sys_login_log`, `sys_audit_log` |
| rule or definition tables | `sys_config_definition`, `sys_sequence_rule` |

### 2.3 View Naming

```
vw_{module_prefix}_{view_name}
mv_{module_prefix}_{view_name}
```

| Type | Examples |
| --- | --- |
| normal views | `vw_sys_user_detail`, `vw_md_staff_summary` |
| aggregate views | `vw_hr_attendance_stat` |
| materialized views | `mv_sys_daily_login_stat`, `mv_hr_monthly_attendance` |

### 2.4 Column Naming

| Scenario | Naming rule | Examples |
| --- | --- | --- |
| primary key | `id` | `id` |
| foreign key reference | `{referenced_entity}_id` | `role_id`, `organization_id`, `app_id` |
| self-reference in tree structures | `parent_id` | `parent_id` |
| boolean field | `is_{adjective_or_state}` | `is_active`, `is_locked`, `is_default` |
| timestamp field | `{action}_at` | `created_at`, `updated_at`, `password_updated_at` |
| sort field | `sort_order` | `sort_order` |
| remark field | `remark` | `remark` |
| metadata field | `meta` | `meta` |

### 2.5 Constraint Naming

All constraints must be explicitly named. Use the format `{constraint_prefix}_{table_name}__{column_name_or_semantics}`. The table name and the column or semantic part must be separated by a double underscore `__`.

| Constraint type | Prefix | Naming format | Examples |
| --- | --- | --- | --- |
| primary key | `pk` | `pk_{table_name}` | `pk_sys_user` |
| unique key | `uk` | `uk_{table_name}__{column_name}` | `uk_sys_user__username`, `uk_sys_user__staff_id` |
| foreign key | `fk` | `fk_{table_name}__{column_name}` | `fk_sys_user__created_by` |
| check constraint | `ck` | `ck_{table_name}__{column_name}` | `ck_sys_user__gender` |

Composite constraints:

- when a constraint spans multiple columns, concatenate the full column names with a single underscore by default
- if the full name would exceed PostgreSQL's identifier length limit of 63 bytes, a semantic abbreviation is allowed

```sql
-- Default: full column names
CONSTRAINT uk_sys_dictionary_item__dictionary_id_code UNIQUE (dictionary_id, code)
CONSTRAINT uk_sys_user_role__user_id_role_id UNIQUE (user_id, role_id)

-- Only when the full name becomes too long: semantic abbreviation
CONSTRAINT uk_md_department__org_code UNIQUE (organization_id, code)
```

### 2.6 Index Naming

Index naming must distinguish index types. Use `idx_` for default B-tree indexes and a type-specific prefix for other index types.

| Index type | Prefix | Naming format | Examples |
| --- | --- | --- | --- |
| B-tree (default) | `idx` | `idx_{table_name}__{column_name}` | `idx_sys_role__created_by` |
| GIN | `gin` | `gin_{table_name}__{column_name}` | `gin_md_medical_code__meta` |
| GiST | `gist` | `gist_{table_name}__{column_name}` | `gist_md_organization__location` |
| BRIN | `brin` | `brin_{table_name}__{column_name}` | `brin_sys_audit_log__created_at` |

For composite indexes, join multiple column names with a single underscore. If the full name is too long, a semantic abbreviation is allowed using the same rule as composite constraints.

```sql
CREATE INDEX idx_sys_audit_log__api_resource_api_action_api_version ON sys_audit_log (api_resource, api_action, api_version);

-- Use a semantic abbreviation only when the full name is too long
CREATE INDEX idx_sys_audit_log__api ON sys_audit_log (api_resource, api_action, api_version);
```

Partial indexes:

- append the `__partial` suffix
- explain the filter condition in a comment

```sql
-- Index only active users
CREATE INDEX idx_sys_user__email__partial ON sys_user (email) WHERE is_active = TRUE;
```

Covering indexes using `INCLUDE`:

- append the `__include` suffix

```sql
-- Covering index to avoid heap lookup
CREATE INDEX idx_sys_user__username__include ON sys_user (username) INCLUDE (name, email);
```

## 3. Data Type Rules

### 3.1 Standard Type Mapping

| Usage | Type | Description |
| --- | --- | --- |
| primary key or foreign key | `VARCHAR(32)` | stores application identifiers; compact `xid` values are used by default, while 32 characters are reserved for future expansion and for system-reserved IDs such as `system`, `anonymous`, and `cron_job` |
| person name | `VARCHAR(16)` | short human names |
| entity name | `VARCHAR(32)` | role name, department name, and similar entity names |
| long name | `VARCHAR(64)` | dictionary item names, code names, and similar values |
| business code | `VARCHAR(32)` | all kinds of business codes |
| email | `VARCHAR(128)` | ordinary email addresses |
| URL or file path | `VARCHAR(512)` | avatar URLs, callback URLs, file paths, and similar values |
| profile or long description | `VARCHAR(2048)` | organization profile, failure reason, and similar long descriptions |
| remark | `VARCHAR(512)` | unified length for remark fields |
| boolean | `BOOLEAN` | pair with `NOT NULL DEFAULT FALSE` |
| sort order | `INTEGER` | pair with `NOT NULL DEFAULT 0` |
| small-range integer | `SMALLINT` | for step values, widths, and similar fields |
| timestamp | `TIMESTAMP` | no time zone; use `LOCALTIMESTAMP` |
| date | `DATE` | date-only scenarios such as birth dates |
| dictionary value or enum | `VARCHAR(8)` | unified length; use with a `CHECK` constraint or dictionary table |
| metadata | `JSONB` | extensible structured metadata |
| large text | `TEXT` | unrestricted text such as configuration values |

### 3.2 Dictionary and Enum Handling

PostgreSQL `ENUM` types must not be used. Dictionary values and enum fields must always use `VARCHAR(8)`, together with a `CHECK` constraint or a dictionary-table constraint. Do not size the `VARCHAR` length according to the individual values.

```sql
-- All dictionary or enum fields use VARCHAR(8)
gender    VARCHAR(8) NOT NULL DEFAULT 'U',
CONSTRAINT ck_sys_user__gender CHECK (gender = ANY (ARRAY['M', 'F', 'U']))
```

```sql
-- Longer enum values still use VARCHAR(8)
overflow_strategy  VARCHAR(8) NOT NULL DEFAULT 'error',
CONSTRAINT ck_sys_sequence_rule__overflow_strategy CHECK (overflow_strategy = ANY (ARRAY['error', 'reset', 'extend']))
```

## 4. Table Structure Templates

### 4.1 Standard Entity Table

Use this template for business entities that require full create, read, update, and delete support, such as users, roles, organizations, and departments:

```sql
CREATE TABLE {module_prefix}_{entity_name} (
    -- 1. Primary key
    id                       VARCHAR(32) NOT NULL,

    -- 2. Audit fields (fixed 4 columns; order must not change)
    created_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    created_by               VARCHAR(32) NOT NULL DEFAULT 'system',
    updated_by               VARCHAR(32) NOT NULL DEFAULT 'system',

    -- 3. Business fields
    ...

    -- 4. Common optional fields (use as needed; order: is_active -> sort_order -> remark -> meta)
    is_active                BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order               INTEGER NOT NULL DEFAULT 0,
    remark                   VARCHAR(512),
    meta                     JSONB,

    -- 5. Constraints (order: PK -> UK -> CK -> FK)
    CONSTRAINT pk_{table_name} PRIMARY KEY (id),
    CONSTRAINT uk_{table_name}__{column_name} UNIQUE (...),
    CONSTRAINT ck_{table_name}__{column_name} CHECK (...),
    CONSTRAINT fk_{table_name}__created_by FOREIGN KEY (created_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_{table_name}__updated_by FOREIGN KEY (updated_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
```

### 4.2 Relation Table (Many-to-Many)

Use this template for many-to-many relationships such as user-role, department-staff, and role-permission tables. These tables must not include `updated_at` or `updated_by`.

```sql
CREATE TABLE {module_prefix}_{entity_a}_{entity_b} (
    -- 1. Primary key
    id                       VARCHAR(32) NOT NULL,

    -- 2. Audit fields (2 columns only)
    created_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    created_by               VARCHAR(32) NOT NULL DEFAULT 'system',

    -- 3. Relation foreign keys
    {entity_a}_id            VARCHAR(32) NOT NULL,
    {entity_b}_id            VARCHAR(32) NOT NULL,

    -- 4. Additional business fields, if any
    ...

    -- 5. Constraints
    CONSTRAINT pk_{table_name} PRIMARY KEY (id),
    CONSTRAINT uk_{table_name}__{entity_a}_id_{entity_b}_id UNIQUE ({entity_a}_id, {entity_b}_id),
    CONSTRAINT fk_{table_name}__{entity_a}_id FOREIGN KEY ({entity_a}_id) REFERENCES {entity_a_table}(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_{table_name}__{entity_b}_id FOREIGN KEY ({entity_b}_id) REFERENCES {entity_b_table}(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_{table_name}__created_by FOREIGN KEY (created_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
```

### 4.3 Log Table (Append-Only)

Use this template for immutable records such as login logs and audit logs. These tables must not include `updated_at` or `updated_by`.

```sql
CREATE TABLE {module_prefix}_{log_name}_log (
    -- 1. Primary key
    id                       VARCHAR(32) NOT NULL,

    -- 2. Audit fields (2 columns only)
    created_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    created_by               VARCHAR(32) NOT NULL DEFAULT 'system',

    -- 3. Business fields
    ...

    -- 4. Constraints
    CONSTRAINT pk_{table_name} PRIMARY KEY (id),
    CONSTRAINT fk_{table_name}__created_by FOREIGN KEY (created_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
```

## 5. Foreign Key Rules

### 5.1 `ON DELETE` / `ON UPDATE` Strategy

| Scenario | ON DELETE | ON UPDATE |
| --- | --- | --- |
| `created_by` or `updated_by` -> `sys_user` | `RESTRICT` | `CASCADE` |
| self-reference through `parent_id` | `RESTRICT` | `CASCADE` |
| strong dependency between business entities | `RESTRICT` | `CASCADE` |
| relation table referencing a main entity | `CASCADE` | `CASCADE` |
| configuration key referencing a definition table | `RESTRICT` | `CASCADE` |

### 5.2 Audit Foreign Keys

All `created_by` and `updated_by` columns, when present, must reference `sys_user(id)`:

```sql
CONSTRAINT fk_{table_name}__created_by FOREIGN KEY (created_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE,
CONSTRAINT fk_{table_name}__updated_by FOREIGN KEY (updated_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE
```

## 6. Index Rules

### 6.1 Automatically Created Indexes

PostgreSQL automatically creates indexes for `PRIMARY KEY` and `UNIQUE` constraints. Those indexes must not be created manually again.

### 6.2 Common B-tree Indexes

| Column | Sort order | Description |
| --- | --- | --- |
| `created_at` | `DESC` | add when reverse chronological lookup is needed |
| `created_by` | default ascending | add only when filtering by creator or "my own" records is a stable query pattern and the table is expected to grow large enough for the index to matter |
| non-unique foreign keys | default ascending | for example `parent_id`, `app_id`, `organization_id` |
| `sort_order` | default ascending | add when the table actually uses explicit sort order |

### 6.3 Special Indexes

| Scenario | Type | Example |
| --- | --- | --- |
| JSONB query | GIN index | `CREATE INDEX gin_md_medical_code__meta ON md_medical_code USING gin (meta);` |
| composite condition query | composite B-tree index | `CREATE INDEX idx_sys_audit_log__api ON sys_audit_log (api_resource, api_action, api_version);` |
| tree path query | B-tree index | `CREATE INDEX idx_md_medical_code_category__tree_path ON md_medical_code_category (tree_path);` |

### 6.4 Partial Indexes

Use partial indexes when queries frequently target a fixed filter condition and the smaller index materially improves performance:

```sql
-- Index only active user emails
CREATE INDEX idx_sys_user__email__partial ON sys_user (email) WHERE is_active = TRUE;
```

Typical scenarios:

- status filters such as `WHERE is_active = TRUE`
- indexing only non-null values such as `WHERE deleted_reason IS NOT NULL`

### 6.5 Covering Indexes With `INCLUDE`

Use `INCLUDE` when a query can be fully covered by index data and you want to avoid heap lookup:

```sql
-- Cover name and email for username lookups to avoid heap access
CREATE INDEX idx_sys_user__username__include ON sys_user (username) INCLUDE (name, email);
```

Note:

- columns listed in `INCLUDE` do not participate in search or ordering
- they are stored only as extra payload for index-only access

## 7. COMMENT Rules

### 7.1 Requirements

- every table must have a `COMMENT ON TABLE`
- every column must have a `COMMENT ON COLUMN`
- this English page renders SQL comments and `COMMENT ON` text in English for readability
- `COMMENT` statements must appear immediately after `CREATE TABLE` and before index creation

### 7.2 Format

```sql
CREATE TABLE sys_user (
    ...
);
COMMENT ON TABLE sys_user IS 'User';
COMMENT ON COLUMN sys_user.id IS 'Primary key';
COMMENT ON COLUMN sys_user.created_at IS 'Created at';
COMMENT ON COLUMN sys_user.updated_at IS 'Updated at';
COMMENT ON COLUMN sys_user.created_by IS 'Created by';
COMMENT ON COLUMN sys_user.updated_by IS 'Updated by';
-- Business column comments ...
CREATE INDEX ...
```

### 7.3 Fixed Audit Field Comments

| Column | Comment |
| --- | --- |
| `id` | `Primary key` |
| `created_at` | `Created at` |
| `updated_at` | `Updated at` |
| `created_by` | `Created by` |
| `updated_by` | `Updated by` |
| `remark` | `Remark` |
| `meta` | `Metadata` |
| `sort_order` | `Sort order` |
| `is_active` | `Is active` |

### 7.4 Foreign Key Column Comments

Foreign key column comments must use the `{referenced_entity_name}ID` pattern. They must not be described as "primary key".

| Column | Correct comment | Incorrect comment |
| --- | --- | --- |
| `user_id` | `User ID` | ~~User primary key~~ |
| `role_id` | `Role ID` | ~~Role primary key~~ |
| `organization_id` | `Organization ID` | ~~Organization primary key~~ |
| `parent_id` | `Parent ID` | ~~Parent primary key~~ |

## 8. Column Definition Formatting

### 8.1 Alignment

Column names and data types must be aligned to keep the same table visually consistent:

```sql
CREATE TABLE sys_user (
    id                       VARCHAR(32) NOT NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    username                 VARCHAR(32) NOT NULL,
    is_active                BOOLEAN NOT NULL DEFAULT FALSE,
    remark                   VARCHAR(512),
    meta                     JSONB,
    ...
);
```

Formatting rules:

- column names are left-aligned
- data types start at the same visual column, with a suggested name area width of 25 characters
- modifiers such as `NOT NULL` and `DEFAULT` follow the type directly

### 8.2 Column Order

Columns inside a table must follow this order:

1. `id`
2. `created_at`, `updated_at`
3. `created_by`, `updated_by`
4. foreign key columns such as `parent_id` and `{entity}_id`
5. core business fields
6. status or switch fields such as `is_active`, `is_locked`, and `status`
7. `sort_order` when present
8. `remark` when present
9. `meta` when present

### 8.3 Constraint Order

`CONSTRAINT` clauses must follow this order:

1. `PRIMARY KEY`
2. `UNIQUE`
3. `CHECK`
4. `FOREIGN KEY`

Among foreign keys, `created_by` and `updated_by` must be placed last.

## 9. Default Value Rules

| Type or usage | Default value |
| --- | --- |
| `created_at` | `LOCALTIMESTAMP` |
| `updated_at` | `LOCALTIMESTAMP` |
| `created_by` | `'system'` |
| `updated_by` | `'system'` |
| boolean fields such as `is_active` | `FALSE` |
| `sort_order` | `0` |
| `gender` | `'U'` |
| business booleans that need to default to enabled | `TRUE`, but the reason must be explicitly documented |

### 9.1 Audit Field Update Mechanism

- `created_at` and `created_by` are filled by database defaults on insert and must not be updated later
- `updated_at` and `updated_by` must be assigned by the application on every update
- database triggers must not be used for this responsibility

### 9.2 System-Reserved User Initialization Prerequisite

The default values `'system'` for `created_by` and `updated_by` reference `sys_user(id)` through foreign keys. Reserved user IDs such as `system`, `anonymous`, and `cron_job` use the same primary-key type as ordinary users. Therefore, database initialization must insert a system user row with `id = 'system'` first. Otherwise, inserts that depend on these defaults will fail on the foreign key check.

```sql
INSERT INTO sys_user (id, username, name) VALUES ('system', 'system', 'System');
```

If your audit flow writes `anonymous`, `cron_job`, or other reserved actors, preload those user rows as well.

## 10. `ALTER TABLE` Rules

When cross-module circular dependencies exist, use `ALTER TABLE` to add foreign keys after the referenced table has been created:

```sql
-- Add the sys_user.staff_id foreign key after md_staff has been created
ALTER TABLE sys_user
    ADD CONSTRAINT fk_sys_user__staff_id FOREIGN KEY (staff_id) REFERENCES md_staff(id) ON DELETE RESTRICT ON UPDATE CASCADE;
```

Such statements must be placed at the end of the DDL file for that dependency chain, before `COMMIT;`.

## 11. Complete Examples

The following example shows a standard entity table:

```sql
CREATE TABLE sys_role (
    id                       VARCHAR(32) NOT NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP,
    created_by               VARCHAR(32) NOT NULL DEFAULT 'system',
    updated_by               VARCHAR(32) NOT NULL DEFAULT 'system',
    name                     VARCHAR(32) NOT NULL,
    is_active                BOOLEAN NOT NULL DEFAULT FALSE,
    remark                   VARCHAR(512),
    meta                     JSONB,

    CONSTRAINT pk_sys_role PRIMARY KEY (id),
    CONSTRAINT fk_sys_role__created_by FOREIGN KEY (created_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_sys_role__updated_by FOREIGN KEY (updated_by) REFERENCES sys_user(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
COMMENT ON TABLE sys_role IS 'Role';
COMMENT ON COLUMN sys_role.id IS 'Primary key';
COMMENT ON COLUMN sys_role.created_at IS 'Created at';
COMMENT ON COLUMN sys_role.updated_at IS 'Updated at';
COMMENT ON COLUMN sys_role.created_by IS 'Created by';
COMMENT ON COLUMN sys_role.updated_by IS 'Updated by';
COMMENT ON COLUMN sys_role.name IS 'Role name';
COMMENT ON COLUMN sys_role.is_active IS 'Is active';
COMMENT ON COLUMN sys_role.remark IS 'Remark';
COMMENT ON COLUMN sys_role.meta IS 'Metadata';
CREATE INDEX idx_sys_role__created_at ON sys_role (created_at DESC);
-- Add the created_by index only when creator-based filtering becomes a stable query pattern on a large table
-- CREATE INDEX idx_sys_role__created_by ON sys_role (created_by);
```

The following example shows a relation table:

```sql
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
COMMENT ON TABLE sys_user_role IS 'User-role relation';
COMMENT ON COLUMN sys_user_role.id IS 'Primary key';
COMMENT ON COLUMN sys_user_role.created_at IS 'Created at';
COMMENT ON COLUMN sys_user_role.created_by IS 'Created by';
COMMENT ON COLUMN sys_user_role.user_id IS 'User ID';
COMMENT ON COLUMN sys_user_role.role_id IS 'Role ID';
COMMENT ON COLUMN sys_user_role.effective_from IS 'Effective from';
COMMENT ON COLUMN sys_user_role.effective_to IS 'Effective to';
-- Add the created_by index only when querying records created by a specific user becomes common and row count is expected to grow
-- CREATE INDEX idx_sys_user_role__created_by ON sys_user_role (created_by);
CREATE INDEX idx_sys_user_role__role_id ON sys_user_role (role_id);
```

## See also

- [Application Project Conventions](./application-project-conventions) for Go code and API naming rules
- [Models](../guide/models) for the interaction between Go field names, JSON tags, and database columns
