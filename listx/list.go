package listx

// Element represents an element of a null terminated (non-circular) intrusive doubly linked list.
// It contains the key and value of the corresponding entry.
type Element[K comparable, V any] struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element[K, V]

	// The key that corresponds to this element in the ordered map.
	Key K

	// The value stored with this element.
	Value V
}

// Next returns the next list element.
//
// Returns:
//   - *Element[K, V]: the next element, or nil if there is none.
func (e *Element[K, V]) Next() *Element[K, V] {
	return e.next
}

// Prev returns the previous list element.
//
// Returns:
//   - *Element[K, V]: the previous element, or nil if there is none.
func (e *Element[K, V]) Prev() *Element[K, V] {
	return e.prev
}

// List represents a null terminated (non-circular) intrusive doubly linked list.
// The list is immediately usable after instantiation without dedicated initialization.
type List[K comparable, V any] struct {
	root Element[K, V] // list head and tail
}

// IsEmpty reports whether the list is empty.
//
// Returns:
//   - bool: true if the list contains no elements, false otherwise.
func (l *List[K, V]) IsEmpty() bool {
	return l.root.next == nil
}

// Front returns the first element of the list.
//
// Returns:
//   - *Element[K, V]: the first element, or nil if the list is empty.
func (l *List[K, V]) Front() *Element[K, V] {
	return l.root.next
}

// Back returns the last element of the list.
//
// Returns:
//   - *Element[K, V]: the last element, or nil if the list is empty.
func (l *List[K, V]) Back() *Element[K, V] {
	return l.root.prev
}

// Remove removes the specified element from the list.
//
// Parameters:
//   - e: the element to remove.
func (l *List[K, V]) Remove(e *Element[K, V]) {
	if e.prev == nil {
		l.root.next = e.next
	} else {
		e.prev.next = e.next
	}
	if e.next == nil {
		l.root.prev = e.prev
	} else {
		e.next.prev = e.prev
	}
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
}

// PushFront inserts a new element with the given key and value at the front of the list.
//
// Parameters:
//   - key: the key for the new element.
//   - value: the value for the new element.
//
// Returns:
//   - *Element[K, V]: the newly inserted element.
func (l *List[K, V]) PushFront(key K, value V) *Element[K, V] {
	e := &Element[K, V]{Key: key, Value: value}
	if l.root.next == nil {
		// It's the first element
		l.root.next = e
		l.root.prev = e
		return e
	}

	e.next = l.root.next
	l.root.next.prev = e
	l.root.next = e
	return e
}

// PushBack inserts a new element with the given key and value at the back of the list.
//
// Parameters:
//   - key: the key for the new element.
//   - value: the value for the new element.
//
// Returns:
//   - *Element[K, V]: the newly inserted element.
func (l *List[K, V]) PushBack(key K, value V) *Element[K, V] {
	e := &Element[K, V]{Key: key, Value: value}
	if l.root.prev == nil {
		// It's the first element
		l.root.next = e
		l.root.prev = e
		return e
	}

	e.prev = l.root.prev
	l.root.prev.next = e
	l.root.prev = e
	return e
}
