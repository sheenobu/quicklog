package ql

// Parser defines the interface for parsing incoming data
type Parser interface {
	// Parse parses the incoming buffer and returns a line object
	Parse(buffer []byte, line *Line) error
}

var (
	parsers map[string]Parser
)

func init() {
	parsers = make(map[string]Parser)
}

// GetParser returns the parser registered under the name
func GetParser(name string) Parser {
	return parsers[name]
}

// RegisterParser registers a new parser
func RegisterParser(name string, parser Parser) {
	parsers[name] = parser
}
