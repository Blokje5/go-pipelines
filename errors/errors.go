package errors

// Optional is a type that helps with error propagation 
// between channels. Instead of creating a channel with values
// and having the channel producer handle any error, by propagating
// the error upstream, consumers can also decide what to do with the errors.
type Optional struct {
	Val interface{}
	Error error
}

func (o Optional) MustGetVal() interface{} {
	if o.Error != nil {
		panic(o.Error.Error())
	}

	return o.Val
}