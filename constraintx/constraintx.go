package constraintx

import "golang.org/x/exp/constraints"

// Numeric is a type constraint that matches any numeric type (both integer and float types).
type Numeric interface {
	constraints.Integer | constraints.Float
}
