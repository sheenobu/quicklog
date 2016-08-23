package ql

import (
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestChain(t *testing.T) {
	var ch Chain
	statusChan := make(chan int, 0)

	ch.Input = InputProcessFunc(func(ctx context.Context, out chan<- Buffer) error {
		go func() {
			out <- Buffer{Data: []byte("hello")}
			statusChan <- 1
		}()
		return nil
	})

	ch.Filter = FilterHandlerFunc(func(ctx context.Context, prev <-chan Line, next chan<- Line, _ map[string]interface{}) error {
		go func() {
			select {
			case line := <-prev:
				line.Data["message"] = strings.ToUpper(line.Data["message"].(string))
				statusChan <- 2
				next <- line
			}
		}()
		return nil
	})

	var outputValue string

	ch.Output = OutputHandlerFunc(func(ctx context.Context, in <-chan Line, _ map[string]interface{}) error {
		go func() {
			line := <-in
			outputValue = line.Data["message"].(string)
			statusChan <- 3
		}()
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case i := <-statusChan:
				if i == 3 {
					if outputValue != "HELLO" {
						t.Errorf("Expected output value to be 'HELLO', was '%s'", outputValue)
					}
					cancel()
					return
				}
			case <-ctx.Done():
				t.Error("Goroutine cancelled before output value got processed")
				return
			}
		}
	}()

	go ch.Execute(ctx)

	select {
	case <-time.After(10 * time.Second):
		t.Error("Timeout waiting for message input")
		cancel()
	case <-ctx.Done():
		return
	}

}

func TestChainNoFilter(t *testing.T) {
	var ch Chain
	statusChan := make(chan int, 0)

	ch.Input = InputProcessFunc(func(ctx context.Context, out chan<- Buffer) error {
		go func() {
			out <- Buffer{Data: []byte("hello")}
			statusChan <- 1
		}()
		return nil
	})

	var outputValue string

	ch.Output = OutputHandlerFunc(func(ctx context.Context, in <-chan Line, _ map[string]interface{}) error {
		go func() {
			line := <-in
			outputValue = line.Data["message"].(string)
			statusChan <- 3
		}()
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case i := <-statusChan:
				if i == 3 {
					if outputValue != "hello" {
						t.Errorf("Expected output value to be 'HELLO', was '%s'", outputValue)
					}
					cancel()
					return
				}
			case <-ctx.Done():
				t.Error("Goroutine cancelled before output value got processed")
				return
			}
		}
	}()

	go ch.Execute(ctx)

	select {
	case <-time.After(10 * time.Second):
		t.Error("Timeout waiting for message input")
		cancel()
	case <-ctx.Done():
		return
	}

}
