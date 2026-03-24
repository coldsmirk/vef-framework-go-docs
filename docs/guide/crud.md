---
sidebar_position: 2
---

# Generic CRUD

The `crud` package is one of the most important user-facing layers in VEF. It turns typed models and typed request structs into reusable API operations with built-in transactions, validation, data permissions, file promotion, and result formatting.

## The Basic Pattern

You usually embed CRUD providers into a resource struct:

```go
type UserResource struct {
	api.Resource

	crud.FindPage[User, UserSearch]
	crud.Create[User, UserParams]
	crud.Update[User, UserParams]
	crud.Delete[User]
}

func NewUserResource() api.Resource {
	return &UserResource{
		Resource: api.NewRPCResource("sys/user"),
		FindPage: crud.NewFindPage[User, UserSearch]().PermToken("sys:user:query"),
		Create:   crud.NewCreate[User, UserParams]().PermToken("sys:user:create"),
		Update:   crud.NewUpdate[User, UserParams]().PermToken("sys:user:update"),
		Delete:   crud.NewDelete[User]().PermToken("sys:user:delete"),
	}
}
```

The framework collects embedded CRUD builders automatically because they implement `api.OperationsProvider`.

## Generic Parameter Meanings

Most CRUD builders only use one of these generic shapes:

| Generic | Meaning | Typical type |
| --- | --- | --- |
| `TModel` | persistence model loaded from or written to the database | `User`, `Role`, `Flow` |
| `TParams` | write-side params decoded from `Request.Params` | `UserParams`, `CreateUserParams` |
| `TSearch` | read-side search params decoded from `Request.Params` | `UserSearch`, `RoleSearch` |

Operation families use them like this:

| Builder family | Generic shape | Meaning |
| --- | --- | --- |
| single-record write builders | `Create[TModel, TParams]`, `Update[TModel, TParams]` | params are copied into a model before persistence |
| batch write builders | `CreateMany[TModel, TParams]`, `UpdateMany[TModel, TParams]` | framework wraps `TParams` into batch params types |
| read builders | `FindOne[TModel, TSearch]`, `FindPage[TModel, TSearch]`, and similar | model defines the query target, search defines filters |
| delete builders | `Delete[TModel]`, `DeleteMany[TModel]` | deletion works from primary-key payloads, so no extra `TParams` type is needed |
| export builder | `Export[TModel, TSearch]` | export runs a read query and then renders the result into a file |
| import builder | `Import[TModel]` | imported rows are decoded directly into models |

## Prebuilt Builder Matrix

| Builder | Default RPC action | Default REST action | Input contract | Output contract | Typical use |
| --- | --- | --- | --- | --- | --- |
| `NewCreate[TModel, TParams]` | `create` | `post /` | `TParams` from `params` | primary-key map | create one record |
| `NewUpdate[TModel, TParams]` | `update` | `put /:id` | `TParams` from `params`, including PK fields | success result | update one record |
| `NewDelete[TModel]` | `delete` | `delete /:id` | raw PK values from `params` | success result | delete one record |
| `NewCreateMany[TModel, TParams]` | `create_many` | `post /many` | `CreateManyParams[TParams]` with `list` | list of primary-key maps | batch create |
| `NewUpdateMany[TModel, TParams]` | `update_many` | `put /many` | `UpdateManyParams[TParams]` with `list` | success result | batch update |
| `NewDeleteMany[TModel]` | `delete_many` | `delete /many` | `DeleteManyParams` with `pks` | success result | batch delete |
| `NewFindOne[TModel, TSearch]` | `find_one` | `get /:id` | `TSearch` from `params` | one model | single-record query |
| `NewFindAll[TModel, TSearch]` | `find_all` | `get /` | `TSearch` from `params` | `[]TModel` | filtered list without paging metadata |
| `NewFindPage[TModel, TSearch]` | `find_page` | `get /page` | `TSearch` from `params` + `page.Pageable` from `meta` | `page.Page[T]` | admin list screen |
| `NewFindOptions[TModel, TSearch]` | `find_options` | `get /options` | `TSearch` from `params` + `DataOptionConfig` from `meta` | `[]DataOption` | dropdown options |
| `NewFindTree[TModel, TSearch](treeBuilder)` | `find_tree` | `get /tree` | `TSearch` from `params` | hierarchical `[]TModel` | tree-structured data |
| `NewFindTreeOptions[TModel, TSearch]` | `find_tree_options` | `get /tree/options` | `TSearch` from `params` + `DataOptionConfig` from `meta` | `[]TreeDataOption` | tree options |
| `NewExport[TModel, TSearch]` | `export` | `get /export` | `TSearch` from `params` + export format from `meta` | file download | Excel or CSV export |
| `NewImport[TModel]` | `import` | `post /import` | multipart file upload + import format from `meta` | `{total: n}` | Excel or CSV import |

