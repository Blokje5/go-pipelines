package main

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FromSlice", func() {
	It("Should return a valid generator from a slice", func() {
		done := make(chan interface{})
	
		s := []interface{}{1,2,3}
		res := make([]interface{}, 0)
	
		go func() {
			defer close(done)
			ctx := context.Background()
			g := FromSlice(s...)
			c := g.Generate(ctx)
			
			for v := range c {
				res = append(res, v)
			}
		}()

		Eventually(done, timeout).Should(BeClosed())
		Expect(res).To(Equal(s))
	}, 1)
})

var _ = Describe("Repeat", func() {
	It("Should create a Generator that repeats the value v n times", func() {
		done := make(chan interface{})
	
		s := []interface{}{1,1,1}
		res := make([]interface{}, 0)


		go func() {
			defer close(done)
			ctx := context.Background()
			g := Repeat(1, 3)
			c := g.Generate(ctx)
			
			for v := range c {
				res = append(res, v)
			}
		}()

		Eventually(done, timeout).Should(BeClosed())
		Expect(res).To(Equal(s))
	})

	It("Should be preemptable", func() {
		done := make(chan interface{})

		go func() {
			defer close(done)
			ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Microsecond)
			defer cancel()
			g := Repeat(1, 3)
			c := g.Generate(ctx)
			<-c
		}()

		Eventually(done, timeout).Should(BeClosed())
	})
})