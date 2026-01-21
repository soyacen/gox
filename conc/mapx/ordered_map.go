package mapx

import (
	"container/list"
	"sync"
)

var _ MapInterface = (*OrderedMap)(nil)

type entry[K comparable, V any] struct {
	Key   K
	Value V
}

type OrderedMap struct {
	mu    sync.RWMutex
	items map[any]*list.Element
	list  *list.List
}

func (m *OrderedMap) Load(key any) (value any, loaded bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	elem, loaded := m.items[key]
	if !loaded {
		return nil, loaded
	}
	return elem.Value.(*entry[any, any]).Value, loaded
}

func (m *OrderedMap) Store(key, value any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, ok := m.items[key]
	if !ok {
		entry := &entry[any, any]{Key: key, Value: value}
		elem = m.list.PushBack(entry)
		m.items[key] = elem
		return
	}
	elem.Value.(*entry[any, any]).Value = value
}

func (m *OrderedMap) LoadOrStore(key, value any) (actual any, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, loaded := m.items[key]
	if !loaded {
		entry := &entry[any, any]{Key: key, Value: value}
		elem = m.list.PushBack(entry)
		m.items[key] = elem
		return value, loaded
	}
	return elem.Value.(*entry[any, any]).Value, loaded
}

func (m *OrderedMap) LoadAndDelete(key any) (value any, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, loaded := m.items[key]
	if !loaded {
		return nil, loaded
	}
	delete(m.items, key)
	m.list.Remove(elem)
	return elem.Value.(*entry[any, any]).Value, loaded
}

func (m *OrderedMap) Delete(key any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, ok := m.items[key]
	if !ok {
		return
	}
	delete(m.items, key)
	m.list.Remove(elem)
}

func (m *OrderedMap) Swap(key, value any) (previous any, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, loaded := m.items[key]
	if !loaded {
		return previous, loaded
	}
	entry := elem.Value.(*entry[any, any])
	previous = entry.Value
	entry.Value = value
	return previous, loaded
}

func (m *OrderedMap) CompareAndSwap(key, old, new any) (swapped bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, loaded := m.items[key]
	if !loaded {
		return false
	}
	entry := elem.Value.(*entry[any, any])
	if entry.Value != old {
		return false
	}
	entry.Value = new
	return true
}

func (m *OrderedMap) CompareAndDelete(key, old any) (deleted bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, loaded := m.items[key]
	if !loaded {
		return false
	}
	entry := elem.Value.(*entry[any, any])
	if entry.Value != old {
		return false
	}
	delete(m.items, key)
	m.list.Remove(elem)
	return true
}

func (m *OrderedMap) Range(f func(key any, value any) (shouldContinue bool)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for elem := m.list.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*entry[any, any])
		if !f(entry.Key, entry.Value) {
			break
		}
	}
}

func (m *OrderedMap) Front() *list.Element {
	return m.list.Front()
}

func (m *OrderedMap) Back() *list.Element {
	return m.list.Back()
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		items: make(map[any]*list.Element),
		list:  list.New(),
	}
}
