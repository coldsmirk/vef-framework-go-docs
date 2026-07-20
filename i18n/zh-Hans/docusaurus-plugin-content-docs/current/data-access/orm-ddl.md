---
sidebar_position: 8
---

# ORM：DDL 与 Surface Map

模式 DDL 构造器与 `orm` 包的完整公开接口面。

## DDL 操作

### 创建表

```go
_, err := db.NewCreateTable().
	Model((*User)(nil)).
	IfNotExists().
	Exec(ctx)
```

### 创建索引

```go
_, err := db.NewCreateIndex().
	Model((*User)(nil)).
	Index("idx_user_email").
	Column("email").
	Unique().
	IfNotExists().
	Exec(ctx)
```

### 其他 DDL

```go
db.NewDropTable().Model((*User)(nil)).IfExists().Exec(ctx)
db.NewTruncateTable().Model((*User)(nil)).Exec(ctx)
db.NewAddColumn().Model((*User)(nil)).
    Column("phone", orm.DataType.VarChar(20)).
    IfNotExists().
    Exec(ctx)
db.NewDropColumn().Model((*User)(nil)).Column("phone").Exec(ctx)
```

## ORM 公开 Surface Map

`orm` 包有意重新导出了应用代码会用到的底层 query-builder 契约：

共享的 `QueryBuilder` contract 公开 `QueryBuilder.DB()`，用于返回创建当前
query 的 VEF `orm.DB`；它还包含 dialect、table metadata、expression-builder、
subquery 和 condition-builder helpers。

| API 组 | 公开 surface |
| --- | --- |
| 数据库入口 | `DB`, `Tx`, `Executor`, `RawQuery`, `SelectQuery`, `InsertQuery`, `UpdateQuery`, `DeleteQuery`, `MergeQuery`, `CreateTableQuery`, `DropTableQuery`, `CreateIndexQuery`, `DropIndexQuery`, `TruncateTableQuery`, `AddColumnQuery`, `DropColumnQuery`, `TableTarget`, `QueryBuilder` |
| Bun/schema alias | `BunSelectQuery`, `BunInsertQuery`, `BunUpdateQuery`, `BunDeleteQuery`, `Table`, `Field`, `Relation`, `Dialect` |
| 模型基类 | `BaseModel`, `Model`, `CreationAuditedModel`, `FullAuditedModel`, `CreationTrackedModel`, `FullTrackedModel` |
| hook | `BeforeSelectHook`, `AfterSelectHook`, `BeforeInsertHook`, `AfterInsertHook`, `BeforeUpdateHook`, `AfterUpdateHook`, `BeforeDeleteHook`, `AfterDeleteHook`, `BeforeScanRowHook`, `AfterScanRowHook` |
| 核心 builder | `ConditionBuilder`, `ExprBuilder`, `OrderBuilder`, `CaseBuilder`, `CaseWhenBuilder`, `ConflictBuilder`, `ConflictAction`, `ConflictUpdateBuilder`, `MergeWhenBuilder`, `MergeUpdateBuilder`, `MergeInsertBuilder`, `RelationSpec` |
| aggregate builder | `CountBuilder`, `SumBuilder`, `AvgBuilder`, `MinBuilder`, `MaxBuilder`, `StringAggBuilder`, `ArrayAggBuilder`, `StdDevBuilder`, `VarianceBuilder`, `JSONObjectAggBuilder`, `JSONArrayAggBuilder`, `BitOrBuilder`, `BitAndBuilder`, `BoolOrBuilder`, `BoolAndBuilder` |
| window builder | `WindowCountBuilder`, `WindowSumBuilder`, `WindowAvgBuilder`, `WindowMinBuilder`, `WindowMaxBuilder`, `WindowStringAggBuilder`, `WindowArrayAggBuilder`, `WindowStdDevBuilder`, `WindowVarianceBuilder`, `WindowJSONObjectAggBuilder`, `WindowJSONArrayAggBuilder`, `WindowBitOrBuilder`, `WindowBitAndBuilder`, `WindowBoolOrBuilder`, `WindowBoolAndBuilder`, `RowNumberBuilder`, `RankBuilder`, `DenseRankBuilder`, `PercentRankBuilder`, `CumeDistBuilder`, `NTileBuilder`, `LagBuilder`, `LeadBuilder`, `FirstValueBuilder`, `LastValueBuilder`, `NthValueBuilder` |
| DDL builder 与类型 | `DataTypeDef`, `ColumnConstraint`, `RawDefault`, `PrimaryKeyBuilder`, `UniqueBuilder`, `CheckBuilder`, `ForeignKeyBuilder`, `ReferenceAction`, `IndexMethod`, `PartitionStrategy` |
| 表达式占位符 | `PlaceholderKeyOperator`, `ExprOperator`, `ExprTableColumns`, `ExprColumns`, `ExprTablePKs`, `ExprPKs`, `ExprTableName`, `ExprTableAlias` |
| 审计常量 | `ColumnID`, `ColumnCreatedAt`, `ColumnUpdatedAt`, `ColumnCreatedBy`, `ColumnUpdatedBy`, `ColumnCreatedByName`, `ColumnUpdatedByName`, `FieldID`, `FieldCreatedAt`, `FieldUpdatedAt`, `FieldCreatedBy`, `FieldUpdatedBy`, `FieldCreatedByName`, `FieldUpdatedByName`, `OperatorSystem`, `OperatorCronJob`, `OperatorAnonymous` |
| enum family | `JoinType`, `FuzzyKind`, `NullsMode`, `FromDirection`, `FrameType`, `FrameBoundKind`, `StatisticalMode`, `DateTimeUnit`, `JoinDefault`, `JoinInner`, `JoinLeft`, `JoinRight`, `JoinFull`, `JoinCross`, `FuzzyStarts`, `FuzzyEnds`, `FuzzyContains`, `NullsDefault`, `NullsRespect`, `NullsIgnore`, `FromDefault`, `FromFirst`, `FromLast`, `FrameDefault`, `FrameRows`, `FrameRange`, `FrameGroups`, `FrameBoundNone`, `FrameBoundUnboundedPreceding`, `FrameBoundUnboundedFollowing`, `FrameBoundCurrentRow`, `FrameBoundPreceding`, `FrameBoundFollowing`, `StatisticalDefault`, `StatisticalPopulation`, `StatisticalSample`, `ConflictDoNothing`, `ConflictDoUpdate`, `UnitYear`, `UnitMonth`, `UnitDay`, `UnitHour`, `UnitMinute`, `UnitSecond`, `ReferenceCascade`, `ReferenceRestrict`, `ReferenceSetNull`, `ReferenceSetDefault`, `ReferenceNoAction`, `IndexBTree`, `IndexHash`, `IndexGIN`, `IndexGiST`, `IndexSPGiST`, `IndexBRIN`, `PartitionRange`, `PartitionList`, `PartitionHash` |
| helper | `Applier`, `ApplyFunc`, `ApplySort`, `DataType`, `PrimaryKey`, `NotNull`, `Nullable`, `Default`, `Unique`, `Check`, `References`, `AutoIncrement`, `PKField`, `ColumnInfo`, `LabelsEqual`, `ValidateLabels`, `ErrInvalidLabel` |

只有 DDL default 必须原样渲染可信 SQL 表达式时，才使用
`orm.RawDefault("CURRENT_TIMESTAMP")`。普通 `orm.Default(value)` 仍让字符串、
布尔值和数字走安全的 literal/bound-value 路径。

## 下一步

- [ORM：查询](./orm-querying) — 查询构造的主参考
- [模型](./models) — 模型定义如何映射到数据表
