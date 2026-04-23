# conc AGENTS.md

**Generated:** 2026-04-23

## OVERVIEW

Concurrency primitives: barriers, channels, goroutine pools, mutex extensions, lazy loading, async batching, and panic recovery.

## STRUCTURE

```
conc/
├── no_copy.go          # NoCopy struct (copylocks vet marker)
├── asyncbatch/         # Async batch processor (size/time triggered)
├── atomicx/            # Atomic operation helpers
├── barrier/            # CyclicBarrier (Group coordinator)
├── brave/              # Panic recovery wrappers
├── chanx/              # Channel ops: FanIn, FanOut, Pipeline, Map, Filter
├── gofer/              # Gofer interface + 5 third-party pool adapters
│   ├── ants/           # github.com/panjf2000/ants
│   ├── grpool/         # github.com/ivpusic/grpool
│   ├── gopgpool/       # github.com/alitto/pond
│   ├── tunny/          # github.com/Jeffail/tunny
│   └── workerpool/     # github.com/gammazero/workerpool
├── goroutinex/         # Goroutine utilities
├── lazyload/           # singleflight-based lazy load (Group)
├── mapx/               # MapInterface + ShardedMap
├── mutexx/             # TryLock via unsafe.Pointer on sync.Mutex
├── oncex/              # sync.Once per key (Group)
├── poolx/              # Object/buffer pools
└── waiter/             # Wait() → channel conversion
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add goroutine pool adapter | `gofer/<pool>/` | Each has independent go.mod; must implement `Gofer` interface |
| Add batch coordinator | `asyncbatch/` | Follow Option + Apply + Correct pattern |
| Add channel primitive | `chanx/chan.go` | File is 673 lines; add alongside existing ops |
| Add map implementation | `mapx/` | Must satisfy `MapInterface` for `ShardedMap` compatibility |
| Panic recovery pattern | `brave/` or `asyncbatch` | `recover()` + `debug.Stack()` |

## CONVENTIONS

- **"Group" naming**: oncex.Group, lazyload.Group, asyncbatch.Group, barrier.Group = coordinator types
- **Functional options**: asyncbatch uses `Option` + `Apply(...)` + `Correct()` chain
- **Gofer interface**: `Go(ctx, f)` / `Close(ctx)` — 5 adapter submodules
- **Generic structs**: `Pool[T]`, `Group[Obj]` use type parameters
- **Duck-typed interfaces**: `WaitNotify(waiter interface{ Wait() })` avoids importing sync
- **Error sentinels**: `ErrNilFunction`, `ErrTaskInvalid`

## ANTI-PATTERNS

- **mutexx/mutex.go**: Uses `unsafe.Pointer` on `sync.Mutex` internals — avoid if possible
- **chanx.Pipe / AsyncPipe**: Deprecated, use `Copy` / `AsyncCopy`
- **Workerpool submodule**: `conc/gofer/workerpool/go.mod` has wrong module path (declares `.../tunny`)
