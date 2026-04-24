package listx

import "errors"

// ErrEmptyStack is returned when attempting to pop or peek from an empty stack.
var ErrEmptyStack = errors.New("stack is empty")

// ErrFullStack is returned when attempting to push onto a full stack.
var ErrFullStack = errors.New("stack is full")

// Stack is a last-in-first-out (LIFO) list.
type Stack interface {

	// Pop removes the object at the top of the stack and returns it.
	// Returns ErrEmptyStack if the stack is empty.
	//
	// Returns:
	//   - any: the object at the top of the stack.
	//   - error: nil on success, or ErrEmptyStack if the stack is empty.
	Pop() (any, error)

	// Push pushes an item onto the top of the stack.
	// Returns ErrFullStack if the stack is full.
	//
	// Parameters:
	//   - v: the item to push.
	//
	// Returns:
	//   - error: nil on success, or ErrFullStack if the stack is full.
	Push(v any) error

	// Peek looks at the object at the top of the stack without removing it.
	// Returns ErrEmptyStack if the stack is empty.
	//
	// Returns:
	//   - any: the object at the top of the stack.
	//   - error: nil on success, or ErrEmptyStack if the stack is empty.
	Peek() (any, error)

	// Empty tests whether the stack is empty.
	Empty()

	// Search returns the 1-based position where an object is in the stack.
	//
	// Parameters:
	//   - v: the object to search for.
	//
	// Returns:
	//   - int: the 1-based position, or 0 if the object is not found.
	Search(v any) int
}
