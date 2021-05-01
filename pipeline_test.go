package main

import (
	"context"
	"math"
	"time"

	"github.com/blokje5/go-pipelines/internals"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pipeline", func() {
	It("Should execute a simple pipeline and return the result correctly", func() {
		ctx := context.Background()
		res := From(SliceGenerator(1, 2, 3, 4)).
			Map(MapFunc(func(_ context.Context, x interface{}) interface{} {return x.(int) * x.(int)})).
			Reduce(ToSlice()).
			Run(ctx)
		
		Expect(res).To(Equal([]interface{}{1, 4, 9, 16}))
	})

	It("Should execute a simple summing pipeline and return the result correctly", func() {
		ctx := context.Background()
		res := From(SliceGenerator(1, 2, 3, 4)).
			Map(MapFunc(func(_ context.Context, x interface{}) interface{} {return x.(int) * x.(int)})).
			Reduce(Accumulate(0, func(a, b interface{}) interface{} { return a.(int) + b.(int)})).
			Run(ctx)
		
		Expect(res).To(Equal(30))
	})


	It("Should be preemtable", func() {
		internals.TestGoroutineClosure(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Microsecond)
			defer cancel()
			_ = From(Repeat(1, math.MaxInt32)).
				Map(MapFunc(func(_ context.Context, x interface{}) interface{} {return x.(int) * x.(int)})).
				Reduce(ToSlice()).
				Run(ctx)
		}, timeout)
	})
})