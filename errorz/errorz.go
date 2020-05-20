package errorz

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
)

// Err ...
type Err struct {
	StackTrace []byte
	Err        error
}

func (e *Err) Error() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Error:\n %s\n", e.Err)
	fmt.Fprintf(&buf, "Trace:\n %s\n", e.StackTrace)
	return buf.String()
}

// NewErr ...
func NewErr(msg string) *Err {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	return &Err{StackTrace: buf, Err: errors.New(msg)}
}
