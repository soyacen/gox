package mutexx

import (
	"testing"
	"time"
)

const (
	callbackTimeout = 1 * time.Second
)

func newGroupMutexes() []*GroupMutex {
	return []*GroupMutex{
		{N: 0},
		{N: 1},
		{N: 2},
		{N: 3},
	}
}

func Test_SingleLock_NoUnlock(t *testing.T) {
	for _, gm := range newGroupMutexes() {
		// Arrange
		key := "fakeid"
		callbackCh := make(chan interface{})

		// Act
		go lockAndCallback(gm, key, callbackCh)

		// Assert
		verifyCallbackHappens(t, callbackCh)
	}
}

func Test_SingleLock_SingleUnlock(t *testing.T) {
	for _, gm := range newGroupMutexes() {
		// Arrange
		key := "fakeid"
		callbackCh := make(chan interface{})

		// Act & Assert
		go lockAndCallback(gm, key, callbackCh)
		verifyCallbackHappens(t, callbackCh)
		gm.Unlock(key)
	}
}

func Test_DoubleLock_DoubleUnlock(t *testing.T) {
	for _, gm := range newGroupMutexes() {
		// Arrange
		key := "fakeid"
		callbackCh1stLock := make(chan interface{})
		callbackCh2ndLock := make(chan interface{})

		// Act & Assert
		go lockAndCallback(gm, key, callbackCh1stLock)
		verifyCallbackHappens(t, callbackCh1stLock)
		go lockAndCallback(gm, key, callbackCh2ndLock)
		verifyCallbackDoesntHappens(t, callbackCh2ndLock)
		gm.Unlock(key)
		verifyCallbackHappens(t, callbackCh2ndLock)
		gm.Unlock(key)
	}
}

func lockAndCallback(gm *GroupMutex, id string, callbackCh chan<- interface{}) {
	gm.Lock(id)
	callbackCh <- true
}

func verifyCallbackHappens(t *testing.T, callbackCh <-chan interface{}) bool {
	select {
	case <-callbackCh:
		return true
	case <-time.After(callbackTimeout):
		t.Fatalf("Timed out waiting for callback.")
		return false
	}
}

func verifyCallbackDoesntHappens(t *testing.T, callbackCh <-chan interface{}) bool {
	select {
	case <-callbackCh:
		t.Fatalf("Unexpected callback.")
		return false
	case <-time.After(callbackTimeout):
		return true
	}
}
