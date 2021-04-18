package main

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ToSlice", func() {
	It("Should reduce a channel into a slice", func() {
		done := make(chan interface{})
		channel := make(chan interface{}, 3)
		for i := 0; i < 3; i++ {
			channel <- i+1
		}
		close(channel)

		res := make([]interface{}, 0)
	
		go func() {
			defer close(done)
			ctx := context.Background()
			g := ToSlice()
			res = g.Reduce(ctx, channel).([]interface{})
		}()

		expected := []interface{}{1,2,3}
		Eventually(done, timeout).Should(BeClosed())
		Expect(res).To(Equal(expected))
	})

	It("Should be preemtable", func() {
		done := make(chan interface{})
		channel := make(chan interface{}, 3)
		for i := 0; i < 3; i++ {
			channel <- i+1
		}
		defer close(channel)
		// never close channel before this function ends
	
		go func() {
			defer close(done)
			ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Microsecond)
			defer cancel()
			g := ToSlice()
			_ = g.Reduce(ctx, channel).([]interface{})
		}()
		
		Eventually(done, timeout).Should(BeClosed())
	})
})