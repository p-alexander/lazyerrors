// Package lazyerrors - contains simple error handling via panic/recover approach.
package lazyerrors

import "fmt"

var (
	// Catch - common catch handler, defaults to CatchAllFunc.
	Catch = CatchAllFunc
	// Try - common try handler, defaults to TryThrowFunc.
	Try = TryThrowFunc
)

// TryThrowFunc - throws an error if any.
func TryThrowFunc(err error) {
	if err != nil {
		panic(err)
	}
}

// CatchErrorFunc - catches thrown errors.
func CatchErrorFunc(ep *error) {
	if ep == nil {
		return
	}

	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			*ep = err

			return
		}

		panic(r)
	}
}

// CatchAllFunc - catches thrown errors as well as panics.
func CatchAllFunc(ep *error) {
	if ep == nil {
		return
	}

	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			*ep = err

			return
		}

		*ep = fmt.Errorf("panic: %v", r)
	}
}
