package main

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/blokje5/go-pipelines/internals"
)

var _ = Describe("ToSlice", func() {
	It("Should reduce a channel into a slice", func() {
		channel := make(chan interface{}, 3)
		for i := 0; i < 3; i++ {
			channel <- i+1
		}
		close(channel)

		res := make([]interface{}, 0)
		internals.TestGoroutineClosure(func() {
			ctx := context.Background()
			g := ToSlice()
			res = g.Reduce(ctx, channel).([]interface{})
		}, timeout)

		expected := []interface{}{1,2,3}
		Expect(res).To(Equal(expected))
	})

	It("Should be preemtable", func() {
		channel := make(chan interface{}, 3)
		for i := 0; i < 3; i++ {
			channel <- i+1
		}
		defer close(channel)
		// never close channel before this function ends
	
		internals.TestGoroutineClosure(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Microsecond)
			defer cancel()
			g := ToSlice()
			_ = g.Reduce(ctx, channel).([]interface{})
		}, timeout)
	})
})