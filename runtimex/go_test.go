package runtimex

import (
	"testing"
)

func TestStack(t *testing.T) {
	stackData := Stack(0)
	t.Log(string(stackData))
}
