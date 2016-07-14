package json

import (
	"bytes"
	"encoding/json"

	"github.com/sheenobu/quicklog/ql"
)

func init() {
	ql.RegisterParser("json", &Parser{})
}

// Parser is a parser for JSON data.
type Parser struct{}

// Parse parses the given buffer and adds it to the line
func (jp *Parser) Parse(buffer []byte, line *ql.Line, _ map[string]interface{}) error {
	err := json.NewDecoder(bytes.NewReader(buffer)).Decode(&line.Data)
	if err != nil {
		return err
	}
	if line.Data["message"] == nil {
		line.Data["message"] = ""
	}
	return nil
}
