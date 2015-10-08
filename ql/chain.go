package ql

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
func (ch *Chain) Execute() {

	inputHandler := ch.Input
	outputHandler := ch.Output

	var chann chan Line

	inputChan := make(chan Line)

	ctx, _ := inputHandler.Handle(inputChan, ch.InputConfig)

	if ch.Filter != nil {
		filterHandler := ch.Filter
		chann = make(chan Line)
		filterHandler.Handle(ctx, inputChan, chann, ch.FilterConfig)
	} else {
		chann = inputChan
	}

	outputHandler.Handle(ctx, chann, ch.OutputConfig)

	<-ctx.Done()
}
