# lazyerrors [![GoDoc](https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat)](https://godoc.org/github.com/p-alexander/lazyerrors)

Golang error handling via panic/recover approach.

## Summary:

Purpose of this package is to show a way to diminish endless chains of:

```go
if err != nil {
	return err
}
```

## Download:

`go get github.com/p-alexander/lazyerrors`

## Usage:

Import:

```go
import "github.com/p-alexander/lazyerrors"
```

 Defer Catch at the beggining of your function and then check errors with Try.

```go
     func foo() (err error) {
             defer lazyerrors.Catch(&err)
             lazyerrors.Try(bar())

             return
     }
```

 Or put Try/Catch inside of your code with an anonymous function.

```go
     var err error

     func() {
             defer lazyerrors.Catch(&err)
             lazyerrors.Try(bar())
     }()
```

 As a result, you'll have 'return on error' behaviour as if standart approach was used.

 What happens inside of Try:
 - On nil error execution will procede normally.
 - On non-nil error it will be wrapped to show the caller, then it will start to panic until Catch.
 - If an error was already wrapped, it won't be wrapped again to preserve the caller.

 Now about Catch:
 - By default Catch can recover from any error or panic.
 - Default behaviour can be changed through assigning Catch with another handler from a given set.
 - If Catch recovers from a panic, it wraps recovered information into LazyErrorFromPanic.
