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