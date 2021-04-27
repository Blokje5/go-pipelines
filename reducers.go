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
	f := func(a, b interface{}) interface{} {
		cast := a.([]interface{})
		return append(cast, b)
	}

	r := Accumulate([]interface{}{}, f)

	return r
}

// MergeFunc is a function that merges a & b
type MergeFunc func(a, b interface{}) interface{}

// Accumulate creates a Reducer that accumulates the results of
// the channel into a result by applying MergeFunc f to each incoming message and
// merging it with the previous result. It will accumulate the result into v
func Accumulate(v interface{}, f MergeFunc) Reducer {
	r := ReducerFunc(func(ctx context.Context, channel <-chan interface{}) interface{} {
		var wg sync.WaitGroup
		wg.Add(1)
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
					v = f(v, msg)
				}
			}
		}()

		wg.Wait()

		return v
	})

	return r
}