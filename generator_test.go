package main

import (
	"context"
	"math"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/blokje5/go-pipelines/internals"
)

var _ = Describe("FromSlice", func() {
	It("Should return a valid generator from a slice", func() {
		s := []interface{}{1,2,3}
		res := make([]interface{}, 0)
	
		internals.TestGoroutineClosure(func() {
			ctx := context.Background()
			g := FromSlice(s...)
			c := g.Generate(ctx)
			
			for v := range c {
				res = append(res, v)
			}
		}, timeout)
		Expect(res).To(Equal(s))
	}, 1)
})

var _ = Describe("Repeat", func() {
	It("Should create a Generator that repeats the value v n times", func() {
		s := []interface{}{1,1,1}
		res := make([]interface{}, 0)

		internals.TestGoroutineClosure(func() {
			ctx := context.Background()
			g := Repeat(1, 3)
			c := g.Generate(ctx)
			
			for v := range c {
				res = append(res, v)
			}
		}, timeout)
		Expect(res).To(Equal(s))
	})

	It("Should be preemptable", func() {
		internals.TestGoroutineClosure(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Microsecond)
			defer cancel()
			g := Repeat(1, math.MaxInt32)
			c := g.Generate(ctx)
			// block until channel is closed, which should happen when timeout occurs
			for {
				select {
				case _, ok := <- c:
					if !ok { return }
				}
			}
		}, timeout)
	})
})