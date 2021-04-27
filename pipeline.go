package main

import "context"

// Pipeline is a builder for go channel pipelines. It is instantiated from a 
// generator with the From method. Pipeline allows multiple sequential calls to Map,
// which ensures the passed Transformers are executed in order.
// Reduce can be called to add a Reducer to the Pipeline. A reducer is mandatory
// before being able to call Run, which executes the pipeline.
type Pipeline interface {
	Map(transformer Transformer) Pipeline
	Reduce(reducer Reducer) Pipeline
	Run(ctx context.Context) interface{}
}

type pipeline struct {
	generator Generator
	transformers []Transformer
	reducer Reducer
}

// From returns a new Pipeline from the generator.
func From(generator Generator) Pipeline {
	p := &pipeline{}
	p.generator = generator
	return p
}

// Map configures a Transformer for the pipeline.
// Each Transformer is called in the order the Map functions where called.
func (p *pipeline) Map(transformer Transformer) Pipeline {
	p.transformers = append(p.transformers, transformer)
	return p
}

// Reduce passes a Reducer to the pipeline. This method must be called
// before Run is called.
func (p *pipeline) Reduce(reducer Reducer) Pipeline {
	p.reducer = reducer
	return p
}

// Run executes the pipeline. It expects Reduce has been called before.
// It returns the reduced results once the pipeline is closed.
func (p *pipeline) Run(ctx context.Context) interface{} {
	if p.reducer == nil {
		panic("No reducer for the pipeline, unable to Run")
	}

	c := p.generator.Generate(ctx)
	var prev <-chan interface{}
	prev = c
	for _, t := range p.transformers {
		prev = t.Transform(ctx, prev)
	}

	return p.reducer.Reduce(ctx, prev)
}