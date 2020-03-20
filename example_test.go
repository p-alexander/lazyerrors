package lazyerrors_test

import (
	"errors"
	"fmt"

	"github.com/p-alexander/lazyerrors"
)

// Example - example of error handling via Try and Catch.
func Example() {
	functions := []func() error{
		func() error { return nil },
		func() error { return errors.New("some error") },
		func() error { panic("some panic") },
	}
	// direct usage.
	someFunc := func() (err error) {
		// defer a Catch function at the beginning of the wrapping function.
		defer lazyerrors.Catch(&err)
		// execute function that panics.
		lazyerrors.Try(functions[2]())

		return
	}
	// upon executing wrapping function it will return an error instead of panic as Catch suppresses panic by default.
	fmt.Println(someFunc())
	// anonymous usage.
	var err error
	// to catch errors inline anonymous function can be used to defer the Catch function in the necessary block of code.
	for _, f := range functions {
		func() {
			defer lazyerrors.Catch(&err)
			lazyerrors.Try(f())
		}()
		// print the contents.
		fmt.Println(err)
	}
}
