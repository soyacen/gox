package heapx

import (
	"container/heap"
	"fmt"
)

// KeyError represents an error that occurred while calling a key function.
type KeyError struct {
	Obj any
	Err error
}

// Error implements the error interface for KeyError.
func (k KeyError) Error() string {
	return fmt.Sprintf("heap: failed to call key func, obj: %v, error : %v", k.Obj, k.Err)
}

// lessFunc is a function that receives two items and returns true if the first
// item should be placed before the second one when the list is sorted.
type lessFunc = func(item1, item2 any) bool

// KeyFunc is a function type to get the key from an object.
// KeyFunc is a function type to get the key from an object.
type KeyFunc func(obj any) (string, error)

type heapItem struct {
	obj   any // The object which is stored in the heap.
	index int // The index of the object's key in the Heap.queue.
}

type itemKeyValue struct {
	key string
	obj any
}

// data is an internal struct that implements the standard heap interface
// and keeps the data stored in the heap.
type data struct {
	// items is a map from key of the objects to the objects and their index.
	// We depend on the property that items in the map are in the queue and vice versa.
	items map[string]*heapItem
	// queue implements a heap data structure and keeps the order of elements
	// according to the heap invariant. The queue keeps the keys of objects stored
	// in "items".
	queue []string

	// keyFunc is used to make the key used for queued item insertion and retrieval, and
	// should be deterministic.
	keyFunc KeyFunc
	// lessFunc is used to compare two objects in the heap.
	lessFunc lessFunc
}

var (
	_ = heap.Interface(&data{}) // heapData is a standard heap
)

// Less compares two objects and returns true if the first one should go
// in front of the second one in the heap.
//
// Parameters:
//   - i: The index of the first item
//   - j: The index of the second item
//
// Returns:
//   - bool: True if the first item should go in front
func (h *data) Less(i, j int) bool {
	if i > len(h.queue) || j > len(h.queue) {
		return false
	}
	itemi, ok := h.items[h.queue[i]]
	if !ok {
		return false
	}
	itemj, ok := h.items[h.queue[j]]
	if !ok {
		return false
	}
	return h.lessFunc(itemi.obj, itemj.obj)
}

// Len returns the number of items in the heap.
//
// Returns:
//   - int: The number of items
func (h *data) Len() int { return len(h.queue) }

// Swap implements swapping of two elements in the heap.
// This is part of the standard heap interface and should never be called directly.
//
// Parameters:
//   - i: The index of the first item
//   - j: The index of the second item
func (h *data) Swap(i, j int) {
	h.queue[i], h.queue[j] = h.queue[j], h.queue[i]
	item := h.items[h.queue[i]]
	item.index = i
	item = h.items[h.queue[j]]
	item.index = j
}

// Push adds an item to the heap.
// It is supposed to be called by heap.Push only.
//
// Parameters:
//   - kv: The key-value pair to push
func (h *data) Push(kv any) {
	keyValue := kv.(*itemKeyValue)
	n := len(h.queue)
	h.items[keyValue.key] = &heapItem{keyValue.obj, n}
	h.queue = append(h.queue, keyValue.key)
}

// Pop removes and returns the last item from the heap.
// It is supposed to be called by heap.Pop only.
//
// Returns:
//   - any: The popped item
func (h *data) Pop() any {
	key := h.queue[len(h.queue)-1]
	h.queue = h.queue[0 : len(h.queue)-1]
	item, ok := h.items[key]
	if !ok {
		// This is an error
		return nil
	}
	delete(h.items, key)
	return item.obj
}

// Peek returns the top item from the heap without removing it.
//
// Returns:
//   - any: The top item, or nil if the heap is empty
func (h *data) Peek() any {
	if len(h.queue) > 0 {
		return h.items[h.queue[0]].obj
	}
	return nil
}

