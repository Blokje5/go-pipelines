package internals

import (
	. "github.com/onsi/gomega"
)

// TestGoroutineClosure runs the given func f in a goroutine and checks
// whether it finished in the given interval using Eventually.
func TestGoroutineClosure(f func(), intervals ...interface{}) {
	done := make(chan interface{})

	go func() {
		defer close(done)
		f()
	}()

	Eventually(done, intervals...).Should(BeClosed())
}

// BlockUntilClose blocks the goroutine until the channel is closed.
// This is used in testing preemptability.
func BlockUntilClose(c <-chan interface{}) {
	for {
		select {
		case _, ok := <- c:
			if !ok { return }
		}
	}
}