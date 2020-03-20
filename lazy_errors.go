// Package lazyerrors - contains simple error handling via panic/recover approach.
//
// Purpose of this package is to show a somewhat dirty way to diminish endless chains of:
//
// if err != nil {
//	 return err
// }
package lazyerrors

import "fmt"

var (
	// Catch - common catch handler, defaults to CatchAllFunc.
	Catch = CatchAllFunc
	// Try - common try handler, defaults to TryThrowFunc.
	Try = TryThrowFunc
)

// TryThrowFunc - throws non-nil error err.
func TryThrowFunc(err error) {
	if err != nil {
		panic(err)
	}
}

// CatchErrorFunc - catches thrown error.
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

// CatchAllFunc - catches thrown error or panic.
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
