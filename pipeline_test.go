package main

import (
	"context"
	"time"

	"github.com/blokje5/go-pipelines/internals"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pipeline", func() {
	It("Should execute a simple pipeline and return the result correctly", func() {
		ctx := context.Background()
		res := From(FromSlice(1, 2, 3, 4)).
			Map(Map(func(x interface{}) interface{} {return x.(int) * x.(int)})).
			Reduce(ToSlice()).
			Run(ctx)
		
		Expect(res).To(Equal([]interface{}{1, 4, 9, 16}))
	})

	It("Should be preemtable", func() {
		internals.TestGoroutineClosure(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Microsecond)
			defer cancel()
			_ = From(FromSlice(1, 2, 3, 4)).
				Map(Map(func(x interface{}) interface{} {return x.(int) * x.(int)})).
				Reduce(ToSlice()).
				Run(ctx)
		}, timeout)
	})
})