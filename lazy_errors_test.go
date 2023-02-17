package lazyerrors

import (
	"errors"
	"fmt"
	"testing"
)

func TestCatchAllWithStackFunc(t *testing.T) {
	if err := testWrapper(TryWrapErrorFunc, CatchAllWithStackFunc, testFuncNoError); err != nil {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}

	if err := testWrapper(TryWrapErrorFunc, CatchAllWithStackFunc, testFuncError); err == nil {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}

	if err := testWrapper(TryWrapErrorFunc, CatchAllWithStackFunc, testFuncPanic); err == nil {
		t.Fatal("unexpected:", err)
	} else {
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

	if err := testWrapper(TryWrapErrorFunc, CatchErrorFunc, testFuncNoError); err != nil {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}

	if err := testWrapper(TryWrapErrorFunc, CatchErrorFunc, testFuncError); err == nil {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}
	// this should panic.
	if err := testWrapper(TryWrapErrorFunc, CatchErrorFunc, testFuncPanic); err != nil {
		t.Fatal("unexpected:", err)
	}
}

func TestCatchLazyErrorFunc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			panic("this should panic")
		} else {
			fmt.Println("as expected:", r)
		}
	}()

	if err := testWrapper(TryWrapErrorFunc, CatchLazyErrorFunc, testFuncNoError); err != nil {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}

	if err := testWrapper(TryWrapErrorFunc, CatchLazyErrorFunc, testFuncError); err == nil {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}
	// this should panic.
	if err := testWrapper(TryWrapErrorFunc, CatchLazyErrorFunc, testFuncPanic); err != nil {
		t.Fatal("unexpected:", err)
	}
}

func BenchmarkDefault(b *testing.B) {
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

func BenchmarkMinimal(b *testing.B) {
	functions := []func() error{
		testFuncNoError,
		testFuncError,
		testFuncPanic,
	}

	var err error

	for i := 0; i < b.N; i++ {
		for _, f := range functions {
			func() {
				defer CatchAllFunc(&err)
				TryErrorFunc(f())
			}()
		}
	}
}

func TestNestedError(t *testing.T) {
	functions := []func() error{
		testFuncNoError,
		testFuncError,
		testFuncPanic,
	}

	var err error

	for _, f := range functions {
		func() {
			defer Catch(&err)
			Try(testWrapper(Try, Catch, f))
		}()

		fmt.Printf("wrapped:\n%v\noriginal:\n%v\n\n", err, errors.Unwrap(err))
	}
}

func TestTryErrorFunc(t *testing.T) {
	if err := testWrapper(TryErrorFunc, CatchAllFunc, testFuncNoError); err != nil {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}

	if err := testWrapper(TryErrorFunc, CatchAllFunc, testFuncError); err == nil {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}

	if err := testWrapper(TryErrorFunc, CatchAllFunc, testFuncPanic); err == nil {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}
}

func TestCatchAllFunc(t *testing.T) {
	if err := testWrapper(TryErrorFunc, CatchAllFunc, testFuncNoError); err != nil {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}

	if err := testWrapper(TryErrorFunc, CatchAllFunc, testFuncError); err == nil {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}

	if err := testWrapper(TryErrorFunc, CatchAllFunc, testFuncPanic); err == nil {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}
}

func TestErrorIs(t *testing.T) {
	customError := errors.New("test error")

	f := func() error {
		return customError
	}

	if err := testWrapper(TryErrorFunc, CatchAllFunc, f); !errors.Is(err, customError) {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}

	if err := testWrapper(TryWrapErrorFunc, CatchAllFunc, f); !errors.Is(err, customError) {
		t.Fatal("unexpected:", err)
	} else {
		fmt.Println(err)
	}
}

func testWrapper(tryFunc func(error), catchFunc func(*error), f func() error) (err error) {
	defer catchFunc(&err)
	tryFunc(f())

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
