package barrier

import (
	"context"
	"sync"
)

// Group implements a reusable barrier similar to Java's CyclicBarrier.
// Create with NewGroup(parties, action). Each of the parties calls
// Wait(ctx). When the last party arrives, optional action is executed
// by that goroutine and all waiters are released. Wait supports context
// cancellation; if any waiter cancels, the barrier becomes broken and
// all waiters (except the canceling caller) receive *BrokenBarrierError.
// Reset breaks the current generation and starts a fresh one.
type Group struct {
	parties int
	action  func()

	mutex      sync.Mutex
	cond       *sync.Cond
	count      int   // remaining arrivals for current generation
	generation int64 // generation number
	broken     bool
}

// BrokenBarrierError indicates the barrier is broken.
type BrokenBarrierError struct{}

func (e *BrokenBarrierError) Error() string { return "barrier: barrier is broken" }

// NewGroup creates a new Group that waits for parties goroutines.
// barrierAction is executed by the last arriving goroutine. parties must be > 0.
func NewGroup(parties int, barrierAction func()) *Group {
	if parties <= 0 {
		panic("barrier: parties must be > 0")
	}
	g := &Group{
		parties: parties,
		action:  barrierAction,
		count:   parties,
	}
	g.cond = sync.NewCond(&g.mutex)
	return g
}

// Wait blocks until all parties have called Wait or until ctx is done.
// If ctx is canceled before release, the barrier is broken and all waiters
// receive *BrokenBarrierError (the caller that canceled receives ctx.Err()).
func (g *Group) Wait(ctx context.Context) error {
	g.mutex.Lock()
	// If already broken, return immediately
	if g.broken {
		g.mutex.Unlock()
		return &BrokenBarrierError{}
	}

	gen := g.generation

	g.count--
	if g.count == 0 {
		// Last arrival: run action and advance generation
		if g.action != nil {
			func() {
				defer func() {
					if r := recover(); r != nil {
						// action panicked -> break barrier and wake everyone
						g.broken = true
						g.cond.Broadcast()
					}
				}()
				g.action()
			}()
			if g.broken {
				g.mutex.Unlock()
				return &BrokenBarrierError{}
			}
		}
		g.nextGeneration()
		g.mutex.Unlock()
		return nil
	}

	// Not last: wait for generation change or broken state or ctx done.
	// We'll park this goroutine on the condition variable, but also respond to ctx cancellation.
	waitCh := make(chan struct{})

	go func(gen int64) {
		select {
		case <-ctx.Done():
			g.mutex.Lock()
			if g.generation == gen && !g.broken {
				g.broken = true
				g.cond.Broadcast()
			}
			g.mutex.Unlock()
		case <-waitCh:
			// released normally
		}
	}(gen)

	for gen == g.generation && !g.broken {
		g.cond.Wait()
	}

	// release helper goroutine
	close(waitCh)

	// If this goroutine's context was canceled, prefer returning ctx.Err()
	if err := ctx.Err(); err != nil {
		g.mutex.Unlock()
		return err
	}

	if g.broken {
		g.mutex.Unlock()
		return &BrokenBarrierError{}
	}

	g.mutex.Unlock()
	return nil
}

// Reset breaks the barrier and starts a fresh generation.
// All waiters in the current generation will observe a broken barrier.
func (g *Group) Reset() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.broken = true
	g.cond.Broadcast()
	g.nextGeneration()
}

func (g *Group) nextGeneration() {
	g.count = g.parties
	g.generation++
	g.broken = false
	g.cond.Broadcast()
}

// GetNumberWaiting returns how many goroutines are currently waiting.
func (g *Group) GetNumberWaiting() int {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.broken {
		return 0
	}
	return g.parties - g.count
}

// IsBroken returns whether the barrier is in a broken state.
func (g *Group) IsBroken() bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.broken
}

// GetParties returns the configured number of parties.
func (g *Group) GetParties() int { return g.parties }
