package ql

var (
	parsers map[string]Parser
)

func init() {
	parsers = make(map[string]Parser)
}

// Parser defines the interface for parsing incoming data
type Parser interface {
	// Parse parses the incoming buffer into the line object
	Parse(buffer []byte, line *Line, config map[string]interface{}) error
}

// ParserFunc is the adaptor for converting a function to a Parser type
type ParserFunc func([]byte, *Line, map[string]interface{}) error

// Parse parses the incoming buffer in the line object
func (pf ParserFunc) Parse(buffer []byte, line *Line, config map[string]interface{}) error {
	return pf(buffer, line, config)
}

// GetParser returns the parser registered under the name
func GetParser(name string) Parser {
	return parsers[name]
}

// RegisterParser registers a new parser
func RegisterParser(name string, parser Parser) {
	parsers[name] = parser
}
