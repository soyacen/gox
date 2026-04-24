package operator

// Pointer stores v in a new value of type E and returns a pointer to it.
//
// Parameters:
//   - v: The value to store
//
// Returns:
//   - *E: A pointer to the stored value
func Pointer[E any](v E) *E {
	return &v
}

// Indirect returns the value pointed to by p.
// If p is nil, it returns the zero value of E.
//
// Parameters:
//   - p: The pointer to dereference
//
// Returns:
//   - E: The value pointed to by p, or the zero value if p is nil
func Indirect[P *E, E any](p P) E {
	if p == nil {
		var e E
		return e
	}
	return *p
}
