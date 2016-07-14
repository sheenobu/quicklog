package otto

import (
	"github.com/sheenobu/quicklog/ql"

	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore" // enable javascript underscore support for otto
)

// Parser runs data through a javascript engine to execute parsing
type Parser struct {
	o *otto.Otto
}

func init() {
	ql.RegisterParser("otto", &Parser{o: otto.New()})
}

// Parse parses the buffer and adds it to the line struct
// the config[otto.script] should be a javascript function which
// returns a hash. Each key in the hash will be added to the line Data
// and will be indexed
func (op *Parser) Parse(buffer []byte, line *ql.Line, config map[string]interface{}) error {

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