## Shared Builder Controls

Every CRUD builder inherits the common controls from `Builder[T]`:

| Method | Effect |
| --- | --- |
| `ResourceKind(kind)` | switches the builder between RPC and REST naming/validation rules |
| `Action(action)` | overrides the default action name |
| `Public()` | marks the operation as unauthenticated |
| `PermToken(token)` | requires a permission token for access |
| `Timeout(duration)` | sets the request timeout |
| `EnableAudit()` | enables audit logging for the operation |
| `RateLimit(max, period)` | applies per-operation rate limiting |

Important detail:

- `Action(...)` is validated according to the current `ResourceKind(...)`
- if you are overriding a REST action, set `ResourceKind(api.KindREST)` first

## Shared Find Controls

All read-oriented builders are built on top of `Find[...]`, so they share a richer set of query-shaping options:

| Method | Purpose |
| --- | --- |
| `WithProcessor(...)` | post-processes the query result before response serialization |
| `WithOptions(...)` | appends reusable low-level `FindOperationOption` values |
| `WithSelect(column)` | adds a column to the select list |
| `WithSelectAs(column, alias)` | adds a selected column with an alias |
| `WithDefaultSort(...)` | sets fallback sorting when no dynamic sort is provided |
| `WithCondition(...)` | adds a `WHERE` condition using `orm.ConditionBuilder` |
| `DisableDataPerm()` | disables automatic data-permission filtering |
| `WithRelation(...)` | adds relation joins through `orm.RelationSpec` |
| `WithAuditUserNames(userModel, nameColumn...)` | joins audit user information to populate creator/updater names |
| `WithQueryApplier(...)` | applies arbitrary query modifications with typed access to `TSearch` |

Runtime defaults for most find-style builders:

- search tags from `TSearch` are applied automatically
- data permission filtering is enabled by default
- sort defaults to primary key descending when the model has a single PK
- if no single PK exists, the fallback sort is `created_at DESC` when available

### Query Parts For Tree Builders

Tree builders use recursive CTEs, so some options can target different query stages:

| Query part | Meaning |
| --- | --- |
| `QueryRoot` | the final outer query |
| `QueryBase` | the starting query inside the recursive CTE |
| `QueryRecursive` | the recursive branch of the CTE |
| `QueryAll` | all query parts |

For `FindTree` and `FindTreeOptions`, several methods intentionally change their defaults:

- `WithCondition(...)` defaults to `QueryBase`
- `WithQueryApplier(...)` defaults to `QueryBase`
- `WithSelect(...)`, `WithSelectAs(...)`, and `WithRelation(...)` default to both `QueryBase` and `QueryRecursive`

## Read Builders

### `FindOne[TModel, TSearch]`

Use `FindOne` when the resource should return one record.

| Aspect | Details |
| --- | --- |
| Generics | `TModel` is the query target model, `TSearch` defines filters |
| Input | `TSearch` from `params`, raw `api.Meta` from `meta` |
| Output | one `TModel` value after optional `WithProcessor(...)` transformation |
| Default behavior | runs a select with model columns and `LIMIT 1` |
| Common configuration | shared find controls such as `WithCondition`, `WithRelation`, `WithQueryApplier`, `WithAuditUserNames` |

Use this when the read still behaves like a query instead of a fixed metadata fetch.

### `FindAll[TModel, TSearch]`

Use `FindAll` when you need a filtered list without paging metadata.

| Aspect | Details |
| --- | --- |
| Generics | `TModel` is the result model, `TSearch` defines filters |
| Input | `TSearch` from `params`, `api.Meta` from `meta` |
| Output | `[]TModel` or the processed slice returned by `WithProcessor(...)` |
| Default behavior | applies a safety limit (`maxQueryLimit`) and returns an empty slice instead of `nil` |
| Common configuration | shared find controls, especially `WithDefaultSort`, `WithCondition`, `WithRelation`, `WithQueryApplier` |

