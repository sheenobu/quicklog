package managed

import (
	"io"
)

// Writer defines the interface for writing system data
type Writer interface {

	// WriteProcess writes the process
	WriteProcess(sys *System, process *Process) error

	// WriteSystem writes the system
	WriteSystem(sys *System) error

	// Child gets the child writer
	Child() Writer
}

// TextWriter returns the plaintext writer
func TextWriter(writer io.Writer, tabs int) Writer {
	return &IoWriter{
		wr:   writer,
		tabs: tabs,
	}
}

// IoWriter is the plaintext writer that write to io.Writer
type IoWriter struct {
	wr   io.Writer
	tabs int
}

// Child returns the child IoWriter with the tab spacing shifted over by 1
func (wr *IoWriter) Child() Writer {
	return &IoWriter{
		wr:   wr.wr,
		tabs: wr.tabs + 1,
	}
}

// WriteProcess writes the process information
func (wr *IoWriter) WriteProcess(sys *System, process *Process) error {
	for i := 0; i != wr.tabs+1; i++ {
		if _, err := wr.wr.Write([]byte("\t")); err != nil {
			return err
		}
	}

	return process.Writer(process, wr.wr)
}

// WriteSystem writes the system information
func (wr *IoWriter) WriteSystem(sys *System) error {
	for i := 0; i != wr.tabs; i++ {
		if _, err := wr.wr.Write([]byte("\t")); err != nil {
			return err
		}
	}

	_, err := wr.wr.Write([]byte("System: " + sys.name + "-\n"))
	return err
}

// WriteTree recursively writes the system(s) and processes to the given writer.
func (sys *System) WriteTree(w Writer) error {
	if err := w.WriteSystem(sys); err != nil {
		return err
	}

	sys.lock.RLock()
	defer sys.lock.RUnlock()

	for _, ch := range sys.ChildrenProcs {
		if err := w.WriteProcess(sys, ch); err != nil {
			return err
		}
	}

	child := w.Child()
	for _, ch := range sys.Children {
		if err := ch.WriteTree(child); err != nil {
			return err
		}
	}

	return nil
}
