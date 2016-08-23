package stdout

import (
	"os"

	"github.com/sheenobu/quicklog/outputs/writer"
	"github.com/sheenobu/quicklog/ql"
)

func init() {
	ql.RegisterOutput("stdout", Process())
}

// Process builds the standard output process
func Process() ql.OutputHandler {
	return &writer.Process{W: os.Stdout, Name: "stdout"}
}
