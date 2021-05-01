package main

import (
	"context"
	"net/http"

	"github.com/blokje5/go-pipelines/errors"
	"github.com/blokje5/go-pipelines/internals"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("http", func() {
	var ts *ghttp.Server

	BeforeEach(func() {
		ts = ghttp.NewServer()
	})

	AfterEach(func() {
		ts.Close()
	})

	Describe("HTTPGenerator", func() {
		BeforeEach(func() {
			ts.AppendHandlers(
				ghttp.VerifyRequest("GET", "/test"),
				ghttp.VerifyRequest("GET", "/test"),
				ghttp.VerifyRequest("GET", "/test"),
			)
		})



		It("Should execute requests against a remote http server", func() {
			res := make([]interface{}, 0)
	
			internals.TestGoroutineClosure(func() {
				ctx := context.Background()
				var reqs []*http.Request
				for i := 0; i < 3; i++ {
					req, _ := http.NewRequest("GET", ts.URL() + "/test", nil)
					reqs = append(reqs, req)
				}

				g := HTTPGenerator(reqs...)
				c := g.Generate(ctx)
				
				for v := range c {
					res = append(res, v)
				}
			}, timeout)

			Expect(ts.ReceivedRequests()).Should(HaveLen(3))
			Expect(res).To(HaveLen(3))
			for _, v := range res {
				Expect(v).To(BeAssignableToTypeOf(errors.Optional{}))
				res := v.(errors.Optional).Val.(*http.Response)
				Expect(res).To(HaveHTTPStatus(200))
			}
		})
	})

	Describe("HTTPTransformer", func() {
		BeforeEach(func() {
			ts.AppendHandlers(
				ghttp.VerifyRequest("GET", "/test"),
				ghttp.VerifyRequest("GET", "/test"),
				ghttp.VerifyRequest("GET", "/test"),
			)
		})

		It("Should execute requests against a remote http server based on the channels on the request", func() {
			reqs := make(chan interface{}, 3)
			for i := 0; i < 3; i++ {
				req, _ := http.NewRequest("GET", ts.URL() + "/test", nil)
				reqs <- req
			}
			close(reqs)

			res := make([]interface{}, 0)

			internals.TestGoroutineClosure(func() {
				ctx := context.Background()
				g := HTTPTransformer()
				c := g.Transform(ctx, reqs)
				
				for v := range c {
					res = append(res, v)
				}
			}, timeout)

			Expect(ts.ReceivedRequests()).Should(HaveLen(3))
			Expect(res).To(HaveLen(3))
			for _, v := range res {
				Expect(v).To(BeAssignableToTypeOf(errors.Optional{}))
				res := v.(errors.Optional).Val.(*http.Response)
				Expect(res).To(HaveHTTPStatus(200))
			}
		})
	})
})
