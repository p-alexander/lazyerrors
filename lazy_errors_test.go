package lazyerrors

import (
	"errors"
	"fmt"
	"testing"
)

func TestCatchAllFunc(t *testing.T) {
	fmt.Println(testWrapper(CatchAllFunc, testFuncNoError))
	fmt.Println(testWrapper(CatchAllFunc, testFuncError))
	fmt.Println(testWrapper(CatchAllFunc, testFuncPanic))

	functions := []func() error{
		testFuncNoError,
		testFuncError,
		testFuncPanic,
	}

	var err error

	for _, f := range functions {
		func() {
			defer Catch(&err)
			Try(f())

			return
		}()

		fmt.Println(err)
	}
}

func TestCatchErrorFunc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			panic("this should panic")
		} else {
			fmt.Println("as expected:", r)
		}
	}()

	fmt.Println(testWrapper(CatchErrorFunc, testFuncNoError))
	fmt.Println(testWrapper(CatchErrorFunc, testFuncError))
	fmt.Println(testWrapper(CatchErrorFunc, testFuncPanic))
}

func BenchmarkCatchAllFunc(b *testing.B) {
	functions := []func() error{
		testFuncNoError,
		testFuncError,
		testFuncPanic,
	}

	var err error

	for i := 0; i < b.N; i++ {
		for _, f := range functions {
			func() {
				defer Catch(&err)
				Try(f())
			}()
		}
	}
}

func testWrapper(catcher func(*error), f func() error) (err error) {
	defer catcher(&err)
	Try(f())

	return
}

func testFuncPanic() error {
	panic("test panic")
}

func testFuncNoError() error {
	return nil
}

func testFuncError() error {
	return errors.New("test error")
}
