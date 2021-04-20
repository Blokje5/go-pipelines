package main

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/blokje5/go-pipelines/internals"
)

var _ = Describe("FuncTransformer", func() {
	var transformer Transformer

	BeforeEach(func() {
		f := func(msg interface{}) interface{} {
			i := msg.(int)
			return i*i
		}

		transformer = FuncTransformer(f)
	})
	It("Should transform a channel by applying f", func() {
		channel := make(chan interface{}, 3)
		for i := 0; i < 3; i++ {
			channel <- i+1
		}
		close(channel)

		res := make([]interface{}, 0)
		internals.TestGoroutineClosure(func() {
			ctx := context.Background()
			c := transformer.Transform(ctx, channel)
			for v := range c {
				res = append(res, v)
			}
		})

		Expect(res).To(Equal([]interface{}{1, 4, 9}))
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
			c := transformer.Transform(ctx, channel)
			internals.BlockUntilClose(c)
		}, timeout)
	})
})