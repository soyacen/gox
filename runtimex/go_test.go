package runtimex

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	stackData := Stack(0)
	t.Log(string(stackData))
}

func TestLine(t *testing.T) {
	fmt.Println(Line())
}
