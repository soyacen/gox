package errorx

import (
	"errors"
	"reflect"
	"testing"
)

func TestEquals(t *testing.T) {
	innerErr := errors.New("inner")
	testCases := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			name:     "both nil",
			err:      nil,
			target:   nil,
			expected: true,
		},
		{
			name:     "one nil one not",
			err:      nil,
			target:   errors.New("test"),
			expected: false,
		},
		{
			name:     "same error",
			err:      errors.New("test"),
			target:   errors.New("test"),
			expected: true,
		},
		{
			name:     "different errors",
			err:      errors.New("test 1"),
			target:   errors.New("test 2"),
			expected: false,
		},
		{
			name:     "wrapped error",
			err:      Join(innerErr),
			target:   innerErr,
			expected: true, // because errors.Is should find it
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Equals(tc.err, tc.target)
			if result != tc.expected {
				t.Errorf("Equals() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		errs     []error
		expected int // number of errors in the joined result
	}{
		{
			name:     "nil error",
			err:      nil,
			errs:     []error{errors.New("1"), errors.New("2")},
			expected: 2,
		},
		{
			name:     "single error",
			err:      errors.New("1"),
			errs:     []error{},
			expected: 1,
		},
		{
			name:     "multiple errors",
			err:      errors.New("1"),
			errs:     []error{errors.New("2"), errors.New("3")},
			expected: 3,
		},
		{
			name:     "already joined error",
			err:      Join(errors.New("1"), errors.New("2")),
			errs:     []error{errors.New("3")},
			expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Join(tc.err, tc.errs...)
			unwrapped := UnwrapMultiErr(result)
			if len(unwrapped) != tc.expected {
				t.Errorf("Join() resulted in %d errors, want %d", len(unwrapped), tc.expected)
			}
		})
	}
}

func TestUnwrapMultiErr(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: 0,
		},
		{
			name:     "single error",
			err:      errors.New("test"),
			expected: 1,
		},
		{
			name:     "joined errors",
			err:      Join(errors.New("1"), errors.New("2")),
			expected: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := UnwrapMultiErr(tc.err)
			if len(result) != tc.expected {
				t.Errorf("UnwrapMultiErr() = %d errors, want %d", len(result), tc.expected)
			}
		})
	}
}

func TestErrorType(t *testing.T) {
	if ErrorType == nil {
		t.Error("ErrorType is nil")
	}

	// Check that ErrorType is the same as reflect.TypeOf((*error)(nil)).Elem()
	err := errors.New("test")
	if !reflect.TypeOf(err).Implements(ErrorType) {
		t.Error("ErrorType should be implemented by error")
	}
}

func TestIgnore(t *testing.T) {
	v, err := testString()
	result := Ignore(v, err)
	if result != v {
		t.Errorf("Ignore() = %v, want %v", result, v)
	}

	// With error
	v, err = testStringError()
	result = Ignore(v, err)
	if result != v {
		t.Errorf("Ignore() with error = %v, want %v", result, v)
	}
}

func TestConcern(t *testing.T) {
	v, err := testString()
	result := Concern(v, err)
	if result != err {
		t.Errorf("Concern() = %v, want %v", result, err)
	}

	// With error
	v, err = testStringError()
	result = Concern(v, err)
	if result != err {
		t.Errorf("Concern() with error = %v, want %v", result, err)
	}
}

func TestQuiet(t *testing.T) {
	// Just test that it doesn't panic
	Quiet(nil)
	Quiet(errors.New("test"))
}

func TestStringfy(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "",
		},
		{
			name:     "non-nil error",
			err:      errors.New("test"),
			expected: "test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Stringfy(tc.err)
			if result != tc.expected {
				t.Errorf("Stringfy() = %v, want %v", result, tc.expected)
			}
		})
	}
}
