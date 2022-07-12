// Code generated by "go-option -type Number"; DO NOT EDIT.
// Install go-option by "go get -u github.com/searKing/golang/tools/go-option"

package main

// A NumberOption sets options.
type NumberOption interface {
	apply(*Number)
}

// EmptyNumberOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyNumberOption struct{}

func (EmptyNumberOption) apply(*Number) {}

// NumberOptionFunc wraps a function that modifies Number into an
// implementation of the NumberOption interface.
type NumberOptionFunc func(*Number)

func (f NumberOptionFunc) apply(do *Number) {
	f(do)
}
func (o *Number) ApplyOptions(options ...NumberOption) *Number {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}

// WithNumberArrayType sets arrayType in Number.
func WithNumberArrayType(v [5]int64) NumberOption {
	return NumberOptionFunc(func(o *Number) {
		o.arrayType = v
	})
}

// WithNumberInterfaceType sets interfaceType in Number.
func WithNumberInterfaceType(v interface{}) NumberOption {
	return NumberOptionFunc(func(o *Number) {
		o.interfaceType = v
	})
}

// WithNumberMapType appends mapType in Number.
func WithNumberMapType(m map[string]int64) NumberOption {
	return NumberOptionFunc(func(o *Number) {
		if o.mapType == nil {
			o.mapType = m
			return
		}
		for k, v := range m {
			o.mapType[k] = v
		}
	})
}

// WithNumberMapTypeReplace sets mapType in Number.
func WithNumberMapTypeReplace(v map[string]int64) NumberOption {
	return NumberOptionFunc(func(o *Number) {
		o.mapType = v
	})
}

// WithNumberSliceType appends sliceType in Number.
func WithNumberSliceType(v ...int64) NumberOption {
	return NumberOptionFunc(func(o *Number) {
		o.sliceType = append(o.sliceType, v...)
	})
}

// WithNumberSliceTypeReplace sets sliceType in Number.
func WithNumberSliceTypeReplace(v ...int64) NumberOption {
	return NumberOptionFunc(func(o *Number) {
		o.sliceType = v
	})
}

// WithNumberName sets name in Number.
func WithNumberName(v string) NumberOption {
	return NumberOptionFunc(func(o *Number) {
		o.name = v
	})
}
