# lazyerrors [![GoDoc](https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat)](https://godoc.org/github.com/p-alexander/lazyerrors)

Golang error handling via panic/recover approach.

import "github.com/p-alexander/lazyerrors"

Example:

```go
import (
    ...
    
    "github.com/p-alexander/lazyerrors"
)

...
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
// wrapping function will return an error instead of panic as Catch suppresses it by default.
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
```
