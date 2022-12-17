//go:build js && wasm && !tiny
// +build js,wasm,!tiny

package main

import (
	"context"
	"syscall/js"
	"time"
)

func multiply(this js.Value, args []js.Value) interface{} {
	if nargs := len(args); nargs != 2 {
		return js.ValueOf(0)
	}
	x := args[0].Int()
	y := args[1].Int()
	return js.ValueOf(x * y)
}

// This calls a JS function from Go.
func main() {
	done := make(chan struct{}, 0)
	// expose functions to js.
	js.Global().Set("multiply", js.FuncOf(multiply))
	// call js function from go
	println("adding promiss two numbers:", addPromissWrap(2, 3))
	println("adding async two numbers:", addAsyncWrap(2, 3))

	<-done
}

// https://stackoverflow.com/a/68427221
func await(awaitable js.Value) ([]js.Value, []js.Value) {
	then := make(chan []js.Value)
	defer close(then)
	thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		then <- args
		return nil
	})
	defer thenFunc.Release()

	catch := make(chan []js.Value)
	defer close(catch)
	catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		catch <- args
		return nil
	})
	defer catchFunc.Release()

	awaitable.Call("then", thenFunc).Call("catch", catchFunc)

	select {
	case result := <-then:
		return result, nil
	case err := <-catch:
		return nil, err
	}
}

func awaitContext(context context.Context, awaitable js.Value) ([]js.Value, []js.Value, error) {
	result := make(chan [2][]js.Value)
	defer close(result)
	go func() {
		then, catch := await(awaitable)
		select {
		case <-context.Done():
			// already done, not send anymore
		default:
			result <- [2][]js.Value{then, catch}
		}
	}()
	select {
	case <-context.Done():
		// NOTE: this case await function may still running...
		// released it later.
		return nil, nil, context.Err()
	case ret := <-result:
		return ret[0], ret[1], nil
	}
}

func addAsyncWrap(x, y int) int {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	addAsync := js.Global().Get("addAsync")
	promiss := addAsync.Invoke(x, y)
	then, catch, err := awaitContext(ctx, promiss)

	// error case
	if err != nil {
		println("await context error:", err)
		return 0
	}
	if catch != nil {
		msg := catch[0].String()
		println("catched await error:", msg)
		return 0
	}
	// truth case
	if then != nil {
		ret := then[0].Int()
		return ret
	}
	return 0
}

func addPromissWrap(x, y int) int {
	ch := make(chan int)
	defer close(ch)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	jsSuccess := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		result := args[0].Int()
		ch <- result
		return nil
	})
	defer jsSuccess.Release()
	addPromiss := js.Global().Get("addPromiss")
	go addPromiss.Invoke(x, y, jsSuccess)

	select {
	case ret := <-ch:
		return ret
	case <-ctx.Done():
		return 0
	}
}
