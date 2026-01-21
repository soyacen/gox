package lazyload

import (
	"errors"
	"sync"

	"golang.org/x/sync/singleflight"
)

var ErrNilFunction = errors.New("lazyload: New function is nil")

// Group 用于实现缓存和懒加载功能
type Group[Obj any] struct {
	// m 用于存储键值对
	m sync.Map

	// g 用于并发控制，确保多个 goroutine 同时请求同一个键时，只有一个 goroutine 执行创建逻辑。
	g singleflight.Group

	// New 一个函数指针，用于创建新的值。在创建时设置，之后不再修改
	New func(key string) (Obj, error)
}

// Load 用于获取键对应的值。如果 key 不存在，则调用 g.New 创建新值。
func (g *Group[Obj]) Load(key string) (Obj, error, bool) {
	return g.LoadOrNew(key, nil)
}

// Delete 删除指定的键值对。
func (g *Group[Obj]) Delete(key string) {
	g.m.Delete(key)
}

// Range 遍历所有键值对，并调用 f 函数。
func (g *Group[Obj]) Range(f func(key string, value Obj) bool) {
	g.m.Range(func(key, value any) bool {
		return f(key.(string), value.(Obj))
	})
}

// LoadOrNew 用于获取键对应的值。如果 key 不存在，则调用 f 或 g.New 创建新值。
func (g *Group[Obj]) LoadOrNew(key string, f func(key string) (Obj, error)) (Obj, error, bool) {
	var obj Obj

	// 1. 初次检查,如果 key 已存在于 sync.Map 中，直接返回该值。
	if value, ok := g.m.Load(key); ok {
		obj = value.(Obj)
		return obj, nil, true
	}

	// 2. 并发控制,g.g.Do 会确保并发控制，即如果有多个 goroutine 同时请求同一个 key，只有一个 goroutine 会执行闭包中的逻辑。
	if f == nil && g.New == nil {
		return obj, ErrNilFunction, false
	}
	if f == nil {
		f = g.New
	}
	value, err, _ := g.g.Do(key, func() (any, error) {
		// 3. 再次检查 key 是否已存在于 sync.Map 中, 如果存在，则直接返回该值。如果不存在，则调用 g.New 创建新值。
		if value, ok := g.m.Load(key); ok {
			return value, nil
		}
		value, err := f(key)
		if err != nil {
			return nil, err
		}

		// 4. 将新值存入 sync.Map 中
		g.m.Store(key, value)
		return value, nil
	})
	if err != nil {
		return obj, err, false
	}

	obj = value.(Obj)
	return obj, nil, false
}
