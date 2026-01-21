# Gofer

Gofer is a simple asynchronous task executor interface and its implementation collection that provides a unified interface for using different goroutine pool libraries.

## Introduction

Gofer defines a simple interface for asynchronously executing tasks and managing the lifecycle of goroutine pools. This project includes encapsulation implementations of various popular goroutine pool libraries.

## Features

- Unified asynchronous task execution interface
- Multiple goroutine pool implementations
- Simple and easy-to-use API design
- Support for graceful shutdown and resource release

## Interface Definition

```go
type Gofer interface {
    Go(f func()) error
    Close(ctx context.Context) error
}
```

- [Go(f func()) error](file:///Users/soyacen/Workspace/github.com/soyacen/gox/conc/gofer/gofer.go#L10-L10): Submit an asynchronous task for execution
- `Close(ctx context.Context) error`: Close the executor and wait for all tasks to complete

## Installation

```bash
go get github.com/soyacen/gox/conc/gofer
```

## Implementation List

### Sample Implementation
A simple implementation based on `sync.WaitGroup` and `atomic` without using third-party libraries. This implementation provides a complete goroutine pool with the following features:

- **Core and Non-Core Threads**: Configurable core thread count and maximum thread count. Core threads will not be recycled even when idle
- **Task Queue**: Supports buffered task queues. When the thread count reaches the core thread count, tasks will enter the queue to wait for execution
- **Dynamic Expansion**: When the task queue is full and the current thread count has not reached the maximum, new non-core threads will be created to handle tasks
- **Thread Recycling**: Non-core threads will be automatically recycled after being idle for a period of time
- **Graceful Shutdown**: Supports context-controlled shutdown timeout, waiting for all tasks to complete
- **Error Handling**: Built-in panic capture and recovery mechanism with customizable error handling functions
- **Thread Safety**: All operations are thread-safe and support concurrent access

#### Configuration Options:
- [CorePoolSize](file:///Users/soyacen/Workspace/github.com/soyacen/gox/conc/gofer/sample/gofer.go#L42-L46): Core thread count - threads will not be recycled even when idle
- [MaximumPoolSize](file:///Users/soyacen/Workspace/github.com/soyacen/gox/conc/gofer/sample/gofer.go#L49-L53): Maximum thread count the pool can accommodate
- [KeepAliveTime](file:///Users/soyacen/Workspace/github.com/soyacen/gox/conc/gofer/sample/gofer.go#L56-L60): Idle survival time for non-core threads
- [WorkQueue](file:///Users/soyacen/Workspace/github.com/soyacen/gox/conc/gofer/sample/gofer.go#L63-L67): Task queue for storing tasks waiting for execution
- [Recover](file:///Users/soyacen/Workspace/github.com/soyacen/gox/conc/gofer/sample/gofer.go#L70-L74): Custom error handling function for handling panics during task execution

#### Usage Example:
```go
// Create a thread pool with default configuration
gofer := sample.New()

// Create a thread pool with custom configuration
gofer := sample.New(
    sample.CorePoolSize(5),
    sample.MaximumPoolSize(10),
    sample.KeepAliveTime(60*time.Second)
)

// Submit a task
gofer.Go(func() {
    // Execute task
})

// Close the thread pool
gofer.Close(context.Background())
```

### Ants Implementation
Implementation based on [github.com/panjf2000/ants/v2](https://github.com/panjf2000/ants/v2).

### Tunny Implementation
Implementation based on [github.com/Jeffail/tunny](https://github.com/Jeffail/tunny).

### WorkerPool Implementation
Implementation based on [github.com/gammazero/workerpool](https://github.com/gammazero/workerpool).

### GRPool Implementation
Implementation based on [github.com/ivpusic/grpool](https://github.com/ivpusic/grpool).

### GoPgPool Implementation
Implementation based on [gopkg.in/go-playground/pool.v3](https://gopkg.in/go-playground/pool.v3).

## License

MIT