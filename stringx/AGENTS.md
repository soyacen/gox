# stringx AGENTS.md

## OVERVIEW

String manipulation utilities — 22 Go files, the largest single package in gox.

## WHERE TO LOOK

| Task | File | Notes |
|------|------|-------|
| Pad, Template, ToArgv | `funcsPZ.go` | 533 lines, hand-written state-machine parser for shell-like tokenization |
| Case conversion | `case.go` | `GoCamelCase`, protobuf-style identifier conversion |
| Builder wrapper | `builder.go` | Wraps `strings.Builder` with nil-check lazy init |
| Builder pooling | `pool.go` | `sync.Pool` for `strings.Builder`, caps at 4KB |
| Empty/blank checks | `empty.go`, `blank.go` | Generic `~string` constraints |
| Remove ops | `remove.go` | `RemoveApostrophe` delegates to `internal/` |
| Join, indexes, replace | `join.go`, `indexes.go`, `replace.go` | Standard string ops |
| Match, fold | `match.go`, `fold.go` | Pattern matching, Unicode folding |
| Internal helper | `internal/remove_apostrophe.go` | Only `internal/` usage in entire repo |

## CONVENTIONS

- **Filter functions**: Dual API — `Func` and `FuncF` variants (e.g. `Pad`/`PadF`, `PadLeft`/`PadLeftF`)
- **Builder lifecycle**: Use `GetBuilder()`/`PutBuilder()` for pooled builders; cap checked on return
- **Test packages**: Mixed `package stringx` and `package stringx_test` (examples in-package, builder tests external)
- **Generic constraints**: `~string` used for empty/blank checks to accept string aliases

## ANTI-PATTERNS

- `Builder` wrapper adds indirection over `strings.Builder`; nil-check on every method call
- `internal/` package contains only `RemoveApostrophe` — consider if it justifies the separation
- `ToArgv` state machine parser is complex; changes require careful testing against shell tokenization edge cases
- Do not leak pooled builders — always `PutBuilder()` after use
