# databasex AGENTS.md

## OVERVIEW

Database utility cluster: pagination, SQL injection detection, and unsafe SQL string building.

## STRUCTURE

```
databasex/
├── pagex/         # Pagination: Page struct, Option pattern, protobuf bridge
├── sqls/          # SQL injection detection via regex pattern matching
└── unsafesql/     # SQL query string builder (fluent API, Must* panic variants)
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add pagination field | `pagex/page.go` | Add to struct + getter + Option func + proto |
| Change SQL injection rules | `sqls/sql_injection.go` | Two regexes: syntax + comment patterns |
| Add SQL operator | `unsafesql/sql.go` | Add safe + Must* pair; no abstraction exists |
| Store page in context | `pagex/context.go` | `NewContext` / `FromContext` with private key |

## CONVENTIONS

- **Option pattern**: `Option func(p *Page)` in pagex; functional options for Page construction
- **Dual API**: Every unsafesql operator has safe (skip empty) + `Must*` (panic) variant
- **Protobuf bridge**: `Page.AsProto()` / `FromProto()` with `timestamppb` conversion
- **sync.Once guard**: `Page.SetTotal` uses `totalOnce` to prevent double-write
- **Context storage**: Private `key struct{}` for context value storage
- **Cross-package dep**: `sqls` imports `stringx.IsBlank` for blank check

## ANTI-PATTERNS

- `unsafesql/sql.go` repeats same pattern 30+ times with zero abstraction
- `sqls/context.go` is entirely commented-out dead code
- `unsafesql` is truly unsafe: string concatenation, no parameterization
- `pagex` pageNum is 1-based; offset calculated as `(pageNum-1)*pageSize`
- `CheckSqlInjection` regex may have false positives/negatives
