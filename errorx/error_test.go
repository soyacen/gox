package errorx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	v := Must(testString())
	assert.NotEmpty(t, v)

	// Test that Must panics when there's an error
	defer func() {
		if r := recover(); r == nil {
			t.Error("Must() should have panicked")
		}
	}()
	Must(testStringError())
}

func testString() (string, error) {
	return "testString", nil
}

func testStringError() (string, error) {
	return "", errors.New("error")
}
