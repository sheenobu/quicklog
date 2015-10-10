package plain

import (
	"github.com/sheenobu/quicklog/ql"
)

func init() {
	ql.RegisterParser("plain", &PlainParser{})
}

type PlainParser struct{}

func (pp *PlainParser) Parse(buffer []byte, line *ql.Line) error {
	line.Data["message"] = string(buffer)
	return nil
}
