package listx

import (
	"sync/atomic"
	"unsafe"
)

// LockFreeQueue is a lock-free queue implementation based on a linked list.
type LockFreeQueue struct {
	// headPtr is the head pointer of the queue.
	headPtr unsafe.Pointer
	// tailPtr is the tail pointer of the queue.
	tailPtr unsafe.Pointer
}

// lockFreeElement represents a node in the linked list.
type lockFreeElement struct {
	value   any
	nextPtr unsafe.Pointer
}

// NewLockFreeQueue creates and returns a new lock-free queue.
//
// Reference: https://www.sobyte.net/post/2021-07/implementing-lock-free-queues-with-go/
//
// Returns:
//   - *LockFreeQueue: a new empty lock-free queue.
func NewLockFreeQueue() *LockFreeQueue {
	elemPtr := unsafe.Pointer(&lockFreeElement{})
	return &LockFreeQueue{headPtr: elemPtr, tailPtr: elemPtr}
}

// Enqueue adds a value to the tail of the queue.
//
// Parameters:
//   - v: the value to enqueue.
func (q *LockFreeQueue) Enqueue(v any) {
	n := &lockFreeElement{value: v}
	for {
		tail := load(&q.tailPtr)
		next := load(&tail.nextPtr)
		// tail is still tail
		if tail == load(&q.tailPtr) {
			// no new element has been enqueued yet
			if next == nil {
				// add to the tail
				if cas(&tail.nextPtr, next, n) {
					// enqueue succeeded, move tail pointer
					_ = cas(&q.tailPtr, tail, n)
					return
				}
			} else {
				// a new element has been added, need to move tail pointer
				_ = cas(&q.tailPtr, tail, next)
			}
		}
	}
}

// Dequeue removes and returns the value at the head of the queue.
// Returns nil if the queue is empty.
//
// Returns:
//   - any: the dequeued value, or nil if the queue is empty.
func (q *LockFreeQueue) Dequeue() any {
	for {
		head := load(&q.headPtr)
		tail := load(&q.tailPtr)
		next := load(&head.nextPtr)
		// head is still head
		if head == load(&q.headPtr) {
			// head and tail are the same
			if head == tail {
				// empty queue
				if next == nil {
					return nil
				}
				// tail pointer hasn't been adjusted yet, try to move it
				_ = cas(&q.tailPtr, tail, next)
			} else {
				// read the value to dequeue
				v := next.value
				// move head pointer to next
				if cas(&q.headPtr, head, next) {
					// Dequeue is done.
					return v
				}
			}
		}
	}
}

func load(p *unsafe.Pointer) *lockFreeElement {
	return (*lockFreeElement)(atomic.LoadPointer(p))
}

func cas(p *unsafe.Pointer, old, new *lockFreeElement) bool {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}