// Heap is a producer/consumer queue that implements a heap data structure.
// It can be used to implement priority queues and similar data structures.
//
// Fields:
//   - data: Stores objects and maintains their ordering according to the heap invariant.
type Heap struct {
	// data stores objects and has a queue that keeps their ordering according
	// to the heap invariant.
	data *data
}

// Add inserts an item into the heap.
// If the item already exists, it is updated.
//
// Parameters:
//   - obj: The object to insert or update
//
// Returns:
//   - error: Error if the key function fails
func (h *Heap) Add(obj any) error {
	key, err := h.data.keyFunc(obj)
	if err != nil {
		return KeyError{Obj: obj, Err: err}
	}
	if _, exists := h.data.items[key]; exists {
		h.data.items[key].obj = obj
		heap.Fix(h.data, h.data.items[key].index)
	} else {
		heap.Push(h.data, &itemKeyValue{key, obj})
	}
	return nil
}

// Update updates an item in the heap.
// If the item does not exist, it is added.
//
// Parameters:
//   - obj: The object to update or add
//
// Returns:
//   - error: Error if the key function fails
func (h *Heap) Update(obj any) error {
	return h.Add(obj)
}

// Delete removes an item from the heap.
//
// Parameters:
//   - obj: The object to remove
//
// Returns:
//   - error: Error if the key function fails or the object is not found
func (h *Heap) Delete(obj any) error {
	key, err := h.data.keyFunc(obj)
	if err != nil {
		return KeyError{Obj: obj, Err: err}
	}
	if item, ok := h.data.items[key]; ok {
		heap.Remove(h.data, item.index)
		return nil
	}
	return fmt.Errorf("heap: object not found")
}

// Peek returns the head of the heap without removing it.
//
// Returns:
//   - any: The head item, or nil if the heap is empty
func (h *Heap) Peek() any {
	return h.data.Peek()
}

// Pop returns the head of the heap and removes it.
//
// Returns:
//   - any: The head item
//   - error: Error if the heap is empty
func (h *Heap) Pop() (any, error) {
	obj := heap.Pop(h.data)
	if obj != nil {
		return obj, nil
	}
	return nil, fmt.Errorf("heap: object was removed from heap data")
}

// Get returns the requested item.
//
// Parameters:
//   - obj: The object to search for
//
// Returns:
//   - any: The found item
//   - bool: True if the item exists
//   - error: Error if the key function fails
func (h *Heap) Get(obj any) (any, bool, error) {
	key, err := h.data.keyFunc(obj)
	if err != nil {
		return nil, false, KeyError{Obj: obj, Err: err}
	}
	return h.GetByKey(key)
}

// GetByKey returns the requested item by its key.
//
// Parameters:
//   - key: The key of the item to retrieve
//
// Returns:
//   - any: The found item
//   - bool: True if the item exists
//   - error: Error if the key function fails
func (h *Heap) GetByKey(key string) (any, bool, error) {
	item, exists := h.data.items[key]
	if !exists {
		return nil, false, nil
	}
	return item.obj, true, nil
}

// List returns a list of all items in the heap.
//
// Returns:
//   - []any: A slice containing all items
func (h *Heap) List() []any {
	list := make([]any, 0, len(h.data.items))
	for _, item := range h.data.items {
		list = append(list, item.obj)
	}
	return list
}

// Len returns the number of items in the heap.
//
// Returns:
//   - int: The number of items
func (h *Heap) Len() int {
	return len(h.data.queue)
}

// New creates a new Heap instance.
//
// Parameters:
//   - keyFn: Function to extract keys from objects
//   - lessFn: Function to compare two objects
//
// Returns:
//   - *Heap: A new heap instance
func New(keyFn KeyFunc, lessFn lessFunc) *Heap {
	return &Heap{
		data: &data{
			items:    map[string]*heapItem{},
			queue:    []string{},
			keyFunc:  keyFn,
			lessFunc: lessFn,
		},
	}
}
