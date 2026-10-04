# databasex AGENTS.md

## OVERVIEW

Database utility cluster: MySQL data types, pagination, SQL injection detection, and unsafe SQL string building.

## STRUCTURE

```
databasex/
├── mysqlx/        # MySQL data types modeled after pgx/v5/pgtype (nullable types, Scanner/Valuer/JSON)
├── pagex/         # Pagination: Page struct, Option pattern, protobuf bridge
├── sqls/          # SQL injection detection via regex pattern matching
└── unsafesql/     # SQL query string builder (fluent API, Must* panic variants)
```

`mysqlx` requires Go 1.27: `mysqlx.UUID` is backed by the standard library
`uuid` package added in that release. It is the only package in the module with
a Go 1.27 floor.

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add a MySQL type | `mysqlx/<family>.go` | Struct with `Valid bool` + Scan/Value/MarshalJSON/UnmarshalJSON |
| Add a MySQL type name | `mysqlx/types.go` | `Type*` string constants |
| Change scan source coercion | `mysqlx/scan.go` | `coerce*` helpers + generic `scanSigned`/`scanUnsigned` |
| Change UUID handling | `mysqlx/uuid.go` | Wraps stdlib `uuid.UUID`; `Value()` writes CHAR(36), `Bin()` writes BINARY(16) |
| Add a pagination field | `pagex/page.go` | Add to struct + getter + Option func + proto |
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
- **Nullable MySQL types**: Every mysqlx type carries `Valid bool`, implements `sql.Scanner`/`driver.Valuer` and JSON, and NULL maps to the zero value
- **Aliases over repetition**: Same-representation MySQL variants are type aliases (`type Char = Text`, `type LineString = Geometry`) to avoid the unsafe-ql style duplication
- **Generics for shared logic**: mysqlx integer logic lives in `scanSigned`/`scanUnsigned`/`signedValue`/`marshalSignedJSON` rather than per-type copies
- **Standard library first**: mysqlx reuses stdlib types where they exist (`uuid.UUID`, `time.Time`, `math/big.Int`) instead of reimplementing them

## ANTI-PATTERNS

- `unsafesql/sql.go` repeats same pattern 30+ times with zero abstraction
- `sqls/context.go` is entirely commented-out dead code
- `unsafesql` is truly unsafe: string concatenation, no parameterization
- `pagex` pageNum is 1-based; offset calculated as `(pageNum-1)*pageSize`
- `CheckSqlInjection` regex may have false positives/negatives
- `mysqlx.Geometry` keeps raw WKB bytes and does not decode geometry primitives; only `Point` is fully decoded
- `mysqlx` has no `Codec`/`Map` registry (unlike pgtype); types are used directly with `database/sql`
- `mysqlx.UUID` scanning guesses the storage form: a 16 byte `[]byte` is read as BINARY(16), anything else as text. A CHAR(36) column is the only unambiguous choice
- MySQL has no native UUID type, so `mysqlx.TypeUUID` is MariaDB-only; MySQL users need CHAR(36) or BINARY(16)
