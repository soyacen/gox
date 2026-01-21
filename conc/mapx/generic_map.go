package mapx

type GenericMap[K comparable, V any] struct {
	MapInterface MapInterface
}

func (m *GenericMap[K, V]) Load(key K) (V, bool) {
	var v V
	load, ok := m.MapInterface.Load(key)
	if !ok {
		return v, false
	}
	v = load.(V)
	return v, true
}

func (m *GenericMap[K, V]) Store(key K, value V) {
	m.MapInterface.Swap(key, value)
}

func (m *GenericMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	var v V
	actual, loaded := m.MapInterface.LoadOrStore(key, value)
	if !loaded {
		return v, false
	}
	v = actual.(V)
	return v, true
}

func (m *GenericMap[K, V]) LoadAndDelete(key K) (V, bool) {
	var v V
	value, loaded := m.MapInterface.LoadAndDelete(key)
	if !loaded {
		return v, false
	}
	v = value.(V)
	return v, true
}

func (m *GenericMap[K, V]) Delete(key K) {
	m.MapInterface.Delete(key)
}

func (m *GenericMap[K, V]) Swap(key K, value V) (V, bool) {
	var v V
	previous, loaded := m.MapInterface.Swap(key, value)
	if !loaded {
		return v, false
	}
	v = previous.(V)
	return v, true
}

func (m *GenericMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	return m.MapInterface.CompareAndSwap(key, old, new)
}

func (m *GenericMap[K, V]) CompareAndDelete(key K, old V) bool {
	return m.MapInterface.CompareAndDelete(key, old)
}

func (m *GenericMap[K, V]) Range(f func(key K, value V) bool) {
	m.MapInterface.Range(func(key, value any) bool {
		k := key.(K)
		v := value.(V)
		return f(k, v)
	})
}
