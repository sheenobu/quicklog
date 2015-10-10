package json

import (
	"bytes"
	"encoding/json"
	"github.com/sheenobu/quicklog/ql"
)

func init() {
	ql.RegisterParser("json", &JSONParser{})
}

type JSONParser struct{}

func (jp *JSONParser) Parse(buffer []byte, line *ql.Line, config map[string]interface{}) error {
	err := json.NewDecoder(bytes.NewReader(buffer)).Decode(&line.Data)
	if err != nil {
		return err
	}
	if line.Data["message"] == nil {
		line.Data["message"] = ""
	}
	return nil
}
