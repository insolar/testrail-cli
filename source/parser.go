package source

import (
	"io"
)

type Parser interface {
	Parse(io.Reader) []TestEvent
}

// TestEvent go test2json event object
type TestEvent struct {
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}
