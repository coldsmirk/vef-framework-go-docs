---
sidebar_position: 8
---

# ORM: DDL & Surface Map

Schema DDL builders and the complete public surface of the `orm` package.

## DDL Operations

### Create Table

```go
_, err := db.NewCreateTable().
	Model((*User)(nil)).
	IfNotExists().
	Exec(ctx)
```

### Create Index

```go
_, err := db.NewCreateIndex().
	Model((*User)(nil)).
	Index("idx_user_email").
	Column("email").
	Unique().
	IfNotExists().
	Exec(ctx)
```

### Other DDL

```go
db.NewDropTable().Model((*User)(nil)).IfExists().Exec(ctx)
db.NewTruncateTable().Model((*User)(nil)).Exec(ctx)
db.NewAddColumn().Model((*User)(nil)).
    Column("phone", orm.DataType.VarChar(20)).
    IfNotExists().
    Exec(ctx)
db.NewDropColumn().Model((*User)(nil)).Column("phone").Exec(ctx)
```

## Public ORM Surface Map

The `orm` package intentionally re-exports the lower-level query-builder
contracts used by application code:

The shared `QueryBuilder` contract exposes `QueryBuilder.DB()` to return the
VEF `orm.DB` that created the query, alongside dialect, table metadata,
expression-builder, subquery, and condition-builder helpers.

| API group | Public surface |
| --- | --- |
| database entry points | `DB`, `Tx`, `Executor`, `RawQuery`, `SelectQuery`, `InsertQuery`, `UpdateQuery`, `DeleteQuery`, `MergeQuery`, `CreateTableQuery`, `DropTableQuery`, `CreateIndexQuery`, `DropIndexQuery`, `TruncateTableQuery`, `AddColumnQuery`, `DropColumnQuery`, `TableTarget`, `QueryBuilder` |
| Bun/schema aliases | `BunSelectQuery`, `BunInsertQuery`, `BunUpdateQuery`, `BunDeleteQuery`, `Table`, `Field`, `Relation`, `Dialect` |
| model bases | `BaseModel`, `Model`, `CreationAuditedModel`, `FullAuditedModel`, `CreationTrackedModel`, `FullTrackedModel` |
| hooks | `BeforeSelectHook`, `AfterSelectHook`, `BeforeInsertHook`, `AfterInsertHook`, `BeforeUpdateHook`, `AfterUpdateHook`, `BeforeDeleteHook`, `AfterDeleteHook`, `BeforeScanRowHook`, `AfterScanRowHook` |
| core builders | `ConditionBuilder`, `ExprBuilder`, `OrderBuilder`, `CaseBuilder`, `CaseWhenBuilder`, `ConflictBuilder`, `ConflictAction`, `ConflictUpdateBuilder`, `MergeWhenBuilder`, `MergeUpdateBuilder`, `MergeInsertBuilder`, `RelationSpec` |
| aggregate builders | `CountBuilder`, `SumBuilder`, `AvgBuilder`, `MinBuilder`, `MaxBuilder`, `StringAggBuilder`, `ArrayAggBuilder`, `StdDevBuilder`, `VarianceBuilder`, `JSONObjectAggBuilder`, `JSONArrayAggBuilder`, `BitOrBuilder`, `BitAndBuilder`, `BoolOrBuilder`, `BoolAndBuilder` |
| window builders | `WindowCountBuilder`, `WindowSumBuilder`, `WindowAvgBuilder`, `WindowMinBuilder`, `WindowMaxBuilder`, `WindowStringAggBuilder`, `WindowArrayAggBuilder`, `WindowStdDevBuilder`, `WindowVarianceBuilder`, `WindowJSONObjectAggBuilder`, `WindowJSONArrayAggBuilder`, `WindowBitOrBuilder`, `WindowBitAndBuilder`, `WindowBoolOrBuilder`, `WindowBoolAndBuilder`, `RowNumberBuilder`, `RankBuilder`, `DenseRankBuilder`, `PercentRankBuilder`, `CumeDistBuilder`, `NTileBuilder`, `LagBuilder`, `LeadBuilder`, `FirstValueBuilder`, `LastValueBuilder`, `NthValueBuilder` |
| DDL builders and types | `DataTypeDef`, `ColumnConstraint`, `RawDefault`, `PrimaryKeyBuilder`, `UniqueBuilder`, `CheckBuilder`, `ForeignKeyBuilder`, `ReferenceAction`, `IndexMethod`, `PartitionStrategy` |
| expression placeholders | `PlaceholderKeyOperator`, `ExprOperator`, `ExprTableColumns`, `ExprColumns`, `ExprTablePKs`, `ExprPKs`, `ExprTableName`, `ExprTableAlias` |
| audit constants | `ColumnID`, `ColumnCreatedAt`, `ColumnUpdatedAt`, `ColumnCreatedBy`, `ColumnUpdatedBy`, `ColumnCreatedByName`, `ColumnUpdatedByName`, `FieldID`, `FieldCreatedAt`, `FieldUpdatedAt`, `FieldCreatedBy`, `FieldUpdatedBy`, `FieldCreatedByName`, `FieldUpdatedByName`, `OperatorSystem`, `OperatorCronJob`, `OperatorAnonymous` |
| enum families | `JoinType`, `FuzzyKind`, `NullsMode`, `FromDirection`, `FrameType`, `FrameBoundKind`, `StatisticalMode`, `DateTimeUnit`, `JoinDefault`, `JoinInner`, `JoinLeft`, `JoinRight`, `JoinFull`, `JoinCross`, `FuzzyStarts`, `FuzzyEnds`, `FuzzyContains`, `NullsDefault`, `NullsRespect`, `NullsIgnore`, `FromDefault`, `FromFirst`, `FromLast`, `FrameDefault`, `FrameRows`, `FrameRange`, `FrameGroups`, `FrameBoundNone`, `FrameBoundUnboundedPreceding`, `FrameBoundUnboundedFollowing`, `FrameBoundCurrentRow`, `FrameBoundPreceding`, `FrameBoundFollowing`, `StatisticalDefault`, `StatisticalPopulation`, `StatisticalSample`, `ConflictDoNothing`, `ConflictDoUpdate`, `UnitYear`, `UnitMonth`, `UnitDay`, `UnitHour`, `UnitMinute`, `UnitSecond`, `ReferenceCascade`, `ReferenceRestrict`, `ReferenceSetNull`, `ReferenceSetDefault`, `ReferenceNoAction`, `IndexBTree`, `IndexHash`, `IndexGIN`, `IndexGiST`, `IndexSPGiST`, `IndexBRIN`, `PartitionRange`, `PartitionList`, `PartitionHash` |
| helpers | `Applier`, `ApplyFunc`, `ApplySort`, `DataType`, `PrimaryKey`, `NotNull`, `Nullable`, `Default`, `Unique`, `Check`, `References`, `AutoIncrement`, `PKField`, `ColumnInfo`, `LabelsEqual` (v0.39), `ValidateLabels` (v0.39), `ErrInvalidLabel` (v0.39) |

Use `orm.RawDefault("CURRENT_TIMESTAMP")` only when a DDL default must render a
trusted SQL expression verbatim. Plain `orm.Default(value)` keeps ordinary
strings, booleans, and numeric values on the safe literal/bound-value path.

## Next Step

- [ORM: Querying](./orm-querying) — the main query-building reference
- [Models](./models) — how model definitions map to tables
