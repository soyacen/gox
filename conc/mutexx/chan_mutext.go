package mutexx

import (
	"context"
	"errors"
)

type ChanMutex struct {
	state chan struct{}
}

func NewChanMutex() *ChanMutex {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return &ChanMutex{
		state: ch,
	}
}

func (m *ChanMutex) Lock() {
	<-m.state
}

func (m *ChanMutex) Unlock() {
	select {
	case m.state <- struct{}{}:
	default:
		panic(errors.New("mutexx: unlock of unlocked mutex"))
	}

}

func (m *ChanMutex) TryLock() bool {
	select {
	case <-m.state:
		return true
	default:
		return false
	}
}

func (m *ChanMutex) LockContext(ctx context.Context) bool {
	select {
	case <-m.state:
		return true
	case <-ctx.Done():
		return false
	}
}

func (m *ChanMutex) IsLocked() bool {
	return len(m.state) == 0
}
