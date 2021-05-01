package main

import (
	"github.com/blokje5/go-pipelines/errors"
	"context"
	"net/http"
)

// HTTPGenerator returns a generator that executes a slice of requests.
// For each request it returns an Optional with the error or the http.Response.
func HTTPGenerator(reqs ...*http.Request) Generator {
	f := GeneratorFunc(func(ctx context.Context) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			for _, req := range reqs {
				select {
				case <- ctx.Done():
					return
				case c <- executeRequestWithContext(ctx, req):
				}
			}
		}()


		return c
	})

	return f
}

// HTTPTransformer takes a channel of *http.Request and executes each request.
// For each request it returns an Optional with the error or the http.Response.
func HTTPTransformer() Transformer {
	f := MapFunc(func(ctx context.Context, msg interface{}) interface{} {
		req := msg.(*http.Request)
		return executeRequestWithContext(ctx, req) 
	})

	return f
}

func executeRequestWithContext(ctx context.Context, req *http.Request) errors.Optional {
	req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	return errors.Optional{
		Val: res,
		Error: err,
	}
}