### `FindPage[TModel, TSearch]`

Use `FindPage` for most admin-style list screens.

| Aspect | Details |
| --- | --- |
| Generics | `TModel` is the item model, `TSearch` defines query filters |
| Input | `TSearch` from `params`, `page.Pageable` from `meta`, plus any extra `api.Meta` |
| Output | `page.Page[T]` |
| Default behavior | paginates, counts total rows, and normalizes page settings |
| Special configuration | `WithDefaultPageSize(size)` sets the fallback page size |

Use this when the caller needs `total`, page number, page size, and item list together.

### `FindOptions[TModel, TSearch]`

Use `FindOptions` for lightweight option lists such as select boxes.

| Aspect | Details |
| --- | --- |
| Generics | `TModel` is the source model, `TSearch` defines filters |
| Input | `TSearch` from `params`, `DataOptionConfig` from `meta` |
| Output | `[]DataOption` |
| Default behavior | maps data into `label`, `value`, `description`, and optional `meta` |
| Special configuration | `WithDefaultColumnMapping(mapping)` sets fallback label/value/description/meta column mapping |

`DataOptionConfig` comes from `meta` and can override:

| Field | Meaning |
| --- | --- |
| `labelColumn` | source column for `label` |
| `valueColumn` | source column for `value` |
| `descriptionColumn` | optional source column for `description` |
| `metaColumns` | additional columns to include in the option `meta` object |

Defaults:

- label column defaults to `name`
- value column defaults to `id`

### `FindTree[TModel, TSearch]`

Use `FindTree` when the domain is hierarchical and the response should contain nested model records.

Constructor shape:

```go
crud.NewFindTree[Category, CategorySearch](tree.Build)
```

| Aspect | Details |
| --- | --- |
| Generics | `TModel` is the tree node model, `TSearch` defines filters |
| Input | `TSearch` from `params`, `api.Meta` from `meta` |
| Output | hierarchical `[]TModel` |
| Default behavior | builds a recursive CTE, loads flat rows, then runs the provided `treeBuilder` function |
| Special configuration | `WithIDColumn(name)` and `WithParentIDColumn(name)` customize the tree columns |

Defaults:

- node ID column defaults to `id`
- parent ID column defaults to `parent_id`

### `FindTreeOptions[TModel, TSearch]`

Use `FindTreeOptions` when you need a hierarchical option tree instead of full model records.

| Aspect | Details |
| --- | --- |
| Generics | `TModel` is the source model, `TSearch` defines filters |
| Input | `TSearch` from `params`, `DataOptionConfig` from `meta` |
| Output | `[]TreeDataOption` |
| Default behavior | builds a recursive CTE and converts the result into nested `TreeDataOption` values |
| Special configuration | `WithDefaultColumnMapping(...)`, `WithIDColumn(...)`, `WithParentIDColumn(...)` |

Use this when the client needs `label`/`value` plus `children`, not the full persistence model.

## Write Builders

### `Create[TModel, TParams]`

Use `Create` for single-record creation.

| Aspect | Details |
| --- | --- |
| Generics | `TModel` is the persistence model, `TParams` is the write params type |
| Input | `TParams` from `params` |
| Output | primary-key map for the created record |
| Default behavior | copies params into a new model, promotes storage references, runs inside a transaction, inserts the record |
| Special configuration | `WithPreCreate(...)`, `WithPostCreate(...)` |

Hook responsibilities:

| Method | Runs when | Typical use |
| --- | --- | --- |
| `WithPreCreate` | before insert, inside the same transaction | normalization, validation, derived fields, extra query shaping |
| `WithPostCreate` | after insert, inside the same transaction | side effects that belong to the same transaction |

### `Update[TModel, TParams]`

Use `Update` for single-record update.

| Aspect | Details |
| --- | --- |
| Generics | `TModel` is the persistence model, `TParams` is the write params type |
| Input | `TParams` from `params`, including primary-key fields |
| Output | success result |
| Default behavior | copies params into a temporary model, validates PK presence, loads the old model, applies data permissions, merges non-empty fields, updates in a transaction |
| Special configuration | `WithPreUpdate(...)`, `WithPostUpdate(...)`, `DisableDataPerm()` |

