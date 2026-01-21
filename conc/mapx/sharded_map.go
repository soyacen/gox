package mapx

var _ MapInterface = (*ShardedMap)(nil)

type ShardedMap struct {
	Segments []MapInterface
	Hash     func(key any) int
}

func (m *ShardedMap) Load(key any) (value any, loaded bool) {
	return m.bucket(key).Load(key)
}

func (m *ShardedMap) Store(key, value any) {
	m.bucket(key).Store(key, value)
}

func (m *ShardedMap) LoadOrStore(key, value any) (actual any, loaded bool) {
	return m.bucket(key).LoadOrStore(key, value)
}

func (m *ShardedMap) LoadAndDelete(key any) (value any, loaded bool) {
	return m.bucket(key).LoadAndDelete(key)
}

func (m *ShardedMap) Delete(key any) {
	m.bucket(key).Delete(key)
}

func (m *ShardedMap) Swap(key, value any) (previous any, loaded bool) {
	return m.bucket(key).Swap(key, value)
}

func (m *ShardedMap) CompareAndSwap(key, old, new any) (swapped bool) {
	return m.bucket(key).CompareAndSwap(key, old, new)
}

func (m *ShardedMap) CompareAndDelete(key, old any) (deleted bool) {
	return m.bucket(key).CompareAndDelete(key, old)
}

func (m *ShardedMap) Range(f func(key any, value any) (shouldContinue bool)) {
	for _, segment := range m.Segments {
		segment.Range(f)
	}
}

func (m *ShardedMap) bucket(k any) MapInterface {
	return m.Segments[m.Hash(k)%len(m.Segments)]
}
