# Gofer

Gofer 是一个简单的异步任务执行器接口及其实现集合，提供了统一的接口来使用不同的 goroutine 池库。

## 简介

Gofer 定义了一个简单的接口，用于异步执行任务并管理 goroutine 池的生命周期。该项目包含了多种流行的 goroutine 池库的封装实现。

## 功能特性

- 统一的异步任务执行接口
- 多种 goroutine 池实现
- 简单易用的 API 设计
- 支持优雅关闭和资源释放

## 接口定义

```go
type Gofer interface {
    Go(f func()) error
    Close(ctx context.Context) error
}
```

- `Go(f func()) error`: 提交一个异步任务执行
- `Close(ctx context.Context) error`: 关闭执行器并等待所有任务完成


## 安装

```bash
go get github.com/soyacen/gox/conc/gofer
```

## 实现列表

### Sample 实现
基于 `sync.WaitGroup` 和 `atomic` 的简单实现，不使用第三方库。该实现提供了一个完整的 goroutine 池，具有以下特性：

- **核心线程与非核心线程**: 可配置核心线程数和最大线程数，核心线程不会因空闲而被回收
- **任务队列**: 支持缓冲任务队列，当线程数达到核心线程数后，任务将进入队列等待执行
- **动态扩展**: 当任务队列满且当前线程数未达到最大线程数时，会创建新的非核心线程处理任务
- **线程回收**: 非核心线程在空闲一段时间后会被自动回收
- **优雅关闭**: 支持通过 context 控制关闭超时，等待所有任务执行完毕
- **错误处理**: 内置 panic 捕获和恢复机制，支持自定义错误处理函数
- **线程安全**: 所有操作都是线程安全的，支持并发访问

#### 配置选项：
- [CorePoolSize](file:///Users/soyacen/Workspace/github.com/soyacen/gox/conc/gofer/sample/gofer.go#L42-L46): 核心线程数，即使线程处于空闲状态也不会被回收
- [MaximumPoolSize](file:///Users/soyacen/Workspace/github.com/soyacen/gox/conc/gofer/sample/gofer.go#L49-L53): 线程池所能容纳的最大线程数
- [KeepAliveTime](file:///Users/soyacen/Workspace/github.com/soyacen/gox/conc/gofer/sample/gofer.go#L56-L60): 非核心线程闲置时的存活时间
- [WorkQueue](file:///Users/soyacen/Workspace/github.com/soyacen/gox/conc/gofer/sample/gofer.go#L63-L67): 任务队列，用于存放等待执行的任务
- [Recover](file:///Users/soyacen/Workspace/github.com/soyacen/gox/conc/gofer/sample/gofer.go#L70-L74): 自定义错误处理函数，用于处理任务执行过程中的 panic

#### 使用示例：
```go
// 创建默认配置的线程池
gofer := sample.New()

// 创建自定义配置的线程池
gofer := sample.New(
    sample.CorePoolSize(5),
    sample.MaximumPoolSize(10),
    sample.KeepAliveTime(60*time.Second)
)

// 提交任务
gofer.Go(func() {
    // 执行任务
})

// 关闭线程池
gofer.Close(context.Background())
```

### Ants 实现
基于 [github.com/panjf2000/ants/v2](https://github.com/panjf2000/ants/v2) 的实现。

### Tunny 实现
基于 [github.com/Jeffail/tunny](https://github.com/Jeffail/tunny) 的实现。

### WorkerPool 实现
基于 [github.com/gammazero/workerpool](https://github.com/gammazero/workerpool) 的实现。

### GRPool 实现
基于 [github.com/ivpusic/grpool](https://github.com/ivpusic/grpool) 的实现。

### GoPgPool 实现
基于 [gopkg.in/go-playground/pool.v3](https://gopkg.in/go-playground/pool.v3) 的实现。

## 许可证

MIT