package ql

import (
	"github.com/sheenobu/quicklog/log"
	"golang.org/x/net/context"

	"time"
)

// A Chain is a series of handlers that process data
type Chain struct {
	Input       InputHandler
	InputConfig map[string]interface{}
	Parser      Parser

	Filter       FilterHandler
	FilterConfig map[string]interface{}

	Output       OutputHandler
	OutputConfig map[string]interface{}
}

// Execute executes the chain and waits for its completion
func (ch *Chain) Execute(ctx context.Context) {

	inputHandler := ch.Input
	outputHandler := ch.Output
	parser := ch.Parser

	if parser == nil {
		parser = GetParser("plain")
	}

	var chann chan Line

	bufferChan := make(chan Buffer)
	inputChan := make(chan Line)

	err := inputHandler.Handle(ctx, bufferChan, ch.InputConfig)
	if err != nil {
		log.Log(ctx).Crit("Error creating input handler", "error", err)
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case buffer := <-bufferChan:
				l := Line{
					Data:      make(map[string]interface{}),
					Timestamp: time.Now(),
				}

				if err := parser.Parse(buffer.data, &l, ch.InputConfig); err != nil {
					log.Log(ctx).Error("Error parsing incoming data", "error", err)
					continue
				}

				for k, v := range buffer.metadata {
					l.Data[k] = v
				}

				inputChan <- l
			}
		}
	}()

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
