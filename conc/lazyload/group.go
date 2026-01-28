package lazyload

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/singleflight"
)

// ErrNilFunction 当 New 函数为 nil 时返回的错误
var ErrNilFunction = errors.New("lazyload: New function is nil")

// entry 包装了实际对象和其清理函数
type entry[Obj any] struct {
	Obj       Obj                                      // 实际存储的对象
	CloseFunc func(ctx context.Context, obj Obj) error // 对象清理函数
}

// Group 实现了并发安全的懒加载缓存，确保相同键的值只创建一次
// 支持自定义创建和清理函数
type Group[Obj any] struct {
	m        sync.Map                                 // 存储键值对，value 为 entry 类型
	g        singleflight.Group                       // 并发控制，防止重复创建
	New      func(key string) (Obj, error)            // 创建新值的函数
	Finalize func(ctx context.Context, obj Obj) error // 清理资源的函数
}

// Option 允许自定义 LoadOrNew 的行为
type Option[Obj any] func(*options[Obj])

// options 存储 LoadOrNew 的配置选项
type options[Obj any] struct {
	factory   func(key string) (Obj, error)            // 自定义创建函数
	finalizer func(ctx context.Context, obj Obj) error // 自定义清理函数
}

// WithFactory 指定自定义的对象创建函数
func WithFactory[Obj any](factory func(key string) (Obj, error)) Option[Obj] {
	return func(opts *options[Obj]) {
		opts.factory = factory
	}
}

// WithFinalizer 指定自定义的对象清理函数
func WithFinalizer[Obj any](finalizer func(ctx context.Context, obj Obj) error) Option[Obj] {
	return func(opts *options[Obj]) {
		opts.finalizer = finalizer
	}
}

// Load 获取键对应的值，不存在时使用默认的 New 函数创建
// 返回值、错误和是否存在标志
func (g *Group[Obj]) Load(key string) (Obj, error, bool) {
	return g.LoadOrNew(key)
}

// LoadOrNew 加载现有值或使用指定选项创建新值
// 返回值、错误和是否存在标志(true=已存在, false=新创建)
func (g *Group[Obj]) LoadOrNew(key string, opts ...Option[Obj]) (Obj, error, bool) {
	// 快速路径：如果值已存在，直接返回
	if value, ok := g.m.Load(key); ok {
		entry := value.(*entry[Obj])
		return entry.Obj, nil, true
	}

	// 应用选项
	opt := &options[Obj]{}
	for _, o := range opts {
		o(opt)
	}

	// 确定创建函数
	factory := opt.factory
	if factory == nil {
		factory = g.New
	}
	if factory == nil {
		var obj Obj
		return obj, ErrNilFunction, false
	}

	// 使用 singleflight 防止重复创建
	value, err, _ := g.g.Do(key, func() (any, error) {
		// 再次检查，防止在等待期间被其他 goroutine 创建
		if value, ok := g.m.Load(key); ok {
			return value, nil
		}
		value, err := factory(key)
		if err != nil {
			return nil, err
		}

		// 使用适当的清理函数存储值
		finalizer := g.Finalize
		if opt.finalizer != nil {
			finalizer = opt.finalizer
		}
		g.m.Store(key, &entry[Obj]{
			Obj:       value,
			CloseFunc: finalizer,
		})
		return value, nil
	})
	if err != nil {
		var obj Obj
		return obj, err, false
	}

	obj := value.(Obj)
	return obj, nil, false
}

// Close 清理所有缓存项并释放资源
// 返回所有清理错误的组合
func (g *Group[Obj]) Close(ctx context.Context) error {
	var errs []error
	// 遍历并清理所有缓存项
	g.m.Range(func(key, value any) bool {
		g.m.Delete(key.(string))
		entry := value.(*entry[Obj])
		if entry.CloseFunc != nil {
			errs = append(errs, entry.CloseFunc(ctx, entry.Obj))
		}
		return true
	})
	return errors.Join(errs...)
}
