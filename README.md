# lazyerrors [![GoDoc](https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat)](https://godoc.org/github.com/p-alexander/lazyerrors)

Golang error handling via panic/recover approach.

import "github.com/p-alexander/lazyerrors"

Example:

```go
import (
    ...
    
    "github.com/p-alexander/lazyerrors"
)

// functions with return types that can be caught. 
functions := []func() error{
	func() error { return nil },
	func() error { return errors.New("some error") },
	func() error { panic("some panic") },
}

// direct usage.
someFunc := func() (err error) {
	// defer a Catch function at the beginning of the wrapping function.
	defer lazyerrors.Catch(&err)
	// execute the function that panics.
	lazyerrors.Try(functions[2]())

	return
}
// wrapping function will return an error instead of panic as Catch suppresses it by default.
fmt.Println(someFunc())

// anonymous function usage.
var err error

for _, f := range functions {
	// define the anonymous function to defer Catch in the necessary block of code.
	func() {
		defer lazyerrors.Catch(&err)
		lazyerrors.Try(f())
	}()
	// print the contents.
	fmt.Println(err)
}
```
