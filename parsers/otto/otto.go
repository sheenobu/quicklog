package otto

import (
	"github.com/sheenobu/quicklog/ql"

	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
)

type OttoParser struct {
	o *otto.Otto
}

func init() {
	ql.RegisterParser("otto", &OttoParser{o: otto.New()})
}

func (op *OttoParser) Parse(buffer []byte, line *ql.Line, config map[string]interface{}) error {

	script := config["otto.script"].(string)

	fn, err := op.o.Run(script)
	if err != nil {
		return err
	}

	this, err := otto.ToValue(nil)
	if err != nil {
		return err
	}

	result, err := fn.Call(this, string(buffer))
	if err != nil {
		return err
	}

	object := result.Object()

	for _, key := range object.Keys() {
		v, _ := object.Get(key)
		if v.IsPrimitive() {
			line.Data[key] = v.String()
		}
	}

	return nil
}
