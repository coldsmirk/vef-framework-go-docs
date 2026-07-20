---
sidebar_position: 4
---

# Code Maps

Contracts standardize *structures*; code maps standardize *values*. A code
map translates the values of one code set (gender, marital status, order
status, ...) between the host's canonical codes and one external system's
codes — the code-level instance of the canonical data model pattern.

## The Model

One `integration.CodeMap` row (`itg_code_map`) binds a system to a code set:

| Field | Meaning |
| --- | --- |
| `systemId` | the owning system |
| `codeSet` | identifier of the translated code set (e.g. `"gender"`). Where the host catalog exposes the same set (mold translate tags, `mold.CodeSetInspector`), the identifiers should agree |
| `name` | display name |
| `entries` | the mapping pairs (below) |
| `onUnmapped` | default behavior for lookups no entry matches: `reject` (default, fail closed), `passthrough`, or `fallback` |
| `fallbackCanonical` / `fallbackExternal` | the value each direction yields for unmapped input under the `fallback` policy |
| `isEnabled` | disabled maps behave as missing |

Each `CodeMapEntry` is one bidirectional pair:

```json
{
  "canonical": "F",
  "external": 2,
  "canonicalAliases": ["female"],
  "externalAliases": ["02"]
}
```

- Lookups match the primary or any alias; translations always emit the
  opposite side's **primary** — aliases are matched, never emitted.
- Values keep their JSON type (string, number, boolean) end to end; lookups
  compare by normalized string form, so `1` and `"1"` address the same entry.
- Save-time validation rejects duplicate lookup values per side across
  primaries and aliases (`ErrInvalidCodeMap`), so every lookup is
  deterministic.

## The `codes` Script Library

Adapter scripts (both directions) consume code maps through the global
`codes` object, scoped to the executing system:

```js
codes.toExternal('gender', input.gender)          // canonical → external
codes.toCanonical('gender', body.sex)             // external → canonical
codes.toExternal('gender', v, { fallback: 'U' })  // per-call unmapped override
codes.toCanonical('status', v, { passthrough: true })
codes.toCanonical('status', v, { reject: true })
codes.entries('gender')                           // raw mapping pairs
```

- `null` and `undefined` pass through untranslated — translating absence is
  not a lookup.
- The per-call options object carries exactly one of `fallback`,
  `passthrough: true`, or `reject: true`, and overrides the map's stored
  policy for that call. A malformed options object fails on every call, not
  only on the first unmapped value.
- A lookup against a code set the system has no enabled map for throws
  `ErrMissingCodeMap` (classified as `config`); an unmapped value under the
  reject policy throws `ErrUnmappedValue`.
- One run sees one snapshot (loaded maps are memoized per execution);
  saves are live on the next invocation. Compiled lookup indexes are shared
  across executions by content hash.

## The Host Code Set Catalog

The mapping editor benefits from knowing which code sets the host defines
and which values each set contains. When the host's
[mold code set registration](../data-tools/mold) also implements
`mold.CodeSetInspector`, the `integration/code_set` resource exposes the
catalog:

```go
type CodeSetInspector interface {
    ListCodeSets(ctx context.Context) ([]mold.CodeSetInfo, error)
    ListCodes(ctx context.Context, codeSet string) ([]mold.CodeInfo, error)
}
```

- The inspector is asserted from the registered `mold.CodeSetLoader` first
  (the common path), then from the `mold.CodeSetResolver` (hosts that
  replace it wholesale).
- With an inspector present, saving a code map also validates that its
  `codeSet` identifier is registered by the host catalog.
- Without one, the catalog operations report `supported: false` and the
  editor degrades to free-text input.

Field-by-field request/response documentation for `integration/code_map` and
`integration/code_set` lives in [RPC Resources](./resources).

## Next Step

[RPC Resources](./resources) — the complete management API reference.
