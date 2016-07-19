package ql

import (
	"github.com/sheenobu/quicklog/log"
	"golang.org/x/net/context"

	"time"
)

// A Chain is a series of handlers that process data
type Chain struct {
	Input       InputProcess
	InputConfig map[string]interface{}

	Parser Parser

	Filter       FilterHandler
	FilterConfig map[string]interface{}

	Output       OutputHandler
	OutputConfig map[string]interface{}
}

func (ch *Chain) parserLoop(ctx context.Context, bufferChan <-chan Buffer, inputChan chan<- Line) {

	parser := ch.Parser

	if parser == nil {
		parser = GetParser("plain")
	}

	for {
		select {
		case <-ctx.Done():
			return
		case buffer := <-bufferChan:
			l := Line{
				Data:      make(map[string]interface{}),
				Timestamp: time.Now(),
			}

			if len(buffer.Data) == 0 {
				continue // skip line
			}

			if err := parser.Parse(buffer.Data, &l, ch.InputConfig); err != nil {
				log.Log(ctx).Error("Error parsing incoming data", "error", err)
				continue
			}

			if buffer.Metadata != nil {
				for k, v := range buffer.Metadata {
					l.Data[k] = v
				}
			}

			inputChan <- l
		}
	}
}

// Execute executes the chain and waits for its completion
func (ch *Chain) Execute(ctx context.Context) {

	outputHandler := ch.Output

	var chann chan Line

	bufferChan := make(chan Buffer)
	inputChan := make(chan Line)

	if ch.InputConfig == nil {
		ch.InputConfig = make(map[string]interface{})
	}
	if ch.OutputConfig == nil {
		ch.OutputConfig = make(map[string]interface{})
	}
	if ch.FilterConfig == nil {
		ch.FilterConfig = make(map[string]interface{})
	}

	if err := ch.Input.Start(ctx, bufferChan); err != nil {
		log.Log(ctx).Crit("Error starting input handler", "error", err)
	}

	go ch.parserLoop(ctx, bufferChan, inputChan)

	if ch.Filter != nil {
		filterHandler := ch.Filter
		chann = make(chan Line)
		if err := filterHandler.Handle(ctx, inputChan, chann, ch.FilterConfig); err != nil {
			log.Log(ctx).Crit("Error creating filter handler", "error", err)
			return
		}

	} else {
		chann = inputChan
	}

	if err := outputHandler.Handle(ctx, chann, ch.OutputConfig); err != nil {
		log.Log(ctx).Crit("Error creating output handler", "error", err)
		return
	}

	<-ctx.Done()
}