Important detail:

- `Update` uses `copier.WithIgnoreEmpty()` when merging the incoming model into the loaded model

### `Delete[TModel]`

Use `Delete` for single-record deletion.

| Aspect | Details |
| --- | --- |
| Generics | `TModel` is the persistence model |
| Input | primary-key values from raw `api.Params` |
| Output | success result |
| Default behavior | validates PK input, loads the model, applies data permissions, deletes in a transaction, then cleans up promoted files |
| Special configuration | `WithPreDelete(...)`, `WithPostDelete(...)`, `DisableDataPerm()` |

### Batch Builders

#### `CreateMany[TModel, TParams]`

| Aspect | Details |
| --- | --- |
| Input contract | `CreateManyParams[TParams]` with a `list` field |
| Output | list of primary-key maps |
| Special configuration | `WithPreCreateMany(...)`, `WithPostCreateMany(...)` |
| Behavior | copies each params item into a model, inserts all models in one transaction |

#### `UpdateMany[TModel, TParams]`

| Aspect | Details |
| --- | --- |
| Input contract | `UpdateManyParams[TParams]` with a `list` field |
| Output | success result |
| Special configuration | `WithPreUpdateMany(...)`, `WithPostUpdateMany(...)`, `DisableDataPerm()` |
| Behavior | validates PKs for every item, loads all old models, merges updates, and executes a bulk update in one transaction |

#### `DeleteMany[TModel]`

| Aspect | Details |
| --- | --- |
| Input contract | `DeleteManyParams` with a `pks` field |
| Output | success result |
| Special configuration | `WithPreDeleteMany(...)`, `WithPostDeleteMany(...)`, `DisableDataPerm()` |
| Behavior | supports single-PK payloads as scalar values and composite-PK payloads as maps |

`DeleteManyParams.pks` rules:

| Model PK shape | Accepted payload shape |
| --- | --- |
| single primary key | `["id1", "id2"]` |
| composite primary key | `[{"user_id":"u1","role_id":"r1"}]` |

## Export And Import Builders

### `Export[TModel, TSearch]`

Use `Export` when the caller should download a query result as an Excel or CSV file.

| Aspect | Details |
| --- | --- |
| Generics | `TModel` is the exported row model, `TSearch` defines query filters |
| Input | `TSearch` from `params`, `format` from `meta` |
| Output | file download |
| Default behavior | runs a find-style query, applies optional pre-export processing, and writes Excel or CSV to the response |
| Special configuration | `WithDefaultFormat(...)`, `WithExcelOptions(...)`, `WithCsvOptions(...)`, `WithPreExport(...)`, `WithFilenameBuilder(...)` |

`format` values:

| Format | Value |
| --- | --- |
| Excel | `excel` |
| CSV | `csv` |

Defaults:

- export format defaults to `excel`
- default filenames are `data.xlsx` and `data.csv`

### `Import[TModel]`

Use `Import` when the caller uploads a CSV or Excel file that should be decoded into models and inserted.

| Aspect | Details |
| --- | --- |
| Generics | `TModel` is the model type imported from the file |
| Input | multipart file upload in `params.file`, plus optional `format` in `meta` |
| Output | `{total: n}` on success |
| Default behavior | requires multipart input, parses rows into models, validates imported rows, inserts them in a transaction |
| Special configuration | `WithDefaultFormat(...)`, `WithExcelOptions(...)`, `WithCsvOptions(...)`, `WithPreImport(...)`, `WithPostImport(...)` |

Important details:

- JSON requests are rejected for import
- if row-level import validation fails, the response contains an `errors` payload instead of partial persistence
- import format defaults to `excel`

## Practical Advice

- start with `FindPage + Create + Update + Delete` for admin resources
- keep write params and search params separate
- add permissions at the builder level
- rely on default data permissions unless you have a specific reason to disable them
- use `FindOptions` or `FindTreeOptions` for UI option payloads instead of overloading full model endpoints
- prefer the standard CRUD vocabulary unless your business action has a stronger domain verb
- it is normal for one resource to combine CRUD builders with a few custom actions when the UI needs both

## Next Step

Read [Custom Handlers](./custom-handlers) when a resource needs operations that do not fit the generic CRUD model.
