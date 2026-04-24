package signalx

import (
	"os"
	"os/signal"
)

// AsyncWatch watches for the specified signals asynchronously.
// It returns a channel that will receive the first signal that is caught.
// The channel will be closed after the first signal is received.
//
// Parameters:
//   - sig: The signals to watch for.
//
// Returns:
//   - <-chan os.Signal: A channel that receives the caught signal.
func AsyncWatch(sig ...os.Signal) <-chan os.Signal {
	signalC := make(chan os.Signal)
	go func() {
		signalC <- SyncWatch(sig...)
		close(signalC)
	}()
	return signalC
}

// SyncWatch watches for the specified signals synchronously.
// It blocks until one of the specified signals is received.
//
// Parameters:
//   - sig: The signals to watch for.
//
// Returns:
//   - os.Signal: The signal that was received.
func SyncWatch(sig ...os.Signal) os.Signal {
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, sig...)
	incomingSignal := <-signalC
	signal.Stop(signalC)
	close(signalC)
	return incomingSignal
}
