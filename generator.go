package main

import "context"

// Generator generates a read-only channel. This is the starting point
// of a pipeline. For example, it can be a simple slice turned into a channel,
// or it could be a function polling an API and turning each response in a message on the channel.
type Generator interface {
	Generate(ctx context.Context) <-chan interface{}
}

// GeneratorFunc type is an adapter allowing any function adhering to the Generator
// interface to be turned into a Generator.
type GeneratorFunc func(ctx context.Context) <-chan interface{}

func (g GeneratorFunc) Generate(ctx context.Context) <-chan interface{} {
	return g(ctx)
}

// SliceGenerator creates a Generator from a slice.
// Each value of the slice will be send over the channel.
func SliceGenerator(s ...interface{}) Generator {
	f := GeneratorFunc(func(ctx context.Context) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			for _, v := range s {
				select {
				case <- ctx.Done():
					return
				case c <- v:
				}
			}
		}()


		return c
	})

	return f
}

// Repeat creates a Generator that will repeat the value v n times.
func Repeat(v interface{}, n int) Generator {
	f := GeneratorFunc(func(ctx context.Context) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			for i := 0; i < n; i++ {
				select {
				case <- ctx.Done():
					return
				case c <- v:
				}
			}
		}()

		return c
	})

	return f
}