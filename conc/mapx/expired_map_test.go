package mapx

import (
	"testing"
	"time"
)

func TestExpiredMap_Store(t *testing.T) {
	key := "testKey"
	value := "testValue"
	expiredMap := NewExpiredMap(ExpireAfter(func(k any) time.Duration {
		if key == k {
			return time.Second
		}
		return 0
	}))

	expiredMap.Store(key, value)

	result, found := expiredMap.Load(key)
	if !found {
		t.Errorf("Expected item to be found, but it was not")
	}
	if result != value {
		t.Errorf("Expected value to be %v, but got %v", value, result)
	}

	time.Sleep(time.Second * 2)
	result, found = expiredMap.Load(key)
	if found {
		t.Errorf("Expected item to be not found, but it was found")
	}
}

func TestExpiredMap_LoadOrStore(t *testing.T) {
	key := "testKey"
	value := "testValue"

	expiredMap := NewExpiredMap(ExpireAfter(func(k any) time.Duration {
		if key == k {
			return time.Second
		}
		return 0
	}))
	result, found := expiredMap.LoadOrStore(key, value)
	if found {
		t.Errorf("Expected item to not be found, but it was")
	}
	if result != value {
		t.Errorf("Expected value to be %v, but got %v", value, result)
	}

	// Store another value with the same key
	newValue := "newValue"
	result, found = expiredMap.LoadOrStore(key, newValue)
	if !found {
		t.Errorf("Expected item to be found, but it was not")
	}
	if result == newValue {
		t.Errorf("Expected value to be %v, but got %v", newValue, result)
	}

	time.Sleep(2 * time.Second)
	result, found = expiredMap.LoadOrStore(key, newValue)
	if found {
		t.Errorf("Expected item to be not found, but it was found")
	}
}

func TestExpiredMap_LoadAndDelete(t *testing.T) {
	key := "testKey"
	value := "testValue"
	expiredMap := NewExpiredMap(ExpireAfter(func(k any) time.Duration {
		if key == k {
			return time.Second
		}
		return 0
	}))
	expiredMap.Store(key, value)
	result, found := expiredMap.LoadAndDelete(key)
	if !found {
		t.Errorf("Expected item to be found, but it was not")
	}
	if result != value {
		t.Errorf("Expected value to be %v, but got %v", value, result)
	}

	expiredMap.Store(key, value)
	time.Sleep(2 * time.Second)
	result, found = expiredMap.LoadAndDelete(key)
	if found {
		t.Errorf("Expected item to be not found, but it was found")
	}

}

func TestExpiredMap_Delete(t *testing.T) {
	key := "testKey"
	value := "testValue"
	expiredMap := NewExpiredMap(ExpireAfter(func(k any) time.Duration {
		if key == k {
			return time.Second
		}
		return 0
	}))
	expiredMap.Store(key, value)
	expiredMap.Delete(key)
	_, found := expiredMap.Load(key)
	if found {
		t.Errorf("Expected item to be deleted, but it was not")
	}

	expiredMap.Store(key, value)
	time.Sleep(2 * time.Second)
	expiredMap.Delete(key)
	_, found = expiredMap.Load(key)
	if found {
		t.Errorf("Expected item to be deleted, but it was not")
	}
}

func TestExpiredMap_Swap(t *testing.T) {
	key := "testKey"
	oldValue := "oldValue"
	newValue := "newValue"
	expiredMap := NewExpiredMap(ExpireAfter(func(k any) time.Duration {
		if key == k {
			return time.Second
		}
		return 0
	}))
	expiredMap.Store(key, oldValue)

	result, found := expiredMap.Swap(key, newValue)
	if !found {
		t.Errorf("Expected item to be found, but it was not")
	}
	if result != oldValue {
		t.Errorf("Expected previous value to be %v, but got %v", oldValue, result)
	}
	if val, _ := expiredMap.Load(key); val != newValue {
		t.Errorf("Expected value to be swapped, but it was not")
	}

	time.Sleep(2 * time.Second)
	result, found = expiredMap.Swap(key, oldValue)
	if found {
		t.Errorf("Expected item to be not found, but it was found")
	}

}

func TestExpiredMap_CompareAndSwap(t *testing.T) {
	key := "testKey"
	oldValue := "oldValue"
	newValue := "newValue"

	expiredMap := NewExpiredMap(ExpireAfter(func(k any) time.Duration {
		if key == k {
			return time.Second
		}
		return 0
	}))
	expiredMap.Store(key, oldValue)

	success := expiredMap.CompareAndSwap(key, oldValue, newValue)
	if !success {
		t.Errorf("Expected CompareAndSwap to succeed, but it failed")
	}
	if val, _ := expiredMap.Load(key); val != newValue {
		t.Errorf("Expected value to be swapped, but it was not")
	}

	success = expiredMap.CompareAndSwap(key, "wrongValue", newValue)
	if success {
		t.Errorf("Expected CompareAndSwap to fail, but it succeeded")
	}

	time.Sleep(2 * time.Second)
	success = expiredMap.CompareAndSwap(key, newValue, oldValue)
	if success {
		t.Errorf("Expected CompareAndSwap to fail, but it succeeded")
	}
}

func TestExpiredMap_CompareAndDelete(t *testing.T) {
	key := "testKey"
	oldValue := "oldValue"
	expiredMap := NewExpiredMap(ExpireAfter(func(k any) time.Duration {
		if key == k {
			return time.Second
		}
		return 0
	}))

	expiredMap.Store(key, oldValue)
	success := expiredMap.CompareAndDelete(key, oldValue)
	if !success {
		t.Errorf("Expected CompareAndDelete to succeed, but it failed")
	}
	_, found := expiredMap.Load(key)
	if found {
		t.Errorf("Expected item to be deleted, but it was not")
	}

	expiredMap.Store(key, oldValue)
	success = expiredMap.CompareAndDelete(key, "wrongValue")
	if success {
		t.Errorf("Expected CompareAndDelete to fail, but it succeeded")
	}

	time.Sleep(2 * time.Second)
	success = expiredMap.CompareAndDelete(key, oldValue)
	if success {
		t.Errorf("Expected CompareAndDelete to failed, but it successed")
	}

}

func TestExpiredMap_Range(t *testing.T) {
	expiredMap := NewExpiredMap()

	keys := []string{"key1", "key2", "key3"}
	values := []string{"value1", "value2", "value3"}

	for i, key := range keys {
		expiredMap.Store(key, values[i])
	}

	count := 0
	expiredMap.Range(func(key any, value any) bool {
		count++
		return true
	})

	if count != len(keys) {
		t.Errorf("Expected Range to iterate over %v items, but it iterated over %v", len(keys), count)
	}
}
