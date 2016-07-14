package plain

import (
	"github.com/sheenobu/quicklog/ql"
)

func init() {
	ql.RegisterParser("plain", &Parser{})
}

// Parser treats every input buffer as a single line
type Parser struct{}

// Parse adds the buffer to the Line data as the 'message' key
func (pp *Parser) Parse(buffer []byte, line *ql.Line, _ map[string]interface{}) error {
	line.Data["message"] = string(buffer)
	return nil
}
