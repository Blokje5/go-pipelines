package main

import (
	"context"
	"sync"
)

// Reducers take a channel and reduce it into an output.
// Most reducers will block until the channel is empty.
type Reducer interface {
	Reduce(ctx context.Context, channel <-chan interface{}) interface{}
}

// ReducerFunc type is an adapter allowing any function adhering to the Reducer
// interface to be turned into a Reducer.
type ReducerFunc func(ctx context.Context, channel <-chan interface{}) interface{}

func (r ReducerFunc) Reduce(ctx context.Context, channel <-chan interface{}) interface{} {
	return r(ctx, channel)
}

// ToSlice returns a reducer that will collect the channel into a slice and return.
// It will block until the channel is closed or the context is closed.
func ToSlice() Reducer {
	f := ReducerFunc(func(ctx context.Context, channel <-chan interface{}) interface{} {
		var wg sync.WaitGroup
		wg.Add(1)
		var res []interface{}
		go func() {
			defer wg.Done()
			for {
				select {
				case <- ctx.Done():
					return
				case msg, ok := <-channel:
					if !ok {
						return
					}
					res = append(res, msg)
				}
			}
		}()

		wg.Wait()

		return res
	})

	return f
}