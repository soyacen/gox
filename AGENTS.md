# PROJECT KNOWLEDGE BASE

**Generated:** 2026-04-23
**Commit:** 6721e96
**Branch:** main

## OVERVIEW

`gox` is a comprehensive Go utility library (monorepo style) providing extensions across concurrency, crypto, strings, collections, HTTP, database, and more. Module: `github.com/soyacen/gox`, Go 1.25.0.

## STRUCTURE

```
.
├── conc/           # Concurrency: pools, barriers, channels, maps, mutexes
├── cryptox/        # Cryptography: AES, RSA, SHA, HMAC, TLS, MD5
├── databasex/      # Database: pagination, SQL injection detection, query builder
├── stringx/        # String utilities: 22 files, heaviest package
├── slicex/         # Slice generics: Max, Min, Filter, Map, Reduce, etc.
├── httpx/          # HTTP client/server tools (mostly deprecated)
├── backoff/        # Backoff algorithms: exponential, fibonacci, linear, constant
├── errorx/         # Error helpers: Must, Ignore, Concern, Quiet, chain, join
├── reflectx/       # Reflection: unsafe accessor, field access
├── timex/          # Time/date utilities
├── slogx/          # Structured logging extensions
├── protox/         # Protobuf: clone, merge, wrapper conversions
├── strconvx/       # String-to-type conversions
├── filex/          # File operations
├── iox/            # IO utilities
├── imagex/          # Image processing
├── gen/            # Code generation tools
└── tools/          # Build tools (go:build tools)
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add slice/collection utility | `slicex/` | Heavy generics, uses `constraintx` |
| Add string function | `stringx/` | Filter func variants (`PadF`), `ToArgv` parser |
| Add concurrency primitive | `conc/` | See `conc/AGENTS.md` |
| Add crypto algorithm | `cryptox/` | See `cryptox/AGENTS.md` |
| Add database tool | `databasex/` | See `databasex/AGENTS.md` |
| Add HTTP helper | `httpx/` | Most APIs deprecated; prefer `netx/httpx/outgoing` |
| Error handling pattern | `errorx/` | `Must[T]`, `Ignore[T]`, `Quiet`, `Break`/`Continue` chains |
| Generic constraint | `constraintx/` | `Numeric = Integer \| Float` |
| Pointer/ternary helper | `operator/` | `Pointer[T]`, `Ternary[T]`, `SwitchCases[T]` |
| Context extension | `contextx/` | Deadline, error wrapping, reflection |
| Retry logic | `retry/` | Fluent: `MaxAttempts(3).Backoff(...).RetryOn(...).Exec(...)` |
| Logging extension | `slogx/` | Context-aware structured logging |
| Protobuf bridge | `protox/` | `Clone`, `Merge`, `Wrapper` slice conversions |
| Goroutine pool adapter | `conc/gofer/<pool>/` | 5 third-party pool adapters |

## CODE MAP

| Symbol | Type | Location | Role |
|--------|------|----------|------|
| `Must` | func | `errorx/error.go` | Generic panic-on-error helper |
| `Ternary` | func | `operator/ternary.go` | Generic ternary operator |
| `BackoffStrategy` | interface | `retry/retry.go` | Fluent retry chaining |
| `Gofer` | interface | `conc/gofer/gofer.go` | Goroutine pool abstraction |
| `Option` | type | `conc/asyncbatch/group.go` | Functional options pattern |
| `Numeric` | constraint | `constraintx/constraintx.go` | Generic numeric constraint |
| `NoCopy` | struct | `conc/no_copy.go` | Copylocks vet checker marker |

## CONVENTIONS

- **Package naming**: All packages use `x` suffix (`slicex`, `stringx`, `mutexx`)
- **Filter functions**: Many packages provide `Func` and `FuncF` variants (e.g., `Pad` / `PadF`)
- **Must* pattern**: `Must[T](v, err)` panics on error; `MustEq` / `MustWhere` in `unsafesql`
- **Dual API**: Crypto packages provide raw bytes + Hex string variants (`Sha256` / `Sha256Hex`)
- **Generic structs**: `Pool[T]`, `Group[Obj]`, `Ring[T]` use type parameters
- **Error variables**: Package-prefixed sentinel errors (`ErrNilFunction`, `ErrTaskInvalid`)
- **Context first**: All async/concurrent APIs accept `context.Context`
- **Duck-typed interfaces**: `WaitNotify(waiter interface{ Wait() })` avoids import coupling

## ANTI-PATTERNS (THIS PROJECT)

- **Deprecated httpx APIs**: `RequestBuilder` and `ResponseHelper` are fully deprecated; use `github.com/soyacen/netx/httpx/outgoing`
- `slicex.Difference` - deprecated
- `conc/chanx.Pipe` / `AsyncPipe` - deprecated, use `Copy` / `AsyncCopy`
- `PooledClient()` / `PooledTransport()` - do NOT use for transient clients/transports
- `strconvx.StringToBytes` / `BytesToString` - use `unsafe`; results must not leak to end users
- `mutexx/mutex.go` - uses `unsafe.Pointer` on `sync.Mutex` internals
- Very minimal `internal/` usage - only `stringx/internal/` exists

## UNIQUE STYLES

- **Chinese comments**: Most source comments and docs are in Chinese
- **Go version mismatch**: Root `go.mod` declares 1.25.0; CI uses 1.20; submodules use 1.20
- **Submodules for pool adapters**: `conc/gofer/*/` each have independent `go.mod`
- **No `pkg/` or `internal/` directories**: Public API lives flat at root
- **Placeholder CI messages**: `.github/workflows/greetings.yml` uses literal placeholder text

## COMMANDS

```bash
# Run tests
make test          # go test -v ./...

# Build all packages
go build -v ./...

# Run tests for a submodule
cd conc/gofer/ants && go test -v ./...
```

## NOTES

- `conc/gofer/workerpool/go.mod` has wrong module path (declares `.../tunny` instead of `.../workerpool`)
- `errorx/chain.go` provides `Break` / `Continue` for sequential error handling (like loop control)
- `stringx/funcsPZ.go` contains a hand-written state-machine parser (`ToArgv`) for shell-like tokenization
- `databasex/unsafesql/sql.go` repeats the same pattern 30+ times (operator + MustOperator) with no abstraction
- `backoff/` and `retry/` are separate packages; `retry` depends on `backoff`
- No fuzz tests, no example tests, no golangci-lint config
