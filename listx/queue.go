package listx

import "errors"

// ErrEmptyQueue is returned when attempting to dequeue or peek from an empty queue.
var ErrEmptyQueue = errors.New("queue is empty")

// ErrFullQueue is returned when attempting to enqueue into a full queue.
var ErrFullQueue = errors.New("queue is full")

// Queue is a first-in-first-out (FIFO) list.
type Queue interface {

	// Enqueue inserts the specified element into the queue.
	// Returns ErrFullQueue if the queue is full.
	//
	// Parameters:
	//   - v: the element to insert.
	//
	// Returns:
	//   - error: nil on success, or ErrFullQueue if the queue is full.
	Enqueue(v any) error

	// Dequeue retrieves and removes the head of the queue.
	// Returns ErrEmptyQueue if the queue is empty.
	//
	// Returns:
	//   - any: the dequeued element.
	//   - error: nil on success, or ErrEmptyQueue if the queue is empty.
	Dequeue() (any, error)

	// Peek retrieves but does not remove the head of the queue.
	// Returns ErrEmptyQueue if the queue is empty.
	//
	// Returns:
	//   - any: the head element.
	//   - error: nil on success, or ErrEmptyQueue if the queue is empty.
	Peek() (any, error)

	// Empty tests whether the queue is empty.
	Empty()

	// Search returns the 1-based position where an object is in the queue.
	//
	// Parameters:
	//   - v: the object to search for.
	//
	// Returns:
	//   - int: the 1-based position, or 0 if the object is not found.
	Search(v any) int
}
