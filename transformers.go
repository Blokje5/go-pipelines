package main

import "context"

// Transformer transforms the contents of a read-only channel, returning a new channel with the new contents.
type Transformer interface {
	Transform(ctx context.Context, channel <-chan interface{}) <-chan interface{}
}

// TranformFunc type is an adapter allowing any function adhering to the Transformer
// interface to be turned into a Transformer.
type TransformFunc func(ctx context.Context, channel <-chan interface{}) <-chan interface{}

func (f TransformFunc) Transform(ctx context.Context, channel <-chan interface{}) <-chan interface{} {
	return f(ctx, channel)
}

// Map: See FuncTransformer
func Map(f func(msg interface{}) interface{}) Transformer {
	return FuncTransformer(f)
}

// FuncTransformer (aka Map) takes a function and returns a Transformer
// which applies that function on each incoming msg on the read-only channel.
func FuncTransformer(f func(msg interface{}) interface{}) Transformer {
	f1 := TransformFunc(func(ctx context.Context, channel <-chan interface{}) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			for {
				select {
				case <- ctx.Done():
					return
				case msg, ok := <-channel:
					if !ok {
						return
					}
					select {
					case c<- f(msg):
					case <- ctx.Done():
						return
					}
				}
			}
		}()

		return c
	})

	return f1
}