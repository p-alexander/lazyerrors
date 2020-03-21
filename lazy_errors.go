// Package lazyerrors - contains simple error handling via panic/recover approach.
//
// Purpose of this package is to show a somewhat dirty way to diminish endless chains of:
//
//     if err != nil {
//             return err
//     }
//
// Defer Catch at the beggining of your function and then check errors with Try.
//
//     func foo() (err error) {
//             defer Catch(&err)
//             Try(bar())
//
//             return
//     }
//
// Or put Try/Catch inside of your code with an anonymous function.
//
//     var err error
//
//     func() {
//             defer Catch(&err)
//             Try(bar())
//     }()
//
// As a result, you'll have 'return on error' behaviour as if standard approach was used.
//
// What happens inside of Try:
//     - On nil error execution will procede normally.
//     - On non-nil error it will be wrapped to show the caller, then it will start to panic until Catch.
//     - If an error was already wrapped, it won't be wrapped again to preserve the caller.
//
// Now about Catch:
//     - By default Catch can recover from any error or panic.
//     - Default behaviour can be changed through assigning Catch with another handler from a given set.
//     - If Catch recovers from a panic, it wraps recovered information into LazyErrorFromPanic.
package lazyerrors

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	// Catch - common catch handler, defaults to CatchAllFunc.
	Catch = CatchAllFunc
	// Try - common try handler, defaults to TryThrowFunc.
	Try = TryThrowFunc
	// ErrPanic - default error wrapped inside of LazyErrorFromPanic for Uwrap consistency.
	ErrPanic = errors.New("panic")
)

type (
	// LazyErrorWithCaller - custom error structure that contains caller information.
	LazyErrorWithCaller struct {
		Err    error
		Caller string
	}
	// LazyErrorFromPanic - custom error structure that contains recover information and stack trace.
	LazyErrorFromPanic struct {
		Recovered interface{}
		Stack     string
	}
)

// Error - error interface implementation.
func (e *LazyErrorWithCaller) Error() string {
	return e.Caller + e.Err.Error()
}

// Unwrap - error interface implementation (1.13).
func (e *LazyErrorWithCaller) Unwrap() error {
	return e.Err
}

// Error - error interface implementation.
func (e *LazyErrorFromPanic) Error() string {
	return fmt.Sprintf("[%v recovered]:\n%v\n[stack]:\n%s", ErrPanic, e.Recovered, e.Stack)
}

// Unwrap - error interface implementation (1.13).
func (e *LazyErrorFromPanic) Unwrap() error {
	return ErrPanic
}

// NewErrorWithCaller - adds caller information to error err and wraps it into LazyErrorWithCaller.
func NewErrorWithCaller(err error) error {
	return &LazyErrorWithCaller{
		Err:    err,
		Caller: caller(),
	}
}

// NewErrorFromPanic - wraps given recovered information and stack trace into LazyErrorFromPanic.
func NewErrorFromPanic(recovered interface{}, stack []byte) error {
	return &LazyErrorFromPanic{
		Recovered: recovered,
		Stack:     string(stack),
	}
}

// caller - returns a caller for ErrorWithCaller.
func caller() string {
	if _, file, line, ok := runtime.Caller(3); ok {
		return fmt.Sprintf("%s:%d: ", file, line)
	}

	return ""
}

// TryThrowFunc - throws non-nil error err.
func TryThrowFunc(err error) {
	if err != nil {
		switch err.(type) {
		// if an error is already wrapped, then return it as is.
		case *LazyErrorFromPanic, *LazyErrorWithCaller:
			panic(err)
		// else - wrap it into ErrorWithCaller.
		default:
			panic(NewErrorWithCaller(err))
		}
	}
}

// CatchLazyErrorFunc - catches only lazy errors.
func CatchLazyErrorFunc(ep *error) {
	if ep == nil {
		return
	}
	// recover from panic.
	if r := recover(); r != nil {
		// panic upon everything execept for LazyErrorFromPanic and LazyErrorWithCaller.
		switch t := r.(type) {
		case *LazyErrorFromPanic:
			*ep = t
		case *LazyErrorWithCaller:
			*ep = t
		default:
			panic(r)
		}
	}
}

// CatchErrorFunc - catches thrown error.
func CatchErrorFunc(ep *error) {
	if ep == nil {
		return
	}
	// recover from panic.
	if r := recover(); r != nil {
		// if an error was thrown, assign it through the pointer and return.
		if err, ok := r.(error); ok {
			*ep = err

			return
		}
		// else continue panicking.
		panic(r)
	}
}

// CatchAllFunc - catches thrown error or panic.
func CatchAllFunc(ep *error) {
	if ep == nil {
		return
	}
	// recover from panic.
	if r := recover(); r != nil {
		// if an error was thrown, assign it through the pointer and return.
		if err, ok := r.(error); ok {
			*ep = err

			return
		}
		// else wrap a panic info into an error.
		*ep = NewErrorFromPanic(r, debug.Stack())
	}
}
