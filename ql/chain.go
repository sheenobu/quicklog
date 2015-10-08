package ql

import (
	"github.com/sheenobu/golibs/log"
	"golang.org/x/net/context"
)

// A Chain is a series of handlers that process data
type Chain struct {
	Input       InputHandler
	InputConfig map[string]interface{}

	Filter       FilterHandler
	FilterConfig map[string]interface{}

	Output       OutputHandler
	OutputConfig map[string]interface{}
}

// Execute executes the chain and waits for its completion
func (ch *Chain) Execute(ctx context.Context) {

	inputHandler := ch.Input
	outputHandler := ch.Output

	var chann chan Line

	inputChan := make(chan Line)

	err := inputHandler.Handle(ctx, inputChan, ch.InputConfig)
	if err != nil {
		log.Log(ctx).Crit("Error creating input handler", "error", err)
		return
	}

	if ch.Filter != nil {
		filterHandler := ch.Filter
		chann = make(chan Line)
		err = filterHandler.Handle(ctx, inputChan, chann, ch.FilterConfig)
		if err != nil {
			log.Log(ctx).Crit("Error creating filter handler", "error", err)
			return
		}

	} else {
		chann = inputChan
	}

	err = outputHandler.Handle(ctx, chann, ch.OutputConfig)
	if err != nil {
		log.Log(ctx).Crit("Error creating output handler", "error", err)
		return
	}

	<-ctx.Done()
}
